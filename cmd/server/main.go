package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/arganaphang/notification/internal/dto"
	"github.com/arganaphang/notification/internal/repository"
	"github.com/arganaphang/notification/internal/service"
	"github.com/arganaphang/notification/pkg/sse"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type App struct {
	SSE      *sse.SSEConn
	Services service.Services
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	gin.Default().SetTrustedProxies(nil)

	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln(err)
	}

	repositories := repository.Repositories{
		NotificationRepository: repository.NewNotification(db),
	}
	services := service.Services{
		NotificationService: service.NewNotification(repositories),
	}
	app := App{
		SSE:      sse.NewSSEConn(),
		Services: services,
	}

	e := gin.New()
	e.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.1:3000"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	e.GET("/healthz", app.health)

	// Notification
	notification := e.Group("/notification")
	notification.GET("", app.getNotifications)
	notification.POST("", app.createNotification)
	notification.GET(":id", app.getNotificationByID)
	notification.PUT(":id", app.readNotificationByID)
	notification.GET("count/:user_id", app.countNotificationsByUserID)
	notification.GET("count-watch/:user_id", app.watchCountNotificationsByUserID)

	logrus.Info("Server running at http://127.0.0.1:8000")
	defer catch()
	e.Run("0.0.0.0:8000")
}
func catch() {
	if r := recover(); r != nil {
		fmt.Println("Error occured", r)
	} else {
		fmt.Println("Application running perfectly")
	}
}

func (a *App) health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

// Notification
func (a *App) getNotifications(ctx *gin.Context) {
	var params dto.GetNotificationsRequest
	if err := ctx.BindQuery(&params); err != nil {
		ctx.JSON(http.StatusExpectationFailed, gin.H{
			"message": "failed to serialize query",
		})
		return
	}
	result, err := a.Services.NotificationService.GetNotifications(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusExpectationFailed, gin.H{
			"message": "failed to get notifications",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "get notifications success",
		"data":    result,
	})
}

func (a *App) createNotification(ctx *gin.Context) {
	var params dto.CreateNotificationRequest
	if err := ctx.ShouldBind(&params); err != nil {
		ctx.JSON(http.StatusExpectationFailed, gin.H{
			"message": "failed to serialize request body",
		})
		return
	}
	err := a.Services.NotificationService.CreateNotification(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusExpectationFailed, gin.H{
			"message": "failed to create notifications",
		})
		return
	}

	a.SSE.Broadcast(params.UserID.String(), "create notification")

	ctx.JSON(http.StatusOK, gin.H{
		"message": "create notifications success",
	})
}

func (a *App) getNotificationByID(ctx *gin.Context) {
	pathID := ctx.Param("id")
	id, err := uuid.Parse(pathID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "please add query param `id`",
		})
		return
	}

	result, err := a.Services.NotificationService.GetNotificationByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusExpectationFailed, gin.H{
			"message": "failed to get notification detail",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "get notification detail success",
		"data":    result,
	})
}

func (a *App) readNotificationByID(ctx *gin.Context) {
	pathID := ctx.Param("id")
	notificationID, err := uuid.Parse(pathID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "please add query param `user_id`",
		})
		return
	}

	err = a.Services.NotificationService.ReadNotificationByID(ctx, notificationID)
	if err != nil {
		ctx.JSON(http.StatusExpectationFailed, gin.H{
			"message": "failed to read notification",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "read notification success",
	})
}

func (a *App) countNotificationsByUserID(ctx *gin.Context) {
	pathID := ctx.Param("user_id")
	userID, err := uuid.Parse(pathID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "please add query param `user_id`",
		})
		return
	}

	result, err := a.Services.NotificationService.CountNotificationsByUserID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusExpectationFailed, gin.H{
			"message": "failed to get count notifications",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "get count notifications success",
		"data":    result,
	})
}

func (a *App) watchCountNotificationsByUserID(ctx *gin.Context) {
	pathID := ctx.Param("user_id")
	userID, err := uuid.Parse(pathID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "please add query param `user_id`",
		})
		return
	}

	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("Transfer-Encoding", "chunked")
	// Logic
	ch := a.SSE.AddClient(userID.String())
	go func() {
		a.SSE.Broadcast(userID.String(), "initial")
	}()

	// Sent initial response
	i := 0
	ctx.Stream(func(w io.Writer) bool {
		for {
			select {
			case message := <-*ch:
				result, err := a.Services.NotificationService.CountNotificationsByUserID(ctx, userID)
				if err != nil {
					// logrus.Info(err.Error()) // SKIP ERROR
					return true
				}
				ctx.SSEvent("message", gin.H{
					"message": fmt.Sprintf("notification: [%s]", message),
					"count":   result,
				})
				i++
				return true
			case <-ctx.Writer.CloseNotify(): // Act as defer function for http handler :man_shrugging: -> ref [https://github.com/gin-gonic/gin/issues/515#issuecomment-176434018]
				a.SSE.RemoveClient(userID.String(), *ch)
				return false
			}
		}
	})
}

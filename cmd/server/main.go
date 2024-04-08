package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/arganaphang/notification/pkg/sse"
	"github.com/gin-gonic/gin"
)

type App struct {
	SSE *sse.SSEConn
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	gin.Default().SetTrustedProxies(nil)
	e := gin.New()

	app := App{
		SSE: sse.NewSSEConn(),
	}

	e.GET("/healthz", app.health)
	e.GET("/get/:id", app.get)
	e.GET("/send/:id", app.send)

	e.Run("0.0.0.0:8000")
}

func (a *App) health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func (a *App) get(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": "please add query param `id`",
		})
		return
	}

	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("Transfer-Encoding", "chunked")
	// Logic
	ch := a.SSE.AddClient(id)
	defer a.SSE.RemoveClient(id, *ch)
	ctx.Stream(func(w io.Writer) bool {
		for {
			select {
			case message := <-*ch:
				ctx.SSEvent("message", message)
				return true
			case <-ctx.Done():
				fmt.Println("Client closed connection")
				return false
			}
		}

	})
}

func (a *App) send(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": "please add query param `id`",
		})
		return
	}
	msg := ctx.Query("msg")
	if msg == "" {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"message": "please add query param `msg`",
		})
		return
	}
	a.SSE.Broadcast(id, msg)
	ctx.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("Message sent to: %s, data %s", id, msg),
	})
}

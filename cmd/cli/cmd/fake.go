package cmd

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/arganaphang/notification/internal/dto"
	"github.com/arganaphang/notification/internal/repository"
	"github.com/arganaphang/notification/internal/service"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"

	"github.com/go-faker/faker/v4"
)

type App struct {
	Services service.Services
}

type Notification struct {
	Title   string    `json:"title" faker:"sentence"`
	Content string    `json:"content" faker:"paragraph"`
	UserID  uuid.UUID `json:"user_id" faker:"-"`
	OrderID int8      `json:"order_id"`
}

var fakeCmd = &cobra.Command{
	Use:   "fake",
	Short: "Create fake notification",
	Long:  `Create fake notification for spesific user_id`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("please add user id")
		}
		_, err := strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			return errors.New("user must be uuid")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		n, _ := strconv.ParseUint(args[1], 10, 64)
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
			Services: services,
		}
		for i := 0; i < int(n); i++ {
			notification := Notification{}
			err := faker.FakeData(&notification)
			if err != nil {
				log.Fatalln(err.Error())
			}
			err = app.Services.NotificationService.CreateNotification(cmd.Context(), dto.CreateNotificationRequest{
				Title:   notification.Title,
				Content: notification.Content,
				OrderID: int(notification.OrderID),
				UserID:  uuid.New(),
			})
			if err != nil {
				log.Fatalln(err.Error())
			}
		}
		log.Println("create fake notification successfully")
	},
}

func init() {
	rootCmd.AddCommand(fakeCmd)
	fakeCmd.Flags().UintP("n", "", 10, "Total Fake Notification")
}

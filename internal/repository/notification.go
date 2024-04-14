package repository

import (
	"context"
	"time"

	"github.com/arganaphang/notification/internal/dto"
	"github.com/arganaphang/notification/internal/model"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type INotificationRepository interface {
	GetNotifications(ctx context.Context, params dto.GetNotificationsRequest) ([]model.Notification, error)
	GetNotificationByID(ctx context.Context, id uuid.UUID) (*model.Notification, error)
	CountNotificationsByUserID(ctx context.Context, userID uuid.UUID) (*uint, error)
	ReadNotificationByID(ctx context.Context, id uuid.UUID) error
	CreateNotification(ctx context.Context, data dto.CreateNotificationRequest) error
}

type NotificationRepository struct {
	DB *sqlx.DB
}

func NewNotification(db *sqlx.DB) INotificationRepository {
	return &NotificationRepository{DB: db}
}

func (r *NotificationRepository) GetNotifications(ctx context.Context, params dto.GetNotificationsRequest) ([]model.Notification, error) {
	sql := goqu.From("notifications")

	if params.UserID != nil {
		sql = sql.Where(goqu.Ex{
			"user_id": params.UserID,
		})
	}

	query, _, err := sql.Order(goqu.I("created_at").Desc()).ToSQL()
	if err != nil {
		return nil, err
	}

	var notifications []model.Notification
	rows, err := r.DB.Queryx(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var notification model.Notification
		err = rows.StructScan(&notification)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}
	return notifications, nil
}

func (r *NotificationRepository) GetNotificationByID(ctx context.Context, id uuid.UUID) (*model.Notification, error) {
	query, _, err := goqu.
		From("notifications").
		Where(goqu.Ex{"id": id}).
		Limit(1).
		ToSQL()
	if err != nil {
		return nil, err
	}

	var notification model.Notification
	err = r.DB.QueryRowx(query).StructScan(&notification)
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

func (r *NotificationRepository) ReadNotificationByID(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	query, _, err := goqu.
		Update("notifications").
		Set(goqu.Record{"is_read": "t", "updated_at": now}).
		Where(goqu.Ex{"id": id}).
		ToSQL()
	if err != nil {
		return err
	}

	_, err = r.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (r *NotificationRepository) CountNotificationsByUserID(ctx context.Context, userID uuid.UUID) (*uint, error) {
	var counter *uint
	query, _, err := goqu.
		From("notifications").
		Select(goqu.COUNT("*")).
		Where(goqu.Ex{"user_id": userID, "is_read": "f"}).
		ToSQL()
	if err != nil {
		return nil, err
	}
	err = r.DB.QueryRow(query).Scan(&counter)
	if err != nil {
		return nil, err
	}
	return counter, nil
}

func (r *NotificationRepository) CreateNotification(ctx context.Context, data dto.CreateNotificationRequest) error {
	_, err := r.DB.Exec(`INSERT INTO notifications (title, content, user_id, order_id) VALUES ($1, $2, $3, $4)`, data.Title, data.Content, data.UserID, data.OrderID)
	if err != nil {
		return err
	}
	return nil
}

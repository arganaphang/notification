package dto

import "github.com/google/uuid"

type CreateNotificationRequest struct {
	Title   string    `json:"title"`
	Content string    `json:"content"`
	UserID  uuid.UUID `json:"user_id"`
	OrderID int       `json:"order_id"`
}

type GetNotificationsRequest struct {
	UserID *string `json:"user_id" form:"user_id"`
}

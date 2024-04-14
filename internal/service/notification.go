package service

import (
	"context"

	"github.com/arganaphang/notification/internal/dto"
	"github.com/arganaphang/notification/internal/model"
	"github.com/arganaphang/notification/internal/repository"
	"github.com/google/uuid"
)

type INotificationService interface {
	GetNotifications(ctx context.Context, params dto.GetNotificationsRequest) ([]model.Notification, error)
	GetNotificationByID(ctx context.Context, id uuid.UUID) (*model.Notification, error)
	CountNotificationsByUserID(ctx context.Context, userID uuid.UUID) (*uint, error)
	ReadNotificationByID(ctx context.Context, id uuid.UUID) error
	CreateNotification(ctx context.Context, data dto.CreateNotificationRequest) error
}

type NotificationService struct {
	repostiory repository.Repositories
}

func NewNotification(repostiory repository.Repositories) INotificationService {
	return &NotificationService{repostiory: repostiory}
}

func (s *NotificationService) GetNotifications(ctx context.Context, params dto.GetNotificationsRequest) ([]model.Notification, error) {
	return s.repostiory.NotificationRepository.GetNotifications(ctx, params)
}

func (s *NotificationService) GetNotificationByID(ctx context.Context, id uuid.UUID) (*model.Notification, error) {
	return s.repostiory.NotificationRepository.GetNotificationByID(ctx, id)
}
func (s *NotificationService) ReadNotificationByID(ctx context.Context, id uuid.UUID) error {
	return s.repostiory.NotificationRepository.ReadNotificationByID(ctx, id)
}

func (s *NotificationService) CountNotificationsByUserID(ctx context.Context, userID uuid.UUID) (*uint, error) {
	return s.repostiory.NotificationRepository.CountNotificationsByUserID(ctx, userID)
}

func (s *NotificationService) CreateNotification(ctx context.Context, data dto.CreateNotificationRequest) error {
	return s.repostiory.NotificationRepository.CreateNotification(ctx, data)
}

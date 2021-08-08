package services

import (
	"story-app-monolith/domain"
	"story-app-monolith/repo"
)

type NotificationService interface {
	GetAllUnreadNotificationByUsername(string) (*[]domain.Notification, error)
}

type DefaultNotificationService struct {
	repo repo.NotificationRepo
}

func (n DefaultNotificationService) GetAllUnreadNotificationByUsername(owner string) (*[]domain.Notification, error) {
	notifications, err := n.repo.GetAllUnreadNotificationByUsername(owner)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func NewNotificationService(repository repo.NotificationRepo) DefaultNotificationService {
	return DefaultNotificationService{repository}
}

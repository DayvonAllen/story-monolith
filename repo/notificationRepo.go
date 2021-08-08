package repo

import "story-app-monolith/domain"

type NotificationRepo interface {
	GetAllUnreadNotificationByUsername(string) (*[]domain.Notification, error)
}


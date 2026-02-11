package services

import "github.com/Sistal/bff-orchestrator/internal/models"

type NotificationService interface {
	GetNotifications(userID string) ([]models.Notification, error)
	MarkAsRead(userID, id string) error
	MarkAllAsRead(userID string) error
}

type MockNotificationService struct{}

func NewMockNotificationService() NotificationService {
	return &MockNotificationService{}
}

func (s *MockNotificationService) GetNotifications(userID string) ([]models.Notification, error) {
	return []models.Notification{
		{
			ID:      "NOT-001",
			Titulo:  "Solicitud Aprobada",
			Mensaje: "Tu solicitud de cambio ha sido aprobada.",
			Leida:   false,
			Fecha:   "2025-02-01",
		},
	}, nil
}

func (s *MockNotificationService) MarkAsRead(userID, id string) error {
	return nil
}

func (s *MockNotificationService) MarkAllAsRead(userID string) error {
	return nil
}

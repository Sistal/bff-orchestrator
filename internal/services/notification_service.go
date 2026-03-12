package services

import "github.com/Sistal/bff-orchestrator/internal/models"

type NotificationService interface {
	GetNotifications(userID string) ([]models.Notification, error)
	MarkAsRead(userID, id string) error
	// MarkAllAsRead retorna el número de notificaciones marcadas como leídas
	MarkAllAsRead(userID string) (int, error)
}

// MockNotificationService — implementación mock para desarrollo/testing
type MockNotificationService struct{}

func NewMockNotificationService() NotificationService {
	return &MockNotificationService{}
}

func (s *MockNotificationService) GetNotifications(userID string) ([]models.Notification, error) {
	return []models.Notification{
		{
			ID:          "NOT-001",
			Type:        "approved",
			Title:       "Solicitud Aprobada",
			Message:     "Tu solicitud de reposición ha sido aprobada.",
			Timestamp:   "2026-03-08T10:00:00Z",
			IsRead:      false,
			ActionLabel: "Ver detalle",
		},
		{
			ID:          "NOT-002",
			Type:        "delivery",
			Title:       "Entrega en Camino",
			Message:     "Tu uniforme está en camino. Código de seguimiento: 8956234712.",
			Timestamp:   "2026-03-07T14:30:00Z",
			IsRead:      true,
			ActionLabel: "Seguir entrega",
		},
	}, nil
}

func (s *MockNotificationService) MarkAsRead(userID, id string) error {
	return nil
}

func (s *MockNotificationService) MarkAllAsRead(userID string) (int, error) {
	return 1, nil // mock: 1 notificación marcada
}

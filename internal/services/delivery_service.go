package services

import "github.com/Sistal/bff-orchestrator/internal/models"

type DeliveryService interface {
	GetDeliveries(userID string) ([]models.DeliverySummary, error)
	GetDeliveryByID(userID, id string) (*models.DeliverySummary, error)
	ConfirmDelivery(userID, id string) (*models.DeliverySummary, error)
	// GetUpcomingDeliveries retorna entregas próximas para el dashboard
	GetUpcomingDeliveries(userID string) ([]models.DeliverySummary, error)
}

// MockDeliveryService — implementación mock para desarrollo/testing
type MockDeliveryService struct{}

func NewMockDeliveryService() DeliveryService {
	return &MockDeliveryService{}
}

func (s *MockDeliveryService) GetDeliveries(userID string) ([]models.DeliverySummary, error) {
	return []models.DeliverySummary{
		{
			ID:           "DEL-001",
			RequestID:    "SOL-2024-1547",
			DispatchDate: "2024-12-13",
			Garments:     "Chaqueta, Polera",
			Address:      "Sucursal Centro",
			Status:       "in-transit",
			TrackingCode: "8956234712",
			Type:         "full-uniform",
		},
	}, nil
}

func (s *MockDeliveryService) GetDeliveryByID(userID, id string) (*models.DeliverySummary, error) {
	return &models.DeliverySummary{
		ID:           id,
		RequestID:    "SOL-2024-1547",
		DispatchDate: "2024-12-13",
		Garments:     "Chaqueta, Polera",
		Address:      "Sucursal Centro",
		Status:       "in-transit",
		TrackingCode: "8956234712",
		Type:         "full-uniform",
		Timeline: []models.TimelineEvent{
			{Status: "Creado", Date: "2024-12-10", Completed: true},
			{Status: "En tránsito", Date: "2024-12-13", Completed: true},
			{Status: "Entregado", Date: "", Completed: false},
		},
	}, nil
}

func (s *MockDeliveryService) ConfirmDelivery(userID, id string) (*models.DeliverySummary, error) {
	return &models.DeliverySummary{
		ID:           id,
		RequestID:    "SOL-2024-1547",
		DispatchDate: "2024-12-13",
		Garments:     "Chaqueta, Polera",
		Address:      "Sucursal Centro",
		Status:       "delivered",
		TrackingCode: "8956234712",
		Type:         "full-uniform",
	}, nil
}

func (s *MockDeliveryService) GetUpcomingDeliveries(userID string) ([]models.DeliverySummary, error) {
	return []models.DeliverySummary{
		{
			ID:           "DEL-002",
			RequestID:    "SOL-2026-0012",
			DispatchDate: "2026-03-15",
			Garments:     "Pantalón, Polera",
			Address:      "Casa Matriz Santiago",
			Status:       "in-transit",
			TrackingCode: "1234567890",
			Type:         "reposicion",
		},
	}, nil
}

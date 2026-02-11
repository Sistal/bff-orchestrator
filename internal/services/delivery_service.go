package services

import "github.com/Sistal/bff-orchestrator/internal/models"

type DeliveryService interface {
	GetDeliveries(userID string) ([]models.DeliverySummary, error)
	GetDeliveryByID(userID, id string) (*models.DeliverySummary, error)
	ConfirmDelivery(userID, id string) (*models.DeliverySummary, error)
}

type MockDeliveryService struct{}

func NewMockDeliveryService() DeliveryService {
	return &MockDeliveryService{}
}

func (s *MockDeliveryService) GetDeliveries(userID string) ([]models.DeliverySummary, error) {
	return []models.DeliverySummary{
		{
			ID:           "DEL-001",
			RequestID:    "SOL-2024-1547",
			Status:       "in-transit",
			TrackingCode: "8956234712",
			Type:         "full-uniform",
			Garments:     "Chaqueta",
		},
	}, nil
}

func (s *MockDeliveryService) GetDeliveryByID(userID, id string) (*models.DeliverySummary, error) {
	return &models.DeliverySummary{
		ID:           id,
		RequestID:    "SOL-2024-1547",
		Status:       "in-transit",
		TrackingCode: "8956234712",
		Type:         "full-uniform",
		Garments:     "Chaqueta",
	}, nil
}

func (s *MockDeliveryService) ConfirmDelivery(userID, id string) (*models.DeliverySummary, error) {
	return &models.DeliverySummary{
		ID:           id,
		RequestID:    "SOL-2024-1547",
		Status:       "Entregado",
		TrackingCode: "8956234712",
		Type:         "full-uniform",
		Garments:     "Chaqueta",
	}, nil
}

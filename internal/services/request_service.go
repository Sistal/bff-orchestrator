package services

import (
	"mime/multipart"

	"github.com/Sistal/bff-orchestrator/internal/models"
)

type RequestService interface {
	GetRequests(userID string, params map[string]string) ([]models.RequestSummary, error)
	GetRequestByID(userID string, id string) (*models.RequestSummary, error)
	GetRecentRequests(userID string) ([]models.RequestSummary, error)
	CreateReplenishmentRequest(userID string, req models.CreateReplenishmentRequest) (*models.RequestSummary, error)
	CreateGarmentChangeRequest(userID string, req models.CreateGarmentChangeRequest) (*models.RequestSummary, error)
	CreateUniformRequest(userID string, req models.CreateUniformRequest) (*models.RequestSummary, error)
	UploadFile(userID string, file *multipart.FileHeader) (*models.FileUploadResponse, error)
}

type MockRequestService struct{}

func NewMockRequestService() RequestService {
	return &MockRequestService{}
}

func (s *MockRequestService) GetRequests(userID string, params map[string]string) ([]models.RequestSummary, error) {
	return []models.RequestSummary{
		{
			ID:     "SOL-001",
			Tipo:   "Reposición",
			Fecha:  "2025-01-15",
			Estado: "En Proceso",
			Items:  []string{"Polera"},
			Motivo: "Desgaste",
		},
	}, nil
}

func (s *MockRequestService) GetRequestByID(userID string, id string) (*models.RequestSummary, error) {
	return &models.RequestSummary{
		ID:     id,
		Tipo:   "Reposición",
		Fecha:  "2025-01-15",
		Estado: "En Proceso",
		Items:  []string{"Polera"},
		Motivo: "Desgaste",
	}, nil
}

func (s *MockRequestService) GetRecentRequests(userID string) ([]models.RequestSummary, error) {
	return []models.RequestSummary{
		{
			ID:     "SOL-001",
			Tipo:   "Reposición",
			Fecha:  "2026-03-01",
			Estado: "En Proceso",
			Items:  []string{"Polera", "Pantalón"},
			Motivo: "Desgaste",
		},
	}, nil
}

func (s *MockRequestService) CreateReplenishmentRequest(userID string, req models.CreateReplenishmentRequest) (*models.RequestSummary, error) {
	return &models.RequestSummary{
		ID:     "SOL-NEW-001",
		Tipo:   "Reposición",
		Fecha:  "2026-02-04",
		Estado: "Creado",
		Items:  req.Items,
		Motivo: req.Reason,
	}, nil
}

func (s *MockRequestService) CreateGarmentChangeRequest(userID string, req models.CreateGarmentChangeRequest) (*models.RequestSummary, error) {
	return &models.RequestSummary{}, nil
}

func (s *MockRequestService) CreateUniformRequest(userID string, req models.CreateUniformRequest) (*models.RequestSummary, error) {
	return &models.RequestSummary{
		ID:     "SOL-NUEVA",
		Tipo:   "Uniforme",
		Fecha:  "2025-01-01",
		Estado: "Pendiente",
	}, nil
}

func (s *MockRequestService) UploadFile(userID string, file *multipart.FileHeader) (*models.FileUploadResponse, error) {
	return &models.FileUploadResponse{
		FileID: "FILE-123",
		URL:    "https://storage.example.com/file-123.jpg",
		Name:   file.Filename,
	}, nil
}

package services

import (
	"github.com/Sistal/bff-orchestrator/internal/clients"
	"github.com/Sistal/bff-orchestrator/internal/models"
)

type BranchService interface {
	GetAllBranches() ([]models.Branch, error)
	GetChangeHistory() ([]models.BranchChangeRequestHistory, error)
	CreateChangeRequest(req models.CreateBranchChangeRequest) (*models.CreateBranchChangeRequest, error)
}

type MockBranchService struct{}

func NewMockBranchService() BranchService {
	return &MockBranchService{}
}

func (s *MockBranchService) GetAllBranches() ([]models.Branch, error) {
	return []models.Branch{
		{ID: 1, Name: "Casa Matriz - Santiago", Region: "RM"},
		{ID: 2, Name: "Sucursal Norte", Region: "II"},
	}, nil
}

func (s *MockBranchService) GetChangeHistory() ([]models.BranchChangeRequestHistory, error) {
	return []models.BranchChangeRequestHistory{
		{
			ID:               1,
			FechaSolicitud:   "2023-11-15",
			FechaEfectiva:    "2024-01-01",
			SucursalAnterior: "Sucursal Norte",
			SucursalNueva:    "Casa Matriz",
			Motivo:           "Traslado personal",
			Estado:           "Aprobado",
		},
	}, nil
}

func (s *MockBranchService) CreateChangeRequest(req models.CreateBranchChangeRequest) (*models.CreateBranchChangeRequest, error) {
	return &req, nil
}

type HTTPBranchService struct {
	client *clients.HRClient
}

func NewHTTPBranchService(client *clients.HRClient) BranchService {
	return &HTTPBranchService{client: client}
}

func (s *HTTPBranchService) GetAllBranches() ([]models.Branch, error) {
	return s.client.GetBranches()
}

func (s *HTTPBranchService) GetChangeHistory() ([]models.BranchChangeRequestHistory, error) {
	return s.client.GetChangeHistory()
}

func (s *HTTPBranchService) CreateChangeRequest(req models.CreateBranchChangeRequest) (*models.CreateBranchChangeRequest, error) {
	err := s.client.CreateBranchChangeRequest(req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

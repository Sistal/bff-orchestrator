package services

import (
	"github.com/Sistal/bff-orchestrator/internal/clients"
	"github.com/Sistal/bff-orchestrator/internal/models"
)

type BranchService interface {
	GetAllBranches() ([]models.Branch, error)
	// employeeID es el id_funcionario resuelto por el middleware desde la cookie.
	GetChangeHistory(employeeID string) ([]models.BranchChangeRequestHistory, error)
	// CreateChangeRequest retorna el historial creado, no el request de entrada.
	// employeeID es el id_funcionario resuelto por el middleware desde la cookie.
	CreateChangeRequest(employeeID string, req models.CreateBranchChangeRequest) (*models.BranchChangeRequestHistory, error)
}

// MockBranchService — implementación mock para desarrollo/testing
type MockBranchService struct{}

func NewMockBranchService() BranchService {
	return &MockBranchService{}
}

func (s *MockBranchService) GetAllBranches() ([]models.Branch, error) {
	return []models.Branch{
		{ID: 1, Name: "Casa Matriz - Santiago", Direccion: "Av. Apoquindo 4800, Las Condes"},
		{ID: 2, Name: "Sucursal Norte", Direccion: "Av. Grecia 750, Antofagasta"},
	}, nil
}

func (s *MockBranchService) GetChangeHistory(employeeID string) ([]models.BranchChangeRequestHistory, error) {
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

func (s *MockBranchService) CreateChangeRequest(employeeID string, req models.CreateBranchChangeRequest) (*models.BranchChangeRequestHistory, error) {
	return &models.BranchChangeRequestHistory{
		ID:               99,
		FechaSolicitud:   "2026-03-12",
		FechaEfectiva:    req.EffectiveDate,
		SucursalAnterior: "",
		SucursalNueva:    "",
		Motivo:           req.Reason,
		Estado:           "Pendiente",
	}, nil
}

// HTTPBranchService — implementación real que delega al ms-funcionario
type HTTPBranchService struct {
	client *clients.HRClient
}

func NewHTTPBranchService(client *clients.HRClient) BranchService {
	return &HTTPBranchService{client: client}
}

func (s *HTTPBranchService) GetAllBranches() ([]models.Branch, error) {
	return s.client.GetBranches()
}

func (s *HTTPBranchService) GetChangeHistory(employeeID string) ([]models.BranchChangeRequestHistory, error) {
	return s.client.GetChangeHistory(employeeID)
}

func (s *HTTPBranchService) CreateChangeRequest(employeeID string, req models.CreateBranchChangeRequest) (*models.BranchChangeRequestHistory, error) {
	return s.client.CreateBranchChangeRequest(employeeID, req)
}

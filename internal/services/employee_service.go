package services

import (
	"github.com/Sistal/bff-orchestrator/internal/clients"
	"github.com/Sistal/bff-orchestrator/internal/models"
)

type EmployeeService interface {
	GetProfile(id string) (*models.EmployeeProfile, error)
	UpdateContact(id string, req models.UpdateContactRequest) (*models.EmployeeProfile, error)
	GetStats(id string) (*models.HomeStats, error)
	GetMeasurements(id string) (*models.BodyMeasurements, error)
	// Added back for interface compliance
	UpdatePreferences(id string, req models.UpdatePreferencesRequest) error
	UpdateSecurity(id string) error
	GetActivity(id string) ([]models.ActivityLog, error)
	RegisterMeasurements(id string, req models.BodyMeasurements) (*models.BodyMeasurements, error)
}

type MockEmployeeService struct{}

func NewMockEmployeeService() EmployeeService {
	return &MockEmployeeService{}
}

func (s *MockEmployeeService) GetProfile(id string) (*models.EmployeeProfile, error) {
	return &models.EmployeeProfile{
		ID:      123,
		Nombre:  "Juan Pérez",
		Email:   "juan.perez@empresa.com",
		Cargo:   "Funcionario",
		Celular: "+56912345678",
	}, nil
}

func (s *MockEmployeeService) UpdateContact(id string, req models.UpdateContactRequest) (*models.EmployeeProfile, error) {
	return &models.EmployeeProfile{
		ID:      123,
		Nombre:  "Juan Pérez",
		Email:   req.Email,
		Cargo:   "Funcionario",
		Celular: req.Celular,
	}, nil
}

func (s *MockEmployeeService) GetStats(id string) (*models.HomeStats, error) {
	return &models.HomeStats{
		SolicitudesPendientes: 2,
		EntregasProximas:      1,
	}, nil
}

func (s *MockEmployeeService) GetMeasurements(id string) (*models.BodyMeasurements, error) {
	return &models.BodyMeasurements{
		EstaturaM: 1.75,
		PechoCm:   100,
	}, nil
}

type HTTPEmployeeService struct {
	hrClient  *clients.HRClient
	opsClient *clients.OpsClient
}

func NewHTTPEmployeeService(hr *clients.HRClient, ops *clients.OpsClient) EmployeeService {
	return &HTTPEmployeeService{hrClient: hr, opsClient: ops}
}

func (s *HTTPEmployeeService) GetProfile(id string) (*models.EmployeeProfile, error) {
	return s.hrClient.GetEmployeeProfile(id)
}

func (s *HTTPEmployeeService) UpdateContact(id string, req models.UpdateContactRequest) (*models.EmployeeProfile, error) {
	err := s.hrClient.UpdateContact(id, req)
	if err != nil {
		return nil, err
	}
	// Return updated profile (could fetch it again or fake it)
	return s.GetProfile(id)
}

func (s *HTTPEmployeeService) GetStats(id string) (*models.HomeStats, error) {
	return s.opsClient.GetDashboardStats(id)
}

func (s *HTTPEmployeeService) GetMeasurements(id string) (*models.BodyMeasurements, error) {
	return s.hrClient.GetMeasurements(id)
}

func (s *HTTPEmployeeService) UpdatePreferences(id string, req models.UpdatePreferencesRequest) error {
	return nil // Not implemented
}

func (s *HTTPEmployeeService) UpdateSecurity(id string) error {
	return nil // Not implemented
}

func (s *HTTPEmployeeService) GetActivity(id string) ([]models.ActivityLog, error) {
	return []models.ActivityLog{}, nil
}

func (s *HTTPEmployeeService) RegisterMeasurements(id string, req models.BodyMeasurements) (*models.BodyMeasurements, error) {
	return &req, nil
}

// Ensure Mock also implements the new methods signature
func (s *MockEmployeeService) UpdatePreferences(id string, req models.UpdatePreferencesRequest) error {
	return nil
}
func (s *MockEmployeeService) UpdateSecurity(id string) error {
	return nil
}
func (s *MockEmployeeService) GetActivity(id string) ([]models.ActivityLog, error) {
	return []models.ActivityLog{}, nil
}
func (s *MockEmployeeService) RegisterMeasurements(id string, req models.BodyMeasurements) (*models.BodyMeasurements, error) {
	return &req, nil
}

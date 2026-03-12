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
	RegisterMeasurements(id string, req models.BodyMeasurements) (*models.BodyMeasurements, error)
	UpdateMeasurements(id string, req models.BodyMeasurements) (*models.BodyMeasurements, error)
	// UpdatePreferences: mock temporal sin persistencia — retorna el objeto enviado
	UpdatePreferences(id string, req models.UpdatePreferencesRequest) (*models.UpdatePreferencesRequest, error)
	// UpdateSecurity: sin persistencia en DDL actual — acepta body y retorna éxito
	UpdateSecurity(id string, req models.UpdateSecurityRequest) error
	GetActivity(id string) ([]models.ActivityLog, error)
}

// MockEmployeeService — implementación mock para desarrollo/testing
type MockEmployeeService struct{}

func NewMockEmployeeService() EmployeeService {
	return &MockEmployeeService{}
}

func (s *MockEmployeeService) GetProfile(id string) (*models.EmployeeProfile, error) {
	return &models.EmployeeProfile{
		ID:              123,
		RutFuncionario:  "12345678-9",
		Nombres:         "Juan",
		ApellidoPaterno: "Pérez",
		ApellidoMaterno: "González",
		Email:           "juan.perez.contacto@empresa.com",
		Celular:         "+56912345678",
		Telefono:        "+5622345678",
		Direccion:       "Av. Providencia 1234, Santiago",
		Cargo: models.CargoRef{
			ID:          1,
			NombreCargo: "Operario",
		},
		Sucursal: models.SucursalRef{
			ID:             1,
			NombreSucursal: "Casa Matriz Santiago",
			Direccion:      "Av. Apoquindo 4800",
		},
		Estado: models.EstadoRef{
			ID:           1,
			NombreEstado: "Activo",
		},
		Preferences: models.UserPreferences{
			Notifications: models.NotificationPreferences{
				Email: true,
				Push:  true,
				SMS:   false,
			},
			Theme: "light",
		},
	}, nil
}

func (s *MockEmployeeService) UpdateContact(id string, req models.UpdateContactRequest) (*models.EmployeeProfile, error) {
	profile, _ := s.GetProfile(id)
	profile.Nombres = req.Nombres
	profile.ApellidoPaterno = req.ApellidoPaterno
	profile.ApellidoMaterno = req.ApellidoMaterno
	profile.Celular = req.Celular
	profile.Telefono = req.Telefono
	profile.Email = req.Email
	profile.Direccion = req.Direccion
	return profile, nil
}

func (s *MockEmployeeService) GetStats(id string) (*models.HomeStats, error) {
	return &models.HomeStats{
		UserID:                123,
		TotalSolicitudes:      5,
		SolicitudesPendientes: 2,
		EntregasProximas:      1,
	}, nil
}

func (s *MockEmployeeService) GetMeasurements(id string) (*models.BodyMeasurements, error) {
	return &models.BodyMeasurements{
		ID:            1,
		FuncionarioID: 123,
		EstaturaM:     1.75,
		PechoCm:       100.0,
		CinturaCm:     85.0,
		CaderaCm:      95.0,
		MangaCm:       60.0,
		FechaInicio:   "2024-01-01",
		FechaFin:      nil,
		Activa:        true,
	}, nil
}

func (s *MockEmployeeService) RegisterMeasurements(id string, req models.BodyMeasurements) (*models.BodyMeasurements, error) {
	req.Activa = req.FechaFin == nil
	return &req, nil
}

func (s *MockEmployeeService) UpdateMeasurements(id string, req models.BodyMeasurements) (*models.BodyMeasurements, error) {
	req.Activa = req.FechaFin == nil
	return &req, nil
}

func (s *MockEmployeeService) UpdatePreferences(id string, req models.UpdatePreferencesRequest) (*models.UpdatePreferencesRequest, error) {
	// Mock temporal: retorna el objeto enviado sin persistencia
	return &req, nil
}

func (s *MockEmployeeService) UpdateSecurity(id string, req models.UpdateSecurityRequest) error {
	return nil
}

func (s *MockEmployeeService) GetActivity(id string) ([]models.ActivityLog, error) {
	return []models.ActivityLog{}, nil
}

// HTTPEmployeeService — implementación real que delega al ms-funcionario
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
	return s.GetProfile(id)
}

func (s *HTTPEmployeeService) GetStats(id string) (*models.HomeStats, error) {
	return s.opsClient.GetDashboardStats(id)
}

func (s *HTTPEmployeeService) GetMeasurements(id string) (*models.BodyMeasurements, error) {
	return s.hrClient.GetMeasurements(id)
}

func (s *HTTPEmployeeService) RegisterMeasurements(id string, req models.BodyMeasurements) (*models.BodyMeasurements, error) {
	return s.hrClient.RegisterMeasurements(id, req)
}

func (s *HTTPEmployeeService) UpdateMeasurements(id string, req models.BodyMeasurements) (*models.BodyMeasurements, error) {
	return s.hrClient.UpdateMeasurements(id, req)
}

func (s *HTTPEmployeeService) UpdatePreferences(id string, req models.UpdatePreferencesRequest) (*models.UpdatePreferencesRequest, error) {
	// Mock temporal: retorna el objeto enviado sin persistencia en DDL
	return &req, nil
}

func (s *HTTPEmployeeService) UpdateSecurity(id string, req models.UpdateSecurityRequest) error {
	return s.hrClient.UpdateSecurity(id, req)
}

func (s *HTTPEmployeeService) GetActivity(id string) ([]models.ActivityLog, error) {
	return []models.ActivityLog{}, nil
}

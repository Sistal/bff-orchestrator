package services

import (
	"github.com/Sistal/bff-orchestrator/internal/clients"
	"github.com/Sistal/bff-orchestrator/internal/logger"
	"github.com/Sistal/bff-orchestrator/internal/models"
	"go.uber.org/zap"
)

type EmployeeService interface {
	GetProfile(id string) (*models.EmployeeProfile, error)
	// CreateEmployee crea un nuevo funcionario asociado a usuario
	CreateEmployee(req models.CreateEmployeeRequest) (*models.EmployeeProfile, error)
	// CreateEmployeeFromLogin crea un funcionario con datos mínimos desde el login
	CreateEmployeeFromLogin(req models.CreateEmployeeRequestFromLogin) (*models.EmployeeProfile, error)
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
	// GetEmployeeByUserID obtiene el ID del funcionario asociado a un usuario
	GetEmployeeByUserID(userID string) (int, error)
}

// HTTPEmployeeService — implementación real que delega al ms-funcionario
type HTTPEmployeeService struct {
	hrClient  *clients.HRClient
	opsClient *clients.OpsClient
}

func NewHTTPEmployeeService(hr *clients.HRClient, ops *clients.OpsClient) EmployeeService {
	return &HTTPEmployeeService{hrClient: hr, opsClient: ops}
}

// GetHRClient permite exponer el cliente a otros servicios del mismo paquete
func (s *HTTPEmployeeService) GetHRClient() *clients.HRClient {
	return s.hrClient
}

func (s *HTTPEmployeeService) GetProfile(id string) (*models.EmployeeProfile, error) {
	log := logger.Get()
	log.Debug("EmployeeService.GetProfile: Consultando perfil", zap.String("id", id))
	return s.hrClient.GetEmployeeProfile(id)
}

func (s *HTTPEmployeeService) CreateEmployee(req models.CreateEmployeeRequest) (*models.EmployeeProfile, error) {
	return s.hrClient.CreateEmployee(req)
}

func (s *HTTPEmployeeService) CreateEmployeeFromLogin(req models.CreateEmployeeRequestFromLogin) (*models.EmployeeProfile, error) {
	return s.hrClient.CreateEmployeeFromLogin(req)
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

func (s *HTTPEmployeeService) GetEmployeeByUserID(userID string) (int, error) {
	return s.hrClient.GetEmployeeByUserID(userID)
}

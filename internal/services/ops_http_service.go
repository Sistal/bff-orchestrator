package services

import (
	"mime/multipart"

	"github.com/Sistal/bff-orchestrator/internal/clients"
	"github.com/Sistal/bff-orchestrator/internal/models"
)

// HTTPOpsService implementa RequestService, DeliveryService y NotificationService
// delegando al ms-operations a través de OpsClient.
type HTTPOpsService struct {
	client *clients.OpsClient
}

func NewHTTPOpsService(client *clients.OpsClient) *HTTPOpsService {
	return &HTTPOpsService{client: client}
}

// ─── RequestService ────────────────────────────────────────────────────────────

func (s *HTTPOpsService) GetRequests(userID string, params map[string]string) ([]models.RequestSummary, error) {
	return s.client.GetRequests(userID, params)
}

func (s *HTTPOpsService) GetRequestByID(userID, id string) (*models.RequestSummary, error) {
	return s.client.GetRequestByID(userID, id)
}

func (s *HTTPOpsService) CreateReplenishmentRequest(userID string, req models.CreateReplenishmentRequest) (*models.RequestSummary, error) {
	return s.client.CreateReplenishmentRequest(userID, req)
}

func (s *HTTPOpsService) CreateGarmentChangeRequest(userID string, req models.CreateGarmentChangeRequest) (*models.RequestSummary, error) {
	return s.client.CreateGarmentChangeRequest(userID, req)
}

func (s *HTTPOpsService) CreateUniformRequest(userID string, req models.CreateUniformRequest) (*models.RequestSummary, error) {
	return s.client.CreateUniformRequest(userID, req)
}

func (s *HTTPOpsService) UploadFile(userID string, file *multipart.FileHeader) (*models.FileUploadResponse, error) {
	return s.client.UploadFile(userID, file)
}

func (s *HTTPOpsService) GetRecentRequests(userID string) ([]models.RequestSummary, error) {
	return s.client.GetRecentRequests(userID)
}

// ─── DeliveryService ───────────────────────────────────────────────────────────

func (s *HTTPOpsService) GetDeliveries(userID string) ([]models.DeliverySummary, error) {
	return s.client.GetDeliveries(userID)
}

func (s *HTTPOpsService) GetDeliveryByID(userID, id string) (*models.DeliverySummary, error) {
	return s.client.GetDeliveryByID(userID, id)
}

func (s *HTTPOpsService) ConfirmDelivery(userID, id string) (*models.DeliverySummary, error) {
	return s.client.ConfirmDelivery(userID, id)
}

func (s *HTTPOpsService) GetUpcomingDeliveries(userID string) ([]models.DeliverySummary, error) {
	return s.client.GetUpcomingDeliveries(userID)
}

// ─── NotificationService ───────────────────────────────────────────────────────

func (s *HTTPOpsService) GetNotifications(userID string) ([]models.Notification, error) {
	return s.client.GetNotifications(userID)
}

func (s *HTTPOpsService) MarkAsRead(userID, id string) error {
	return s.client.MarkAsRead(userID, id)
}

func (s *HTTPOpsService) MarkAllAsRead(userID string) (int, error) {
	return s.client.MarkAllAsRead(userID)
}

package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Sistal/bff-orchestrator/internal/models"
)

type OpsClient struct {
	BaseURL string
	Client  *http.Client
}

func NewOpsClient() *OpsClient {
	baseURL := os.Getenv("MS_OPERATIONS_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8083" // Default port for Ops Service
	}
	return &OpsClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

// Request Operations

func (c *OpsClient) GetRequests(employeeID string) ([]models.RequestSummary, error) {
	url := fmt.Sprintf("%s/api/v1/peticiones?id_funcionario=%s", c.BaseURL, employeeID)
	resp, err := c.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get requests: %d", resp.StatusCode)
	}

	var requests []models.RequestSummary
	if err := json.NewDecoder(resp.Body).Decode(&requests); err != nil {
		return nil, err
	}
	return requests, nil
}

func (c *OpsClient) CreateReplenishmentRequest(employeeID string, req models.CreateReplenishmentRequest) error {
	body, _ := json.Marshal(req)
	// We might need to pass employeeID in headers or body in the real service
	// For now assuming headers or auth context is handled by proxy or param
	url := fmt.Sprintf("%s/api/v1/peticiones", c.BaseURL)
	reqHTTP, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")
	reqHTTP.Header.Set("X-User-ID", employeeID) // Pass context

	resp, err := c.Client.Do(reqHTTP)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create replenishment: %d", resp.StatusCode)
	}
	return nil
}

// Delivery / Logistics Operations

func (c *OpsClient) GetDeliveries(employeeID string) ([]models.DeliverySummary, error) {
	url := fmt.Sprintf("%s/api/v1/despachos?id_funcionario=%s", c.BaseURL, employeeID)
	resp, err := c.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get deliveries: %d", resp.StatusCode)
	}

	var deliveries []models.DeliverySummary
	if err := json.NewDecoder(resp.Body).Decode(&deliveries); err != nil {
		return nil, err
	}
	return deliveries, nil
}

// Notification Operations

func (c *OpsClient) GetNotifications(employeeID string) ([]models.Notification, error) {
	url := fmt.Sprintf("%s/api/v1/notificaciones?id_usuario=%s", c.BaseURL, employeeID)
	resp, err := c.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get notifications: %d", resp.StatusCode)
	}

	var notifs []models.Notification
	if err := json.NewDecoder(resp.Body).Decode(&notifs); err != nil {
		return nil, err
	}
	return notifs, nil
}

// Stats for Dashboard

func (c *OpsClient) GetDashboardStats(employeeID string) (*models.HomeStats, error) {
	// Call a specialized endpoint on Ops Service that aggregates counts
	url := fmt.Sprintf("%s/api/v1/operaciones/dashboard?id_usuario=%s", c.BaseURL, employeeID)
	resp, err := c.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get dashboard stats: %d", resp.StatusCode)
	}

	var stats models.HomeStats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

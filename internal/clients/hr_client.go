package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Sistal/bff-orchestrator/internal/models"
)

type HRClient struct {
	BaseURL string
	Client  *http.Client
}

func NewHRClient() *HRClient {
	baseURL := os.Getenv("MS_FUNCIONARIO_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8082" // Default port for HR Service
	}
	return &HRClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

func (c *HRClient) GetEmployeeProfile(id string) (*models.EmployeeProfile, error) {
	resp, err := c.Client.Get(fmt.Sprintf("%s/api/v1/funcionarios/%s", c.BaseURL, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get employee profile: %d", resp.StatusCode)
	}

	var profile models.EmployeeProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, err
	}
	return &profile, nil
}

func (c *HRClient) UpdateContact(id string, req models.UpdateContactRequest) error {
	body, _ := json.Marshal(req)
	reqHTTP, _ := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/funcionarios/%s/contacto", c.BaseURL, id), bytes.NewBuffer(body))
	reqHTTP.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(reqHTTP)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update contact: %d", resp.StatusCode)
	}
	return nil
}

func (c *HRClient) GetBranches() ([]models.Branch, error) {
	resp, err := c.Client.Get(c.BaseURL + "/api/v1/sucursales")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get branches: %d", resp.StatusCode)
	}

	var branches []models.Branch
	if err := json.NewDecoder(resp.Body).Decode(&branches); err != nil {
		return nil, err
	}
	return branches, nil
}

func (c *HRClient) GetMeasurements(id string) (*models.BodyMeasurements, error) {
	resp, err := c.Client.Get(fmt.Sprintf("%s/api/v1/funcionarios/%s/medidas", c.BaseURL, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get measurements: %d", resp.StatusCode)
	}

	var m models.BodyMeasurements
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func (c *HRClient) GetChangeHistory() ([]models.BranchChangeRequestHistory, error) {
	// Call /api/v1/transferencias
	return nil, nil // Placeholder
}

func (c *HRClient) CreateBranchChangeRequest(req models.CreateBranchChangeRequest) error {
	// POST /api/v1/transferencias
	return nil
}

package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Sistal/bff-orchestrator/internal/models"
)

type CatalogClient struct {
	BaseURL string
	Client  *http.Client
}

func NewCatalogClient() *CatalogClient {
	baseURL := os.Getenv("MS_CATALOG_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8084" // Default port for Catalog Service
	}
	return &CatalogClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

func (c *CatalogClient) GetMasterData() (*models.MasterDataResponse, error) {
	resp, err := c.Client.Get(c.BaseURL + "/api/v1/master-data")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get master data: %d", resp.StatusCode)
	}

	var data models.MasterDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}

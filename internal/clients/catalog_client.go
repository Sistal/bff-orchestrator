package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Sistal/bff-orchestrator/internal/logger"
	"github.com/Sistal/bff-orchestrator/internal/models"
	"go.uber.org/zap"
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

func (c *CatalogClient) GetCompanies() ([]models.Empresa, error) {
	log := logger.Get()
	fullURL := c.BaseURL + "/api/v1/empresas"
	log.Info("CatalogClient: Solicitando empresas", zap.String("url", fullURL))

	resp, err := c.Client.Get(fullURL)
	if err != nil {
		log.Error("CatalogClient: Error al solicitar empresas", zap.Error(err), zap.String("url", fullURL))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warn("CatalogClient: Respuesta fallida al obtener empresas", zap.Int("status", resp.StatusCode), zap.String("url", fullURL))
		return nil, fmt.Errorf("failed to get companies: %d", resp.StatusCode)
	}

	var result models.CatalogResponse[models.Empresa]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Error("CatalogClient: Error decodificando respuesta de empresas", zap.Error(err), zap.String("url", fullURL))
		return nil, err
	}

	log.Info("CatalogClient: Empresas obtenidas exitosamente", zap.Int("count", len(result.Data)), zap.String("url", fullURL))
	return result.Data, nil
}

func (c *CatalogClient) GetSegments(companyID int) ([]models.Segmento, error) {
	log := logger.Get()
	url := fmt.Sprintf("%s/api/v1/segmentos/%d", c.BaseURL, companyID)
	log.Info("CatalogClient: Solicitando segmentos", zap.Int("company_id", companyID), zap.String("url", url))

	resp, err := c.Client.Get(url)
	if err != nil {
		log.Error("CatalogClient: Error al solicitar segmentos", zap.Error(err), zap.Int("company_id", companyID), zap.String("url", url))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warn("CatalogClient: Respuesta fallida al obtener segmentos", zap.Int("status", resp.StatusCode), zap.Int("company_id", companyID), zap.String("url", url))
		return nil, fmt.Errorf("failed to get segments: %d", resp.StatusCode)
	}

	var result models.CatalogResponse[models.Segmento]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Error("CatalogClient: Error decodificando respuesta de segmentos", zap.Error(err), zap.Int("company_id", companyID), zap.String("url", url))
		return nil, err
	}

	log.Info("CatalogClient: Segmentos obtenidos exitosamente", zap.Int("count", len(result.Data)), zap.Int("company_id", companyID), zap.String("url", url))
	return result.Data, nil
}

func (c *CatalogClient) GetBranches(companyID int) ([]models.Sucursal, error) {
	log := logger.Get()
	url := fmt.Sprintf("%s/api/v1/sucursales/%d", c.BaseURL, companyID)
	log.Info("CatalogClient: Solicitando sucursales", zap.Int("company_id", companyID), zap.String("url", url))

	resp, err := c.Client.Get(url)
	if err != nil {
		log.Error("CatalogClient: Error al solicitar sucursales", zap.Error(err), zap.Int("company_id", companyID), zap.String("url", url))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warn("CatalogClient: Respuesta fallida al obtener sucursales", zap.Int("status", resp.StatusCode), zap.Int("company_id", companyID), zap.String("url", url))
		return nil, fmt.Errorf("failed to get branches: %d", resp.StatusCode)
	}

	var result models.CatalogResponse[models.Sucursal]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Error("CatalogClient: Error decodificando respuesta de sucursales", zap.Error(err), zap.Int("company_id", companyID), zap.String("url", url))
		return nil, err
	}

	log.Info("CatalogClient: Sucursales obtenidas exitosamente", zap.Int("count", len(result.Data)), zap.Int("company_id", companyID), zap.String("url", url))
	return result.Data, nil
}

func (c *CatalogClient) GetUniformsBySegment(segmentID int) ([]models.Uniform, error) {
	log := logger.Get()
	url := fmt.Sprintf("%s/api/v1/segmentos/%d/uniformes", c.BaseURL, segmentID)
	log.Info("CatalogClient: Solicitando uniformes por segmento", zap.Int("segmentID", segmentID), zap.String("url", url))

	resp, err := c.Client.Get(url)
	if err != nil {
		log.Error("CatalogClient: Error al solicitar uniformes", zap.Error(err), zap.Int("segmentID", segmentID))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warn("CatalogClient: Respuesta fallida al obtener uniformes", zap.Int("status", resp.StatusCode), zap.Int("segmentID", segmentID))
		return nil, fmt.Errorf("failed to get uniforms: %d", resp.StatusCode)
	}

	var result models.CatalogResponse[models.Uniform]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Error("CatalogClient: Error decodificando respuesta de uniformes", zap.Error(err))
		return nil, err
	}

	log.Info("CatalogClient: Uniformes obtenidos exitosamente", zap.Int("count", len(result.Data)), zap.Int("segmentID", segmentID))
	return result.Data, nil
}

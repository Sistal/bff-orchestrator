package services

import (
	"github.com/Sistal/bff-orchestrator/internal/clients"
	"github.com/Sistal/bff-orchestrator/internal/models"
)

type CatalogService interface {
	GetSizes() ([]models.CatalogItem, error)
	GetChangeReasons() ([]models.CatalogItem, error)
	GetGarmentTypes() ([]models.CatalogItem, error)
	GetActiveCampaign() (*models.Campaign, error)
	GetMasterData() (*models.MasterDataResponse, error)
}

type MockCatalogService struct{}

func NewMockCatalogService() CatalogService {
	return &MockCatalogService{}
}

func (s *MockCatalogService) GetSizes() ([]models.CatalogItem, error) {
	return []models.CatalogItem{
		{ID: "S", Label: "S - Pequeño"},
		{ID: "M", Label: "M - Mediano"},
		{ID: "L", Label: "L - Grande"},
	}, nil
}

func (s *MockCatalogService) GetChangeReasons() ([]models.CatalogItem, error) {
	return []models.CatalogItem{
		{ID: "TALLA", Label: "Talla incorrecta"},
		{ID: "DEFECTO", Label: "Producto defectuoso"},
	}, nil
}

func (s *MockCatalogService) GetGarmentTypes() ([]models.CatalogItem, error) {
	return []models.CatalogItem{
		{ID: "POLERA", Label: "Polera Institucional"},
		{ID: "PANTALON", Label: "Pantalón Cargo"},
	}, nil
}

func (s *MockCatalogService) GetActiveCampaign() (*models.Campaign, error) {
	return &models.Campaign{
		ID:          "CAM-2025",
		Nombre:      "Temporada Invierno 2025",
		FechaInicio: "2025-05-01",
		FechaFin:    "2025-08-31",
		Activa:      true,
	}, nil
}

func (s *MockCatalogService) GetMasterData() (*models.MasterDataResponse, error) {
	sizes, _ := s.GetSizes()
	types, _ := s.GetGarmentTypes()
	reasons, _ := s.GetChangeReasons()
	camp, _ := s.GetActiveCampaign()
	return &models.MasterDataResponse{
		Sizes:         sizes,
		GarmentTypes:  types,
		ChangeReasons: reasons,
		Campaign:      *camp,
	}, nil
}

type HTTPCatalogService struct {
	client *clients.CatalogClient
}

func NewHTTPCatalogService(client *clients.CatalogClient) CatalogService {
	return &HTTPCatalogService{client: client}
}

func (s *HTTPCatalogService) GetSizes() ([]models.CatalogItem, error) {
	data, err := s.client.GetMasterData()
	if err != nil {
		return nil, err
	}
	return data.Sizes, nil
}

func (s *HTTPCatalogService) GetChangeReasons() ([]models.CatalogItem, error) {
	data, err := s.client.GetMasterData()
	if err != nil {
		return nil, err
	}
	return data.ChangeReasons, nil
}

func (s *HTTPCatalogService) GetGarmentTypes() ([]models.CatalogItem, error) {
	data, err := s.client.GetMasterData()
	if err != nil {
		return nil, err
	}
	return data.GarmentTypes, nil
}

func (s *HTTPCatalogService) GetActiveCampaign() (*models.Campaign, error) {
	data, err := s.client.GetMasterData()
	if err != nil {
		return nil, err
	}
	return &data.Campaign, nil
}

func (s *HTTPCatalogService) GetMasterData() (*models.MasterDataResponse, error) {
	return s.client.GetMasterData()
}

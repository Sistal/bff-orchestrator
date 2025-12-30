package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Sistal/bff-orchestrator/internal/domain/model"
)

// ProductProvider implements the ProductProvider port
type ProductProvider struct {
	baseURL    string
	httpClient *http.Client
}

// NewProductProvider creates a new product provider
func NewProductProvider(baseURL string) *ProductProvider {
	return &ProductProvider{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetProductByID fetches a product by ID from the external API
func (p *ProductProvider) GetProductByID(ctx context.Context, productID int) (*model.Product, error) {
	url := fmt.Sprintf("%s/products/%d", p.baseURL, productID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var product model.Product
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, fmt.Errorf("failed to decode product: %w", err)
	}

	return &product, nil
}

// GetProducts fetches products from the external API with a limit
func (p *ProductProvider) GetProducts(ctx context.Context, limit int) ([]*model.Product, error) {
	url := fmt.Sprintf("%s/products?limit=%d", p.baseURL, limit)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var products []*model.Product
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return nil, fmt.Errorf("failed to decode products: %w", err)
	}

	return products, nil
}

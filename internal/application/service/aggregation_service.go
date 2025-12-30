package service

import (
	"context"

	"github.com/Sistal/bff-orchestrator/internal/domain/model"
	"github.com/Sistal/bff-orchestrator/internal/ports"
)

// AggregationService orchestrates data from multiple providers
type AggregationService struct {
	userProvider    ports.UserProvider
	productProvider ports.ProductProvider
}

// NewAggregationService creates a new aggregation service
func NewAggregationService(userProvider ports.UserProvider, productProvider ports.ProductProvider) *AggregationService {
	return &AggregationService{
		userProvider:    userProvider,
		productProvider: productProvider,
	}
}

// GetDashboard aggregates user and product data for the dashboard
func (s *AggregationService) GetDashboard(ctx context.Context, userID int) (*model.Dashboard, error) {
	// Fetch user data
	user, err := s.userProvider.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Fetch products (limit to 10 for dashboard)
	products, err := s.productProvider.GetProducts(ctx, 10)
	if err != nil {
		return nil, err
	}

	return &model.Dashboard{
		User:     user,
		Products: products,
	}, nil
}

package ports

import (
	"context"

	"github.com/Sistal/bff-orchestrator/internal/domain/model"
)

// UserProvider defines the port for fetching user data from external services
type UserProvider interface {
	GetUserByID(ctx context.Context, userID int) (*model.User, error)
	GetUsers(ctx context.Context) ([]*model.User, error)
}

// ProductProvider defines the port for fetching product data from external services
type ProductProvider interface {
	GetProductByID(ctx context.Context, productID int) (*model.Product, error)
	GetProducts(ctx context.Context, limit int) ([]*model.Product, error)
}

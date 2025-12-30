package ports

import (
	"context"

	"github.com/Sistal/bff-orchestrator/internal/domain/model"
)

// AggregationService defines the port for aggregating data from multiple sources
type AggregationService interface {
	GetDashboard(ctx context.Context, userID int) (*model.Dashboard, error)
}

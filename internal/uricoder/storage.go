package uricoder

import (
	"context"

	"github.com/yury-kuznetsov/shortener/internal/models"
)

// Storage is an interface that defines methods for interacting with a storage system.
type Storage interface {
	Get(ctx context.Context, code string, userID int) (string, error)
	Set(ctx context.Context, uri string, userID int) (string, error)
	GetByUser(ctx context.Context, userID int) ([]models.GetByUserResponse, error)
	SoftDelete(ctx context.Context, messages []models.RmvUrlsMsg) error
	HealthCheck(ctx context.Context) error
}

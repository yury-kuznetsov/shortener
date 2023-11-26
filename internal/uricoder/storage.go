package uricoder

import (
	"context"
	"github.com/yury-kuznetsov/shortener/internal/models"
)

type Storage interface {
	Get(ctx context.Context, code string, userID int) (string, error)
	Set(ctx context.Context, uri string, userID int) (string, error)
	GetByUser(ctx context.Context, userID int) ([]models.GetByUserResponse, error)
	HealthCheck(ctx context.Context) error
}

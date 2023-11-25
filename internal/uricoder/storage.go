package uricoder

import "context"

type Storage interface {
	Get(ctx context.Context, code string) (string, error)
	Set(ctx context.Context, uri string) (string, error)
	HealthCheck(ctx context.Context) error
}

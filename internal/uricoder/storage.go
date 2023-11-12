package uricoder

type Storage interface {
	Get(code string) (string, error)
	Set(uri string) (string, error)
	HealthCheck() error
}

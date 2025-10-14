package service

type metricRepo[T int64 | float64] interface {
	Get(name string) (T, bool)
	Set(name string, value T)
	Update(name string, delta T)
}

// Service provides metric-related operations.
type Service struct {
	gaugeRepo   metricRepo[float64]
	counterRepo metricRepo[int64]
}

// NewService creates a new Service instance.
func NewService(gaugeRepo metricRepo[float64], counterRepo metricRepo[int64]) *Service {
	return &Service{
		gaugeRepo:   gaugeRepo,
		counterRepo: counterRepo,
	}
}

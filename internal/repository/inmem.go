package repository

// MetricInMemRepo is an in-memory repository for metrics.
type MetricInMemRepo[T int64 | float64] struct {
	storage map[string]T
}

// NewMetricInMemRepo creates a new instance of MetricInMemRepo.
func NewMetricInMemRepo[T int64 | float64]() *MetricInMemRepo[T] {
	return &MetricInMemRepo[T]{
		storage: make(map[string]T),
	}
}

// Get returns the value of a metric by its name and a bool flag to check if it exists.
func (r *MetricInMemRepo[T]) Get(name string) (T, bool) {
	value, exists := r.storage[name]
	return value, exists
}

// Set sets the value of a metric by its name.
func (r *MetricInMemRepo[T]) Set(name string, value T) {
	r.storage[name] = value
}

// Update updates the value of a metric by adding the delta to the current value.
func (r *MetricInMemRepo[T]) Update(name string, delta T) {
	if currentValue, exists := r.storage[name]; exists {
		r.storage[name] = currentValue + delta
	} else {
		r.storage[name] = delta
	}
}

package pkg

import (
	"fmt"
	"os"
)

type envParams interface {
	string | int | int64 | float64 | bool
}

// GetEnv returns the value of the environment variable or default value.
func GetEnv[T envParams](key string, defaultVal T) T {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultVal
	}

	var result T
	switch any(defaultVal).(type) {
	case string:
		result = any(value).(T)
	case int:
		fmt.Sscanf(value, "%d", &result)
	case int64:
		fmt.Sscanf(value, "%d", &result)
	case float64:
		fmt.Sscanf(value, "%f", &result)
	case bool:
		fmt.Sscanf(value, "%t", &result)
	}

	return result
}

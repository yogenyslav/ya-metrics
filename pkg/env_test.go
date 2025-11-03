package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	t.Run("envParam string", func(t *testing.T) {
		const key = "TEST_ENV_STRING"
		const defaultVal = "default"
		const expectedVal = "value"

		var v string

		v = GetEnv(key, defaultVal)
		assert.Equal(t, defaultVal, v)

		t.Setenv(key, expectedVal)
		v = GetEnv(key, defaultVal)
		assert.Equal(t, expectedVal, v)
	})

	t.Run("envParam int", func(t *testing.T) {
		const key = "TEST_ENV_INT"
		const defaultVal = 10
		const expectedVal = 20

		var v int

		v = GetEnv(key, defaultVal)
		assert.Equal(t, defaultVal, v)

		t.Setenv(key, "20")
		v = GetEnv(key, defaultVal)
		assert.Equal(t, expectedVal, v)
	})

	t.Run("envParam bool", func(t *testing.T) {
		const key = "TEST_ENV_BOOL"
		const defaultVal = false
		const expectedVal = true

		var v bool

		v = GetEnv(key, defaultVal)
		assert.Equal(t, defaultVal, v)

		t.Setenv(key, "true")
		v = GetEnv(key, defaultVal)
		assert.Equal(t, expectedVal, v)
	})

	t.Run("envParam float64", func(t *testing.T) {
		const key = "TEST_ENV_FLOAT"
		const defaultVal = 1.5
		const expectedVal = 2.5

		var v float64

		v = GetEnv(key, defaultVal)
		assert.Equal(t, defaultVal, v)

		t.Setenv(key, "2.5")
		v = GetEnv(key, defaultVal)
		assert.Equal(t, expectedVal, v)
	})
}

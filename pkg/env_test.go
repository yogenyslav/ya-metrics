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

func Test_parseEnv(t *testing.T) {
	t.Parallel()

	t.Run("parse string", func(t *testing.T) {
		t.Parallel()

		value := "asdfg"
		result := parseEnv[string](value)
		assert.Equal(t, value, result)
	})

	t.Run("parse int", func(t *testing.T) {
		t.Parallel()

		value := "12345"
		want := 12345
		result := parseEnv[int](value)
		assert.Equal(t, want, result)
	})

	t.Run("parse bool", func(t *testing.T) {
		t.Parallel()

		value := "true"
		want := true
		result := parseEnv[bool](value)
		assert.Equal(t, want, result)
	})

	t.Run("parse float64", func(t *testing.T) {
		t.Parallel()

		value := "1.23"
		want := 1.23
		result := parseEnv[float64](value)
		assert.Equal(t, want, result)
	})
}

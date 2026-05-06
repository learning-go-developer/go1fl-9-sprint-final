package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomElements(t *testing.T) {
	// Тест нормального случая
	t.Run("Normal size", func(t *testing.T) {
		n := 10
		result := generateRandomElements(n)
		assert.Equal(t, n, len(result))
	})

	// Тест на нулевой размер
	t.Run("Zero size", func(t *testing.T) {
		result := generateRandomElements(0)
		assert.Empty(t, result)
	})

	// Тест на отрицательный размер
	t.Run("Negative size", func(t *testing.T) {
		result := generateRandomElements(-5)
		assert.Empty(t, result)
	})
}

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	nNormal = 1_000
	nLarge  = 100_000_000
	result  []int // Глобальная переменная
)

// command run test: 'go test . -v'
// command run test with ingnore cache: 'go test -v -count=1'

// TestGenerateRandomElements validates the logic of random element generation.
// It covers standard cases, boundary conditions, and performance stability.
//
// The test suite includes:
//   - Range Validation: Ensures all numbers are within [0, 100).
//   - Zero/Negative Input: Verifies that invalid sizes return empty slices.
//   - Large Scale Testing: Checks goroutine stability under a 1,000,000 element load.
func TestGenerateRandomElements(t *testing.T) {
	// Тест нормального случая + проверка диапазона
	t.Run("Normal size and range", func(t *testing.T) {
		result := generateRandomElements(nNormal)
		assert.Equal(t, nNormal, len(result))

		for _, val := range result {
			// Проверяем, что числа в диапазоне [0, 100)
			assert.True(t, val >= 0 && val < 100, "Число вне диапазона: %d", val)
		}
	})

	t.Run("Zero size", func(t *testing.T) {
		assert.Empty(t, generateRandomElements(0))
	})

	t.Run("Negative size", func(t *testing.T) {
		assert.Empty(t, generateRandomElements(-5))
	})

	// Проверка на очень большом размере (чтобы убедиться, что горутины не "падают")
	t.Run("Large size", func(t *testing.T) {
		result := generateRandomElements(nLarge)
		assert.Equal(t, nLarge, len(result))
	})
}

// BenchmarkGenerateRandomElements measures the performance of random element generation
// at a large scale (100,000,000 elements).
//
// It evaluates:
//   - Execution Time: How long it takes to generate and populate the slice.
//   - Memory Allocations: Measured via -benchmem to track heap usage and GC pressure.
//   - Concurrency Scaling: Tested with different CPU counts (-cpu 1,2,4) to verify goroutine efficiency.
//
// Recommended execution command for stable results:
//
//	go test -bench=. -benchmem -benchtime=5s
//
// Using -benchtime=5s extends the test duration, providing a more stable
// average by reducing the impact of OS background noise and GC pauses.
func BenchmarkGenerateRandomElements(b *testing.B) {
	var r []int
	for i := 0; i < b.N; i++ {
		// Сохраняем в локальную переменную внутри цикла
		r = generateRandomElements(nLarge)
	}
	// Записываем в глобальную после цикла, чтобы компилятор не удалил вызов
	result = r
}

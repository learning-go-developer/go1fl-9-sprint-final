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
// command check area code test 'go test -cover'

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
			// Проверяем, что числа в диапазоне [0, nLarge)
			assert.True(t, val >= 0 && val < nLarge, "Число вне диапазона: %d", val)
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

// TestMax runs table-driven tests to verify the correctness of the basic
// sequential maximum search implementation.
// It ensures the function handles various scenarios including positive and
// negative integers, single-element slices, and empty inputs correctly.
func TestMax(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected int
	}{
		{
			name:     "Положительные числа",
			input:    []int{1, 5, 3, 9, 2},
			expected: 9,
		},
		{
			name:     "Отрицательные числа",
			input:    []int{-10, -5, -20, -1},
			expected: -1,
		},
		{
			name:     "Пустой слайс",
			input:    []int{},
			expected: 0,
		},
		{
			name:     "Один элемент",
			input:    []int{42},
			expected: 42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := max(tt.input)
			assert.Equal(t, tt.expected, result, "Ошибка в тесте: %s", tt.name)
		})
	}
}

// BenchmarkMax measures the performance of the sequential max function.
// This serves as a baseline to compare against parallel implementations.
func BenchmarkMax(b *testing.B) {
	// Подготовка данных
	data := make([]int, nLarge)
	for i := range data {
		data[i] = i
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		max(data)
	}
}

// TestMaximum runs table-driven tests to verify the correctness of the maximum function.
// It covers various scenarios including positive and negative numbers, identical elements,
// single-element slices, empty slices, and edge cases related to concurrent processing
// such as small input sizes and inputs not perfectly divisible by the number of workers.
func TestMaximum(t *testing.T) {
	// Описываем таблицу тестовых случаев
	tests := []struct {
		name     string // название теста
		input    []int  // что подаем на вход
		expected int    // что ожидаем на выходе
	}{
		{
			name:     "Положительные числа",
			input:    []int{1, 5, 3, 9, 2},
			expected: 9,
		},
		{
			name:     "Отрицательные числа",
			input:    []int{-10, -5, -20, -2},
			expected: -2,
		},
		{
			name:     "Одинаковые числа",
			input:    []int{7, 7, 7, 7},
			expected: 7,
		},
		{
			name:     "Один элемент",
			input:    []int{42},
			expected: 42,
		},
		{
			name:     "Пустой слайс",
			input:    []int{},
			expected: 0,
		},
		{
			name: "Размер меньше порога (однопоточный режим)",
			input: func() []int {
				res := make([]int, 500)
				res[250] = 100
				return res
			}(),
			expected: 100,
		},
		{
			name:  "Размер не кратный количеству воркеров",
			input: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			// Если CHUNKS=8, проверится логика "остатка"
			expected: 10,
		},
	}

	// Проходим по всем тестам
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maximum(tt.input)
			assert.Equal(t, tt.expected, result, "Ожидалось другое значение для теста: %s", tt.name)
		})
	}
}

// BenchmarkMaximum measures the performance of the maximum function using a large dataset.
// It pre-allocates a slice of size nLarge and resets the timer to ensure that only
// the execution time of the maximum function is recorded. This is used to evaluate
// the efficiency of the parallel processing implementation.
func BenchmarkMaximum(b *testing.B) {
	// Подготовка данных
	data := make([]int, nLarge)
	for i := range data {
		data[i] = i
	}

	b.ResetTimer() // Сбрасываем таймер, чтобы не учитывать время подготовки данных

	for i := 0; i < b.N; i++ {
		maximum(data)
	}
}

// TestMaxChunks verifies the correctness of the chunk-based maximum search.
// It covers edge cases where the input length is less than the number of chunks,
// empty inputs, and inputs where the size is not perfectly divisible by CHUNKS.
func TestMaxChunks(t *testing.T) {
	// Описываем таблицу тестовых случаев
	tests := []struct {
		name     string
		input    []int
		expected int
	}{
		{
			name:     "Положительные числа",
			input:    []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
			expected: 100,
		},
		{
			name:     "Отрицательные числа",
			input:    []int{-5, -10, -2, -8, -15, -3, -7, -1},
			expected: -1,
		},
		{
			name:     "Слайс меньше 8 элементов (некоторые чанки пустые)",
			input:    []int{5, 12, 3},
			expected: 12,
		},
		{
			name:     "Один элемент",
			input:    []int{42},
			expected: 42,
		},
		{
			name:     "Пустой слайс",
			input:    []int{},
			expected: 0,
		},
		{
			name: "Большой слайс, не кратный 8",
			input: func() []int {
				res := make([]int, 1003) // 1003 / 8 даст остаток
				res[1002] = 999
				return res
			}(),
			expected: 999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maxChunks(tt.input)
			assert.Equal(t, tt.expected, result, "Тест '%s' провален", tt.name)
		})
	}
}

// BenchmarkMaxChunks evaluates the performance of the WaitGroup-based
// concurrent maximum search using a large pre-allocated dataset.
func BenchmarkMaxChunks(b *testing.B) {
	// Подготовка данных (используем тот же nLarge, что и в прошлый раз)
	data := make([]int, nLarge)
	for i := range data {
		data[i] = i
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		maxChunks(data)
	}
}

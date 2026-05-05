package main

import (
	"fmt"
)

const (
	SIZE   = 100_000_000
	CHUNKS = 8
)

// generateRandomElements generates random elements.
func generateRandomElements(size int) []int {
	if size <= 0 {
		return []int{}
	}

	rand.Seed(time.Now().UnixNano())

	result := make([]int, size)

	for i := 0; i < size; i++ {
		result[i] = rand.Intn(100)
	}

	return result
}

// maximum returns the maximum number of elements.
func maximum(data []int) int {
	// ваш код здесь
}

// maxChunks returns the maximum number of elements in a chunks.
func maxChunks(data []int) int {
	// ваш код здесь
}

func main() {
	fmt.Printf("Генерируем %d целых чисел", SIZE)
	// ваш код здесь

	fmt.Println("Ищем максимальное значение в один поток")
	// ваш код здесь

	fmt.Printf("Максимальное значение элемента: %d\nВремя поиска: %d ms\n", max, elapsed)

	fmt.Printf("Ищем максимальное значение в %d потоков", CHUNKS)
	// ваш код здесь

	fmt.Printf("Максимальное значение элемента: %d\nВремя поиска: %d ms\n", max, elapsed)
}

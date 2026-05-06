package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

const (
	SIZE   = 100_000_000
	CHUNKS = 8
)

// generateRandomElements creates a slice of random integers using parallel processing.
// Each integer is in the range [0, 100).
//
// The function automatically partitions the workload across multiple goroutines
// based on the number of available CPU cores (runtime.NumCPU).
//
// Efficiency notes:
//   - It pre-allocates the result slice to avoid multiple memory reallocations.
//   - It uses a sync.WaitGroup to coordinate concurrent writes to different
//     segments of the slice.
//   - It utilizes math/rand/v2 for improved performance and modern random
//     number generation algorithms.
//
// Parameters:
//   - size: The total number of elements to generate.
//
// Returns:
//   - A slice of integers if size > 0; otherwise, returns nil.
func generateRandomElements(size int) []int {
	if size <= 0 {
		return nil
	}

	result := make([]int, size)

	numWorkers := CHUNKS
	if size < numWorkers {
		numWorkers = 1
	}

	var wg sync.WaitGroup
	// Считаем размер порции данных для одной горутины
	chunkSize := (size + numWorkers - 1) / numWorkers

	for i := range numWorkers {
		start := i * chunkSize
		if start >= size {
			break
		}

		end := start + chunkSize
		if end > size {
			end = size
		}

		// last gorutine set remainder for end
		if i == numWorkers-1 {
			end = size
		}

		wg.Add(1)
		go func(s, e int) {
			defer wg.Done()

			for j := s; j < e; j++ {
				result[j] = rand.IntN(100) // use rand from v2 version math package with new algoritms
			}
		}(start, end)
	}

	wg.Wait()
	return result
}

/*
// maximum returns the maximum number of elements.
func maximum(data []int) int {
	// ваш код здесь
}

// maxChunks returns the maximum number of elements in a chunks.
func maxChunks(data []int) int {
	// ваш код здесь
}
*/

func main() {
	start := time.Now()

	fmt.Printf("Генерируем %d целых чисел", SIZE)
	// ваш код здесь
	data := generateRandomElements(SIZE)
	// for test control
	fmt.Printf("Готово! Сгенерировано %d элементов\n", len(data))
	fmt.Printf("Время выполнения: %v\n", time.Since(start))

	/*
		fmt.Println("Ищем максимальное значение в один поток")
		// ваш код здесь

		fmt.Printf("Максимальное значение элемента: %d\nВремя поиска: %d ms\n", max, elapsed)

		fmt.Printf("Ищем максимальное значение в %d потоков", CHUNKS)
		// ваш код здесь

		fmt.Printf("Максимальное значение элемента: %d\nВремя поиска: %d ms\n", max, elapsed)
	*/
}

package main

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"sync"
	"time"
)

const (
	SIZE         = 100_000_000
	CHUNKS       = 8
	minBatchSize = 1_000
)

// generateRandomElements creates a slice of the given size and populates it
// with pseudo-random integers in the range [0, size) using math/rand/v2.
//
// If size is less than or equal to 0, it returns nil.
// The function uses a sequential loop to ensure optimal cache performance
// and minimal memory overhead.
func generateRandomElements(size int) []int {
	if size <= 0 {
		return nil
	}

	result := make([]int, size)

	for i := range size {
		result[i] = rand.IntN(size)
	}

	return result
}

// maximum finds the largest integer in the slice using an optimized hybrid approach.
//
// The function selects the most efficient execution path based on input size:
//   - If the slice is empty, it returns 0 as a safe default.
//   - For small datasets (less than 1,000 elements), it uses a fast sequential
//     search via slices.Max to avoid goroutine orchestration overhead.
//   - For large datasets, it delegates the task to maxChunks for high-performance
//     parallel processing across multiple CPU cores.
//
// This strategy ensures minimal latency for small inputs while maintaining
// maximum throughput for large-scale data processing.
func maximum(data []int) int {
	if len(data) == 0 {
		return 0
	}

	if len(data) < minBatchSize {
		return slices.Max(data)
	}

	return maxChunks(data)
}

// maxChunks finds the maximum value in a slice by dividing it into a fixed
// number of segments (CHUNKS) and processing them concurrently.
//
// The function splits the input slice into CHUNKS segments, spawns a goroutine
// for each segment to find its local maximum, and then identifies the overall
// maximum among these local results.
//
// Concurrency control:
//   - Uses sync.WaitGroup to ensure all concurrent workers finish before returning.
//   - Results from each goroutine are stored in a pre-allocated slice to avoid race conditions.
//
// Edge cases:
//   - If the input slice is empty, it returns 0.
//   - If the input size is smaller than CHUNKS, it gracefully handles empty segments.
//   - Correctly handles negative numbers by initializing intermediate results
//     with the first element of the input data.
//
// Performance note:
// This approach is effective for very large slices where the overhead of
// goroutine creation is outweighed by the benefits of parallel CPU utilization.
func maxChunks(data []int) int {
	size := len(data)
	if size == 0 {
		return 0
	}

	chunkMaxes := make([]int, CHUNKS)

	for i := range chunkMaxes {
		chunkMaxes[i] = data[0]
	}

	chunkSize := size / CHUNKS

	if chunkSize == 0 {
		chunkSize = 1
	}

	var wg sync.WaitGroup

	for i := range CHUNKS {
		start := i * chunkSize

		if start >= size {
			continue
		}

		end := start + chunkSize

		if i == CHUNKS-1 || end > size {
			end = size
		}

		segment := data[start:end]

		wg.Add(1)
		go func(s []int, index int) {
			defer wg.Done()
			if len(s) > 0 {
				chunkMaxes[index] = slices.Max(s)
			}
		}(segment, i)
	}

	wg.Wait()

	return slices.Max(chunkMaxes)
}

func main() {
	fmt.Printf("Генерируем %d целых чисел\n", SIZE)
	data := generateRandomElements(SIZE)

	fmt.Println("Ищем максимальное значение в один поток")
	startMax := time.Now()
	maxVal := maximum(data)
	durationMax := time.Since(startMax)
	fmt.Printf("Максимальное значение: %d\nВремя поиска: %d мкс\n", maxVal, durationMax.Microseconds())

	fmt.Printf("Ищем максимальное значение в %d потоков\n", CHUNKS)
	startChunks := time.Now()
	maxChunksVal := maxChunks(data)
	durationChunks := time.Since(startChunks)
	fmt.Printf("Максимальное значение: %d\nВремя поиска: %d мкс\n", maxChunksVal, durationChunks.Microseconds())
}

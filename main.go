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

// Maximum returns the largest element in the provided slice of integers.
// If the slice is empty, it returns 0.
func maximum(data []int) int {
	if len(data) <= 0 {
		return 0
	}

	tempCheck := data[0]
	for _, v := range data {
		if v > tempCheck {
			tempCheck = v
		}
	}

	return tempCheck
}

// maxChunks calculates the maximum value in the data slice by splitting it
// into multiple chunks and processing them concurrently using goroutines.
// It returns 0 if the input slice is empty.
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
				chunkMaxes[index] = maximum(s)
			}
		}(segment, i)
	}

	wg.Wait()

	return maximum(chunkMaxes)
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

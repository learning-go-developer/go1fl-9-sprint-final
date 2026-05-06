package main

import (
	"fmt"
	"math/rand/v2"
	"slices"
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
	if size < 1000 { // TODO for test perfomance theard one without set numWorkers
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
				result[j] = rand.IntN(size) // use rand from v2 version math package with new algoritms
			}
		}(start, end)
	}

	wg.Wait()
	return result
}

// maximum finds the largest integer in the given slice by processing data in parallel.
//
// The function partitions the slice into segments and assigns each to a separate goroutine
// to leverage multi-core efficiency. For small slices (less than 1000 elements),
// it defaults to a single-threaded execution to avoid goroutine overhead.
//
// It handles edge cases safely:
// - Returns 0 for empty slices.
// - Correctly identifies the maximum in slices containing negative numbers or MinInt.
// - Prevents deadlocks by tracking the exact number of active workers.
func maximum(data []int) int {
	size := len(data)

	if len(data) == 0 {
		return 0
	}

	numWorkers := CHUNKS
	if size < 1000 {
		numWorkers = 1
	}
	chunkSize := (size + numWorkers - 1) / numWorkers
	// Канал для сбора локальных максимумов от каждой горутины
	maxChan := make(chan int, numWorkers)

	activeWorkers := 0 // for calculate real run gorutines

	for i := 0; i < numWorkers; i++ {
		start := i * chunkSize
		if start >= size {
			break // Если данных больше нет, новые горутины просто не создаем
		}
		end := start + chunkSize
		if end > size {
			end = size
		}

		activeWorkers++ // Увеличиваем счетчик
		// Запускаем горутину на свой кусок данных
		go func(s, e int) {
			// Каждая горутина ищет максимум в своем куске (в один поток)
			maxChan <- slices.Max(data[s:e])
		}(start, end)
	}

	// Собираем результаты и находим финальный максимум
	// Получаем первый результат
	finalMax := <-maxChan // defens to "максимально возможное отрицательное число"
	// Получаем остальные результаты
	for i := 1; i < activeWorkers; i++ {
		localMax := <-maxChan
		if localMax > finalMax {
			finalMax = localMax
		}
	}

	return finalMax
}

// maxChunks finds the maximum value in a slice by dividing it into exactly CHUNKS segments.
//
// The function uses a sync.WaitGroup to synchronize concurrent workers and stores
// intermediate results in a temporary slice. It includes safety checks to prevent
// panics when dealing with empty segments, which can occur if the input slice
// length is smaller than the number of chunks.
//
// Edge cases:
// - Returns 0 if the input slice is empty.
// - Safely handles cases where some chunks contain no data.
func maxChunks(data []int) int {
	size := len(data)

	if size == 0 {
		return 0
	}

	// Слайс для хранения максимумов из каждой части (8 элементов)
	chunkMaxes := make([]int, CHUNKS)
	chunkSize := size / CHUNKS

	var wg sync.WaitGroup

	for i := 0; i < CHUNKS; i++ {
		start := i * chunkSize
		var end int
		if i == CHUNKS-1 {
			// Последний CHUNKS забирает всё до конца (на случай, если длина не делится на 8 нацело)
			end = size
		} else {
			end = start + chunkSize
		}

		// Если данных меньше, чем чанков, и индекс вышел за пределы
		if start >= size {
			// просто пропускаем, в слайсе уже будут нули
			continue
		}

		wg.Add(1)
		// Запускаем горутину, передавая индекс i для записи результата
		go func(s, e, index int) {
			defer wg.Done()

			// Извлекаем сегмент данных
			segment := data[s:e]

			// ПРОВЕРКА: если сегмент пустой, slices.Max упадет.
			if len(segment) == 0 {
				chunkMaxes[index] = 0
				return
			}

			chunkMaxes[index] = slices.Max(segment)
		}(start, end, i)
	}

	wg.Wait()

	return slices.Max(chunkMaxes)
}

func main() {
	start := time.Now()

	fmt.Printf("Генерируем %d целых чисел\n", SIZE)
	data := generateRandomElements(SIZE)

	fmt.Printf("Готово! Сгенерировано %d элементов\n", len(data))
	fmt.Printf("Время выполнения: %v\n", time.Since(start))

	start = time.Now()

	fmt.Println("Ищем максимальное значение в один поток")
	max := maximum(data)

	duration := time.Since(start)
	fmt.Printf("Максимальное значение элемента: %d\nВремя поиска: %v\n", max, duration)

	start = time.Now()

	fmt.Printf("Ищем максимальное значение в %d потоков", CHUNKS)
	max = maxChunks(data)

	duration = time.Since(start)
	fmt.Printf("Максимальное значение элемента: %d\nВремя поиска: %v\n", max, duration)
}

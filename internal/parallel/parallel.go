// Package parallel запускает работу для многих элементов одновременно.
package parallel

import "sync"

// Map вызывает fn для каждого элемента, не больше jobs одновременно.
// Результаты идут в том же порядке, что и items.
func Map[T, R any](items []T, jobs int, fn func(T) R) []R {
	results := make([]R, len(items))

	if jobs < 1 {
		jobs = 1
	}

	var wg sync.WaitGroup

	// Не больше jobs горутин одновременно.
	slots := make(chan struct{}, jobs)

	for i, item := range items {
		slots <- struct{}{}

		wg.Go(func() {
			defer func() { <-slots }()

			results[i] = fn(item)
		})
	}

	wg.Wait()

	return results
}

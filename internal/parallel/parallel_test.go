package parallel

import (
	"slices"
	"sync/atomic"
	"testing"
	"time"
)

func TestMapKeepsOrder(t *testing.T) {
	items := []int{5, 1, 4, 2, 3}

	// Заканчивают в другом порядке, а результаты должны идти как в items.
	got := Map(items, len(items), func(n int) int {
		time.Sleep(time.Duration(n) * time.Millisecond)
		return n * 10
	})

	if want := []int{50, 10, 40, 20, 30}; !slices.Equal(got, want) {
		t.Errorf("Map() = %v, want %v", got, want)
	}
}

func TestMapLimitsConcurrency(t *testing.T) {
	const jobs = 3

	var running, maxRunning atomic.Int32

	Map(make([]int, 20), jobs, func(int) struct{} {
		now := running.Add(1)
		defer running.Add(-1)

		for {
			seen := maxRunning.Load()
			if now <= seen || maxRunning.CompareAndSwap(seen, now) {
				break
			}
		}

		time.Sleep(5 * time.Millisecond)

		return struct{}{}
	})

	if got := maxRunning.Load(); got != jobs {
		t.Errorf("max concurrent calls = %d, want %d", got, jobs)
	}
}

func TestMapJobsBelowOne(t *testing.T) {
	var running, maxRunning atomic.Int32

	got := Map([]int{1, 2, 3}, 0, func(n int) int {
		now := running.Add(1)
		defer running.Add(-1)

		if now > maxRunning.Load() {
			maxRunning.Store(now)
		}

		return n
	})

	if !slices.Equal(got, []int{1, 2, 3}) {
		t.Errorf("Map() = %v, want [1 2 3]", got)
	}

	if maxRunning.Load() != 1 {
		t.Errorf("jobs=0 ran %d calls at once, want 1", maxRunning.Load())
	}
}

func TestMapEmpty(t *testing.T) {
	got := Map([]int{}, 4, func(n int) int { return n })

	if len(got) != 0 {
		t.Errorf("Map() of empty input = %v, want empty", got)
	}
}

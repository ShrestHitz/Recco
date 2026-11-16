package workerpool_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/recco-app/recco/internal/platform/workerpool"
)

func TestPool_BoundedConcurrency(t *testing.T) {
	maxWorkers := 3
	pool := workerpool.New(maxWorkers, 20)
	defer pool.Close()

	var (
		activeWorkers int32
		maxConcurrent int32
		tasksFinished int32
		wg            sync.WaitGroup
	)

	totalTasks := 15
	for i := 0; i < totalTasks; i++ {
		wg.Add(1)
		err := pool.Submit(func(ctx context.Context) {
			defer wg.Done()

			current := atomic.AddInt32(&activeWorkers, 1)
			// Track peak concurrent workers
			for {
				oldMax := atomic.LoadInt32(&maxConcurrent)
				if current <= oldMax || atomic.CompareAndSwapInt32(&maxConcurrent, oldMax, current) {
					break
				}
			}

			time.Sleep(30 * time.Millisecond)
			atomic.AddInt32(&activeWorkers, -1)
			atomic.AddInt32(&tasksFinished, 1)
		})

		if err != nil {
			t.Fatalf("failed submitting task: %v", err)
		}
	}

	wg.Wait()

	if atomic.LoadInt32(&tasksFinished) != int32(totalTasks) {
		t.Errorf("expected %d tasks finished, got %d", totalTasks, tasksFinished)
	}

	peak := atomic.LoadInt32(&maxConcurrent)
	if peak > int32(maxWorkers) {
		t.Errorf("expected max concurrent workers <= %d, but observed peak of %d", maxWorkers, peak)
	}
}

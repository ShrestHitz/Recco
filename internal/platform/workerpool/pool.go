package workerpool

import (
	"context"
	"errors"
	"sync"
)

var (
	// ErrPoolClosed is returned when submitting tasks to a stopped pool.
	ErrPoolClosed = errors.New("worker pool is closed")
)

// Task represents a unit of concurrent work.
type Task func(ctx context.Context)

// Pool orchestrates a fixed number of worker goroutines consuming jobs from a buffered queue.
type Pool struct {
	workers int
	taskCh  chan Task
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
	once    sync.Once
}

// New creates and starts a worker pool with the given number of worker goroutines and buffer capacity.
func New(workers, bufferSize int) *Pool {
	ctx, cancel := context.WithCancel(context.Background())
	p := &Pool{
		workers: workers,
		taskCh:  make(chan Task, bufferSize),
		ctx:     ctx,
		cancel:  cancel,
	}

	p.start()
	return p
}

func (p *Pool) start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for {
				select {
				case <-p.ctx.Done():
					// Drain remaining queued tasks before exiting
					for task := range p.taskCh {
						task(context.Background())
					}
					return
				case task, ok := <-p.taskCh:
					if !ok {
						return
					}
					task(p.ctx)
				}
			}
		}()
	}
}

// Submit queues a task for execution by an available worker.
func (p *Pool) Submit(task Task) error {
	select {
	case <-p.ctx.Done():
		return ErrPoolClosed
	default:
	}

	select {
	case p.taskCh <- task:
		return nil
	case <-p.ctx.Done():
		return ErrPoolClosed
	}
}

// Close gracefully signals workers to finish active jobs and terminates the pool.
func (p *Pool) Close() {
	p.once.Do(func() {
		close(p.taskCh)
		p.cancel()
		p.wg.Wait()
	})
}

package concurrency

import (
	"context"
	"fmt"
	"sync"
)

const DefaultWorkerCount = 5

// ParallelTask represents a named task that accepts context and returns an error
type ParallelTask struct {
	Name string
	Task func(context.Context) error
}

// ParallelOptions defines options for running parallel tasks
type ParallelOptions struct {
	WorkerCount int // Number of concurrent workers
}

// RunParallelTasks executes the given tasks in parallel and returns all errors encountered
// Context cancellation will stop all running tasks
func RunParallelTasks(ctx context.Context, opts ParallelOptions, tasks ...ParallelTask) []error {
	workerCount := opts.WorkerCount
	if workerCount <= 0 {
		workerCount = DefaultWorkerCount // Default worker count
	}

	var (
		wg       sync.WaitGroup
		errChan  = make(chan error, len(tasks)) // Buffered channel for errors
		taskChan = make(chan ParallelTask)      // Channel for distributing tasks to workers
	)

	// Start worker goroutines
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskChan {
				// Check for context cancellation
				select {
				case <-ctx.Done():
					errChan <- fmt.Errorf("task %q canceled: %w", task.Name, ctx.Err())
					continue
				default:
				}

				// Run task with context and collect error
				if err := task.Task(ctx); err != nil {
					errChan <- fmt.Errorf("task %q failed: %w", task.Name, err)
				}
			}
		}()
	}

	// Send tasks to the task channel
	go func() {
		for _, task := range tasks {
			taskChan <- task
		}
		close(taskChan)
	}()

	// Wait for all workers to complete
	wg.Wait()
	close(errChan)

	// Collect errors from the channel
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	return errors
}

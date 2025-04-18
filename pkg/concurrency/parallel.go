package concurrency

import (
	"context"
	"fmt"
	"sync"

	"github.com/hibare/GoCommon/v2/pkg/maps"
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

type ErrorMap map[string]error

// RunParallelTasks executes the given tasks in parallel and returns a map of task names to errors
// Context cancellation will stop all running tasks
func RunParallelTasks(ctx context.Context, opts ParallelOptions, tasks ...ParallelTask) ErrorMap {
	workerCount := opts.WorkerCount
	if workerCount <= 0 {
		workerCount = DefaultWorkerCount
	}

	var (
		wg       sync.WaitGroup
		errorMap sync.Map
		taskChan = make(chan ParallelTask)
	)

	// Start worker goroutines
	for range workerCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskChan {
				// Check for context cancellation
				select {
				case <-ctx.Done():
					errorMap.Store(task.Name, fmt.Errorf("task canceled: %w", ctx.Err()))
					continue
				default:
				}

				// Run task with context and collect error
				if err := task.Task(ctx); err != nil {
					errorMap.Store(task.Name, err)
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

	// Convert sync.Map to regular map
	result := maps.MapFromSyncMap[string, error](&errorMap)

	return result
}

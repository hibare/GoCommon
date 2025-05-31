# Concurrency Package Documentation

## Overview

The `concurrency` package provides utilities for running tasks in parallel using worker pools. It is designed to simplify concurrent execution and error collection for multiple tasks, with support for context cancellation.

---

## Key Types and Functions

- **ParallelTask**: Represents a named task that accepts a context and returns an error.
- **ParallelOptions**: Options for running parallel tasks (e.g., number of workers).
- **ErrorMap**: A map of task names to errors.
- **RunParallelTasks(ctx, opts, tasks...)**: Executes the given tasks in parallel using a worker pool. Returns a map of task names to errors. Cancels all running tasks if the context is canceled.

---

## Example Usage

```go
import (
    "context"
    "github.com/hibare/GoCommon/v2/pkg/concurrency"
)

tasks := []concurrency.ParallelTask{
    {Name: "task1", Task: func(ctx context.Context) error { /* ... */ return nil }},
    {Name: "task2", Task: func(ctx context.Context) error { /* ... */ return nil }},
}
opts := concurrency.ParallelOptions{WorkerCount: 2}
errs := concurrency.RunParallelTasks(context.Background(), opts, tasks...)
for name, err := range errs {
    if err != nil {
        fmt.Printf("%s failed: %v\n", name, err)
    }
}
```

---

## Notes

- The default worker count is 5 if not specified.
- Errors are collected per task and returned as a map.
- Context cancellation is respected and will stop all running tasks.

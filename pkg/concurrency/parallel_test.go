package concurrency

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRunParallelTasks(t *testing.T) {
	tests := []struct {
		name        string
		opts        ParallelOptions
		tasks       []ParallelTask
		wantErrCnt  int
		wantErrText string
		cancelCtx   bool
	}{
		{
			name: "successful tasks",
			opts: ParallelOptions{WorkerCount: 2},
			tasks: []ParallelTask{
				{Name: "task1", Task: func(_ context.Context) error { return nil }},
				{Name: "task2", Task: func(_ context.Context) error { return nil }},
			},
			wantErrCnt: 0,
		},
		{
			name: "failing tasks",
			opts: ParallelOptions{WorkerCount: 2},
			tasks: []ParallelTask{
				{Name: "task1", Task: func(_ context.Context) error { return errors.New("failed task: error1") }},
				{Name: "task2", Task: func(_ context.Context) error { return errors.New("failed task: error2") }},
			},
			wantErrCnt:  2,
			wantErrText: "failed",
		},
		{
			name: "context cancellation",
			opts: ParallelOptions{WorkerCount: 1},
			tasks: []ParallelTask{
				{
					Name: "slow-task",
					Task: func(_ context.Context) error {
						time.Sleep(100 * time.Millisecond)
						return nil
					},
				},
			},
			cancelCtx:   true,
			wantErrCnt:  1,
			wantErrText: "canceled",
		},
		{
			name: "default worker count",
			opts: ParallelOptions{WorkerCount: 0},
			tasks: []ParallelTask{
				{Name: "task1", Task: func(_ context.Context) error { return nil }},
			},
			wantErrCnt: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			if tt.cancelCtx {
				cancel()
				time.Sleep(50 * time.Millisecond) // Give some time for tasks to start
			}

			errsMap := RunParallelTasks(ctx, tt.opts, tt.tasks...)

			require.Len(t, errsMap, tt.wantErrCnt)
			if tt.wantErrText != "" {
				for _, err := range errsMap {
					require.Contains(t, err.Error(), tt.wantErrText)
				}
			}
		})
	}
}

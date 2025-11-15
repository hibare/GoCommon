package datetime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSortDateTimes(t *testing.T) {
	expectedStr := []string{"20230722053000", "20230721053000", "20230720053000"}
	dateTimeStr := []string{"20230721053000", "20230720053000", "20230722053000"}

	sorted := SortDateTimes(dateTimeStr)
	require.Equal(t, expectedStr, sorted)
}

func TestHumanizeTime(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{
			name:     "zero time",
			input:    time.Time{},
			expected: "",
		},
		{
			name:     "seconds ago",
			input:    now.Add(-30 * time.Second),
			expected: "30 seconds ago",
		},
		{
			name:     "1 min ago",
			input:    now.Add(-90 * time.Second), // 1.5 minutes
			expected: "1 min ago",
		},
		{
			name:     "minutes ago",
			input:    now.Add(-10 * time.Minute),
			expected: "10 mins ago",
		},
		{
			name:     "1 hour ago",
			input:    now.Add(-90 * time.Minute),
			expected: "1 hour ago",
		},
		{
			name:     "hours ago",
			input:    now.Add(-5 * time.Hour),
			expected: "5 hours ago",
		},
		{
			name:     "yesterday",
			input:    now.Add(-36 * time.Hour),
			expected: "yesterday",
		},
		{
			name:     "days ago",
			input:    now.Add(-5 * 24 * time.Hour),
			expected: "5 days ago",
		},
		{
			name:     "1 month ago",
			input:    now.Add(-45 * 24 * time.Hour),
			expected: "1 month ago",
		},
		{
			name:     "months ago",
			input:    now.Add(-90 * 24 * time.Hour),
			expected: "3 months ago",
		},
		{
			name:     "1 year ago",
			input:    now.Add(-370 * 24 * time.Hour),
			expected: "1 year ago",
		},
		{
			name:     "years ago",
			input:    now.Add(-3 * 365 * 24 * time.Hour),
			expected: "3 years ago",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := HumanizeTime(tc.input)

			require.NotNil(t, result)
			assert.Equal(t, tc.expected, result)
		})
	}
}

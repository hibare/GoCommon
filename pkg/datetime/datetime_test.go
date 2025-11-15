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
	t.Parallel()

	now := time.Now()

	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{"zero time", time.Time{}, ""},
		{"seconds ago", now.Add(-30 * time.Second), "30 seconds ago"},
		{"1 min ago", now.Add(-90 * time.Second), "1 min ago"},
		{"minutes ago", now.Add(-10 * time.Minute), "10 mins ago"},
		{"1 hour ago", now.Add(-90 * time.Minute), "1 hour ago"},
		{"hours ago", now.Add(-5 * time.Hour), "5 hours ago"},
		{"yesterday", now.Add(-36 * time.Hour), "yesterday"},
		{"days ago", now.Add(-5 * 24 * time.Hour), "5 days ago"},
		{"1 month ago", now.Add(-45 * 24 * time.Hour), "1 month ago"},
		{"months ago", now.Add(-90 * 24 * time.Hour), "3 months ago"},
		{"1 year ago", now.Add(-370 * 24 * time.Hour), "1 year ago"},
		{"years ago", now.Add(-3 * 365 * 24 * time.Hour), "3 years ago"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := HumanizeTime(tt.input)
			require.NotNil(t, result)
			assert.Equal(t, tt.expected, result)
		})
	}
}

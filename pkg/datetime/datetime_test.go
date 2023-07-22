package datetime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortDateTimes(t *testing.T) {
	expectedStr := []string{"20230722053000", "20230721053000", "20230720053000"}
	dateTimeStr := []string{"20230721053000", "20230720053000", "20230722053000"}

	sorted := SortDateTimes(dateTimeStr)
	assert.Equal(t, expectedStr, sorted)
}

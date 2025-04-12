package datetime

import (
	"sort"
	"time"

	"github.com/hibare/GoCommon/v2/pkg/constants"
)

func SortDateTimes(dt []string) []string {
	// Convert the strings to time.Time objects
	var times []time.Time
	for _, dt := range dt {
		t, _ := time.Parse(constants.DefaultDateTimeLayout, dt)
		times = append(times, t)
	}

	// Define a sorting function
	sortFn := func(i, j int) bool {
		return times[i].After(times[j])
	}

	// Sort the slice of time.Time objects
	sort.Slice(times, sortFn)

	// Convert the sorted time.Time objects back to strings
	var sorted []string
	for _, t := range times {
		sorted = append(sorted, t.Format(constants.DefaultDateTimeLayout))
	}

	return sorted
}

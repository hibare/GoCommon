// Package datetime provides utilities for working with date and time strings.
package datetime

import (
	"fmt"
	"sort"
	"time"

	"github.com/hibare/GoCommon/v2/pkg/constants"
)

// SortDateTimes sorts a slice of date-time strings.
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

const (
	hourInDay   = 24
	daysInMonth = 30
	daysInYear  = 365
)

func HumanizeTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return fmt.Sprintf("%d seconds ago", int(diff.Seconds()))
	case diff < 2*time.Minute:
		return "1 min ago"
	case diff < time.Hour:
		return fmt.Sprintf("%d mins ago", int(diff.Minutes()))
	case diff < 2*time.Hour:
		return "1 hour ago"
	case diff < hourInDay*time.Hour:
		return fmt.Sprintf("%d hours ago", int(diff.Hours()))
	case diff < 2*hourInDay*time.Hour:
		return "yesterday"
	case diff < 30*hourInDay*time.Hour:
		return fmt.Sprintf("%d days ago", int(diff.Hours()/hourInDay))
	case diff < 60*hourInDay*time.Hour:
		return "1 month ago"
	case diff < daysInYear*hourInDay*time.Hour:
		return fmt.Sprintf("%d months ago", int(diff.Hours()/(hourInDay*daysInMonth)))
	default:
		years := int(diff.Hours() / (hourInDay * daysInYear))
		if years == 1 {
			return "1 year ago"
		}
		return fmt.Sprintf("%d years ago", years)
	}
}

package env

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {
	// Set Sample environment variables
	t.Setenv("STRING_ENV", "test_string")
	t.Setenv("BOOL_ENV", "true")
	t.Setenv("INT_ENV", "100")
	t.Setenv("DURATION_ENV", "24h")
	t.Setenv("SLICE_ENV", "1,2,3")
	t.Setenv("PREFIXED_ENV_1", "value1")
	t.Setenv("PREFIXED_ENV_2", "value2")

	assert.Equal(t, "test_string", MustString("STRING_ENV", ""))
	assert.True(t, MustBool("BOOL_ENV", false))
	assert.Equal(t, 100, MustInt("INT_ENV", 0))
	assert.Equal(t, 24*time.Hour, MustDuration("DURATION_ENV", time.Duration(0)))
	assert.Equal(t, []string{"1", "2", "3"}, MustStringSlice("SLICE_ENV", []string{}))

	assert.Equal(t, "default_string", MustString("DEFAULT_STRING_ENV", "default_string"))
	assert.Equal(t, false, MustBool("DEFAULT_BOOL_ENV", false))
	assert.Equal(t, 1, MustInt("DEFAULT_INT_ENV", 1))
	assert.Equal(t, time.Duration(1), MustDuration("DEFAULT_DURATION_ENV", time.Duration(1)))
	assert.Equal(t, []string{"1"}, MustStringSlice("DEFAULT_SLICE_ENV", []string{"1"}))

	prefixed := GetPrefixed("PREFIXED_ENV_")
	assert.Equal(t, map[string]string{
		"PREFIXED_ENV_1": "value1",
		"PREFIXED_ENV_2": "value2",
	}, prefixed)
}

func TestEnv_EdgeCases(t *testing.T) {
	// Invalid bool, int, duration
	t.Setenv("INVALID_BOOL_ENV", "notabool")
	t.Setenv("INVALID_INT_ENV", "notanint")
	t.Setenv("INVALID_DURATION_ENV", "notaduration")

	assert.Equal(t, false, MustBool("INVALID_BOOL_ENV", false))
	assert.Equal(t, 42, MustInt("INVALID_INT_ENV", 42))
	assert.Equal(t, 5*time.Second, MustDuration("INVALID_DURATION_ENV", 5*time.Second))

	// Empty string slice
	t.Setenv("EMPTY_SLICE_ENV", "")
	assert.Equal(t, []string{"fallback"}, MustStringSlice("EMPTY_SLICE_ENV", []string{"fallback"}))

	// Malformed slice (should just split as usual)
	t.Setenv("MALFORMED_SLICE_ENV", ",,a,,b,,")
	assert.Equal(t, []string{"", "", "a", "", "b", "", ""}, MustStringSlice("MALFORMED_SLICE_ENV", []string{"fallback"}))

	// GetPrefixed with no matches
	unsetPrefix := "UNSET_PREFIX_"
	assert.Equal(t, map[string]string{}, GetPrefixed(unsetPrefix))
}

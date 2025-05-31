package env

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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

	require.Equal(t, "test_string", MustString("STRING_ENV", ""))
	require.True(t, MustBool("BOOL_ENV", false))
	require.Equal(t, 100, MustInt("INT_ENV", 0))
	require.Equal(t, 24*time.Hour, MustDuration("DURATION_ENV", time.Duration(0)))
	require.Equal(t, []string{"1", "2", "3"}, MustStringSlice("SLICE_ENV", []string{}))

	require.Equal(t, "default_string", MustString("DEFAULT_STRING_ENV", "default_string"))
	require.False(t, MustBool("DEFAULT_BOOL_ENV", false))
	require.Equal(t, 1, MustInt("DEFAULT_INT_ENV", 1))
	require.Equal(t, time.Duration(1), MustDuration("DEFAULT_DURATION_ENV", time.Duration(1)))
	require.Equal(t, []string{"1"}, MustStringSlice("DEFAULT_SLICE_ENV", []string{"1"}))

	prefixed := GetPrefixed("PREFIXED_ENV_")
	require.Equal(t, map[string]string{
		"PREFIXED_ENV_1": "value1",
		"PREFIXED_ENV_2": "value2",
	}, prefixed)
}

func TestEnv_EdgeCases(t *testing.T) {
	// Invalid bool, int, duration
	t.Setenv("INVALID_BOOL_ENV", "notabool")
	t.Setenv("INVALID_INT_ENV", "notanint")
	t.Setenv("INVALID_DURATION_ENV", "notaduration")

	require.False(t, MustBool("INVALID_BOOL_ENV", false))
	require.Equal(t, 42, MustInt("INVALID_INT_ENV", 42))
	require.Equal(t, 5*time.Second, MustDuration("INVALID_DURATION_ENV", 5*time.Second))

	// Empty string slice
	t.Setenv("EMPTY_SLICE_ENV", "")
	require.Equal(t, []string{"fallback"}, MustStringSlice("EMPTY_SLICE_ENV", []string{"fallback"}))

	// Malformed slice (should just split as usual)
	t.Setenv("MALFORMED_SLICE_ENV", ",,a,,b,,")
	require.Equal(t, []string{"", "", "a", "", "b", "", ""}, MustStringSlice("MALFORMED_SLICE_ENV", []string{"fallback"}))

	// GetPrefixed with no matches
	unsetPrefix := "UNSET_PREFIX_"
	require.Equal(t, map[string]string{}, GetPrefixed(unsetPrefix))
}

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHostname(t *testing.T) {
	hostname := GetHostname()
	assert.NotEmpty(t, hostname)
}

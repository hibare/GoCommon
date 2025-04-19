package http

import "time"

const (
	// DefaultServerReadTimeout is the maximum duration for reading the entire request.
	DefaultServerReadTimeout = 15 * time.Second

	// DefaultServerWriteTimeout is the maximum duration before timing out writes of the response.
	DefaultServerWriteTimeout = 15 * time.Second

	// DefaultServerIdleTimeout is the maximum amount of time to wait for the next request.
	DefaultServerIdleTimeout = 60 * time.Second

	// DefaultServerShutdownGracePeriod is the duration to wait for server shutdown.
	DefaultServerShutdownGracePeriod = 60 * time.Second

	// DefaultServerTimeout is the overall timeout for server operations.
	DefaultServerTimeout = 60 * time.Second

	// DefaultServerMaxHeaderBytes is the maximum size of request headers.
	HTTPDefaultRequestSize = 1024 * 1024 * 5 // 5MB
)

// API endpoint paths.
const (
	// PingPath is the endpoint for health check pings.
	DefaultPingPath = "/ping"
)

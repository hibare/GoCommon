# HTTP Package Documentation

## Overview

The `http` package provides utilities for HTTP server and client operations, including constants, middleware, handlers, and a customizable HTTP client. It is designed to simplify building and testing HTTP APIs in Go applications.

---

## Subpackages

- **client**: Provides a customizable HTTP client interface and implementation.
- **handler**: Contains HTTP handlers, such as health checks.
- **middleware**: Provides middleware for authentication, logging, and security headers.

---

## Key Types and Functions

- **DefaultServerReadTimeout, DefaultServerWriteTimeout, ...**: Constants for server/client timeouts and limits.
- **WriteJSONResponse(w, statusCode, data)**: Writes a JSON response.
- **WriteErrorResponse(w, statusCode, err)**: Writes a JSON error response.

### Client

- **Client**: Interface for HTTP clients (Do method).
- **NewDefaultClient()**: Returns a default HTTP client with sensible timeouts.

### Middleware

- **TokenAuth(next, tokens)**: Middleware for token-based authentication.
- **RequestLogger(next)**: Middleware for logging requests and responses.
- **BasicSecurity(next, sizeBytes)**: Middleware for adding security headers and request size limits.

### Handler

- **HealthCheck(w, r)**: Simple health check handler.

---

## Example Usage

```go
import (
    "net/http"
    "github.com/hibare/GoCommon/v2/pkg/http/handler"
)

http.HandleFunc("/ping", handler.HealthCheck)
```

---

## Notes

- Middleware and handlers are designed for easy integration with the standard `net/http` package.
- The client and middleware are mockable for testing.

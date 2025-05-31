# Context Package Documentation

## Overview

The `context` package provides a custom context implementation for request-scoped values, extending Go's standard `context.Context` with additional helpers for request IDs and value management.

---

## Key Types and Functions

- **Context**: Struct wrapping a standard `context.Context` and allowing storage of additional values.
- **NewContext()**: Creates a new custom context.
- **WithValue(key, value)**: Adds a key-value pair to the context.
- **WithRequestID(requestID)**: Adds a request ID to the context.
- **WithContext(ctx)**: Sets the underlying context.
- **Get(key)**: Retrieves a value by key.
- **GetContext()**: Returns the underlying `context.Context`.
- **GetRequestID()**: Retrieves the request ID from the context.

---

## Example Usage

```go
import (
    "github.com/hibare/GoCommon/v2/pkg/context"
)

ctx := context.NewContext().WithRequestID("abc-123")
requestID := ctx.GetRequestID()
```

---

## Notes

- Useful for propagating request-scoped values (like request IDs) in web applications.
- Fully compatible with Go's standard context API.

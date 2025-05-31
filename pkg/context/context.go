// Package context provides a custom context implementation for request-scoped values.
package context

import (
	"context"
)

// ContextKey is a type for context keys used in the custom Context.
type ContextKey string

// Context wraps a standard context and allows storing additional values.
type Context struct {
	Context context.Context
}

// NewContext creates a new custom Context.
func NewContext() *Context {
	return &Context{
		Context: context.Background(),
	}
}

// WithValue returns a copy of the Context with the provided key-value pair.
func (c *Context) WithValue(key, value any) *Context {
	c.Context = context.WithValue(c.Context, key, value)
	return c
}

// WithRequestID returns a copy of the Context with the provided request ID.
func (c *Context) WithRequestID(requestID string) *Context {
	return c.WithValue("request_id", requestID)
}

// WithContext sets the underlying context.Context.
func (c *Context) WithContext(ctx context.Context) *Context {
	c.Context = ctx
	return c
}

// Get retrieves a value by key from the Context.
func (c *Context) Get(key any) any {
	return c.Context.Value(key)
}

// GetContext returns the underlying context.Context.
func (c *Context) GetContext() context.Context {
	return c.Context
}

// GetRequestID returns the request ID from the Context.
func (c *Context) GetRequestID() string {
	return c.Context.Value("request_id").(string)
}

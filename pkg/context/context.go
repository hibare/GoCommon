package context

import (
	"context"
)

type ContextKey string

type Context struct {
	Context context.Context
}

func NewContext() *Context {
	return &Context{
		Context: context.Background(),
	}
}

func (c *Context) WithValue(key, value any) *Context {
	c.Context = context.WithValue(c.Context, key, value)
	return c
}

func (c *Context) WithRequestID(requestID string) *Context {
	return c.WithValue("request_id", requestID)
}

func (c *Context) WithContext(ctx context.Context) *Context {
	c.Context = ctx
	return c
}

func (c *Context) Get(key any) any {
	return c.Context.Value(key)
}

func (c *Context) GetContext() context.Context {
	return c.Context
}

func (c *Context) GetRequestID() string {
	return c.Context.Value("request_id").(string)
}

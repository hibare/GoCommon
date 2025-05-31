package context

import (
	"context"
	"testing"
)

func TestNewContext(t *testing.T) {
	ctx := NewContext()
	if ctx.Context == nil {
		t.Error("Expected non-nil context")
	}
}

func TestWithValue(t *testing.T) {
	ctx := NewContext().WithValue("key", "value")
	if ctx.Get("key") != "value" {
		t.Errorf("Expected value 'value', got %v", ctx.Get("key"))
	}
}

func TestWithRequestID(t *testing.T) {
	requestID := "12345"
	ctx := NewContext().WithRequestID(requestID)
	if ctx.GetRequestID() != requestID {
		t.Errorf("Expected request ID '%s', got %s", requestID, ctx.GetRequestID())
	}
}

func TestWithContext(t *testing.T) {
	key := Key("key")

	baseCtx := context.WithValue(t.Context(), key, "value")
	ctx := NewContext().WithContext(baseCtx)
	if ctx.Get(key) != "value" {
		t.Errorf("Expected value 'value', got %v", ctx.Get(key))
	}
}

func TestGet(t *testing.T) {
	ctx := NewContext().WithValue("key", "value")
	if ctx.Get("key") != "value" {
		t.Errorf("Expected value 'value', got %v", ctx.Get("key"))
	}
}

func TestGetContext(t *testing.T) {
	ctx := NewContext()
	if ctx.GetContext() == nil {
		t.Error("Expected non-nil context")
	}
}

func TestGetRequestID(t *testing.T) {
	requestID := "12345"
	ctx := NewContext().WithRequestID(requestID)
	if ctx.GetRequestID() != requestID {
		t.Errorf("Expected request ID '%s', got %s", requestID, ctx.GetRequestID())
	}
}

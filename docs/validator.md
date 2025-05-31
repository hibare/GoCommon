# Validator Package Documentation

## Overview

The `validator` package provides utilities for validating Go structs using [go-playground/validator](https://github.com/go-playground/validator). It supports custom error messages via struct tags and can use JSON tags for error reporting.

---

## Key Functions

- **ValidateStructErrors(obj, validate, useJSONTag) error**: Validates a struct and returns a list of errors, using custom error messages from struct tags if present.

---

## Struct Tag Support

- `validate`: Specifies validation rules (e.g., `required,gt=0`).
- `validate_errs`: Specifies custom error messages for each rule.
- `json`: Used for error field names if `useJSONTag` is true.

---

## Example Usage

```go
import (
    "github.com/go-playground/validator/v10"
    "github.com/hibare/GoCommon/v2/pkg/validator"
)

type MyStruct struct {
    Name string `json:"name" validate:"required" validate_errs:"Name is required"`
}
validate := validator.New()
err := validator.ValidateStructErrors[MyStruct](MyStruct{}, validate, true)
```

---

## Notes

- Returns detailed error messages, including custom messages from struct tags.
- Useful for API request validation and form validation.

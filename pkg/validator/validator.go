// Package validator provides utilities for validating structs.
package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/hibare/GoCommon/v2/pkg/slice"
)

const (
	tagJSON            = "json"
	tagValidateErrMsgs = "validate_errs"
	tagValidate        = "validate"
)

func extractTagAsSlice(field reflect.StructField, tagName string) []string {
	var tagSlice []string
	tag := field.Tag.Get(tagName)
	if tag != "" {
		tagSlice = strings.Split(tag, ",")
		for i, t := range tagSlice {
			tagSlice[i] = strings.TrimSpace(t)
		}
	}

	return tagSlice
}

func getFieldOrTag(field reflect.StructField, useJSON bool) string {
	tag := field.Tag.Get(tagJSON)
	if useJSON && tag != "" && tag != "-" {
		return tag
	}
	return field.Name
}

// ValidateStructErrors validates a struct and returns a list of errors.
func ValidateStructErrors[T any](obj any, validate *validator.Validate, useJSONTag bool) (errs error) {
	defer func() {
		if r := recover(); r != nil {
			errs = fmt.Errorf("unable to validate %+v", r)
		}
	}()

	err := validate.Struct(obj)
	if err == nil {
		return nil
	}

	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		errs = fmt.Errorf("unexpected error during validation: %w", err)
	}
	for _, e := range validationErrors {
		if field, found := reflect.TypeOf(obj).FieldByName(e.Field()); found {
			fieldTag := getFieldOrTag(field, useJSONTag)
			validateTags := extractTagAsSlice(field, tagValidate)
			errMsgs := extractTagAsSlice(field, tagValidateErrMsgs)

			if len(errMsgs) == 0 {
				errs = errors.Join(errs, fmt.Errorf("%s: %w", fieldTag, e))
				continue
			}

			validateTagIndex := slice.IndexOf(e.Tag(), validateTags)
			if validateTagIndex == -1 {
				errs = errors.Join(errs, fmt.Errorf("%s: %s (%s)", fieldTag, strings.Join(errMsgs, ", "), e.Tag()))
			} else {
				errs = errors.Join(errs, fmt.Errorf("%s: %s (%s)", fieldTag, errMsgs[validateTagIndex], e.Tag()))
			}
		}
	}

	return errs
}

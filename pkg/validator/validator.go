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
	tagJson            = "json"
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

func getFieldOrTag(field reflect.StructField, useJson bool) string {
	tag := field.Tag.Get(tagJson)
	if useJson && tag != "" && tag != "-" {
		return tag
	}
	return field.Name
}

func ValidateStructErrors[T any](obj any, validate *validator.Validate, useJsonTag bool) (errs error) {
	defer func() {
		if r := recover(); r != nil {
			errs = fmt.Errorf("Unable to validate %+v", r)
		}
	}()

	err := validate.Struct(obj)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, e := range validationErrors {
				if field, found := reflect.TypeOf(obj).FieldByName(e.Field()); found {
					fieldTag := getFieldOrTag(field, useJsonTag)
					validateTags := extractTagAsSlice(field, tagValidate)
					errMsgs := extractTagAsSlice(field, tagValidateErrMsgs)

					if len(errMsgs) == 0 {
						errs = errors.Join(errs, fmt.Errorf("%s: %w", fieldTag, e))
					} else {
						validateTagIndex := slice.SliceIndexOf(e.Tag(), validateTags)
						if validateTagIndex == -1 {
							errs = errors.Join(errs, fmt.Errorf("%s: %s (%s)", fieldTag, strings.Join(errMsgs, ", "), e.Tag()))
						} else {
							errs = errors.Join(errs, fmt.Errorf("%s: %s (%s)", fieldTag, errMsgs[validateTagIndex], e.Tag()))
						}
					}
				}
			}
		} else {
			errs = fmt.Errorf("unexpected error during validation: %w", err)
		}
	}

	return errs
}

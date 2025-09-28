// Package structs provides utilities for working with Go structs.
package structs

import (
	"errors"
	"reflect"
)

// Copy copies the contents of src struct to dst struct.
func Copy(src, dst any) error {
	srcVal := reflect.ValueOf(src)
	dstVal := reflect.ValueOf(dst)

	if srcVal.Kind() != reflect.Ptr || dstVal.Kind() != reflect.Ptr {
		return errors.New("both src and dst must be pointers to structs")
	}

	srcElem := srcVal.Elem()
	dstElem := dstVal.Elem()

	if srcElem.Kind() != reflect.Struct || dstElem.Kind() != reflect.Struct {
		return errors.New("both src and dst must be pointers to structs")
	}

	dstElem.Set(srcElem)
	return nil
}

// Compare compares two structs for equality.
func Compare(a, b any) (bool, error) {
	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)

	if aVal.Kind() != reflect.Struct || bVal.Kind() != reflect.Struct {
		return false, errors.New("both a and b must be structs")
	}

	return reflect.DeepEqual(a, b), nil
}

// ToMap converts a struct to a map.
func ToMap(s any) (map[string]any, error) {
	val := reflect.ValueOf(s)
	if val.Kind() != reflect.Struct {
		return nil, errors.New("input must be a struct")
	}

	result := make(map[string]any)
	typ := val.Type()

	for i := range val.NumField() {
		field := val.Field(i)
		fieldName := typ.Field(i).Name
		result[fieldName] = field.Interface()
	}

	return result, nil
}

// ContainsField checks if a struct contains a field with the given name.
func ContainsField(s any, fieldName string) (bool, error) {
	val := reflect.ValueOf(s)
	if val.Kind() != reflect.Struct {
		return false, errors.New("input must be a struct")
	}

	// Check if the struct has the field
	field := val.FieldByName(fieldName)
	return field.IsValid(), nil
}

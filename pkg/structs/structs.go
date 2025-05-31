// Package structs provides utilities for working with Go structs.
package structs

import (
	"errors"
	"reflect"
)

// CopyStruct copies the contents of src struct to dst struct.
// Deprecated: Use StructCopy instead.
func CopyStruct(src, dst interface{}) error {
	return StructCopy(src, dst)
}

// StructCopy copies the contents of src struct to dst struct.
func StructCopy(src, dst interface{}) error {
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

// StructCompare compares two structs for equality.
func StructCompare(a, b interface{}) (bool, error) {
	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)

	if aVal.Kind() != reflect.Struct || bVal.Kind() != reflect.Struct {
		return false, errors.New("both a and b must be structs")
	}

	return reflect.DeepEqual(a, b), nil
}

// StructToMap converts a struct to a map.
func StructToMap(s interface{}) (map[string]interface{}, error) {
	val := reflect.ValueOf(s)
	if val.Kind() != reflect.Struct {
		return nil, errors.New("input must be a struct")
	}

	result := make(map[string]interface{})
	typ := val.Type()

	for i := range val.NumField() {
		field := val.Field(i)
		fieldName := typ.Field(i).Name
		result[fieldName] = field.Interface()
	}

	return result, nil
}

// StructContainsField checks if a struct contains a field with the given name.
func StructContainsField(s interface{}, fieldName string) (bool, error) {
	val := reflect.ValueOf(s)
	if val.Kind() != reflect.Struct {
		return false, errors.New("input must be a struct")
	}

	// Check if the struct has the field
	field := val.FieldByName(fieldName)
	return field.IsValid(), nil
}

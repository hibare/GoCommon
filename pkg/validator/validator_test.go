package validator

import (
	"reflect"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

type TestStruct struct {
	Name     string `json:"name" validate:"required"      validate_errs:"Name is required"`
	Age      int    `json:"age"  validate:"required,gt=0" validate_errs:"Age is required,Age must be greater than 0"`
	Email    string `            validate:"email"`
	Internal string `json:"-"`
}

func Test_extractTagAsSlice(t *testing.T) {
	type args struct {
		field   reflect.StructField
		tagName string
	}

	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Empty tag",
			args: args{
				field:   reflect.TypeOf(TestStruct{}).Field(2),
				tagName: "validate_errs",
			},
			want: nil,
		},
		{
			name: "Single value tag",
			args: args{
				field:   reflect.TypeOf(TestStruct{}).Field(0),
				tagName: "validate_errs",
			},
			want: []string{"Name is required"},
		},
		{
			name: "Multiple value tag",
			args: args{
				field:   reflect.TypeOf(TestStruct{}).Field(1),
				tagName: "validate_errs",
			},
			want: []string{"Age is required", "Age must be greater than 0"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractTagAsSlice(tt.args.field, tt.args.tagName)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_getFieldOrTag(t *testing.T) {
	type args struct {
		field   reflect.StructField
		useJSON bool
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Use field name when useJson is false",
			args: args{
				field:   reflect.TypeOf(TestStruct{}).Field(0),
				useJSON: false,
			},
			want: "Name",
		},
		{
			name: "Use json tag when useJson is true",
			args: args{
				field:   reflect.TypeOf(TestStruct{}).Field(0),
				useJSON: true,
			},
			want: "name",
		},
		{
			name: "Use field name when json tag is -",
			args: args{
				field:   reflect.TypeOf(TestStruct{}).Field(3),
				useJSON: true,
			},
			want: "Internal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getFieldOrTag(tt.args.field, tt.args.useJSON)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestValidateStructErrors(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name      string
		obj       TestStruct
		useJSON   bool
		wantError bool
	}{
		{
			name: "Valid struct",
			obj: TestStruct{
				Name:  "John",
				Age:   25,
				Email: "john@example.com",
			},
			useJSON:   true,
			wantError: false,
		},
		{
			name:      "Invalid struct - missing required fields",
			obj:       TestStruct{},
			useJSON:   true,
			wantError: true,
		},
		{
			name: "Invalid email",
			obj: TestStruct{
				Name:  "John",
				Age:   25,
				Email: "invalid-email",
			},
			useJSON:   false,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStructErrors[TestStruct](tt.obj, validate, tt.useJSON)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

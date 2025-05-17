package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Age:    20,
				Email:  "test@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{},
		},
		{
			in: User{
				ID:     "12345678",
				Age:    17,
				Email:  "invalid-email",
				Role:   "user",
				Phones: []string{"12345"},
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: fmt.Errorf("invalid string length: 12345678, should be: 36 but got: 8")},
				{Field: "Age", Err: fmt.Errorf("invalid value: 17, should be greater or equal: 18")},
				{Field: "Email", Err: fmt.Errorf("invalid string value: invalid-email, should match regexp: ^\\w+@\\w+\\.\\w+$")},
				{Field: "Role", Err: fmt.Errorf("invalid value: user, should be in range: [admin stuff]")},
				{Field: "Phones", Err: fmt.Errorf("invalid string length: 12345, should be: 11 but got: 5")},
			},
		},
		{
			in: App{
				Version: "12345",
			},
			expectedErr: ValidationErrors{},
		},
		{
			in: App{
				Version: "1234",
			},
			expectedErr: ValidationErrors{
				{Field: "Version", Err: fmt.Errorf("invalid string length: 1234, should be: 5 but got: 4")},
			},
		},
		{
			in: Response{
				Code: 200,
			},
			expectedErr: ValidationErrors{},
		},
		{
			in: Response{
				Code: 403,
			},
			expectedErr: ValidationErrors{
				{Field: "Code", Err: fmt.Errorf("invalid value: 403, should be in range: [200 404 500]")},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if tt.expectedErr == nil {
				assert.NoError(t, err, "unexpected error: %s", err)
			} else {
				assert.True(t, reflect.DeepEqual(tt.expectedErr, err), "expected error: %v, but got: %v", tt.expectedErr, err)
			}
		})
	}
}

package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	UserS struct {
		ID      string `json:"id" validate:"len:36"`
		Name    string
		Age     int      `validate:"min:18|max:50"`
		Email   string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role    UserRole `validate:"in:admin,stuff"`
		Phones  []string `validate:"len:11"`
		Version *App     `validate:"dive"`
	}
	App struct {
		Version string `validate:"len:5"`
	}

	IncorrectReg struct {
		Value string `validate:"regexp:["`
	}

	InvalidTagMethod struct {
		Value string `validate:"ulen:5"`
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
		in           interface{}
		expectedErrs []error
	}{
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Penny",
				Age:    14,
				Email:  "penny85@gmailcom",
				Role:   "stuf",
				Phones: []string{"7926123456"},
				meta:   []byte(`{"gender":"female"}`),
			},
			expectedErrs: []error{ErrMin, ErrRegexp, ErrInclude, ErrLen},
		},
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Penny",
				Age:    18,
				Email:  "penny85@gmailcom",
				Role:   "stuf",
				Phones: []string{"7926123456"},
				meta:   []byte(`{"gender":"female"}`),
			},
			expectedErrs: []error{ErrRegexp, ErrInclude, ErrLen},
		},
		{
			in: App{
				Version: "123456789012345678901234567890123456",
			},
			expectedErrs: []error{ErrLen},
		},
		{
			in: &App{
				Version: "123456789012345678901234567890123456",
			},
			expectedErrs: []error{ErrLen},
		},
		{
			in: App{
				Version: "12345",
			},
			expectedErrs: nil,
		},
		{
			in: &App{
				Version: "12345",
			},
			expectedErrs: nil,
		},
		{
			in:           []string{"1", "2"},
			expectedErrs: []error{ErrNotStruct},
		},
		{
			in: Token{
				Header:    []byte("12345"),
				Payload:   []byte("12345"),
				Signature: []byte("12345"),
			},
			expectedErrs: nil,
		},
		{
			in: Response{
				Code: 303,
				Body: `{"gender":"female"}`,
			},
			expectedErrs: []error{ErrInclude},
		},
		{
			in: Response{
				Code: 404,
				Body: `{"gender":"female"}`,
			},
			expectedErrs: nil,
		},
		{
			in: IncorrectReg{
				Value: "123456789012345678901234567890123456",
			},
			expectedErrs: []error{ErrRegexpCompile},
		},
		{
			in: InvalidTagMethod{
				Value: "123456789012345678901234567890123456",
			},
			expectedErrs: []error{ErrValidMethod},
		},
		{
			in: UserS{
				ID:      "123456789012345678901234567890123456",
				Name:    "Penny",
				Age:     19,
				Email:   "penny85@gmail.com",
				Role:    "stuff",
				Phones:  []string{"79261234567"},
				Version: &App{Version: "12345667"},
			},
			expectedErrs: []error{ErrLen},
		},
		// ...
		// Place your code here.
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			_ = tt

			if tt.expectedErrs == nil {
				require.NoError(t, err)
			} else {
				var e *ValidationErrors
				if errors.As(err, &e) {
					for i := 0; i < len(e.Unwrap()); i++ {
						require.EqualError(t, e.Unwrap()[i], tt.expectedErrs[i].Error())
					}
				} else {
					require.EqualError(t, errors.Unwrap(err), tt.expectedErrs[0].Error())
				}
			}
		})
	}
}

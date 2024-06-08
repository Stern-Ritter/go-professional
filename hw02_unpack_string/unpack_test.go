package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	type want struct {
		str string
		err error
	}
	tests := []struct {
		input string
		want  want
	}{
		{input: "a4bc2d5e", want: want{str: "aaaabccddddde"}},
		{input: "abccd", want: want{str: "abccd"}},
		{input: "", want: want{str: ""}},
		{input: "aaa0b", want: want{str: "aab"}},
		// uncomment if task with asterisk completed
		{input: `qwe\4\5`, want: want{str: `qwe45`}},
		{input: `qwe\45`, want: want{str: `qwe44444`}},
		{input: `qwe\\5`, want: want{str: `qwe\\\\\`}},
		{input: `qwe\\\3`, want: want{str: `qwe\3`}},
		// additional from the description of the task on the otus.ru
		{input: "abcd", want: want{str: "abcd"}},
		{input: "3abc", want: want{str: "", err: ErrInvalidString}},
		{input: "45", want: want{str: "", err: ErrInvalidString}},
		{input: "aaa10b", want: want{str: "", err: ErrInvalidString}},
		{input: "aaa0b", want: want{str: "aab"}},
		// additional from readme.md
		{input: "qw\ne", want: want{str: "", err: ErrInvalidString}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.input, func(t *testing.T) {
			str, err := Unpack(tt.input)
			require.ErrorIs(t, tt.want.err, err)
			require.Equal(t, tt.want.str, str)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}

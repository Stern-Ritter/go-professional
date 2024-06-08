package hw02unpackstring

import (
	"errors"
	"strconv"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(value string) (string, error) {
	chars := []rune(value)
	acc := make([]rune, 0, len(chars))
	lastSlashes := make([]rune, 0)
	lastChars := make([]rune, 0)

	for _, char := range chars {
		switch {
		case char == '\\':
			acc = append(acc, lastChars...)
			lastChars = lastChars[:0]

			if len(lastSlashes) != 0 {
				lastChars = append(lastChars, char)
				lastSlashes = lastSlashes[:0]
			} else {
				lastSlashes = append(lastSlashes, char)
			}

		case unicode.IsLetter(char):
			acc = append(acc, lastChars...)
			lastChars = lastChars[:0]

			lastChars = append(lastChars, char)

		case unicode.IsNumber(char):
			if len(lastSlashes) == 0 && len(lastChars) == 0 {
				return "", ErrInvalidString
			}

			if len(lastSlashes) != 0 {
				lastChars = append(lastChars, char)
				lastSlashes = lastSlashes[:0]
			} else {
				number, _ := strconv.Atoi(string(char))
				for i := 0; i < number; i++ {
					acc = append(acc, lastChars...)
				}
				lastChars = lastChars[:0]
			}

		default:
			return "", ErrInvalidString
		}
	}

	acc = append(acc, lastChars...)
	return string(acc), nil
}

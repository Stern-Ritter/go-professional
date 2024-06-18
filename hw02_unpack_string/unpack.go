package hw02unpackstring

import (
	"errors"
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/rivo/uniseg" //nolint:depguard
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(value string) (string, error) {
	acc := make([]rune, 0, utf8.RuneCountInString(value))
	lastSlashes := make([]rune, 0)
	lastChars := make([]rune, 0)

	graphemes := uniseg.NewGraphemes(value)
	for graphemes.Next() {
		chars := graphemes.Runes()
		switch {
		case len(chars) == 1:
			char := chars[0]
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
				acc = append(acc, lastChars...)
				lastChars = lastChars[:0]

				lastChars = append(lastChars, char)
			}

		default:
			acc = append(acc, lastChars...)
			lastChars = lastChars[:0]

			lastChars = append(lastChars, chars...)
		}
	}

	acc = append(acc, lastChars...)
	return string(acc), nil
}

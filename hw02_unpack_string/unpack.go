package hw02unpackstring

import (
	"errors"
	"strings"
)

const (
	minDigit          = '0'
	maxDigit          = '9'
	emptyLetter  rune = 0
	escapeLetter rune = '\\'
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var sb strings.Builder
	lastLetter := emptyLetter
	isCharEscaped := false

	for _, ch := range str {
		switch {
		case isCharEscaped:
			lastLetter = ch
			isCharEscaped = false

		case ch == escapeLetter:
			if lastLetter != emptyLetter {
				sb.WriteRune(lastLetter)
			}
			lastLetter = emptyLetter
			isCharEscaped = true

		case ch >= minDigit && ch <= maxDigit:
			if lastLetter == emptyLetter {
				return "", ErrInvalidString
			}

			count := int(ch - minDigit)
			for i := 0; i < count; i++ {
				sb.WriteRune(lastLetter)
			}

			lastLetter = emptyLetter

		default:
			if lastLetter != emptyLetter {
				sb.WriteRune(lastLetter)
			}
			lastLetter = ch
		}
	}

	if lastLetter != emptyLetter {
		sb.WriteRune(lastLetter)
	}

	return sb.String(), nil
}

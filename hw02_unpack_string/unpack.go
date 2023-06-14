package hw02unpackstring

import (
	"errors"
	"strconv"
	str "strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) { //nolint:gocognit
	var b str.Builder
	var pr string
	var isDigit bool
	for _, v := range s {
		iv, err := strconv.Atoi(string(v))
		isDigit = err == nil
		switch {
		case pr == "" && isDigit:
			return "", ErrInvalidString
		case pr == "" && !isDigit:
			pr = string(v)
		case pr == "\\" && string(v) == "\\":
			pr = "\\\\"
		case pr == "\\" && isDigit:
			pr = string(v)
		case pr == "\\" && !isDigit:
			return "", ErrInvalidString
		case pr == "\\\\" && isDigit:
			b.WriteString(str.Repeat("\\", iv))
			pr = ""
		case pr == "\\\\" && !isDigit:
			b.WriteString(str.Repeat("\\", 1))
			pr = string(v)
		case pr != "" && isDigit:
			b.WriteString(str.Repeat(pr, iv))
			pr = ""
		case pr != "" && !isDigit:
			b.WriteString(str.Repeat(pr, 1))
			pr = string(v)
		}
	}
	switch {
	case pr == "\\\\":
		b.WriteString(str.Repeat("\\", 1))
	case pr != "":
		b.WriteString(str.Repeat(pr, 1))
	}
	return b.String(), nil
}

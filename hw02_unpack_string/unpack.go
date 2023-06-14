package hw02unpackstring

import (
	"errors"
	"strconv"
	str "strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var b str.Builder
	var pr string
	for _, v := range s {
		iv, err := strconv.Atoi(string(v))
		// fmt.Printf("pr: %q ... v: %q\n", pr, v)
		switch {
		case pr == "" && err == nil:
			return "", ErrInvalidString
		case pr == "" && err != nil:
			pr = string(v)
		case pr == "\\" && string(v) == "\\":
			pr = "\\\\"
		case pr == "\\" && err == nil:
			pr = string(v)
		case pr == "\\" && err != nil:
			return "", ErrInvalidString
		case pr == "\\\\" && err == nil:
			b.WriteString(str.Repeat("\\", iv))
			pr = ""
		case pr == "\\\\" && err != nil:
			b.WriteString(str.Repeat("\\", 1))
			pr = string(v)
		case pr != "" && err == nil:
			b.WriteString(str.Repeat(pr, iv))
			pr = ""
		case pr != "" && err != nil:
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

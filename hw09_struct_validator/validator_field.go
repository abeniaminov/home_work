package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrLen               = errors.New("wrong len")
	ErrRegexp            = errors.New("value is not match to regexp")
	ErrInclude           = errors.New("wrong value in enumeration")
	ErrMin               = errors.New("value is less then min")
	ErrMax               = errors.New("value is more then max")
	ErrUnsupportableType = errors.New("unsuppotable type of value")
)

func checkLen(field string, v interface{}, length string) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.String {
		return &ValidationError{
			Field: field,
			Err:   fmt.Errorf("%w - actual type %s, expected type string", ErrUnsupportableType, rv.Kind().String()),
		}
	}

	l, err := strconv.Atoi(length)
	if err != nil {
		return ErrConvertionStrToInt
	}
	if len(rv.String()) != l {
		return &ValidationError{
			Field: field,
			Err:   fmt.Errorf("%w - actual length %d", ErrLen, len(rv.String())),
		}
	}
	return nil
}

func checkRegexp(field string, v interface{}, r string) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.String {
		return &ValidationError{
			Field: field,
			Err:   fmt.Errorf("%w - actual type %s, expected type string", ErrUnsupportableType, rv.Kind().String()),
		}
	}

	re, err := regexp.Compile(r)
	if err != nil {
		return fmt.Errorf("%w - regexp=`%s`", ErrRegexpCompile, r)
	}
	if !re.MatchString(rv.String()) {
		return &ValidationError{
			Field: field,
			Err:   fmt.Errorf("%w - value=%s, regexp=`%s`", ErrRegexp, rv, r),
		}
	}
	return nil
}

func checkIn(field string, v interface{}, inc string) error {
	rv := reflect.ValueOf(v)
	if (rv.Kind() != reflect.String) && (rv.Kind() != reflect.Int) {
		return &ValidationError{
			Field: field,
			Err:   fmt.Errorf("%w - actual type %s, expected type string", ErrUnsupportableType, rv.Kind().String()),
		}
	}

	if rv.Kind() == reflect.String {
		s := strings.Split(inc, ",")
		dict := make(map[string]struct{}, len(s))
		for _, val := range s {
			dict[val] = struct{}{}
		}
		if _, ok := dict[rv.String()]; !ok {
			return &ValidationError{
				Field: field,
				Err:   fmt.Errorf("%w - value=%s, enum=%s", ErrInclude, rv.String(), inc),
			}
		}
	}
	if rv.Kind() == reflect.Int {
		s := strings.Split(inc, ",")
		dict := make(map[int]struct{}, len(s))
		for _, val := range s {
			i, err := strconv.Atoi(val)
			if err != nil {
				return fmt.Errorf("%w - value=`%s`", ErrConvertionStrToInt, val)
			}
			dict[i] = struct{}{}
		}
		if _, ok := dict[int(rv.Int())]; !ok {
			return &ValidationError{
				Field: field,
				Err:   fmt.Errorf("%w - value=%d, enum=%s", ErrInclude, rv.Int(), inc),
			}
		}
	}
	return nil
}

func minVal(field string, v interface{}, min string) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Int {
		return &ValidationError{
			Field: field,
			Err:   fmt.Errorf("%w - actual type %s, expected type string", ErrUnsupportableType, rv.Kind().String()),
		}
	}

	i, err := strconv.Atoi(min)
	if err != nil {
		return ErrConvertionStrToInt
	}
	if int(rv.Int()) < i {
		return &ValidationError{
			Field: field,
			Err:   fmt.Errorf("%w - value=%d, min=%d", ErrMin, rv.Int(), i),
		}
	}
	return nil
}

func maxVal(field string, v interface{}, max string) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Int {
		return &ValidationError{
			Field: field,
			Err:   fmt.Errorf("%w - actual type %s, expected type string", ErrUnsupportableType, rv.Kind().String()),
		}
	}

	i, err := strconv.Atoi(max)
	if err != nil {
		return ErrConvertionStrToInt
	}
	if int(rv.Int()) > i {
		return &ValidationError{
			Field: field,
			Err:   fmt.Errorf("%w - value=%d, max=%d", ErrMax, rv.Int(), i),
		}
	}
	return nil
}

var vMap = map[string](func(string, interface{}, string) error){
	"in":     checkIn,
	"len":    checkLen,
	"regexp": checkRegexp,
	"min":    minVal,
	"max":    maxVal,
}

func validateField(field, tag string, v interface{}) error {
	var vErr ValidationErrors
	rules := strings.Split(tag, tagAndSeparator)

	for _, rule := range rules {
		r := strings.Split(rule, tagKeyValueSeparator)
		method := vMap[r[0]]
		if method == nil {
			return fmt.Errorf("%w - method=`%s`", ErrValidMethod, r[0])
		}
		err := method(field, v, r[1])

		var e *ValidationError
		if err != nil {
			if errors.As(err, &e) {
				vErr = append(vErr, *e)
			} else {
				return err
			}
		}
	}
	return &vErr
}

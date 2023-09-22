package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	tagName              = "validate"
	tagAndSeparator      = "|"
	tagKeyValueSeparator = ":"
	tagValueSeparator    = ","
)

var (
	ErrNotStruct          = errors.New("input value is not a structure")
	ErrRegexpCompile      = errors.New("incorrect regexp")
	ErrConvertionStrToInt = errors.New("error convert string to int")
	ErrValidMethod        = errors.New("validation method is not supported")
)

type ValidationData struct {
	Field string
	Tag   string
	Value interface{}
}

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (vErrs *ValidationErrors) Error() string {
	var b strings.Builder
	for _, vErr := range *vErrs {
		b.WriteString(fmt.Sprintf("Field: %s, Error: %v", vErr.Field, vErr.Err))
	}
	return b.String()
}

func (vErrs *ValidationErrors) Unwrap() []error {
	e := make([]error, 0, len(*vErrs))
	for _, vErr := range *vErrs {
		e = append(e, vErr.Unwrap())
	}
	return e
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("Field: %s, Error: %v\n", v.Field, v.Err)
}

func (v *ValidationError) Unwrap() error {
	return errors.Unwrap(v.Err)
}

func Validate(v interface{}) error {
	var validationErrors ValidationErrors

	valuesData, err1 := getValidationData(v)
	if err1 != nil {
		return err1
	}

	for _, val := range *valuesData {
		err := validateField(val.Field, val.Tag, val.Value)
		if err != nil {
			var e *ValidationErrors
			if errors.As(err, &e) {
				validationErrors = append(validationErrors, *e...)
			} else {
				return err
			}
		}
	}

	if len(validationErrors) > 0 {
		return &validationErrors
	}
	return nil
}

func getValidationData(v interface{}) (*[]ValidationData, error) {
	var valuesData []ValidationData
	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%w", ErrNotStruct)
	}

	t := val.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fv := val.Field(i)
		if !fv.CanInterface() {
			continue
		}
		tag := field.Tag.Get(tagName)
		if tag == "" {
			continue
		}
		switch fv.Type().Kind().String() {
		case "int":
			valuesData = append(valuesData, ValidationData{field.Name, tag, int(fv.Int())})
		case "string":
			valuesData = append(valuesData, ValidationData{field.Name, tag, fv.String()})
		case "slice":
			if fv.Type().String() == "[]string" {
				strSlice := fv.Interface().([]string)
				for i := 0; i < len(strSlice); i++ {
					valuesData = append(valuesData, ValidationData{field.Name, tag, strSlice[i]})
				}
			}
			if fv.Type().String() == "[]int" {
				intSlice := fv.Interface().([]int)
				for i := 0; i < len(intSlice); i++ {
					valuesData = append(valuesData, ValidationData{field.Name, tag, intSlice[i]})
				}
			}
		case "struct", "ptr":
			if tag != "dive" {
				continue
			}
			data, err := getValidationData(fv.Interface())
			if err != nil {
				return nil, err
			}
			valuesData = append(valuesData, *data...)
		}
	}
	return &valuesData, nil
}

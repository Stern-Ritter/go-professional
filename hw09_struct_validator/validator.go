package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	validateTagName = "validate"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Validation Errors: %d\n", len(v)))
	for _, err := range v {
		sb.WriteString(fmt.Sprintf("\t%s: %s\n", err.Field, err.Err))
	}
	return sb.String()
}

func Validate(obj interface{}) error {
	errs := make(ValidationErrors, 0)

	v := reflect.ValueOf(obj)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldInfo := t.Field(i)
		fieldType := field.Kind()

		tagValue := fieldInfo.Tag.Get(validateTagName)
		if tagValue == "" {
			continue
		}

		//exhaustive:ignore
		switch fieldType {
		case reflect.Int:
			intValue := int(field.Int())
			validationErrs := validateInt(fieldInfo.Name, intValue, tagValue)
			errs = append(errs, validationErrs...)

		case reflect.String:
			stringValue := field.String()
			validationErrs := validateString(fieldInfo.Name, stringValue, tagValue)
			errs = append(errs, validationErrs...)

		case reflect.Slice:
			fieldElemType := field.Type().Elem().Kind()
			for j := 0; j < field.Len(); j++ {
				elem := field.Index(j)
				//exhaustive:ignore
				switch fieldElemType {
				case reflect.Int:
					intValue := int(elem.Int())
					validationErrs := validateInt(fieldInfo.Name, intValue, tagValue)
					errs = append(errs, validationErrs...)
				case reflect.String:
					stringValue := elem.String()
					validationErrs := validateString(fieldInfo.Name, stringValue, tagValue)
					errs = append(errs, validationErrs...)
				default:
					continue
				}
			}
		default:
			continue
		}
	}

	return errs
}

func validateString(fieldName string, value string, validation string) []ValidationError {
	errs := make([]ValidationError, 0)
	conditions := strings.Split(validation, "|")
	for _, condition := range conditions {
		conditionName, conditionValue, found := strings.Cut(condition, ":")
		if !found {
			continue
		}

		switch conditionName {
		case "len":
			expectedLength, err := strconv.Atoi(conditionValue)
			length := len(value)
			if err != nil {
				errs = append(errs, newValidationError(fieldName,
					fmt.Errorf("invalid expected length: %s: %w", conditionValue, err)))
			} else if length != expectedLength {
				errs = append(errs, newValidationError(fieldName,
					fmt.Errorf("invalid string length: %s, should be: %d but got: %d", value, expectedLength, length)))
			}

		case "regexp":
			regex, err := regexp.Compile(conditionValue)
			if err != nil {
				errs = append(errs, newValidationError(fieldName,
					fmt.Errorf("invalid regexp: %s: %w", conditionValue, err)))
			} else if !regex.MatchString(value) {
				errs = append(errs, newValidationError(fieldName,
					fmt.Errorf("invalid string value: %s, should match regexp: %s", value, conditionValue)))
			}

		case "in":
			expectedValues := strings.Split(conditionValue, ",")
			if !contains(expectedValues, value) {
				errs = append(errs, newValidationError(fieldName,
					fmt.Errorf("invalid value: %s, should be in range: %v", value, expectedValues)))
			}
		}
	}

	return errs
}

func validateInt(fieldName string, value int, validation string) []ValidationError {
	errs := make([]ValidationError, 0)
	conditions := strings.Split(validation, "|")
	for _, condition := range conditions {
		conditionName, conditionValue, found := strings.Cut(condition, ":")
		if !found {
			continue
		}

		switch conditionName {
		case "min":
			minValue, err := strconv.Atoi(conditionValue)
			if err != nil {
				errs = append(errs, newValidationError(fieldName,
					fmt.Errorf("invalid min value: %s: %w", conditionValue, err)))
			} else if value < minValue {
				errs = append(errs, newValidationError(fieldName,
					fmt.Errorf("invalid value: %d, should be greater or equal: %d", value, minValue)))
			}

		case "max":
			maxValue, err := strconv.Atoi(conditionValue)
			if err != nil {
				errs = append(errs, newValidationError(fieldName,
					fmt.Errorf("invalid max value: %s: %w", conditionValue, err)))
			} else if value > maxValue {
				errs = append(errs, newValidationError(fieldName,
					fmt.Errorf("invalid value: %d, should be less or equal: %d", value, maxValue)))
			}

		case "in":
			expectedValues, err := stringsToNumbers(strings.Split(conditionValue, ","))
			if err != nil {
				errs = append(errs, newValidationError(fieldName,
					fmt.Errorf("invalid expected values: %v: %w", conditionValue, err)))
			} else if !contains(expectedValues, value) {
				errs = append(errs, newValidationError(fieldName,
					fmt.Errorf("invalid value: %d, should be in range: %v", value, expectedValues)))
			}
		}
	}

	return errs
}

func newValidationError(field string, err error) ValidationError {
	return ValidationError{
		Field: field,
		Err:   err,
	}
}

func stringsToNumbers(strings []string) ([]int, error) {
	numbers := make([]int, 0, len(strings))
	for _, s := range strings {
		num, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		numbers = append(numbers, num)
	}
	return numbers, nil
}

func contains[T comparable](slice []T, element T) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

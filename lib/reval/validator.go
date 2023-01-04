package reval

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

const defaultKey = "reval"

// ErrValidating will appear in case of validation error.
var ErrValidating = errors.New("validation error")

// ErrUnexpected made in case of panic.
var ErrUnexpected = errors.New("unexpected error occurred")

// ValidationError could occur during unsuccessful validation
// Error message will be rendered dynamically according
// to the first inadequacy of struct field validation pattern that noted in
// struct tag field according to regexp notation.
type ValidationError struct {
	entity   string
	property string
	isZero   bool
	code     string
	pattern  string
}

// Error implements error interface
// and can distinct if non-zero value was provided.
func (vErr *ValidationError) Error() string {
	if !vErr.isZero {
		vErr.code = "reval" + " "
	}

	return strings.ToLower(fmt.Sprintf("%s has to have %s%s according to the pattern: `%s`",
		vErr.entity,
		vErr.code,
		vErr.property,
		vErr.pattern,
	))
}

// ValidateStruct validates struct fields
// according to given regex tag.
func ValidateStruct(src any) (err error) {
	// check if src is a struct
	srcValue, err := inspectSource(src)
	if err != nil {
		return err
	}
	// despite it is abs not necessary, in case of panic we will return ErrUnexpected
	defer func() {
		if recover() != nil {
			err = ErrUnexpected
		}
	}()

	// top level struct name (in case we are using nested structs)
	var structName string
	if structName == "" {
		structName = srcValue.Type().Name()
	}
	// iterate  all over struct fields
	for i := 0; i < srcValue.NumField(); i++ {
		fieldValue := srcValue.Field(i)
		fieldName := srcValue.Type().Field(i).Name
		tagValue := string(srcValue.Type().Field(i).Tag)
		// check presence of regex tag (.Tag.Lookup() would not work here)
		if pattern, ok := getTagValue(tagValue, defaultKey); ok {
			if fieldValue.IsZero() {
				return fmt.Errorf("%v: %w", ErrValidating,
					&ValidationError{
						entity:   structName,
						property: fieldName,
						isZero:   true},
				)
			}
			// field validation according to pattern
			if !regexp.MustCompile(pattern).MatchString(fmt.Sprintf("%v", fieldValue)) {
				return fmt.Errorf("%s: %w", ErrValidating,
					&ValidationError{
						entity:   structName,
						property: fieldName,
						pattern:  pattern},
				)
			}
		}
		// recursive call for nested structs
		if fieldValue.Type().Kind() == reflect.Struct {
			if err := ValidateStruct(fieldValue.Interface()); err != nil {
				return fmt.Errorf("error validating nested struct: %w", err)
			}
		}
	}

	return nil
}

// getTagValue address the problem of looking
// for <key>/<val> pairs of struct tag fields
// which is not solved by reflect.Tag.Lookup())
func getTagValue(tag string, key string) (string, bool) {
	tagStr := fmt.Sprintf("%v", reflect.StructTag(tag))
	tagValue := fmt.Sprintf(`(?s)(?i)\s*(?P<key>%s):\"(?P<value>[^\"]+)\"`, key)

	if match := regexp.MustCompile(tagValue).
		FindStringSubmatch(tagStr); match != nil {
		return match[2], true
	}

	return "", false
}

func inspectSource(src any) (srcValue *reflect.Value, err error) {
	defer func() {
		if recover() != nil {
			err = ErrUnexpected
		}
	}()

	*srcValue = reflect.Indirect(reflect.ValueOf(src))

	if srcType := srcValue.Kind(); srcType != reflect.Struct {
		return nil,
			fmt.Errorf("input value must be struct, got: %v", srcType)
	}

	return srcValue, nil
}

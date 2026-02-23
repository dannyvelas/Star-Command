package app

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var (
	errNotPointer      = errors.New("argument must be a pointer")
	errEmptyEnvVar     = errors.New("environment variable is set but empty")
	errEmptyInput      = errors.New("value cannot be empty")
	errUnsupportedType = errors.New("unsupported field type")
)

// promptSensitiveFields fills fields tagged `sensitive:"true"` on the given
// struct pointer. For each such field, it first checks for a matching
// environment variable (case-insensitive, prefixed with "IAC_" + json tag or
// field name). If no env var is found, it prompts the user interactively using
// the `prompt` tag value, or the field name if no prompt tag is set.
//
// Returns an error if a sensitive field has a type other than string, int, float, or bool.
// Also returns an error if v is not a pointer, or if the user enters an empty value.
func promptSensitiveFields(v any, r io.Reader, w io.Writer) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer {
		return fmt.Errorf("promptSensitiveFields: %w", errNotPointer)
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("promptSensitiveFields: argument must be a pointer to a struct")
	}

	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		fieldVal := rv.Field(i)

		if field.Tag.Get("sensitive") != "true" {
			continue
		}

		// Determine the env var key: IAC_ + json tag name, or IAC_ + field name
		envPrefix := "IAC_"
		envSuffix := field.Name
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			envSuffix = strings.Split(jsonTag, ",")[0]
		}
		envKey := envPrefix + envSuffix

		// Case-insensitive search in environment
		value := ""
		found := false
		for _, entry := range os.Environ() {
			parts := strings.SplitN(entry, "=", 2)
			if len(parts) == 2 && strings.EqualFold(parts[0], envKey) {
				value = parts[1]
				found = true
				break
			}
		}

		if found {
			if value == "" {
				return fmt.Errorf("environment variable %q: %w", envKey, errEmptyEnvVar)
			}
		} else {
			// Fall back to interactive prompt
			promptMsg := field.Tag.Get("prompt")
			if promptMsg == "" {
				promptMsg = field.Name
			}

			_, _ = fmt.Fprintf(w, "Enter a value for %q: ", promptMsg)
			reader := bufio.NewReader(r)
			input, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("error reading input for %q: %v", promptMsg, err)
			}
			value = strings.TrimSpace(input)
			if value == "" {
				return fmt.Errorf("field %q: %w", promptMsg, errEmptyInput)
			}
		}

		if err := setSensitiveField(fieldVal, field, value); err != nil {
			return err
		}
	}

	return nil
}

func setSensitiveField(fieldVal reflect.Value, field reflect.StructField, value string) error {
	switch fieldVal.Kind() {
	case reflect.String:
		fieldVal.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("field %q: cannot parse %q as int: %w", field.Name, value, err)
		}
		fieldVal.SetInt(n)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("field %q: cannot parse %q as float: %w", field.Name, value, err)
		}
		fieldVal.SetFloat(f)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("field %q: cannot parse %q as bool: %w", field.Name, value, err)
		}
		fieldVal.SetBool(b)
	default:
		return fmt.Errorf("promptSensitiveFields: field %q: %w: %s", field.Name, errUnsupportedType, fieldVal.Kind())
	}
	return nil
}

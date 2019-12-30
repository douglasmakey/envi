package envi

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	// ErrIsRequired env is required
	ErrIsRequired            = errors.New("env is required")
	errNotAPointer           = errors.New("expected a pointer")
	errTagOptionNotSupported = errors.New("tag option is not supported")
	errInvalidMapItem        = errors.New("invalid map item")
)

const (
	envSeparator = "envSeparator"
	required     = "required"
)

// EnvError struct define to error
type EnvError struct {
	KeyName   string
	FieldType string
	Value     string
	Err       error
}

func (e *EnvError) Error() string {
	return fmt.Sprintf("Key[%s] of type %s, Value: %s - Error: %s", e.KeyName, e.FieldType, e.Value, e.Err.Error())
}

// Parse env using reflect
func Parse(val interface{}) error {
	ptrValue := reflect.ValueOf(val)
	if ptrValue.Kind() != reflect.Ptr {
		return errNotAPointer
	}
	refValue := ptrValue.Elem()
	if refValue.Kind() != reflect.Struct {
		return errNotAPointer
	}

	return do(refValue)
}

// do set type and value for each env
func do(val reflect.Value) error {
	var err error
	refType := val.Type()

	// With refType.Kind() get kind represents the specific kind of type that a Type represents.
	// With refType.NumField() obtain the number of fields of the struct
	// With refType.Field(position) obtain a struct type's in position
	for i := 0; i < refType.NumField(); i++ {
		value, err := getValue(refType.Field(i))
		if err != nil {
			return &EnvError{
				KeyName:   refType.Field(i).Name,
				FieldType: refType.Field(i).Type.String(),
				Value:     value,
				Err:       err,
			}
		}
		if value == "" {
			continue
		}
		separator := refType.Field(i).Tag.Get(envSeparator)
		if err := setValue(val.Field(i), value, separator); err != nil {
			return &EnvError{
				KeyName:   refType.Field(i).Name,
				FieldType: refType.Field(i).Type.String(),
				Value:     value,
				Err:       err,
			}
		}
	}

	return err

}

// getValue get value or default value of key
func getValue(sf reflect.StructField) (string, error) {
	var (
		value string
		err   error
	)

	key, options := parseKeyForOption(sf.Tag.Get("env"))

	// Get default value if exists
	defaultValue := sf.Tag.Get("envDefault")
	value = getValueOrDefault(key, defaultValue)

	if len(options) > 0 {
		for _, option := range options {
			// TODO: Implement others options supported
			// For now only option supported is "required".
			switch option {
			case "":
				break
			case required:
				value, err = getRequired(key)
			default:
				err = errTagOptionNotSupported
			}
		}
	}

	return value, err
}

func parseKeyForOption(k string) (string, []string) {
	opts := strings.Split(k, ",")
	return opts[0], opts[1:]
}

func getValueOrDefault(k, defValue string) string {
	// Retrieves the value of the environment variable named by the key.
	// If the variable is present in the environment, return value
	value, ok := os.LookupEnv(k)
	if ok {
		return value
	}
	return defValue
}

func getRequired(k string) (string, error) {
	// Retrieves the value of the environment variable named by the key.
	// If the variable is present in the environment, return value and nil for error
	if value, ok := os.LookupEnv(k); ok {
		return value, nil
	}
	return "", ErrIsRequired
}

func setValue(field reflect.Value, value string, separator string) error {
	refType := field.Type()

	if refType.Kind() == reflect.Ptr {
		refType = refType.Elem()
		if field.IsNil() {
			field.Set(reflect.New(refType))
		}
		field = field.Elem()
	}

	switch refType.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var (
			val int64
			err error
		)
		if field.Kind() == reflect.Int64 && refType.PkgPath() == "time" && refType.Name() == "Duration" {
			var td time.Duration
			td, err = time.ParseDuration(value)
			val = int64(td)
		} else {
			val, err = strconv.ParseInt(value, 0, refType.Bits())
		}
		if err != nil {
			return err
		}

		field.SetInt(val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(value, 0, refType.Bits())
		if err != nil {
			return err
		}
		field.SetUint(val)
	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(value, refType.Bits())
		if err != nil {
			return err
		}
		field.SetFloat(val)
	case reflect.Bool:
		val, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(val)
	case reflect.Slice:
		// Validate separator and set default
		if separator == "" {
			separator = ","
		}
		values := strings.Split(value, separator)
		newSlice := reflect.MakeSlice(refType, len(values), len(values))
		for i, val := range values {
			if err := setValue(newSlice.Index(i), val, ""); err != nil {
				return err
			}
		}
		field.Set(newSlice)
	case reflect.Map:
		newMap := reflect.MakeMap(refType)
		if len(strings.TrimSpace(value)) != 0 {
			pairs := strings.Split(value, ",")
			for _, pair := range pairs {
				kPair := strings.Split(pair, ":")
				if len(kPair) != 2 {
					return errInvalidMapItem
				}
				k := reflect.New(refType.Key()).Elem()
				if err := setValue(k, kPair[0], ""); err != nil {
					return err
				}

				v := reflect.New(refType.Elem()).Elem()
				if err := setValue(v, kPair[1], ""); err != nil {
					return err
				}
				newMap.SetMapIndex(k, v)
			}
		}
		field.Set(newMap)
	}

	return nil
}

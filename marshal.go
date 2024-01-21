package urlvalues

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/vcraescu/go-urlvalues/internal"
)

// Marshaler is an interface that defines the method for encoding Go values into URL query parameters.
type Marshaler interface {
	EncodeValues(key string, values *url.Values) error
}

// MarshalerOptions defines optional parameters for the Marshal function.
// These options only apply to maps, for structs the tag options are used.
type MarshalerOptions struct {
	// ArrayBrackets, when set to true, adds brackets to array/slice keys in the URL parameters.
	ArrayBrackets bool

	// ArrayDelimiter specifies a custom delimiter for joining slice values in the URL parameters.
	ArrayDelimiter string

	// IntBool when set to true, represents boolean values as integers (0 or 1).
	IntBool bool

	// TimeUnix time encoded in seconds.
	TimeUnix bool

	// TimeUnixMilli time encoded in milliseconds.
	TimeUnixMilli bool

	// TimeUnixNano time encoded in nanoseconds.
	TimeUnixNano bool

	// TimeLayout allows specifying a custom layout for time values if provided.
	TimeLayout string
}

// Marshal encodes the given value into URL query parameters with the default options.
func Marshal(v any) (url.Values, error) {
	return (MarshalerOptions{}).Marshal(v)
}

// Marshal encodes the given value into URL query parameters.
func (m MarshalerOptions) Marshal(v any) (url.Values, error) {
	values := url.Values{}

	if v == nil {
		return values, nil
	}

	val := internal.Indirect(reflect.ValueOf(v))

	if internal.IsZero(val) {
		return nil, nil
	}

	switch val.Kind() {
	case reflect.String:
		return url.ParseQuery(val.Interface().(string))

	case reflect.Map, reflect.Struct:
		if err := m.marshalAny(values, val, ""); err != nil {
			return nil, err
		}

		return values, nil

	default:
		return nil, fmt.Errorf(
			"expected %s/%s/%s but got %q", reflect.String, reflect.Map, reflect.Struct, val.Kind())
	}
}

func (m MarshalerOptions) marshalAny(urlValues url.Values, val reflect.Value, scope string) error {
	if m, ok := val.Interface().(Marshaler); ok {
		if err := m.EncodeValues(scope, &urlValues); err != nil {
			return fmt.Errorf("EncodeValues: %w", err)
		}

		return nil
	}

	switch val.Kind() {
	case reflect.Map:
		if err := m.marshalMap(urlValues, val, scope); err != nil {
			return fmt.Errorf("marshal map: %w", err)
		}

	case reflect.Array:
		fallthrough
	case reflect.Slice:
		if err := m.marshalSlice(urlValues, val, scope); err != nil {
			return fmt.Errorf("marshal slice: %w", err)
		}

	case reflect.Struct:
		if err := m.marshalStruct(urlValues, val, scope); err != nil {
			return fmt.Errorf("marshal struct: %w", err)
		}

	case reflect.Ptr:
		if err := m.marshalPtr(urlValues, val, scope); err != nil {
			return fmt.Errorf("marshal pointer: %w", err)
		}

	default:
		if err := m.marshalBasicTypes(urlValues, val, scope); err != nil {
			return fmt.Errorf("marshal other: %w", err)
		}
	}

	return nil
}

func (m MarshalerOptions) marshalMap(urlValues url.Values, val reflect.Value, scope string) error {
	for k, v := range internal.StringValueMap(val) {
		if scope != "" {
			k = scope + "[" + k + "]"
		}

		if err := m.marshalAny(urlValues, v, k); err != nil {
			return err
		}
	}

	return nil
}

func (m MarshalerOptions) marshalStruct(urlValues url.Values, val reflect.Value, scope string) error {
	if _, ok := val.Interface().(time.Time); ok {
		if scope != "" {
			urlValues.Add(scope, m.encodeTime(val))
		}

		return nil
	}

	structValues, err := query.Values(val.Interface())
	if err != nil {
		return err
	}

	internal.MergeURLValues(urlValues, structValues, scope)

	return nil
}

func (m MarshalerOptions) marshalSlice(urlValues url.Values, val reflect.Value, scope string) error {
	if m.ArrayBrackets {
		scope += "[]"
	}

	if m.ArrayDelimiter == "" {
		for _, v := range internal.SliceValues(val) {
			if err := m.marshalAny(urlValues, v, scope); err != nil {
				return err
			}
		}

		return nil
	}

	tmp := url.Values{}

	for _, v := range internal.SliceValues(val) {
		if err := m.marshalAny(tmp, v, scope); err != nil {
			return err
		}
	}

	if _, ok := tmp[scope]; ok {
		urlValues.Add(scope, strings.Join(tmp[scope], m.ArrayDelimiter))
	}

	return nil
}

func (m MarshalerOptions) marshalPtr(urlValues url.Values, val reflect.Value, scope string) error {
	val = internal.Indirect(val)

	return m.marshalAny(urlValues, val, scope)
}

func (m MarshalerOptions) marshalBasicTypes(urlValues url.Values, val reflect.Value, scope string) error {
	if scope == "" {
		return nil
	}

	var s string

	switch val.Kind() {
	case reflect.Bool:
		s = m.encodeBool(val)
	default:
		s = encodeAsString(val)
	}

	urlValues.Add(scope, s)

	return nil
}

func (m MarshalerOptions) encodeBool(val reflect.Value) string {
	if !m.IntBool {
		return fmt.Sprintf("%v", val.Interface())
	}

	if val.Bool() {
		return "1"
	}

	return "0"
}

func (m MarshalerOptions) encodeTime(val reflect.Value) string {
	t := val.Interface().(time.Time)

	if m.TimeUnix {
		return fmt.Sprint(t.Unix())
	}

	if m.TimeUnixMilli {
		return fmt.Sprint(t.UnixMilli())
	}

	if m.TimeUnixNano {
		return fmt.Sprint(t.UnixNano())
	}

	layout := time.RFC3339

	if m.TimeLayout != "" {
		layout = m.TimeLayout
	}

	return t.Format(layout)
}

func encodeAsString(val reflect.Value) string {
	return internal.ValueString(val)
}

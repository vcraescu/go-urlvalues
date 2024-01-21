package internal

import (
	"fmt"
	"reflect"
)

type zeroable interface {
	IsZero() bool
}

func Indirect(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return v
		}

		v = v.Elem()
	}

	return v
}

func IsZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	default:
		if !v.IsValid() {
			return true
		}

		if a, ok := v.Interface().(zeroable); ok {
			return a.IsZero()
		}
	}

	return false
}

func StringValueMap(v reflect.Value) map[string]reflect.Value {
	iter := v.MapRange()
	out := make(map[string]reflect.Value)

	for iter.Next() {
		mv := reflect.ValueOf(iter.Value().Interface())

		if !mv.IsValid() {
			continue
		}

		out[ValueString(iter.Key())] = mv
	}

	return out
}

func ValueString(v reflect.Value) string {
	if IsZero(v) {
		return ""
	}

	return fmt.Sprint(v.Interface())
}

func SliceValues(v reflect.Value) []reflect.Value {
	out := make([]reflect.Value, 0, v.Len())

	for i := 0; i < v.Len(); i++ {
		out = append(out, v.Index(i))
	}

	return out
}

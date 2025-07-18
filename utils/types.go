package utils

import (
	"fmt"
	"reflect"
)

type Int interface {
	int8 | int16 | int32 | int64 | int
}

type UInt interface {
	uint8 | uint16 | uint32 | uint64 | uint
}

type Integer interface {
	Int | UInt
}

type Float interface {
	float32 | float64
}

type Number interface {
	Integer | Float
}

type Primitive interface {
	bool | Number | string
}

// IsTypePrimitive uses reflection to test if a given type parameter is considered
// a `Primitive`.
func IsTypePrimitive[T any]() bool {
	switch reflect.TypeFor[T]().Kind() {
	case reflect.Bool, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int,
		reflect.Float32, reflect.Float64,
		reflect.String:
		return true
	}
	return false
}

// IsPrimitive performs a type check on the given value and returns true if it
// is considered a `Primitive`.
func IsPrimitive(val any) bool {
	switch val.(type) {
	case bool, uint8, uint16, uint32, uint64, uint,
		int8, int16, int32, int64, int,
		float32, float64,
		string:
		return true
	}
	return false
}

// String converts any value into a string representation using [fmt.Sprintf]. It
// performs a type check in order to pick a good verb. If the value is not a `Primitive`
// then the "%v" is chosen by default.
func String(val any) string {
	verb := "%v"
	switch val.(type) {
	case bool:
		verb = "%t"
	case uint8, uint16, uint32, uint64, uint,
		int8, int16, int32, int64, int:
		verb = "%d"
	case float32, float64:
		verb = "%g"
	case string:
		verb = "%s"
	}
	return fmt.Sprintf(verb, val)
}

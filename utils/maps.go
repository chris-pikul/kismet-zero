package utils

import "strings"

// TransformMap takes an input map, and runs a given transformer callback on all
// the entries in order to construct a new map.
func TransformMap[IK comparable, IV any, OK comparable, OV any](input map[IK]IV, transformer func(IK, IV) (OK, OV)) map[OK]OV {
	output := make(map[OK]OV, len(input))
	var nextK OK
	var nextV OV
	for k, v := range input {
		nextK, nextV = transformer(k, v)
		output[nextK] = nextV
	}
	return output
}

// TransformMapKeys takes an input map, and runs a given transformer callback on
// all the keys found in order to construct a new map.
func TransformMapKeys[IK comparable, OK comparable, V any](input map[IK]V, transformer func(IK) OK) map[OK]V {
	output := make(map[OK]V, len(input))
	for k, v := range input {
		output[transformer(k)] = v
	}
	return output
}

// TransformMapValues takes an input map, and runs a given transformer callback on
// all the values found in order to construct a new map.
func TransformMapValues[K comparable, IV any, OV any](input map[K]IV, transformer func(IV) OV) map[K]OV {
	output := make(map[K]OV, len(input))
	for k, v := range input {
		output[k] = transformer(v)
	}
	return output
}

// MapContainsKey checks if the given map contains a key. This is just syntactic
// sugar.
func MapContainsKey[K comparable, V any](input map[K]V, query K) bool {
	_, exists := input[query]
	return exists
}

// MapContainsValue checks if the given map contains a value. This requires that
// the value be comparable.
func MapContainsValue[K comparable, V comparable](input map[K]V, query V) bool {
	for _, v := range input {
		if v == query {
			return true
		}
	}
	return false
}

// MapContainsValueFunc checks if the given map contains a value. It accepts a query
// value and a equality checking function. The equality checker will be ran on each
// value, being supplied with the given `query` as well to help reduce state. First
// time the equality checker returns true, so does this function.
func MapContainsValueFunc[K comparable, V any](input map[K]V, query V, matcher func(V, V) bool) bool {
	for _, v := range input {
		if matcher(query, v) {
			return true
		}
	}
	return false
}

// JoinMap turns a map into a string by combining all the key-value pairs and
// performing stringification.
//
// The `colSep` parameter decides the string in-between the key and the value.
//
// The `rowSep` parameter sets the separator between pairs of key-values.
//
// This function uses the [String] function for stringification of both the
// key and values.
func JoinMap[K comparable, V any](m map[K]V, colSep, rowSep string) string {
	var sb strings.Builder

	first := true
	for k, v := range m {
		if !first {
			sb.WriteString(rowSep)
		}
		sb.WriteString(String(k))
		sb.WriteString(colSep)
		sb.WriteString(String(v))
	}

	return sb.String()
}

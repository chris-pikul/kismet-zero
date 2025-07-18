package utils

import "net/url"

// QueryString converts a given map into a URL safe query string. It does not
// prefix it with the "?" character.
func QueryString[K comparable, V any](query map[K]V) string {
	return JoinMap(TransformMapValues(query, func(v V) string {
		return url.QueryEscape(String(v))
	}), "=", "&")
}

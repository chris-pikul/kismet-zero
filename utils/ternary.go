package utils

// If acts as a simple ternary operation and returns `t` if the condition `cond`
// is true, otherwise it returns `f`.
func If[T any](cond bool, t T, f T) T {
	if cond {
		return t
	}
	return f
}

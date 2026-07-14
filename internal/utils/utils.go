package utils

func If[T any](cond bool, t T, f T) T {
	if cond {
		return t
	}
	return f
}

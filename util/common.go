package util

// Ternary 三目运算符
func Ternary[T any](expression bool, a, b T) T {
	if expression {
		return a
	} else {
		return b
	}
}

package util

// Ternary 三目运算符
func Ternary[T any](expression bool, a, b T) T {
	if expression {
		return a
	} else {
		return b
	}
}

// IgnoreError 忽略错误
func IgnoreError[T any](value T, err error) T {
	return value
}

// IgnoreValue 忽略返回值
func IgnoreValue[T any](value T, err error) error {
	return err
}

package utils

func GetStringPointer(s string) *string {
	return &s
}

func GetPointer[T any](val T) *T {
	return &val
}

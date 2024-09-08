package utils

func GetStringPointer(s string) *string {
	return &s
}

func GetPointer[T any](ptr *T) *T {
	return ptr
}

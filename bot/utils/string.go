package utils

// ObterPonteiro obtém o ponteiro de qualquer valor
func ObterPonteiro[T any](v T) *T {
	return &v
}

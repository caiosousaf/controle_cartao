package utils

import "github.com/pkg/errors"

// Wrap retorna um erro com uma mensagem adicional
func Wrap(err error, message string) error {
	return errors.Wrap(err, message)
}

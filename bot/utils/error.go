package utils

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

const ErroPdfVazio = "Bad Request: file must be non-empty"

// Error representa um tipo de erro personalizado com uma mensagem de erro, código de erro e código de status HTTP.
type Error struct {
	// A mensagem de erro.
	Msg string
	// O erro subjacente.
	Err error
	// O código de erro personalizado.
	Code int
	// O código de status HTTP.
	StatusCode int
}

// Wrap retorna um erro com uma mensagem adicional
func Wrap(err error, message string) error {
	return errors.Wrap(err, message)
}

// defaultCode é o código padrão de erro para retorno
const defaultCode = 400

// Error retorna a mensagem de erro para o erro personalizado.
func (e *Error) Error() string {
	return e.Msg
}

// NewErr cria um novo erro personalizado com a mensagem fornecida e define o código de erro padrão e o código de status HTTP.
func NewErr(message string) error {
	return &Error{
		Msg:        message,
		Err:        errors.New(fmt.Sprintf("mensagem de erro inline: '%s'. Consulte o rastreamento de pilha do erro para obter informações adicionais", message)),
		Code:       defaultCode,
		StatusCode: http.StatusBadRequest,
	}
}

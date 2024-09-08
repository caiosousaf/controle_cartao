package middlewares

import (
	"controle_cartao/utils"
	"golang.org/x/crypto/bcrypt"
)

// HashSenha Gera o hash da senha
func HashSenha(senha *string) (*string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*senha), bcrypt.DefaultCost)
	return utils.GetStringPointer(string(bytes)), err
}

// VerificarSenha valida a senha com o hash
func VerificarSenha(senha, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(senha))
	return err == nil
}

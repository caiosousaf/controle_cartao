package middlewares

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"strconv"
	"time"
)

var jwtKey = []byte("jwt_controle_cartao")

// GerarJWT é responstável por gerar o token JWT
func GerarJWT(nome string, idUsuario *uuid.UUID) (string, error) {
	dataAtual := time.Now()

	// Create claims with multiple fields populated
	claims := CustomClaims{
		"controle_cartao",
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    strconv.Itoa(int(dataAtual.Unix())),
			Subject:   nome,
			ID:        idUsuario.String(),
			Audience:  []string{nome + "-" + idUsuario.String()},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(jwtKey)

	return ss, err
}

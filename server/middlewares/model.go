package middlewares

import "github.com/golang-jwt/jwt/v5"

// CustomClaims é a estrutura que define as claims personalizadas
type CustomClaims struct {
	Foo string `json:"foo"`
	jwt.RegisteredClaims
}

package models

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

package models

import "github.com/golang-jwt/jwt/v5"

type CustomClaim struct {
	User int
	jwt.RegisteredClaims
}

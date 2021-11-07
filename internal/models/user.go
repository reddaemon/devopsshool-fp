package models

import (
	jwt "github.com/dgrijalva/jwt-go/v4"
	"github.com/google/uuid"
)

type Token struct {
	UserID uuid.UUID
	jwt.StandardClaims
}

type User struct {
	Email string `json:"email"`
	Password string `json:"password"`
	Token string `json:"token"`
}
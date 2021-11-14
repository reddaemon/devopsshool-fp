package models

import (
	jwt "github.com/dgrijalva/jwt-go/v4"
	"github.com/google/uuid"
)

type Token struct {
	UserID uuid.UUID
	jwt.StandardClaims
}

type TokenDetails struct {
	AccessToken string
	RefreshToken string
	AccessUuid string
	RefreshUuid string
	AtExpires int64
	RtExpires int64
}

type User struct {
	ID uuid.UUID
	Email string `json:"email"`
	Password string `json:"password"`
	Token string `json:"token"`
}
package models

import (
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type AccessDetails struct {
	AccessUuid string
	UserId     uint64
}

type Token struct {
	UserID uuid.UUID
	jwt.StandardClaims
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type User struct {
	ID       float64
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

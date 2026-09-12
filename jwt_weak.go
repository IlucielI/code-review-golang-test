package main

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Vulnerable: Weak JWT secret (predictable, short)
var jwtSecret = []byte("secret")

func generateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func validateToken(tokenString string) (*jwt.Token, error) {
	// Using the same weak secret
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
}

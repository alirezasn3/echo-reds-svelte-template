package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateJWT(user User, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"phoneNumber": user.PhoneNumber,
			"role":        user.Role,
			"firstName":   user.FirstName,
			"lastName":    user.LastName,
			"id":          user.ID,
			"exp":         time.Now().Add(time.Hour * 72).Unix(),
		},
	)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func DecodeJWT(tokenString, secret string) (*User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid jwt")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return &User{
			ID:          claims["id"].(string),
			FirstName:   claims["firstName"].(string),
			LastName:    claims["lastName"].(string),
			PhoneNumber: claims["phoneNumber"].(string),
			Role:        claims["role"].(string),
		}, nil
	} else {
		return nil, fmt.Errorf("failed to parse jwt")
	}
}

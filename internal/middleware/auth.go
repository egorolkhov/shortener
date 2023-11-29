package middleware

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"math/rand"
	"strconv"
	"time"
)

const TOKEN_EXP = time.Hour * 3

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

func BuidToken(secretKey string) (string, error) {
	userID := rand.Int()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP))},
		UserID: strconv.Itoa(userID),
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GetUserID(tokenString, secretKey string) string {
	claims := &Claims{}

	token, _ := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if !token.Valid {
		fmt.Println("Token is not valid")
		return "error"
	}

	fmt.Println("Token is valid")
	return claims.UserID
}

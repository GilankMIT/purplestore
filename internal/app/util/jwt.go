package util

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("test1234")

func BuildJWT(email string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_email": email,
		"exp_at":     time.Now().Add(time.Hour * 2).Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Print("error when building jwt", err)
		return ""
	}

	return tokenString
}

func VerifyJWT(accessToken string) error {

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		
		expiredAt := int64(claims["exp_at"].(float64))
		if time.Now().Unix() > expiredAt {
			return errors.New("token already expired")
		}
	} else {
		return errors.New("failed to extract token")
	}

	return nil
}

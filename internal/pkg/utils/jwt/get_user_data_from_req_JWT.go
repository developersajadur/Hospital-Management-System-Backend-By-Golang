package jwt

import (
	"errors"
	"fmt"
	"hospital_management_system/config"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	UserID string
	Role   string
	Exp    int64
	Iat    int64
}

func GetUserDataFromReqJWT(r *http.Request) (*UserClaims, error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return nil, errors.New("missing token")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.ENV.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := ""
		switch v := claims["userId"].(type) {
		case float64:
			userID = fmt.Sprintf("%v", int64(v))
		case string:
			userID = v
		default:
			return nil, errors.New("invalid userId in token")
		}

		role, _ := claims["role"].(string)
		exp, _ := claims["exp"].(float64)

		return &UserClaims{
			UserID: userID,
			Role:   role,
			Exp:    int64(exp),
			Iat:    int64(claims["iat"].(float64)),
		}, nil
	}

	return nil, errors.New("invalid token")
}

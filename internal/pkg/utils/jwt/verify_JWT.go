package jwt

import (
	"hospital_management_system/config"
	"hospital_management_system/internal/pkg/helpers"

	"github.com/golang-jwt/jwt/v5"
)

// getSecret returns the JWT secret at runtime
func getSecret() []byte {
	if config.ENV == nil || config.ENV.JWTSecret == "" {
		panic("JWT_SECRET is not set. Make sure config.Init() is called before using JWT utils")
	}
	return []byte(config.ENV.JWTSecret)
}


// VerifyJWT verifies the token and returns the claims
func VerifyJWT(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return getSecret(), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	 return nil, helpers.NewAppError(401, "Invalid or expired token")
}

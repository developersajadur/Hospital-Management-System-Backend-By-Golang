package utils

import (
	"hospital_management_system/config"
	"hospital_management_system/internal/pkg/helpers"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims structure
type JWTClaims struct {
	UserID string `json:"userId"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// getSecret returns the JWT secret at runtime
func getSecret() []byte {
	if config.ENV == nil || config.ENV.JWTSecret == "" {
		panic("JWT_SECRET is not set. Make sure config.Init() is called before using JWT utils")
	}
	return []byte(config.ENV.JWTSecret)
}

// GenerateJWT generates a JWT token for a user
func GenerateJWT(userID, email, role string, expiry time.Duration) (string, error) {
	claims := &JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getSecret())
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

// GetDataFromJWT returns userID and role from token string
func GetDataFromJWT(tokenStr string) (string, string, error) {
	claims, err := VerifyJWT(tokenStr)
	if err != nil {
		return "", "", err
	}
	return claims.UserID, claims.Role, nil
}

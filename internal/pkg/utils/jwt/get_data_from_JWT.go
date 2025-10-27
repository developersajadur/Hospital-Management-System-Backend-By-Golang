package jwt

// GetDataFromJWT returns userID and role from token string
func GetDataFromJWT(tokenStr string) (string, string, error) {
	claims, err := VerifyJWT(tokenStr)
	if err != nil {
		return "", "", err
	}
	return claims.UserID, claims.Role, nil
}


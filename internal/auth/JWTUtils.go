package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type JWTUtils struct {
	JwtKey []byte
}

func (j *JWTUtils) GenerateJWT(userID uuid.UUID, exp int64, genTime int64, jti string) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID.String(),
        "exp": exp,
		"iat": time.Now().Unix(),
		"jti": jti,
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
    return token.SignedString(j.JwtKey)
}

func (j *JWTUtils) ValidateJWT(tokenString string) (*uuid.UUID, string, error) {
    token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return j.JwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, "", fmt.Errorf("invalid claims")
	}

	userIDRaw := claims["user_id"]
	userIDStr, ok := userIDRaw.(string)
	if !ok {
		return nil, "", fmt.Errorf("user_id is not a string")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, "", fmt.Errorf("invalid UUID format: %w", err)
	}

	accessJTI := claims["jti"].(string)

    return &userID, accessJTI, nil
}
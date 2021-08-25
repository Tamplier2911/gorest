package access

import (
	"fmt"

	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Token represents authentication token.
type Token struct {
	jwt.StandardClaims

	UserID   uuid.UUID       `json:"userId"`
	UserRole models.UserRole `json:"userRole"`
}

// EncodeToken is used to encode Token to string.
func EncodeToken(token *Token, hmacSecret string) (string, error) {
	// create jwt token object
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, token)

	// encode jwt object to string
	tokenString, err := jwtToken.SignedString([]byte(hmacSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// DecodeToken is used to verify and decode token from string.
func DecodeToken(tokenString string, hmacSecret string) (*Token, error) {
	// parse token
	jwtToken, err := jwt.ParseWithClaims(tokenString, &Token{}, func(token *jwt.Token) (interface{}, error) {
		// validate signing algorithm
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// return validation secret
		return []byte(hmacSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("malformed token")
	}

	// check if token is valid
	if !jwtToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// parse claims
	token, ok := jwtToken.Claims.(*Token)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return token, nil
}

// SaveTokenToContext is used to save auth token to gin context.
func SaveTokenToContext(c echo.Context, token *Token) {
	c.Set("token", token)
}

// GetTokenFromContext is used to get auth token from gin context.
func GetTokenFromContext(c echo.Context) *Token {
	return c.Get("token").(*Token)
}

package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const tokenTTL = 24 * time.Hour

type TokenPayload struct {
	UserID   int    `json:"user_id"`
	UserType string `json:"user_type"`
	Exp      int64  `json:"exp"`
}

type claims struct {
	UserID   int    `json:"user_id"`
	UserType string `json:"user_type"`
	jwt.RegisteredClaims
}

func secret() []byte {
	s := os.Getenv("AUTH_SECRET")
	if s == "" {
		s = "default-secret-change-in-production"
	}
	return []byte(s)
}

// HashPassword returns a bcrypt hash of the password.
func HashPassword(password string) string {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(fmt.Sprintf("bcrypt: %v", err))
	}
	return string(b)
}

// CheckPassword verifies a plain password against a bcrypt hash.
func CheckPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// GenerateToken creates a signed HS256 JWT for the given user.
func GenerateToken(userID int, userType string) (string, error) {
	c := claims{
		UserID:   userID,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(secret())
}

// ValidateToken parses and validates a JWT, returning the payload.
func ValidateToken(tokenStr string) (*TokenPayload, error) {
	var c claims
	token, err := jwt.ParseWithClaims(tokenStr, &c, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret(), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return &TokenPayload{
		UserID:   c.UserID,
		UserType: c.UserType,
		Exp:      c.ExpiresAt.Unix(),
	}, nil
}

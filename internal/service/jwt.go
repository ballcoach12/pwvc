package service

import (
	"errors"
	"os"
	"time"

	"pairwise/internal/domain"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrInvalidToken  = errors.New("invalid token")
	ErrExpiredToken  = errors.New("token has expired")
	ErrMissingClaims = errors.New("token missing required claims")
)

// JWTService handles JWT token generation and validation
type JWTService struct {
	secretKey string
	issuer    string
}

// NewJWTService creates a new JWT service instance
func NewJWTService() *JWTService {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "default-secret-key-change-in-production"
	}

	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		issuer = "pairwise-app"
	}

	return &JWTService{
		secretKey: secretKey,
		issuer:    issuer,
	}
}

// GenerateToken creates a new JWT token for the given user
func (j *JWTService) GenerateToken(userID uint, username string, roles []string) (string, error) {
	now := time.Now()
	expirationTime := now.Add(24 * time.Hour) // Token valid for 24 hours

	claims := &domain.JWTClaims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    j.issuer,
			Subject:   string(rune(userID)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates and parses a JWT token
func (j *JWTService) ValidateToken(tokenString string) (*domain.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*domain.JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// Validate required claims
	if claims.UserID == 0 || claims.Username == "" {
		return nil, ErrMissingClaims
	}

	return claims, nil
}

// ExtractRoles extracts role names from JWT claims
func (j *JWTService) ExtractRoles(claims *domain.JWTClaims) []string {
	if claims == nil || claims.Roles == nil {
		return []string{}
	}
	return claims.Roles
}

// HasRole checks if the user has a specific role
func (j *JWTService) HasRole(claims *domain.JWTClaims, role string) bool {
	roles := j.ExtractRoles(claims)
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

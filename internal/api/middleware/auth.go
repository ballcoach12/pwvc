package middleware

import (
	"net/http"
	"strings"

	"pairwise/internal/domain"
	"pairwise/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware provides JWT authentication middleware
type AuthMiddleware struct {
	jwtService *service.JWTService
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(jwtService *service.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

// JWTAuth middleware validates JWT tokens from cookies or Authorization header
func (m *AuthMiddleware) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to get token from cookie first
		tokenString, err := c.Cookie("auth")

		// If not in cookie, try Authorization header
		if err != nil || tokenString == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Authentication required",
				})
				c.Abort()
				return
			}

			// Extract token from "Bearer <token>" format
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid authorization header format",
				})
				c.Abort()
				return
			}
			tokenString = parts[1]
		}

		// Validate token
		claims, err := m.jwtService.ValidateToken(tokenString)
		if err != nil {
			var statusCode int
			var message string

			switch err {
			case service.ErrExpiredToken:
				statusCode = http.StatusUnauthorized
				message = "Token has expired"
			case service.ErrInvalidToken:
				statusCode = http.StatusUnauthorized
				message = "Invalid token"
			case service.ErrMissingClaims:
				statusCode = http.StatusUnauthorized
				message = "Token missing required claims"
			default:
				statusCode = http.StatusUnauthorized
				message = "Authentication failed"
			}

			c.JSON(statusCode, gin.H{
				"error": message,
			})
			c.Abort()
			return
		}

		// Store user information in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("roles", claims.Roles)
		c.Set("jwt_claims", claims)

		c.Next()
	}
}

// RequireRole middleware checks if the authenticated user has the required role
func (m *AuthMiddleware) RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get JWT claims from context (should be set by JWTAuth middleware)
		claimsInterface, exists := c.Get("jwt_claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		claims, ok := claimsInterface.(*domain.JWTClaims)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid authentication data",
			})
			c.Abort()
			return
		}

		// Check if user has the required role
		if !m.jwtService.HasRole(claims, role) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyRole middleware checks if the authenticated user has any of the required roles
func (m *AuthMiddleware) RequireAnyRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get JWT claims from context
		claimsInterface, exists := c.Get("jwt_claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		claims, ok := claimsInterface.(*domain.JWTClaims)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid authentication data",
			})
			c.Abort()
			return
		}

		// Check if user has any of the required roles
		for _, role := range roles {
			if m.jwtService.HasRole(claims, role) {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error": "Insufficient permissions",
		})
		c.Abort()
	}
}

// GetCurrentUser helper to extract user information from context
func GetCurrentUser(c *gin.Context) (*domain.JWTClaims, error) {
	claimsInterface, exists := c.Get("jwt_claims")
	if !exists {
		return nil, service.ErrInvalidToken
	}

	claims, ok := claimsInterface.(*domain.JWTClaims)
	if !ok {
		return nil, service.ErrInvalidToken
	}

	return claims, nil
}

// GetCurrentUserID helper to extract user ID from context
func GetCurrentUserID(c *gin.Context) (uint, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, service.ErrInvalidToken
	}

	id, ok := userID.(uint)
	if !ok {
		return 0, service.ErrInvalidToken
	}

	return id, nil
}

// GetCurrentUserRoles helper to extract user roles from context
func GetCurrentUserRoles(c *gin.Context) ([]string, error) {
	roles, exists := c.Get("roles")
	if !exists {
		return nil, service.ErrInvalidToken
	}

	roleSlice, ok := roles.([]string)
	if !ok {
		return nil, service.ErrInvalidToken
	}

	return roleSlice, nil
}

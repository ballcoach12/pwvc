# Quickstart Guide: JWT Authentication and Authorization System

**Feature**: JWT Authentication and Authorization System  
**Date**: October 28, 2025  
**Branch**: 001-jwt-auth-system

## Overview

This guide provides step-by-step instructions to implement JWT-based authentication and authorization with role-based access control for the pairwise prioritization web application.

## Prerequisites

- Go 1.23.3+ installed
- Node.js and npm for frontend development
- PostgreSQL database running
- Basic familiarity with Gin framework and React

## Implementation Phases

### Phase 1: Database Setup

#### 1.1 Install JWT Dependencies

```bash
# Backend dependencies
cd /workspaces/pwvc2
go get github.com/golang-jwt/jwt/v4
go get golang.org/x/crypto/bcrypt

# Frontend dependencies
cd web/
npm install js-cookie
npm install @types/js-cookie  # if using TypeScript
```

#### 1.2 Set Up GORM Models for Auto-Migration

Since the application uses SQLite with GORM AutoMigrate, no manual migrations are needed. The database schema will be automatically managed by GORM.

```bash
# No migration files needed - GORM handles schema automatically
# Schema updates happen when you run the application with new struct definitions
cd /workspaces/pwvc2

# The following models will be auto-migrated:
# - Existing User struct (extended with password_hash, is_active)
# - New Role struct
# - New UserRole junction table struct
```

The database schema will be automatically updated when the application starts with the new GORM model definitions.### Phase 2: Backend Implementation

#### 2.1 Create Domain Models

**File**: `internal/domain/auth.go`

```go
package domain

import (
    "time"
    "github.com/golang-jwt/jwt/v4"
)

// JWT Claims structure
type JWTClaims struct {
    UserID   uint     `json:"sub"`
    Username string   `json:"username"`
    Roles    []string `json:"roles"`
    jwt.RegisteredClaims
}

// Login request structure
type LoginRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Password string `json:"password" binding:"required,min=8"`
}
```

**File**: `internal/domain/user.go` (extend existing)

```go
// Add to existing User struct
type User struct {
    // ... existing fields
    PasswordHash string `gorm:"size:255" json:"-"`
    IsActive     bool   `gorm:"default:true" json:"is_active"`

    // Relationships
    Roles []Role `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}
```

#### 2.2 Create JWT Service

**File**: `internal/service/jwt.go`

```go
package service

import (
    "fmt"
    "time"
    "github.com/golang-jwt/jwt/v4"
    "pairwise/internal/domain"
)

type JWTService struct {
    secretKey []byte
    issuer    string
}

func NewJWTService(secretKey, issuer string) *JWTService {
    return &JWTService{
        secretKey: []byte(secretKey),
        issuer:    issuer,
    }
}

func (s *JWTService) GenerateToken(user *domain.User) (string, error) {
    // Extract role names
    roleNames := make([]string, len(user.Roles))
    for i, role := range user.Roles {
        roleNames[i] = role.Name
    }

    // Create claims
    claims := &domain.JWTClaims{
        UserID:   user.ID,
        Username: user.Username,
        Roles:    roleNames,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    s.issuer,
            Subject:   fmt.Sprintf("%d", user.ID),
        },
    }

    // Create and sign token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(s.secretKey)
}

func (s *JWTService) ValidateToken(tokenString string) (*domain.JWTClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &domain.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        return s.secretKey, nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*domain.JWTClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, fmt.Errorf("invalid token")
}
```

#### 2.3 Create Authentication Service

**File**: `internal/service/auth.go`

```go
package service

import (
    "errors"
    "golang.org/x/crypto/bcrypt"
    "pairwise/internal/domain"
    "pairwise/internal/repository"
)

type AuthService struct {
    userRepo   repository.UserRepository
    jwtService *JWTService
}

func NewAuthService(userRepo repository.UserRepository, jwtService *JWTService) *AuthService {
    return &AuthService{
        userRepo:   userRepo,
        jwtService: jwtService,
    }
}

func (s *AuthService) Login(req *domain.LoginRequest) (string, *domain.User, error) {
    // Get user by username
    user, err := s.userRepo.GetUserByUsername(req.Username)
    if err != nil {
        return "", nil, errors.New("invalid credentials")
    }

    // Check if user is active
    if !user.IsActive {
        return "", nil, errors.New("account is inactive")
    }

    // Verify password
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        return "", nil, errors.New("invalid credentials")
    }

    // Load user roles
    if err := s.userRepo.LoadUserRoles(user); err != nil {
        return "", nil, err
    }

    // Generate JWT token
    token, err := s.jwtService.GenerateToken(user)
    if err != nil {
        return "", nil, err
    }

    return token, user, nil
}

func (s *AuthService) HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}
```

#### 2.4 Create JWT Middleware

**File**: `internal/api/middleware/auth.go`

```go
package middleware

import (
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
    "pairwise/internal/service"
)

func JWTAuth(jwtService *service.JWTService) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get token from cookie
        tokenString, err := c.Cookie("auth")
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "success": false,
                "error":   "Authentication required",
                "code":    "AUTH_REQUIRED",
            })
            c.Abort()
            return
        }

        // Validate token
        claims, err := jwtService.ValidateToken(tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "success": false,
                "error":   "Invalid token",
                "code":    "INVALID_TOKEN",
            })
            c.Abort()
            return
        }

        // Store claims in context
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("roles", claims.Roles)

        c.Next()
    }
}

func RequireRole(role string) gin.HandlerFunc {
    return func(c *gin.Context) {
        roles, exists := c.Get("roles")
        if !exists {
            c.JSON(http.StatusForbidden, gin.H{
                "success": false,
                "error":   "Access denied",
                "code":    "ACCESS_DENIED",
            })
            c.Abort()
            return
        }

        userRoles := roles.([]string)
        for _, userRole := range userRoles {
            if userRole == role {
                c.Next()
                return
            }
        }

        c.JSON(http.StatusForbidden, gin.H{
            "success": false,
            "error":   "Insufficient permissions",
            "code":    "INSUFFICIENT_PERMISSIONS",
        })
        c.Abort()
    }
}
```

#### 2.5 Create API Handlers

**File**: `internal/api/auth.go` (extend existing)

```go
func (h *Handler) Login(c *gin.Context) {
    var req domain.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Invalid request format",
            "code":    "INVALID_REQUEST",
        })
        return
    }

    token, user, err := h.authService.Login(&req)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "success": false,
            "error":   err.Error(),
            "code":    "LOGIN_FAILED",
        })
        return
    }

    // Set secure cookie
    c.SetCookie("auth", token, 86400, "/", "", true, true) // 24 hours, secure, httpOnly

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Login successful",
        "user":    user,
    })
}

func (h *Handler) Logout(c *gin.Context) {
    // Clear cookie
    c.SetCookie("auth", "", -1, "/", "", true, true)

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Logout successful",
    })
}

func (h *Handler) GetCurrentUser(c *gin.Context) {
    userID := c.GetUint("user_id")
    username := c.GetString("username")
    roles := c.Get("roles")

    c.JSON(http.StatusOK, gin.H{
        "id":       userID,
        "username": username,
        "roles":    roles,
    })
}
```

#### 2.6 Set Up Routes

**File**: `cmd/server/main.go` (extend existing)

```go
// Add JWT service and middleware
jwtService := service.NewJWTService(os.Getenv("JWT_SECRET"), "pairwise-app")
authService := service.NewAuthService(userRepo, jwtService)

// Authentication routes (public)
auth := r.Group("/api/v1/auth")
{
    auth.POST("/login", handler.Login)
    auth.POST("/logout", handler.Logout)
}

// Protected routes
protected := r.Group("/api/v1")
protected.Use(middleware.JWTAuth(jwtService))
{
    protected.GET("/auth/me", handler.GetCurrentUser)
}

// Admin routes
admin := protected.Group("/admin")
admin.Use(middleware.RequireRole("admin"))
{
    admin.GET("/users", handler.ListUsers)
    admin.POST("/users", handler.CreateUser)
    admin.GET("/users/:id", handler.GetUser)
    admin.PUT("/users/:id", handler.UpdateUser)
    admin.GET("/roles", handler.ListRoles)
}
```

### Phase 3: Frontend Implementation

#### 3.1 Create Authentication Context

**File**: `web/src/contexts/AuthContext.jsx`

```jsx
import React, { createContext, useContext, useState, useEffect } from "react";
import Cookies from "js-cookie";
import { authAPI } from "../services/auth";

const AuthContext = createContext();

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    checkAuthStatus();
  }, []);

  const checkAuthStatus = async () => {
    try {
      const token = Cookies.get("auth");
      if (token) {
        const userData = await authAPI.getCurrentUser();
        setUser(userData);
      }
    } catch (error) {
      // Token might be expired, clear it
      Cookies.remove("auth");
    } finally {
      setLoading(false);
    }
  };

  const login = async (username, password) => {
    const response = await authAPI.login(username, password);
    setUser(response.user);
    return response;
  };

  const logout = async () => {
    await authAPI.logout();
    setUser(null);
    Cookies.remove("auth");
  };

  const hasRole = (role) => {
    return user?.roles?.includes(role) || false;
  };

  const value = {
    user,
    loading,
    login,
    logout,
    hasRole,
    isAuthenticated: !!user,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
```

#### 3.2 Create API Service

**File**: `web/src/services/auth.js`

```javascript
import axios from "axios";

const API_BASE_URL = "/api/v1";

export const authAPI = {
  login: async (username, password) => {
    const response = await axios.post(`${API_BASE_URL}/auth/login`, {
      username,
      password,
    });
    return response.data;
  },

  logout: async () => {
    const response = await axios.post(`${API_BASE_URL}/auth/logout`);
    return response.data;
  },

  getCurrentUser: async () => {
    const response = await axios.get(`${API_BASE_URL}/auth/me`);
    return response.data;
  },
};

export const adminAPI = {
  getUsers: async (page = 1, limit = 20) => {
    const response = await axios.get(`${API_BASE_URL}/admin/users`, {
      params: { page, limit },
    });
    return response.data;
  },

  createUser: async (userData) => {
    const response = await axios.post(`${API_BASE_URL}/admin/users`, userData);
    return response.data;
  },

  updateUser: async (userId, userData) => {
    const response = await axios.put(
      `${API_BASE_URL}/admin/users/${userId}`,
      userData
    );
    return response.data;
  },

  getRoles: async () => {
    const response = await axios.get(`${API_BASE_URL}/admin/roles`);
    return response.data;
  },
};
```

#### 3.3 Create Login Component

**File**: `web/src/components/auth/LoginForm.jsx`

```jsx
import React, { useState } from "react";
import {
  TextField,
  Button,
  Paper,
  Typography,
  Alert,
  Box,
} from "@mui/material";
import { useAuth } from "../../contexts/AuthContext";

const LoginForm = () => {
  const [credentials, setCredentials] = useState({
    username: "",
    password: "",
  });
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const { login } = useAuth();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      await login(credentials.username, credentials.password);
    } catch (err) {
      setError(err.response?.data?.error || "Login failed");
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (field) => (e) => {
    setCredentials({ ...credentials, [field]: e.target.value });
  };

  return (
    <Box
      display="flex"
      justifyContent="center"
      alignItems="center"
      minHeight="100vh"
    >
      <Paper elevation={3} sx={{ padding: 4, maxWidth: 400, width: "100%" }}>
        <Typography variant="h4" component="h1" gutterBottom align="center">
          Login
        </Typography>

        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}

        <form onSubmit={handleSubmit}>
          <TextField
            fullWidth
            label="Username"
            margin="normal"
            required
            value={credentials.username}
            onChange={handleChange("username")}
            disabled={loading}
          />

          <TextField
            fullWidth
            label="Password"
            type="password"
            margin="normal"
            required
            value={credentials.password}
            onChange={handleChange("password")}
            disabled={loading}
          />

          <Button
            type="submit"
            fullWidth
            variant="contained"
            sx={{ mt: 3, mb: 2 }}
            disabled={loading}
          >
            {loading ? "Logging in..." : "Login"}
          </Button>
        </form>
      </Paper>
    </Box>
  );
};

export default LoginForm;
```

#### 3.4 Create Protected Route Component

**File**: `web/src/components/auth/ProtectedRoute.jsx`

```jsx
import React from "react";
import { Navigate } from "react-router-dom";
import { useAuth } from "../../contexts/AuthContext";
import { CircularProgress, Box } from "@mui/material";

const ProtectedRoute = ({ children, requiredRole = null }) => {
  const { user, loading, isAuthenticated, hasRole } = useAuth();

  if (loading) {
    return (
      <Box
        display="flex"
        justifyContent="center"
        alignItems="center"
        minHeight="200px"
      >
        <CircularProgress />
      </Box>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  if (requiredRole && !hasRole(requiredRole)) {
    return <Navigate to="/unauthorized" replace />;
  }

  return children;
};

export default ProtectedRoute;
```

### Phase 4: Testing

#### 4.1 Backend Tests

**File**: `internal/service/auth_test.go`

```go
package service_test

import (
    "testing"
    "pairwise/internal/service"
    "pairwise/internal/domain"
)

func TestAuthService_Login(t *testing.T) {
    // Mock repository and test login functionality
    // Test cases: valid login, invalid credentials, inactive user
}

func TestJWTService_GenerateAndValidateToken(t *testing.T) {
    // Test JWT token generation and validation
    // Test cases: valid token, expired token, invalid signature
}
```

#### 4.2 Frontend Tests

**File**: `web/src/components/auth/LoginForm.test.jsx`

```jsx
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { AuthProvider } from "../../contexts/AuthContext";
import LoginForm from "./LoginForm";

test("renders login form correctly", () => {
  render(
    <AuthProvider>
      <LoginForm />
    </AuthProvider>
  );

  expect(screen.getByLabelText(/username/i)).toBeInTheDocument();
  expect(screen.getByLabelText(/password/i)).toBeInTheDocument();
  expect(screen.getByRole("button", { name: /login/i })).toBeInTheDocument();
});

test("handles login submission", async () => {
  // Test login form submission with valid and invalid credentials
});
```

## Environment Configuration

### Backend Environment Variables

```bash
# Add to .env file
JWT_SECRET=your-super-secret-jwt-signing-key-here-make-it-long-and-random
JWT_ISSUER=pairwise-app
```

### Frontend Environment Variables

```bash
# Add to web/.env
VITE_API_BASE_URL=http://localhost:8080
```

## Deployment Checklist

- [ ] JWT secret is properly configured in production
- [ ] Cookies are set with secure flags in production (HTTPS)
- [ ] GORM models are updated for auto-migration
- [ ] Initial admin user is created
- [ ] Rate limiting is configured for auth endpoints
- [ ] Logging is configured for security events
- [ ] HTTPS is enabled for production deployment

## Common Issues and Troubleshooting

### Cookie Not Being Set

- Ensure secure flag is only set in HTTPS environments
- Check SameSite policy configuration
- Verify domain settings for cookie

### JWT Token Validation Failing

- Verify JWT secret matches between generation and validation
- Check token expiration times
- Ensure proper claim validation

### Role-Based Access Not Working

- Verify roles are properly loaded and included in JWT
- Check middleware order in route configuration
- Validate role names match between frontend and backend

## Next Steps

After implementing the basic authentication system:

1. Add password reset functionality
2. Implement account lockout after failed attempts
3. Add two-factor authentication
4. Implement session management and concurrent login controls
5. Add comprehensive audit logging
6. Set up monitoring and alerting for security events

## Documentation References

- [API Specification](contracts/auth-api.yaml)
- [Data Model](data-model.md)
- [Research Notes](research.md)
- [Feature Specification](spec.md)

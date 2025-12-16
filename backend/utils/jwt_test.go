package utils

import (
	"testing"
	"time"
    "github.com/golang-jwt/jwt/v5"
)

func TestGenerateAndValidateToken(t *testing.T) {
	// 1. Generate Token
	token, err := GenerateToken(1, "admin")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// 2. Validate Token
	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	// 3. Check Claims
	if claims.UserID != 1 {
		t.Errorf("Expected UserID 1, got %d", claims.UserID)
	}
	if claims.Role != "admin" {
		t.Errorf("Expected Role admin, got %s", claims.Role)
	}
}

func TestExpiredToken(t *testing.T) {
    // Manually create expired token
    expirationTime := time.Now().Add(-1 * time.Hour)
    claims := &Claims{
        UserID: 1,
        Role: "user",
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, _ := token.SignedString(jwtKey)

    _, err := ValidateToken(tokenString)
    if err == nil {
        t.Error("Expected error for expired token, got nil")
    }
}

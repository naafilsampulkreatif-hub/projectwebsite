package utils // Package utils

import (
	"time" // Time package
	"os"   // OS package for env vars
	"github.com/golang-jwt/jwt/v5" // JWT library
)

var jwtKey = []byte(os.Getenv("JWT_SECRET")) // Secret key for signing tokens. Should be in env.

// Claims struct defines the data encoded in the token
type Claims struct {
	UserID int    `json:"user_id"` // User's ID
	Role   string `json:"role"`    // User's Role
	jwt.RegisteredClaims // Standard JWT claims (exp, iss, etc.)
}

// GenerateToken creates a new JWT token for a user
func GenerateToken(userID int, role string) (string, error) {
	if len(jwtKey) == 0 { // Check if key is set
		jwtKey = []byte("secret_key_change_me") // Fallback key (unsafe for prod)
	}

	// Create expiration time (24 hours)
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create claims
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime), // Set expiration
		},
	}

	// Create the token using HS256 algorithm
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtKey)

	return tokenString, err // Return the token string and any error
}

// ValidateToken parses and validates a token string
func ValidateToken(tokenString string) (*Claims, error) {
	if len(jwtKey) == 0 { // Check if key is set
		jwtKey = []byte("secret_key_change_me") // Fallback key
	}

	// Parse the claims
	claims := &Claims{}

	// Parse the token
	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil // Return the key for verification
	})

	if err != nil { // If parsing failed
		return nil, err
	}

	if !tkn.Valid { // If token is invalid
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil // Return the valid claims
}

package middleware // Package middleware

import (
	"context"                 // Context package
	"ecommerce-backend/utils" // Utils for JWT
	"net/http"                // HTTP package
	"strings"                 // String manipulation
)

// AuthMiddleware checks for a valid JWT token (Strict)
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" { // If missing
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Check if it starts with "Bearer "
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		// Validate the token
		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString)
		if err != nil { // If invalid
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user info to context (normalize role to lowercase)
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "role", strings.ToLower(claims.Role))

		// Call the next handler with the new context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuthMiddleware checks for a valid JWT token but proceeds even if missing (Guest)
func OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		// If header exists and looks correct, try to validate
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 {
				tokenString := parts[1]
				claims, err := utils.ValidateToken(tokenString)
				if err == nil {
					// Valid token, add to context
					// Normalize role to lowercase for consistent checks
					ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
					ctx = context.WithValue(ctx, "role", strings.ToLower(claims.Role))
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}
		}

		// Proceed without user_id in context (Guest mode)
		next.ServeHTTP(w, r)
	})
}

// AdminMiddleware checks if the user has admin role
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get role from context (set by AuthMiddleware)
		role, ok := r.Context().Value("role").(string)
		if !ok || role != "admin" { // If not admin
			http.Error(w, "Forbidden: Admin access required", http.StatusForbidden)
			return
		}

		// Proceed
		next.ServeHTTP(w, r)
	})
}

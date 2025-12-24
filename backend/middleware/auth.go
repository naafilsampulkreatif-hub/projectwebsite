package middleware // Package middleware

import (
	"context"                 // Context package
	"ecommerce-backend/utils" // Utils for JWT
	"encoding/json"           // JSON encoding
	"net/http"                // HTTP package
	"strings"                 // String manipulation
)

// AuthMiddleware checks for a valid JWT token (Strict)
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" { // If missing
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Authorization header required"})
			return
		}

		// Check if it starts with "Bearer "
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid authorization format"})
			return
		}

		// Validate the token
		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString)
		if err != nil { // If invalid
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid token"})
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
		w.Header().Set("Content-Type", "application/json")

		// Get role from context (set by AuthMiddleware)
		role, ok := r.Context().Value("role").(string)
		if !ok || role != "admin" { // If not admin
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": "Forbidden: Admin access required"})
			return
		}

		// Proceed
		next.ServeHTTP(w, r)
	})
}

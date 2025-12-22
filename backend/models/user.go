package models // Defines the package name 'models'

import "time" // Import time package for timestamp fields

// User struct represents a user in the system
type User struct {
	ID           int       `json:"id"`            // Unique identifier for the user
	Name         string    `json:"name"`          // User's full name
	Email        string    `json:"email"`         // User's email address
	Password     string    `json:"password"`      // User's hashed password (should not be returned in JSON usually, but kept here for struct)
	Role         string    `json:"role"`          // User's role (admin or customer)
	Phone        string    `json:"phone"`         // User's phone number
	Address      string    `json:"address"`       // User's address
	Province     string    `json:"province"`      // User's province
	City         string    `json:"city"`          // User's city
	PostalCode   string    `json:"postal_code"`   // User's postal code
	ProfileImage string    `json:"profile_image"` // User's profile image URL
	CreatedAt    time.Time `json:"created_at"`    // Timestamp when the user was created
}

// LoginRequest struct for login payload
type LoginRequest struct {
	Email    string `json:"email"`    // Email provided during login
	Password string `json:"password"` // Password provided during login
}

// RegisterRequest struct for registration payload
type RegisterRequest struct {
	Name     string `json:"name"`     // Name provided during registration
	Email    string `json:"email"`    // Email provided during registration
	Password string `json:"password"` // Password provided during registration
}

// LoginResponse struct for successful login
type LoginResponse struct {
	Token string `json:"token"` // JWT token
	User  User   `json:"user"`  // User details
}

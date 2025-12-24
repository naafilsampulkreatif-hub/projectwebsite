package handlers // Package handlers

import (
	"database/sql"             // Database interface
	"ecommerce-backend/db"     // Import our DB package
	"ecommerce-backend/models" // Import our models
	"ecommerce-backend/utils"  // Import utils
	"encoding/json"            // JSON encoding/decoding
	"log"
	"net/http" // HTTP server
	"strings"

	"golang.org/x/crypto/bcrypt" // Bcrypt for password hashing
)

// Register handler for creating a new user
func Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest // Declare request struct

	// Decode the JSON body into the struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	// Insert user into database
	// Role defaults to 'customer'. Admin must be set manually in DB for now or via special endpoint
	query := "INSERT INTO users (name, email, password, role) VALUES (?, ?, ?, ?)"
	_, err = db.DB.Exec(query, req.Name, req.Email, string(hashedPassword), "customer")

	if err != nil {
		// Check for duplicate entry (assuming email is unique)
		http.Error(w, "Error registering user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)                                                       // Set status to 201 Created
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"}) // Return success JSON
}

// Login handler for authenticating a user
func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest // Declare request struct

	// Decode the JSON body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Trim and normalise inputs to avoid accidental whitespace mismatches
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Password = strings.TrimSpace(req.Password)

	// Query the user by email
	var user models.User
	var profileImage sql.NullString                               // Handle potential NULL
	var phone, address, province, city, postalCode sql.NullString // Handle potential NULLs

	query := "SELECT id, name, email, password, role, phone, address, province, city, postal_code, profile_image FROM users WHERE email = ?"
	err := db.DB.QueryRow(query, req.Email).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.Role,
		&phone, &address, &province, &city, &postalCode, &profileImage,
	)

	if err == sql.ErrNoRows { // If no user found
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	} else if err != nil { // Other DB error
		log.Printf("Login database error for email=%s: %v", req.Email, err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Convert NULL fields to empty strings
	if phone.Valid {
		user.Phone = phone.String
	}
	if address.Valid {
		user.Address = address.String
	}
	if province.Valid {
		user.Province = province.String
	}
	if city.Valid {
		user.City = city.String
	}
	if postalCode.Valid {
		user.PostalCode = postalCode.String
	}
	if profileImage.Valid {
		user.ProfileImage = profileImage.String
	}

	// Compare the stored hashed password with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil { // If passwords don't match
		// Log for debugging: do not include raw password
		log.Printf("Login failed for email=%s: bcrypt compare error=%v, stored_hash_len=%d", req.Email, err, len(user.Password))
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Return the token and user info
	response := models.LoginResponse{
		Token: token,
		User:  user,
	}

	// Don't send the password back!
	response.User.Password = ""

	w.Header().Set("Content-Type", "application/json") // Set content type
	json.NewEncoder(w).Encode(response)                // Encode response
}

package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"ecommerce-backend/db"
	"ecommerce-backend/models"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// GetUsers returns list of users (admin only)
func GetUsers(w http.ResponseWriter, r *http.Request) {
	query := "SELECT id, name, email, role, created_at FROM users ORDER BY created_at DESC"
	rows, err := db.DB.Query(query)
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt); err != nil {
			continue
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetCustomerInfo returns customer info for users who checked out
func GetCustomerInfo(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetCustomerInfo called by user: %v", r.Context().Value("user_id"))
	
	query := `SELECT id, user_id, session_id, full_name, email, phone, address, province, city, postal_code, created_at 
	          FROM customer_info ORDER BY created_at DESC`
	rows, err := db.DB.Query(query)
	if err != nil {
		log.Printf("GetCustomerInfo query error: %v", err)
		http.Error(w, "Failed to fetch customer info", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type CustomerInfo struct {
		ID         int    `json:"id"`
		UserID     *int   `json:"user_id"`
		SessionID  string `json:"session_id"`
		FullName   string `json:"full_name"`
		Email      string `json:"email"`
		Phone      string `json:"phone"`
		Address    string `json:"address"`
		Province   string `json:"province"`
		City       string `json:"city"`
		PostalCode string `json:"postal_code"`
		CreatedAt  string `json:"created_at"`
	}

	var customers []CustomerInfo
	for rows.Next() {
		var c CustomerInfo
		if err := rows.Scan(&c.ID, &c.UserID, &c.SessionID, &c.FullName, &c.Email, &c.Phone, &c.Address, &c.Province, &c.City, &c.PostalCode, &c.CreatedAt); err != nil {
			log.Printf("GetCustomerInfo scan error: %v", err)
			continue
		}
		customers = append(customers, c)
	}

	log.Printf("GetCustomerInfo returning %d customers", len(customers))
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

// UpdateUserRole updates a user's role (admin only)
func UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var req struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	query := "UPDATE users SET role = ? WHERE id = ?"
	_, err := db.DB.Exec(query, req.Role, userID)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User role updated"})
}

// DeleteUser deletes a user (admin only)
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	query := "DELETE FROM users WHERE id = ?"
	_, err := db.DB.Exec(query, userID)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted"})
}

// GetAdminProfile returns admin's profile info
func GetAdminProfile(w http.ResponseWriter, r *http.Request) {
	// Get admin ID from context (set by middleware)
	adminID := r.Context().Value("user_id")

	query := "SELECT id, name, email, role, created_at FROM users WHERE id = ? AND role = 'admin'"
	var user models.User
	err := db.DB.QueryRow(query, adminID).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt)
	if err == sql.ErrNoRows {
		http.Error(w, "Admin not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UpdateAdminProfile updates admin's name and email
func UpdateAdminProfile(w http.ResponseWriter, r *http.Request) {
	adminID := r.Context().Value("user_id")

	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	query := "UPDATE users SET name = ?, email = ? WHERE id = ? AND role = 'admin'"
	_, err := db.DB.Exec(query, req.Name, req.Email, adminID)
	if err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated"})
}

// ChangeAdminPassword changes admin's password
func ChangeAdminPassword(w http.ResponseWriter, r *http.Request) {
	adminID := r.Context().Value("user_id")

	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Get current password hash
	var currentHash string
	query := "SELECT password FROM users WHERE id = ? AND role = 'admin'"
	err := db.DB.QueryRow(query, adminID).Scan(&currentHash)
	if err != nil {
		http.Error(w, "Admin not found", http.StatusNotFound)
		return
	}

	// Verify old password
	err = bcrypt.CompareHashAndPassword([]byte(currentHash), []byte(req.OldPassword))
	if err != nil {
		http.Error(w, "Invalid old password", http.StatusUnauthorized)
		return
	}

	// Hash new password
	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	// Update password
	updateQuery := "UPDATE users SET password = ? WHERE id = ? AND role = 'admin'"
	_, err = db.DB.Exec(updateQuery, string(newHash), adminID)
	if err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Password changed successfully"})
}

// GetActivityLog returns admin activity log
func GetActivityLog(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "50"
	}

	query := `SELECT id, admin_id, action, entity_type, entity_id, details, created_at 
	          FROM admin_activity_log ORDER BY created_at DESC LIMIT ?`
	rows, err := db.DB.Query(query, limit)
	if err != nil {
		http.Error(w, "Failed to fetch activity log", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type ActivityLog struct {
		ID         int    `json:"id"`
		AdminID    int    `json:"admin_id"`
		Action     string `json:"action"`
		EntityType string `json:"entity_type"`
		EntityID   *int   `json:"entity_id"`
		Details    string `json:"details"`
		CreatedAt  string `json:"created_at"`
	}

	var logs []ActivityLog
	for rows.Next() {
		var log ActivityLog
		if err := rows.Scan(&log.ID, &log.AdminID, &log.Action, &log.EntityType, &log.EntityID, &log.Details, &log.CreatedAt); err != nil {
			continue
		}
		logs = append(logs, log)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

// LogAdminActivity logs an admin action (called internally)
func LogAdminActivity(adminID int, action, entityType string, entityID *int, details map[string]interface{}) error {
	detailsJSON, _ := json.Marshal(details)
	query := "INSERT INTO admin_activity_log (admin_id, action, entity_type, entity_id, details) VALUES (?, ?, ?, ?, ?)"
	_, err := db.DB.Exec(query, adminID, action, entityType, entityID, string(detailsJSON))
	return err
}

// GetUserStats returns general user stats for dashboard
func GetUserStats(w http.ResponseWriter, r *http.Request) {
	var stats struct {
		TotalUsers     int `json:"total_users"`
		AdminUsers     int `json:"admin_users"`
		CustomerUsers  int `json:"customer_users"`
		TotalOrders    int `json:"total_orders"`
		TotalCustomers int `json:"total_customers"`
	}

	db.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&stats.TotalUsers)
	db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'admin'").Scan(&stats.AdminUsers)
	db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'customer'").Scan(&stats.CustomerUsers)
	db.DB.QueryRow("SELECT COUNT(*) FROM orders").Scan(&stats.TotalOrders)
	db.DB.QueryRow("SELECT COUNT(DISTINCT user_id) FROM customer_info WHERE user_id IS NOT NULL").Scan(&stats.TotalCustomers)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

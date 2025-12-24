package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"ecommerce-backend/db"

	"github.com/gorilla/mux"
)

// GetShippingMethods returns all active shipping methods
func GetShippingMethods(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := "SELECT id, name, description, cost, is_active, created_at FROM shipping_methods WHERE is_active = 1 ORDER BY created_at ASC"
	rows, err := db.DB.Query(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch shipping methods"})
		return
	}
	defer rows.Close()

	type ShippingMethod struct {
		ID          int     `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Cost        float64 `json:"cost"`
		IsActive    bool    `json:"is_active"`
		CreatedAt   string  `json:"created_at"`
	}

	var methods []ShippingMethod
	for rows.Next() {
		var m ShippingMethod
		if err := rows.Scan(&m.ID, &m.Name, &m.Description, &m.Cost, &m.IsActive, &m.CreatedAt); err != nil {
			continue
		}
		methods = append(methods, m)
	}

	if methods == nil {
		methods = []ShippingMethod{}
	}

	json.NewEncoder(w).Encode(methods)
}

// AdminGetShippingMethods returns all shipping methods (for admin)
func AdminGetShippingMethods(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := "SELECT id, name, description, cost, is_active, created_at FROM shipping_methods ORDER BY created_at ASC"
	rows, err := db.DB.Query(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch shipping methods"})
		return
	}
	defer rows.Close()

	type ShippingMethod struct {
		ID          int     `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Cost        float64 `json:"cost"`
		IsActive    bool    `json:"is_active"`
		CreatedAt   string  `json:"created_at"`
	}

	var methods []ShippingMethod
	for rows.Next() {
		var m ShippingMethod
		if err := rows.Scan(&m.ID, &m.Name, &m.Description, &m.Cost, &m.IsActive, &m.CreatedAt); err != nil {
			continue
		}
		methods = append(methods, m)
	}

	if methods == nil {
		methods = []ShippingMethod{}
	}

	json.NewEncoder(w).Encode(methods)
}

// CreateShippingMethod creates a new shipping method (admin only)
func CreateShippingMethod(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Cost        float64 `json:"cost"`
		IsActive    bool    `json:"is_active"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}

	if req.Name == "" || req.Cost < 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Name and cost are required"})
		return
	}

	query := "INSERT INTO shipping_methods (name, description, cost, is_active) VALUES (?, ?, ?, ?)"
	res, err := db.DB.Exec(query, req.Name, req.Description, req.Cost, req.IsActive)
	if err != nil {
		log.Printf("Create shipping method error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create shipping method"})
		return
	}

	id, _ := res.LastInsertId()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id, "message": "Shipping method created successfully"})
}

// UpdateShippingMethod updates a shipping method (admin only)
func UpdateShippingMethod(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Cost        float64 `json:"cost"`
		IsActive    *bool   `json:"is_active"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}

	if req.IsActive != nil {
		query := "UPDATE shipping_methods SET name = ?, description = ?, cost = ?, is_active = ? WHERE id = ?"
		_, err := db.DB.Exec(query, req.Name, req.Description, req.Cost, *req.IsActive, id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update shipping method"})
			return
		}
	} else {
		query := "UPDATE shipping_methods SET name = ?, description = ?, cost = ? WHERE id = ?"
		_, err := db.DB.Exec(query, req.Name, req.Description, req.Cost, id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update shipping method"})
			return
		}
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Shipping method updated successfully"})
}

// DeleteShippingMethod deletes a shipping method (admin only)
func DeleteShippingMethod(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	idInt, err := strconv.Atoi(id)
	if err != nil || idInt == 1 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Cannot delete default shipping method"})
		return
	}

	query := "DELETE FROM shipping_methods WHERE id = ?"
	_, err = db.DB.Exec(query, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete shipping method"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Shipping method deleted successfully"})
}

// GetShippingMethodCost returns cost of a specific shipping method
func GetShippingMethodCost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	var cost float64
	query := "SELECT cost FROM shipping_methods WHERE id = ? AND is_active = 1"
	err := db.DB.QueryRow(query, id).Scan(&cost)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Shipping method not found"})
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		return
	}

	json.NewEncoder(w).Encode(map[string]float64{"cost": cost})
}

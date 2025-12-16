package handlers // Package handlers

import (
	"database/sql" // Database interface
	"encoding/json" // JSON encoding
	"net/http" // HTTP server
	"ecommerce-backend/db" // DB package
	"ecommerce-backend/models" // Models package
	"github.com/gorilla/mux" // Router package
)

// GetProducts returns a list of products
func GetProducts(w http.ResponseWriter, r *http.Request) {
	// Query to select all products
	rows, err := db.DB.Query("SELECT id, name, slug, description, price, stock, image_url, category_id FROM products")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close() // Close rows when done

	var products []models.Product // Slice to hold products

	// Iterate over the rows
	for rows.Next() {
		var p models.Product // Temp product
		// Scan columns into struct fields
		if err := rows.Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.Price, &p.Stock, &p.ImageURL, &p.CategoryID); err != nil {
			continue // Skip if scan fails
		}
		products = append(products, p) // Add to slice
	}

	w.Header().Set("Content-Type", "application/json") // Set header
	json.NewEncoder(w).Encode(products) // Encode result
}

// GetProduct returns a single product by ID or Slug
func GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get path variables
	id := vars["id"] // Get ID

	var p models.Product // Product struct
	// Query by ID
	err := db.DB.QueryRow("SELECT id, name, slug, description, price, stock, image_url, category_id FROM products WHERE id = ?", id).
		Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.Price, &p.Stock, &p.ImageURL, &p.CategoryID)

	if err == sql.ErrNoRows { // Not found
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	} else if err != nil { // Error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json") // Set header
	json.NewEncoder(w).Encode(p) // Encode result
}

// CreateProduct creates a new product (Admin only)
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var p models.Product // Product struct
	// Decode JSON
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Insert into DB
	stmt, err := db.DB.Prepare("INSERT INTO products (name, slug, description, price, stock, image_url, category_id) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close() // Close statement

	// Execute statement
	res, err := stmt.Exec(p.Name, p.Slug, p.Description, p.Price, p.Stock, p.ImageURL, p.CategoryID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := res.LastInsertId() // Get new ID
	p.ID = int(id) // Set ID

	w.WriteHeader(http.StatusCreated) // 201 Created
	json.NewEncoder(w).Encode(p) // Return created product
}

// UpdateProduct updates an existing product
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get vars
	id := vars["id"] // Get ID

	var p models.Product // Product struct
	// Decode body
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Update DB
	_, err := db.DB.Exec("UPDATE products SET name=?, slug=?, description=?, price=?, stock=?, image_url=?, category_id=? WHERE id=?",
		p.Name, p.Slug, p.Description, p.Price, p.Stock, p.ImageURL, p.CategoryID, id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json") // Header
	json.NewEncoder(w).Encode(map[string]string{"message": "Product updated"}) // Response
}

// DeleteProduct deletes a product
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get vars
	id := vars["id"] // Get ID

	// Delete from DB
	_, err := db.DB.Exec("DELETE FROM products WHERE id=?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

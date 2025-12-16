package handlers // Package handlers

import (
	"encoding/json" // JSON encoding
	"net/http" // HTTP
	"ecommerce-backend/db" // DB
	"ecommerce-backend/models" // Models
)

// AddToCart adds an item to the user's cart
func AddToCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int) // Get user ID from context

	var item models.CartItem // Cart item struct
	// Decode JSON
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Insert or Update (upsert)
	// If item exists for user, increment quantity
	// MySQL: ON DUPLICATE KEY UPDATE quantity = quantity + VALUES(quantity)
	query := `INSERT INTO cart_items (user_id, product_id, quantity) VALUES (?, ?, ?)
	          ON DUPLICATE KEY UPDATE quantity = quantity + VALUES(quantity)`

	_, err := db.DB.Exec(query, userID, item.ProductID, item.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK) // 200 OK
	json.NewEncoder(w).Encode(map[string]string{"message": "Added to cart"})
}

// GetCart retrieves the user's cart
func GetCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int) // Get user ID

	// Join with products table to get product details
	query := `SELECT c.id, c.user_id, c.product_id, c.quantity,
	                 p.id, p.name, p.price, p.image_url
	          FROM cart_items c
	          JOIN products p ON c.product_id = p.id
	          WHERE c.user_id = ?`

	rows, err := db.DB.Query(query, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []models.CartItem
	for rows.Next() {
		var item models.CartItem
		var p models.Product
		// Scan
		if err := rows.Scan(&item.ID, &item.UserID, &item.ProductID, &item.Quantity,
			&p.ID, &p.Name, &p.Price, &p.ImageURL); err != nil {
			continue
		}
		item.Product = p
		items = append(items, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

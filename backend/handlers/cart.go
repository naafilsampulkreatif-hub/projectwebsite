package handlers // Package handlers

import (
	"database/sql"
	"encoding/json" // JSON encoding
	"net/http" // HTTP
	"ecommerce-backend/db" // DB
	"ecommerce-backend/models" // Models
)

// getIdentity returns (userID *int, sessionID string)
func getIdentity(r *http.Request) (*int, string) {
	// Check Context for UserID
	if val := r.Context().Value("user_id"); val != nil {
		uid := val.(int)
		return &uid, ""
	}
	// Check Header for SessionID
	sessID := r.Header.Get("X-Session-ID")
	return nil, sessID
}

// AddToCart adds an item to the user's cart or guest session cart
func AddToCart(w http.ResponseWriter, r *http.Request) {
	userID, sessionID := getIdentity(r)

	if userID == nil && sessionID == "" {
		http.Error(w, "User ID or Session ID required", http.StatusBadRequest)
		return
	}

	var item models.CartItem // Cart item struct
	// Decode JSON
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Logic: Upsert based on (user_id, product_id) OR (session_id, product_id)
	// Since unique key is (user_id, session_id, product_id) effectively.
	// We handle them separately to avoid complexity.

	var query string
	var err error

	if userID != nil {
		// Logged in user
		query = `INSERT INTO cart_items (user_id, product_id, quantity) VALUES (?, ?, ?)
		         ON DUPLICATE KEY UPDATE quantity = quantity + VALUES(quantity)`
		_, err = db.DB.Exec(query, *userID, item.ProductID, item.Quantity)
	} else {
		// Guest - Fix for MySQL NULL unique constraint behavior
		// Check if item exists first
		var existingQty int
		checkQuery := `SELECT quantity FROM cart_items WHERE session_id = ? AND product_id = ? AND user_id IS NULL`
		err = db.DB.QueryRow(checkQuery, sessionID, item.ProductID).Scan(&existingQty)

		if err == sql.ErrNoRows {
			// Insert
			insertQuery := `INSERT INTO cart_items (session_id, product_id, quantity) VALUES (?, ?, ?)`
			_, err = db.DB.Exec(insertQuery, sessionID, item.ProductID, item.Quantity)
		} else if err == nil {
			// Update
			updateQuery := `UPDATE cart_items SET quantity = quantity + ? WHERE session_id = ? AND product_id = ? AND user_id IS NULL`
			_, err = db.DB.Exec(updateQuery, item.Quantity, sessionID, item.ProductID)
		}
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK) // 200 OK
	json.NewEncoder(w).Encode(map[string]string{"message": "Added to cart"})
}

// GetCart retrieves the user's or guest's cart
func GetCart(w http.ResponseWriter, r *http.Request) {
	userID, sessionID := getIdentity(r)

	if userID == nil && sessionID == "" {
		// Return empty list if no identity
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]models.CartItem{})
		return
	}

	var rows *sql.Rows
	var err error

	if userID != nil {
		query := `SELECT c.id, c.user_id, c.session_id, c.product_id, c.quantity,
		                 p.id, p.name, p.price, p.image_url
		          FROM cart_items c
		          JOIN products p ON c.product_id = p.id
		          WHERE c.user_id = ?`
		rows, err = db.DB.Query(query, *userID)
	} else {
		query := `SELECT c.id, c.user_id, c.session_id, c.product_id, c.quantity,
		                 p.id, p.name, p.price, p.image_url
		          FROM cart_items c
		          JOIN products p ON c.product_id = p.id
		          WHERE c.session_id = ? AND c.user_id IS NULL`
		rows, err = db.DB.Query(query, sessionID)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	items := []models.CartItem{}
	for rows.Next() {
		var item models.CartItem
		var p models.Product
		var uid sql.NullInt64 // temporary scanner
		var sid sql.NullString // temporary scanner

		// Scan
		if err := rows.Scan(&item.ID, &uid, &sid, &item.ProductID, &item.Quantity,
			&p.ID, &p.Name, &p.Price, &p.ImageURL); err != nil {
			continue
		}

		if uid.Valid {
			id := int(uid.Int64)
			item.UserID = &id
		}
		if sid.Valid {
			item.SessionID = sid.String
		}

		item.Product = p
		items = append(items, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// ValidateStock checks if the items in the cart have sufficient stock
func ValidateStock(w http.ResponseWriter, r *http.Request) {
	userID, sessionID := getIdentity(r)

	if userID == nil && sessionID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var rows *sql.Rows
	var err error

	// Query to get cart quantity and actual product stock
	query := `SELECT c.product_id, c.quantity, p.name, p.stock
	          FROM cart_items c
	          JOIN products p ON c.product_id = p.id
	          WHERE `

	if userID != nil {
		query += `c.user_id = ?`
		rows, err = db.DB.Query(query, *userID)
	} else {
		query += `c.session_id = ? AND c.user_id IS NULL`
		rows, err = db.DB.Query(query, sessionID)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type StockIssue struct {
		ProductName string `json:"product_name"`
		Requested   int    `json:"requested"`
		Available   int    `json:"available"`
	}

	issues := []StockIssue{}
	valid := true

	for rows.Next() {
		var pid, reqQty, stock int
		var name string
		if err := rows.Scan(&pid, &reqQty, &name, &stock); err != nil {
			continue
		}

		if reqQty > stock {
			valid = false
			issues = append(issues, StockIssue{
				ProductName: name,
				Requested:   reqQty,
				Available:   stock,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"valid":  valid,
		"issues": issues,
	})
}

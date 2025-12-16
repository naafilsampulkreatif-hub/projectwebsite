package handlers // Package handlers

import (
	"database/sql"
	"encoding/json" // JSON
	"net/http" // HTTP
	"ecommerce-backend/db" // DB
	"ecommerce-backend/models" // Models
)

type CheckoutRequest struct {
	GuestInfo map[string]interface{} `json:"guest_info"` // Generic map for guest info
}

// Checkout creates an order from the cart
func Checkout(w http.ResponseWriter, r *http.Request) {
	userID, sessionID := getIdentity(r)

	if userID == nil && sessionID == "" {
		http.Error(w, "User ID or Session ID required", http.StatusBadRequest)
		return
	}

	var req CheckoutRequest
	// For guests, we expect guest_info in body. For users, it might be optional or pre-filled.
	// We decode body regardless.
	json.NewDecoder(r.Body).Decode(&req)

	// Start a transaction
	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// 1. Get Cart Items
	var rows *sql.Rows
	if userID != nil {
		query := `SELECT c.product_id, c.quantity, p.price
		          FROM cart_items c
		          JOIN products p ON c.product_id = p.id
		          WHERE c.user_id = ?`
		rows, err = tx.Query(query, *userID)
	} else {
		query := `SELECT c.product_id, c.quantity, p.price
		          FROM cart_items c
		          JOIN products p ON c.product_id = p.id
		          WHERE c.session_id = ? AND c.user_id IS NULL`
		rows, err = tx.Query(query, sessionID)
	}

	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []models.OrderItem
	var totalAmount float64 = 0

	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(&item.ProductID, &item.Quantity, &item.Price); err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		items = append(items, item)
		totalAmount += item.Price * float64(item.Quantity)
	}
	rows.Close() // Close before next query

	if len(items) == 0 {
		tx.Rollback()
		http.Error(w, "Cart is empty", http.StatusBadRequest)
		return
	}

	// 2. Create Order
	var res sql.Result
	if userID != nil {
		res, err = tx.Exec("INSERT INTO orders (user_id, total_amount, status) VALUES (?, ?, ?)", *userID, totalAmount, "pending")
	} else {
		guestInfoJSON, _ := json.Marshal(req.GuestInfo)
		res, err = tx.Exec("INSERT INTO orders (session_id, guest_info, total_amount, status) VALUES (?, ?, ?, ?)", sessionID, string(guestInfoJSON), totalAmount, "pending")
	}

	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	orderID, _ := res.LastInsertId()

	// 3. Create Order Items
	stmt, err := tx.Prepare("INSERT INTO order_items (order_id, product_id, quantity, price) VALUES (?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	for _, item := range items {
		if _, err := stmt.Exec(orderID, item.ProductID, item.Quantity, item.Price); err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// 4. Clear Cart
	if userID != nil {
		_, err = tx.Exec("DELETE FROM cart_items WHERE user_id = ?", *userID)
	} else {
		_, err = tx.Exec("DELETE FROM cart_items WHERE session_id = ?", sessionID)
	}

	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		http.Error(w, "Transaction commit failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "Order placed", "order_id": orderID})
}

// GetOrders retrieves orders for the user
func GetOrders(w http.ResponseWriter, r *http.Request) {
	userID, sessionID := getIdentity(r)

	if userID == nil && sessionID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var rows *sql.Rows
	var err error

	if userID != nil {
		rows, err = db.DB.Query("SELECT id, total_amount, status, created_at FROM orders WHERE user_id = ? ORDER BY created_at DESC", *userID)
	} else {
		// Guests can only see orders for their current session
		rows, err = db.DB.Query("SELECT id, total_amount, status, created_at FROM orders WHERE session_id = ? ORDER BY created_at DESC", sessionID)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.TotalAmount, &o.Status, &o.CreatedAt); err != nil {
			continue
		}
		orders = append(orders, o)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

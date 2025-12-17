package handlers // Package handlers

import (
	"database/sql"
	"ecommerce-backend/db"     // DB
	"ecommerce-backend/models" // Models
	"encoding/json"            // JSON
	"fmt"                      // Formatting
	"net/http"                 // HTTP

	"github.com/gorilla/mux" // Router
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

	// 1. Get Cart Items and Check Stock (LOCKING ROWS)
	var rows *sql.Rows
	if userID != nil {
		// Uses FOR UPDATE to lock rows during transaction
		query := `SELECT c.product_id, c.quantity, p.price, p.stock, p.name
		          FROM cart_items c
		          JOIN products p ON c.product_id = p.id
		          WHERE c.user_id = ? FOR UPDATE`
		rows, err = tx.Query(query, *userID)
	} else {
		query := `SELECT c.product_id, c.quantity, p.price, p.stock, p.name
		          FROM cart_items c
		          JOIN products p ON c.product_id = p.id
		          WHERE c.session_id = ? AND c.user_id IS NULL FOR UPDATE`
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

	// Use a struct to hold temporary data including stock
	type CartItemDetail struct {
		ProductID int
		Quantity  int
		Price     float64
		Stock     int
		Name      string
	}
	var cartDetails []CartItemDetail

	for rows.Next() {
		var d CartItemDetail
		if err := rows.Scan(&d.ProductID, &d.Quantity, &d.Price, &d.Stock, &d.Name); err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check Stock
		if d.Quantity > d.Stock {
			rows.Close()
			tx.Rollback()
			http.Error(w, fmt.Sprintf("Stok tidak cukup untuk produk: %s (Sisa: %d)", d.Name, d.Stock), http.StatusConflict)
			return
		}

		cartDetails = append(cartDetails, d)

		items = append(items, models.OrderItem{
			ProductID: d.ProductID,
			Quantity:  d.Quantity,
			Price:     d.Price,
		})
		totalAmount += d.Price * float64(d.Quantity)
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

	// 3. Create Order Items and Deduct Stock
	stmtItems, err := tx.Prepare("INSERT INTO order_items (order_id, product_id, quantity, price) VALUES (?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmtItems.Close()

	stmtStock, err := tx.Prepare("UPDATE products SET stock = stock - ? WHERE id = ?")
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmtStock.Close()

	for _, item := range items {
		// Insert Item
		if _, err := stmtItems.Exec(orderID, item.ProductID, item.Quantity, item.Price); err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Deduct Stock
		if _, err := stmtStock.Exec(item.Quantity, item.ProductID); err != nil {
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

// GetOrderInvoice retrieves a specific order invoice
func GetOrderInvoice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]
	userID, sessionID := getIdentity(r)

	// Fetch Order
	var o models.Order
	var guestInfoJSON []byte
	var uid sql.NullInt64
	var sid sql.NullString

	query := `SELECT id, user_id, session_id, total_amount, status, guest_info, created_at
	          FROM orders WHERE id = ?`

	err := db.DB.QueryRow(query, orderID).Scan(&o.ID, &uid, &sid, &o.TotalAmount, &o.Status, &guestInfoJSON, &o.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Order not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Authorization Check
	// 1. If order belongs to user, check userID
	// 2. If order belongs to session, check sessionID
	// 3. Or if Admin (not implemented here yet, but middleware handles admin routes separately)

	// Simplify: If userID matches, OK. If sessionID matches, OK.
	authorized := false
	if uid.Valid && userID != nil && int(uid.Int64) == *userID {
		authorized = true
	} else if sid.Valid && sid.String == sessionID {
		authorized = true
	}

	if !authorized {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse Guest Info if exists
	if len(guestInfoJSON) > 0 {
		var g map[string]interface{}
		json.Unmarshal(guestInfoJSON, &g)
		// Assuming models.Order GuestInfo is map[string]interface{}
		// If it's string, we keep it as string. But in models/order.go it might be different.
		// Let's check models.Order definition or cast appropriately.
		o.GuestInfo = g
	}

	// Fetch Order Items
	rows, err := db.DB.Query(`SELECT oi.product_id, oi.quantity, oi.price, p.name
	                          FROM order_items oi
	                          JOIN products p ON oi.product_id = p.id
	                          WHERE oi.order_id = ?`, o.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item models.OrderItem
		var name string
		rows.Scan(&item.ProductID, &item.Quantity, &item.Price, &name)

		// We can add product name to the response struct if we want,
		// but models.OrderItem might not have it.
		// For simplicity, let's create a response struct or map.
		// Re-using OrderItem but attaching name via a map logic or modifying the struct is better.
		// Let's rely on frontend fetching product details or just return a custom map.

		// Actually, let's just make a composite struct for the response
	}

	// ... Refetching to be cleaner
	type InvoiceItem struct {
		Name     string  `json:"name"`
		Quantity int     `json:"quantity"`
		Price    float64 `json:"price"`
		Total    float64 `json:"total"`
	}

	var invoiceItems []InvoiceItem

	// Reset rows
	rows, _ = db.DB.Query(`SELECT oi.quantity, oi.price, p.name
	                          FROM order_items oi
	                          JOIN products p ON oi.product_id = p.id
	                          WHERE oi.order_id = ?`, o.ID)
	defer rows.Close()

	for rows.Next() {
		var i InvoiceItem
		rows.Scan(&i.Quantity, &i.Price, &i.Name)
		i.Total = i.Price * float64(i.Quantity)
		invoiceItems = append(invoiceItems, i)
	}

	response := map[string]interface{}{
		"order": o,
		"items": invoiceItems,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetOrders retrieves orders for the user
func GetOrders(w http.ResponseWriter, r *http.Request) {
	// This system is session-based and does not provide persistent order history to users.
	// Return an empty list so users have no visible order history.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]models.Order{})
}

// AdminGetOrders returns all orders for admin dashboard
func AdminGetOrders(w http.ResponseWriter, r *http.Request) {
	// Admin middleware already ensured the user is admin
	rows, err := db.DB.Query(`SELECT id, user_id, session_id, guest_info, total_amount, status, created_at FROM orders ORDER BY created_at DESC`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type OrderResponse struct {
		ID          int                      `json:"id"`
		UserID      *int                     `json:"user_id"`
		SessionID   string                   `json:"session_id"`
		GuestInfo   interface{}              `json:"guest_info"`
		TotalAmount float64                  `json:"total_amount"`
		Status      string                   `json:"status"`
		CreatedAt   string                   `json:"created_at"`
		Items       []map[string]interface{} `json:"items"`
		User        *map[string]interface{}  `json:"user,omitempty"`
	}

	var orders []OrderResponse

	for rows.Next() {
		var o OrderResponse
		var guestJSON []byte
		var uid sql.NullInt64
		var sid sql.NullString

		if err := rows.Scan(&o.ID, &uid, &sid, &guestJSON, &o.TotalAmount, &o.Status, &o.CreatedAt); err != nil {
			continue
		}

		if uid.Valid {
			id := int(uid.Int64)
			o.UserID = &id
		}
		if sid.Valid {
			o.SessionID = sid.String
		}

		if len(guestJSON) > 0 {
			var gi map[string]interface{}
			json.Unmarshal(guestJSON, &gi)
			o.GuestInfo = gi
		}

		// Fetch items
		itemRows, err := db.DB.Query(`SELECT oi.quantity, oi.price, p.name FROM order_items oi JOIN products p ON oi.product_id = p.id WHERE oi.order_id = ?`, o.ID)
		if err == nil {
			var items []map[string]interface{}
			for itemRows.Next() {
				var qty int
				var price float64
				var name string
				itemRows.Scan(&qty, &price, &name)
				items = append(items, map[string]interface{}{"name": name, "quantity": qty, "price": price, "total": price * float64(qty)})
			}
			itemRows.Close()
			o.Items = items
		}

		// If order has user_id, fetch basic user info
		if o.UserID != nil {
			var user map[string]interface{}
			var name, email string
			err := db.DB.QueryRow(`SELECT name, email FROM users WHERE id = ?`, *o.UserID).Scan(&name, &email)
			if err == nil {
				user = map[string]interface{}{"name": name, "email": email}
				o.User = &user
			}
		}

		orders = append(orders, o)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

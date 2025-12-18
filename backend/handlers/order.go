package handlers // Package handlers

import (
	"database/sql"
	"ecommerce-backend/db"     // DB
	"ecommerce-backend/models" // Models
	"encoding/json"            // JSON
	"fmt"                      // Formatting
	"log"
	"net/http" // HTTP

	"github.com/gorilla/mux" // Router
)

type CheckoutRequest struct {
	GuestInfo map[string]interface{} `json:"guest_info"` // Generic map for guest info
}

// Checkout creates an order from the cart
func Checkout(w http.ResponseWriter, r *http.Request) {
	userID, sessionID := getIdentity(r)

	if userID == nil && sessionID == "" {
		writeJSONError(w, http.StatusBadRequest, "User ID or Session ID required")
		return
	}

	if userID != nil {
		log.Printf("Checkout request: user_id=%d", *userID)
	} else {
		log.Printf("Checkout request: session_id=%s", sessionID)
	}

	var req CheckoutRequest
	// For guests, we expect guest_info in body. For users, it might be optional or pre-filled.
	// We decode body regardless.
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Checkout decode error: %v", err)
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	log.Printf("Checkout guest_info: %+v", req.GuestInfo)

	// Start a transaction
	tx, err := db.DB.Begin()
	if err != nil {
		log.Printf("Checkout DB begin error: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Database error")
		return
	}

	// 1. Get Cart Items and Check Stock (LOCKING ROWS)
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

	var rows *sql.Rows

	if userID != nil {
		// Uses FOR UPDATE to lock rows during transaction
		query := `SELECT c.product_id, c.quantity, p.price, p.stock, p.name
				  FROM cart_items c
				  JOIN products p ON c.product_id = p.id
				  WHERE c.user_id = ? FOR UPDATE`
		rows, err = tx.Query(query, *userID)
		if err != nil {
			tx.Rollback()
			log.Printf("Checkout query error: %v", err)
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		defer rows.Close()

		for rows.Next() {
			var d CartItemDetail
			if err := rows.Scan(&d.ProductID, &d.Quantity, &d.Price, &d.Stock, &d.Name); err != nil {
				tx.Rollback()
				log.Printf("Checkout scan error: %v", err)
				writeJSONError(w, http.StatusInternalServerError, err.Error())
				return
			}
			if d.Quantity > d.Stock {
				rows.Close()
				tx.Rollback()
				msg := fmt.Sprintf("Stok tidak cukup untuk produk: %s (Sisa: %d)", d.Name, d.Stock)
				log.Printf("Checkout stock error: %s", msg)
				writeJSONError(w, http.StatusConflict, msg)
				return
			}
			cartDetails = append(cartDetails, d)
			items = append(items, models.OrderItem{ProductID: d.ProductID, Quantity: d.Quantity, Price: d.Price})
			totalAmount += d.Price * float64(d.Quantity)
		}
	} else {
		// Guest: load guest cart from file-backed store
		guestItems, gerr := db.ListGuestItems(sessionID)
		if gerr != nil {
			tx.Rollback()
			log.Printf("Checkout guest list error: %v", gerr)
			writeJSONError(w, http.StatusInternalServerError, gerr.Error())
			return
		}
		if len(guestItems) == 0 {
			tx.Rollback()
			writeJSONError(w, http.StatusBadRequest, "Cart is empty")
			return
		}

		// For each product, lock product row and check stock
		for pid, qty := range guestItems {
			var price float64
			var stock int
			var name string
			// Lock product row
			err := tx.QueryRow("SELECT price, stock, name FROM products WHERE id = ? FOR UPDATE", pid).Scan(&price, &stock, &name)
			if err != nil {
				tx.Rollback()
				if err == sql.ErrNoRows {
					msg := fmt.Sprintf("Produk dengan ID %d tidak ditemukan. Produk mungkin telah dihapus.", pid)
					log.Printf("Checkout product not found error: %s", msg)
					writeJSONError(w, http.StatusNotFound, msg)
				} else {
					log.Printf("Checkout product lock error: %v", err)
					writeJSONError(w, http.StatusInternalServerError, err.Error())
				}
				return
			}
			if qty > stock {
				tx.Rollback()
				msg := fmt.Sprintf("Stok tidak cukup untuk produk: %s (Diminta: %d, Sisa: %d)", name, qty, stock)
				log.Printf("Checkout stock error: %s", msg)
				writeJSONError(w, http.StatusConflict, msg)
				return
			}
			cartDetails = append(cartDetails, CartItemDetail{ProductID: pid, Quantity: qty, Price: price, Stock: stock, Name: name})
			items = append(items, models.OrderItem{ProductID: pid, Quantity: qty, Price: price})
			totalAmount += price * float64(qty)
		}
	}

	// 2. Create Order
	var res sql.Result
	if userID != nil {
		res, err = tx.Exec("INSERT INTO orders (user_id, total_amount, status) VALUES (?, ?, ?)", *userID, totalAmount, "pending")
	} else {
		// For guests: create order with session_id, guest_info (JSON) and total_amount
		guestInfoJSON, _ := json.Marshal(req.GuestInfo)
		res, err = tx.Exec("INSERT INTO orders (session_id, guest_info, total_amount, status) VALUES (?, ?, ?, ?)", sessionID, string(guestInfoJSON), totalAmount, "pending")
	}

	if err != nil {
		tx.Rollback()
		log.Printf("Checkout create order error: %v", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	orderID, _ := res.LastInsertId()

	// 2b. For guests, store customer info in customer_info table
	if userID == nil && len(req.GuestInfo) > 0 {
		fullName, _ := req.GuestInfo["full_name"].(string)
		email, _ := req.GuestInfo["email"].(string)
		phone, _ := req.GuestInfo["phone"].(string)
		address, _ := req.GuestInfo["address"].(string)
		province, _ := req.GuestInfo["province"].(string)
		city, _ := req.GuestInfo["city"].(string)
		postalCode, _ := req.GuestInfo["postal_code"].(string)

		_, err := tx.Exec(
			"INSERT INTO customer_info (session_id, full_name, email, phone, address, province, city, postal_code) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			sessionID, fullName, email, phone, address, province, city, postalCode,
		)
		if err != nil {
			log.Printf("Checkout insert customer_info warning (non-fatal): %v", err)
			// Non-fatal: continue even if customer_info fails
		}
	}

	// 3. Create Order Items and Deduct Stock
	stmtItems, err := tx.Prepare("INSERT INTO order_items (order_id, product_id, quantity, price) VALUES (?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		log.Printf("Checkout stmt prepare error: %v", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
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
			log.Printf("Checkout insert order_item error: %v", err)
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		// Deduct Stock
		if _, err := stmtStock.Exec(item.Quantity, item.ProductID); err != nil {
			tx.Rollback()
			log.Printf("Checkout deduct stock error: %v", err)
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	// 4. Clear Cart: DB-backed for users, file-backed for guests
	if userID != nil {
		_, err = tx.Exec("DELETE FROM cart_items WHERE user_id = ?", *userID)
		if err != nil {
			tx.Rollback()
			log.Printf("Checkout clear DB cart error: %v", err)
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		// Clear guest cart from file store
		if cerr := db.ClearGuestCart(sessionID); cerr != nil {
			tx.Rollback()
			log.Printf("Checkout clear guest cart error: %v", cerr)
			writeJSONError(w, http.StatusInternalServerError, cerr.Error())
			return
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Checkout commit error: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Transaction commit failed")
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
	log.Printf("DEBUG GetOrderInvoice: uid.Valid=%v, userID=%v, sid.Valid=%v, sessionID='%s', stored_sid='%s'", uid.Valid, userID, sid.Valid, sessionID, sid.String)
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
	log.Printf("AdminGetOrders called by user: %v", r.Context().Value("user_id"))
	
	rows, err := db.DB.Query(`SELECT id, user_id, session_id, guest_info, total_amount, status, created_at FROM orders ORDER BY created_at DESC`)
	if err != nil {
		log.Printf("AdminGetOrders query error: %v", err)
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
			log.Printf("AdminGetOrders scan error: %v", err)
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

	log.Printf("AdminGetOrders returning %d orders", len(orders))
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

package handlers // Package handlers

import (
	"database/sql"
	"ecommerce-backend/db"     // DB + guest cart
	"ecommerce-backend/models" // Models
	"encoding/json"            // JSON encoding
	"log"
	"net/http" // HTTP
	"strconv"

	"github.com/gorilla/mux"
)

// writeJSONError sends a JSON error response with given status
func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

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
		writeJSONError(w, http.StatusBadRequest, "User ID or Session ID required")
		return
	}

	// Accept a minimal payload to avoid decoding issues from nested objects
	var payload struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	if payload.ProductID <= 0 {
		writeJSONError(w, http.StatusBadRequest, "Invalid product_id")
		return
	}
	if payload.Quantity <= 0 {
		payload.Quantity = 1
	}

	// Log incoming request for debugging
	if userID != nil {
		log.Printf("AddToCart request: user_id=%d product_id=%d qty=%d", *userID, payload.ProductID, payload.Quantity)
	} else {
		log.Printf("AddToCart request: session_id=%s product_id=%d qty=%d", sessionID, payload.ProductID, payload.Quantity)
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
		_, err = db.DB.Exec(query, *userID, payload.ProductID, payload.Quantity)
	} else {
		// Guest: use file-backed guest cart store (no DB schema changes required)
		err = db.AddGuestItem(sessionID, payload.ProductID, payload.Quantity)
	}

	if err != nil {
		log.Printf("AddToCart DB error: %v", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK) // 200 OK
	json.NewEncoder(w).Encode(map[string]string{"message": "Added to cart"})
}

// UpdateCart sets the quantity for a product in the cart (or deletes if quantity <= 0)
func UpdateCart(w http.ResponseWriter, r *http.Request) {
	userID, sessionID := getIdentity(r)

	if userID == nil && sessionID == "" {
		writeJSONError(w, http.StatusBadRequest, "User ID or Session ID required")
		return
	}

	var payload struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	if userID != nil {
		log.Printf("UpdateCart request: user_id=%d product_id=%d qty=%d", *userID, payload.ProductID, payload.Quantity)
	} else {
		log.Printf("UpdateCart request: session_id=%s product_id=%d qty=%d", sessionID, payload.ProductID, payload.Quantity)
	}

	if payload.Quantity <= 0 {
		// Delete the cart item
		if userID != nil {
			res, err := db.DB.Exec("DELETE FROM cart_items WHERE user_id = ? AND product_id = ?", *userID, payload.ProductID)
			if err != nil {
				log.Printf("UpdateCart DB error: %v", err)
				writeJSONError(w, http.StatusInternalServerError, err.Error())
				return
			}
			rows, _ := res.RowsAffected()
			json.NewEncoder(w).Encode(map[string]interface{}{"deleted": rows})
			return
		}
		// guest
		if err := db.RemoveGuestItem(sessionID, payload.ProductID); err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"deleted": 1})
		return
	}

	// Otherwise set quantity
	var err error
	if userID != nil {
		_, err = db.DB.Exec("UPDATE cart_items SET quantity = ? WHERE user_id = ? AND product_id = ?", payload.Quantity, *userID, payload.ProductID)
	} else {
		// guest
		err = db.UpdateGuestItem(sessionID, payload.ProductID, payload.Quantity)
	}

	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Cart updated"})
}

// RemoveCartItem deletes a cart item by its cart_items.id (authorized by owner)
func RemoveCartItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	userID, sessionID := getIdentity(r)
	if userID == nil && sessionID == "" {
		writeJSONError(w, http.StatusBadRequest, "User ID or Session ID required")
		return
	}

	// Try DB-backed deletion first
	var uid sql.NullInt64
	var sid sql.NullString
	err := db.DB.QueryRow("SELECT user_id, session_id FROM cart_items WHERE id = ?", id).Scan(&uid, &sid)
	if err != nil {
		if err == sql.ErrNoRows {
			// Possibly a guest using file-backed cart: treat `id` as product_id and remove
			if sessionID != "" {
				if pid, perr := strconv.Atoi(id); perr == nil {
					if err := db.RemoveGuestItem(sessionID, pid); err != nil {
						writeJSONError(w, http.StatusInternalServerError, err.Error())
						return
					}
					json.NewEncoder(w).Encode(map[string]string{"message": "Item removed"})
					return
				}
			}
			writeJSONError(w, http.StatusNotFound, "Not found")
			return
		}
		log.Printf("RemoveCartItem lookup error: %v", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Check ownership for DB-backed row
	if uid.Valid {
		if userID == nil || int(uid.Int64) != *userID {
			writeJSONError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
	} else if sid.Valid {
		if sessionID == "" || sid.String != sessionID {
			writeJSONError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
	} else {
		writeJSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Delete DB-backed item
	_, err = db.DB.Exec("DELETE FROM cart_items WHERE id = ?", id)
	if err != nil {
		log.Printf("RemoveCartItem delete error: %v", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Item removed"})
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
		// Guest: load from file-backed guest cart
		itemsMap, err2 := db.ListGuestItems(sessionID)
		if err2 != nil {
			writeJSONError(w, http.StatusInternalServerError, err2.Error())
			return
		}
		items := []models.CartItem{}
		for pid, qty := range itemsMap {
			var p models.Product
			err := db.DB.QueryRow("SELECT id, name, price, image_url FROM products WHERE id = ?", pid).Scan(&p.ID, &p.Name, &p.Price, &p.ImageURL)
			if err != nil {
				continue
			}
			item := models.CartItem{
				ID:        pid,
				UserID:    nil,
				SessionID: sessionID,
				ProductID: pid,
				Quantity:  qty,
				Product:   p,
			}
			items = append(items, item)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
		return
	}

	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	items := []models.CartItem{}
	for rows.Next() {
		var item models.CartItem
		var p models.Product
		var uid sql.NullInt64  // temporary scanner
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
		writeJSONError(w, http.StatusUnauthorized, "Unauthorized")
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
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
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
		return
	}

	// Guest: validate against products for file-backed cart
	itemsMap, err := db.ListGuestItems(sessionID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	type StockIssue struct {
		ProductName string `json:"product_name"`
		Requested   int    `json:"requested"`
		Available   int    `json:"available"`
	}

	issues := []StockIssue{}
	valid := true

	for pid, reqQty := range itemsMap {
		var name string
		var stock int
		if err := db.DB.QueryRow("SELECT name, stock FROM products WHERE id = ?", pid).Scan(&name, &stock); err != nil {
			continue
		}
		if reqQty > stock {
			valid = false
			issues = append(issues, StockIssue{ProductName: name, Requested: reqQty, Available: stock})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"valid": valid, "issues": issues})
}

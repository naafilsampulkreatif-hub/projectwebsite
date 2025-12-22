package handlers

import (
	"database/sql"
	"ecommerce-backend/db"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type SupportMessage struct {
	ID        int       `json:"id"`
	UserID    *int      `json:"user_id"`
	SessionID string    `json:"session_id"`
	Message   string    `json:"message"`
	Type      string    `json:"type"` // 'user' or 'support'
	CreatedAt time.Time `json:"created_at"`
}

// CreateSupportMessage creates a new support message
func CreateSupportMessage(w http.ResponseWriter, r *http.Request) {
	userID, sessionID := getIdentity(r)

	if userID == nil && sessionID == "" {
		writeJSONError(w, http.StatusBadRequest, "User ID or Session ID required")
		return
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	message, ok := req["message"].(string)
	if !ok || message == "" {
		writeJSONError(w, http.StatusBadRequest, "Message is required")
		return
	}

	msgType, ok := req["type"].(string)
	if !ok {
		msgType = "user"
	}

	// For now, only allow 'user' messages from users/guests
	if msgType != "user" {
		writeJSONError(w, http.StatusBadRequest, "Invalid message type")
		return
	}

	var res sql.Result
	var err error

	if userID != nil {
		res, err = db.DB.Exec(
			"INSERT INTO support_messages (user_id, message, type, created_at) VALUES (?, ?, ?, ?)",
			*userID, message, msgType, time.Now(),
		)
	} else {
		res, err = db.DB.Exec(
			"INSERT INTO support_messages (session_id, message, type, created_at) VALUES (?, ?, ?, ?)",
			sessionID, message, msgType, time.Now(),
		)
	}

	if err != nil {
		log.Printf("CreateSupportMessage error: %v", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	msgID, _ := res.LastInsertId()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Message created",
		"id":      msgID,
	})
}

// GetSupportMessages retrieves support messages for a user/guest
func GetSupportMessages(w http.ResponseWriter, r *http.Request) {
	userID, sessionID := getIdentity(r)

	if userID == nil && sessionID == "" {
		writeJSONError(w, http.StatusBadRequest, "User ID or Session ID required")
		return
	}

	var rows *sql.Rows
	var err error

	if userID != nil {
		rows, err = db.DB.Query(
			"SELECT id, user_id, session_id, message, type, created_at FROM support_messages WHERE user_id = ? ORDER BY created_at ASC",
			*userID,
		)
	} else {
		rows, err = db.DB.Query(
			"SELECT id, user_id, session_id, message, type, created_at FROM support_messages WHERE session_id = ? ORDER BY created_at ASC",
			sessionID,
		)
	}

	if err != nil {
		log.Printf("GetSupportMessages error: %v", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var messages []SupportMessage

	for rows.Next() {
		var m SupportMessage
		var uid sql.NullInt64
		var sid sql.NullString

		if err := rows.Scan(&m.ID, &uid, &sid, &m.Message, &m.Type, &m.CreatedAt); err != nil {
			log.Printf("GetSupportMessages scan error: %v", err)
			continue
		}

		if uid.Valid {
			id := int(uid.Int64)
			m.UserID = &id
		}
		if sid.Valid {
			m.SessionID = sid.String
		}

		messages = append(messages, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

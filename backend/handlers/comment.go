package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"ecommerce-backend/db"

	"github.com/gorilla/mux"
)

// Toxic words filter list
var toxicWords = []string{
	"anjing", "babi", "goblok", "bodoh", "idiot", "tolol", "bangsat",
	"monyet", "tukang", "geblek", "edan", "sinting", "gila",
	"jelek", "buruk", "rusa", "kerbau", "kontol", "memek", "pussy",
	"fuck", "shit", "damn", "bastard", "asshole", "dick",
}

// filterToxicWords checks and filters toxic language from comment
func filterToxicWords(text string) (bool, string) {
	lowerText := strings.ToLower(text)
	foundToxic := false

	for _, word := range toxicWords {
		// Use word boundaries for more accurate matching
		re := regexp.MustCompile(`\b` + regexp.QuoteMeta(word) + `\b`)
		if re.MatchString(lowerText) {
			foundToxic = true
			text = re.ReplaceAllString(text, strings.Repeat("*", len(word)))
		}
	}

	return foundToxic, text
}

// CreateComment saves a comment (auto-approved with toxic filter)
func CreateComment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Text  string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Filter toxic words
	hasToxic, cleanedText := filterToxicWords(req.Text)

	// Auto-approve comments (unless they contain toxic language - then filter and still approve)
	status := "approved"
	if hasToxic {
		// Still approve, but with filtered text
		req.Text = cleanedText
	}

	query := "INSERT INTO comments (name, email, text, status) VALUES (?, ?, ?, ?)"
	result, err := db.DB.Exec(query, req.Name, req.Email, req.Text, status)
	if err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	commentID, _ := result.LastInsertId()

	// Fetch the created comment to return full object to frontend
	var c struct {
		ID          int     `json:"id"`
		Name        string  `json:"name"`
		Email       string  `json:"email"`
		Text        string  `json:"text"`
		Status      string  `json:"status"`
		AdminNotes  *string `json:"admin_notes"`
		ModeratedBy *int    `json:"moderated_by"`
		CreatedAt   string  `json:"created_at"`
	}
	err = db.DB.QueryRow("SELECT id, name, email, text, status, admin_notes, moderated_by, created_at FROM comments WHERE id = ?", commentID).
		Scan(&c.ID, &c.Name, &c.Email, &c.Text, &c.Status, &c.AdminNotes, &c.ModeratedBy, &c.CreatedAt)
	if err != nil {
		// If fetch fails, still return basic response with 201 Created
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Komentar dipublikasikan",
			"id":      commentID,
			"flagged": hasToxic,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Komentar dipublikasikan",
		"id":      commentID,
		"flagged": hasToxic,
		"comment": c,
	})
}

// GetComments returns comments (approved ones for frontend, all for admin)
func GetComments(w http.ResponseWriter, r *http.Request) {
	// Check if admin is requesting (has auth context)
	adminID := r.Context().Value("user_id")
	isAdmin := adminID != nil

	var query string
	if isAdmin {
		// Admin sees all comments with status
		query = "SELECT id, name, email, text, status, admin_notes, moderated_by, created_at FROM comments ORDER BY created_at DESC"
	} else {
		// Public sees only approved
		query = "SELECT id, name, email, text, status, admin_notes, moderated_by, created_at FROM comments WHERE status = 'approved' ORDER BY created_at DESC"
	}

	rows, err := db.DB.Query(query)
	if err != nil {
		http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Comment struct {
		ID          int     `json:"id"`
		Name        string  `json:"name"`
		Email       string  `json:"email"`
		Text        string  `json:"text"`
		Status      string  `json:"status"`
		AdminNotes  *string `json:"admin_notes"`
		ModeratedBy *int    `json:"moderated_by"`
		CreatedAt   string  `json:"created_at"`
	}

	var comments []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Text, &c.Status, &c.AdminNotes, &c.ModeratedBy, &c.CreatedAt); err != nil {
			continue
		}
		comments = append(comments, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

// ModerateComment approves or rejects a comment (admin only)
func ModerateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID := vars["id"]
	adminID := r.Context().Value("user_id")

	var req struct {
		Status     string `json:"status"` // 'approved' or 'rejected'
		AdminNotes string `json:"admin_notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Status != "approved" && req.Status != "rejected" {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	query := "UPDATE comments SET status = ?, admin_notes = ?, moderated_by = ?, moderated_at = NOW() WHERE id = ?"
	_, err := db.DB.Exec(query, req.Status, req.AdminNotes, adminID, commentID)
	if err != nil {
		http.Error(w, "Failed to moderate comment", http.StatusInternalServerError)
		return
	}

	// Log admin activity
	cID, _ := strconv.Atoi(commentID)
	LogAdminActivity(adminID.(int), "moderate_comment", "comment", &cID, map[string]interface{}{
		"status":      req.Status,
		"admin_notes": req.AdminNotes,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Comment moderated"})
}

// DeleteComment deletes a comment (admin only)
func DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID := vars["id"]
	adminID := r.Context().Value("user_id")

	query := "DELETE FROM comments WHERE id = ?"
	_, err := db.DB.Exec(query, commentID)
	if err != nil {
		http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
		return
	}

	// Log admin activity
	cID, _ := strconv.Atoi(commentID)
	LogAdminActivity(adminID.(int), "delete_comment", "comment", &cID, map[string]interface{}{})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Comment deleted"})
}

// GetPendingCommentsCount returns count of pending comments (for dashboard badge)
func GetPendingCommentsCount(w http.ResponseWriter, r *http.Request) {
	var count int
	query := "SELECT COUNT(*) FROM comments WHERE status = 'pending'"
	db.DB.QueryRow(query).Scan(&count)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"pending_count": count})
}

package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 1. Get current session ID from frontend localStorage (simulate)
	sessionID := fmt.Sprintf("sess_test_%d", time.Now().Unix())
	fmt.Printf("🔄 Testing checkout with Session ID: %s\n\n", sessionID)

	// 2. Prepare checkout request
	checkoutReq := map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"product_id": 1,
				"quantity":   1,
			},
		},
		"guest_info": map[string]interface{}{
			"full_name":   "Test Customer",
			"email":       "test@example.com",
			"phone":       "08123456789",
			"address":     "Jl. Test 123",
			"province":    "Jawa Barat",
			"city":        "Bandung",
			"postal_code": "40123",
		},
		"shipping_method_id": 1,
	}

	body, _ := json.Marshal(checkoutReq)

	// 3. Make checkout request
	fmt.Println("📤 Sending checkout request to http://localhost:8080/api/checkout")
	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("POST", "http://localhost:8080/api/checkout", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Session-ID", sessionID)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("✅ Response status: %d\n", resp.StatusCode)
	fmt.Printf("✅ Response body: %s\n\n", string(respBody))

	// 4. Check database for customer_info
	time.Sleep(1 * time.Second)
	fmt.Println("📊 Checking database for customer_info...")

	db, _ := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/ecommerce")
	defer db.Close()
	db.Ping()

	rows, _ := db.Query("SELECT id, session_id, full_name, email, phone, city FROM customer_info ORDER BY created_at DESC LIMIT 5")
	defer rows.Close()

	fmt.Println("📋 Recent customer_info records:")
	found := false
	for rows.Next() {
		var id int
		var sessionID, fullName, email, phone, city string
		rows.Scan(&id, &sessionID, &fullName, &email, &phone, &city)
		fmt.Printf("   ID: %d | Session: %s | Name: %s | Email: %s | Phone: %s | City: %s\n",
			id, sessionID, fullName, email, phone, city)
		found = true
	}

	if !found {
		fmt.Println("   ❌ No customer_info records found!")
	}

	fmt.Println("\n✅ Test complete!")
}

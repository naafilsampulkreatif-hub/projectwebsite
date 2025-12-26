package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Connect to MySQL
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/ecommerce")
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer db.Close()

	// Test connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping:", err)
	}

	log.Println("✅ Connected to database")

	// Insert test customer data
	_, err = db.Exec(`
		INSERT INTO customer_info (session_id, full_name, email, phone, address, province, city, postal_code) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, "sess_test_001", "Budi Santoso", "budi@example.com", "08123456789", "Jl. Merdeka 123", "Jawa Barat", "Bandung", "40123")

	if err != nil {
		log.Fatal("Insert failed:", err)
	}

	log.Println("✅ Test data inserted successfully!")

	// Query to verify
	rows, err := db.Query("SELECT id, session_id, full_name, email, phone, address, province, city, postal_code, created_at FROM customer_info")
	if err != nil {
		log.Fatal("Query failed:", err)
	}
	defer rows.Close()

	log.Println("\n=== Customer Info Table ===")
	for rows.Next() {
		var id int
		var sessionID, fullName, email, phone, address, province, city, postalCode, createdAt string
		err := rows.Scan(&id, &sessionID, &fullName, &email, &phone, &address, &province, &city, &postalCode, &createdAt)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("ID: %d | Name: %s | Email: %s | Phone: %s | City: %s | Created: %s", id, fullName, email, phone, city, createdAt)
	}
}

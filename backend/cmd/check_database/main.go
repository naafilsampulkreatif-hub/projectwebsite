package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Try to connect to existing ecommerce database
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/ecommerce")
	if err != nil {
		log.Fatal("Failed to open connection:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("❌ DATABASE ECOMMERCE NOT FOUND OR CONNECTION FAILED")
		fmt.Println("Error:", err)

		// Try connecting to MySQL without db to create it
		db2, _ := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/")
		defer db2.Close()
		if db2.Ping() == nil {
			fmt.Println("✅ MySQL server is running")
			fmt.Println("\n📋 Available databases:")
			rows, _ := db2.Query("SHOW DATABASES;")
			defer rows.Close()
			for rows.Next() {
				var dbName string
				rows.Scan(&dbName)
				fmt.Printf("   - %s\n", dbName)
			}
		}
		return
	}

	fmt.Println("✅ Connected to ecommerce database!")

	// Check what tables exist
	fmt.Println("\n📋 Existing tables:")
	rows, err := db.Query("SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_SCHEMA='ecommerce';")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		rows.Scan(&tableName)

		// Count rows in each table
		var count int
		db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)).Scan(&count)
		fmt.Printf("   - %s (%d rows)\n", tableName, count)
	}

	// Show customer_info details
	fmt.Println("\n📊 Customer Info Records:")
	custRows, _ := db.Query("SELECT id, session_id, full_name, email, phone, city, created_at FROM customer_info ORDER BY created_at DESC LIMIT 10")
	defer custRows.Close()

	for custRows.Next() {
		var id int
		var sessionID, fullName, email, phone, city, createdAt string
		custRows.Scan(&id, &sessionID, &fullName, &email, &phone, &city, &createdAt)
		fmt.Printf("✅ ID: %d | Name: %s | Email: %s | Phone: %s | City: %s | Session: %s | Created: %s\n",
			id, fullName, email, phone, city, sessionID, createdAt)
	}

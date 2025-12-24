package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Database credentials
	dbUser := "root"
	dbPass := ""
	dbHost := "127.0.0.1"
	dbPort := "3306"
	dbName := "ecommerce"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	log.Println("Database connected!")

	// Clear existing users (optional - comment out if you want to keep them)
	_, err = db.Exec("DELETE FROM users")
	if err != nil {
		log.Printf("Error clearing users: %v", err)
	}

	// Test users to create
	users := []struct {
		name     string
		email    string
		password string
		role     string
	}{
		{"Test User", "user@example.com", "password123", "customer"},
		{"Admin User", "admin@example.com", "admin123", "admin"},
	}

	for _, u := range users {
		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password for %s: %v", u.email, err)
			continue
		}

		// Insert user
		query := "INSERT INTO users (name, email, password, role) VALUES (?, ?, ?, ?)"
		_, err = db.Exec(query, u.name, u.email, string(hashedPassword), u.role)
		if err != nil {
			log.Printf("Error inserting user %s: %v", u.email, err)
			continue
		}

		log.Printf("Created user: %s (%s) with password: %s", u.email, u.role, u.password)
	}

	log.Println("User seeding complete!")
}

package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// DB Connection details
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" { dbUser = "root" }
	dbPass := os.Getenv("DB_PASS")
	if dbPass == "" { dbPass = "password" }
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" { dbHost = "127.0.0.1" }
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" { dbPort = "3306" }
	dbName := os.Getenv("DB_NAME")
	if dbName == "" { dbName = "ecommerce" }

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true", dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	fmt.Println("Connected to database.")

	// Read schema.sql
	content, err := ioutil.ReadFile("schema.sql")
	if err != nil {
		log.Fatal("Error reading schema.sql: ", err)
	}

	sqlContent := string(content)

	// Execute the entire script (requires multiStatements=true in DSN)
	_, err = db.Exec(sqlContent)
	if err != nil {
		log.Fatal("Error executing schema: ", err)
	}

	fmt.Println("Schema applied successfully!")
}

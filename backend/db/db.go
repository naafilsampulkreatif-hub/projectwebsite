package db // Defines the package name 'db'

import ( // Import the necessary libraries
	"database/sql" // Standard library for SQL database interactions
	"fmt"          // Standard library for formatting strings and printing
	"log"          // Standard library for logging errors and messages
	"os"           // Standard library for interacting with the operating system (env vars)

	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver anonymously so it registers itself with database/sql
)

var DB *sql.DB // Global variable to hold the database connection pool

// InitDB initializes the database connection
func InitDB() {
	// Get the database username from environment variables, default to 'root'
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" { // Check if the variable is empty
		dbUser = "root" // Set default value
	}

	// Get the database password from environment variables, default to 'password'
	dbPass := os.Getenv("DB_PASS")
	if dbPass == "" { // Check if the variable is empty
		dbPass = "" // Set default value
	}

	// Get the database host from environment variables, default to 'localhost'
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" { // Check if the variable is empty
		dbHost = "127.0.0.1" // Set default value
	}

	// Get the database port from environment variables, default to '3306'
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" { // Check if the variable is empty
		dbPort = "3306" // Set default value
	}

	// Get the database name from environment variables, default to 'ecommerce'
	dbName := os.Getenv("DB_NAME")
	if dbName == "" { // Check if the variable is empty
		dbName = "ecommerce" // Set default value
	}

	// Construct the Data Source Name (DSN) string
	// Format: username:password@tcp(host:port)/dbname?parseTime=true
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)

	var err error // Declare error variable
	// Open a connection to the database
	DB, err = sql.Open("mysql", dsn)
	if err != nil { // Check if there was an error opening the connection object
		log.Fatal("Error opening database: ", err) // Log fatal error and exit
	}

	// Ping the database to verify the connection is actually alive
	err = DB.Ping()
	if err != nil { // Check if the ping failed
		log.Fatal("Error connecting to the database: ", err) // Log fatal error and exit
	}

	// Print success message
	fmt.Println("Database connected successfully!")

	// Ensure a default 'Uncategorized' category exists to avoid FK issues
	_, err = DB.Exec("INSERT INTO categories (name, slug) SELECT ?, ? FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM categories WHERE slug = ?)", "Uncategorized", "uncategorized", "uncategorized")
	if err != nil {
		log.Println("Warning: could not ensure default category 'uncategorized':", err)
	} else {
		log.Println("Ensured default category 'uncategorized' exists or was already present")
	}
}

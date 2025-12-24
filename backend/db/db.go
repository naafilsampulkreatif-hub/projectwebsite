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

	// Ensure user_id in orders table allows NULL (for guest orders)
	_, err = DB.Exec("ALTER TABLE orders MODIFY user_id INT NULL")
	if err != nil {
		log.Println("Warning: could not modify orders.user_id to NULL:", err)
	}

	// Ensure session_id column exists in orders table (for guest tracking)
	var sessionIdCount int
	q := `SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 'orders' AND COLUMN_NAME = 'session_id'`
	err = DB.QueryRow(q, dbName).Scan(&sessionIdCount)
	if err != nil || sessionIdCount == 0 {
		_, err = DB.Exec("ALTER TABLE orders ADD COLUMN session_id VARCHAR(255) AFTER user_id")
		if err != nil {
			log.Println("Warning: could not add session_id column to orders:", err)
		} else {
			log.Println("Added session_id column to orders table")
		}
	}

	// Ensure guest_info column exists in orders table
	var guestInfoCount int
	q = `SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 'orders' AND COLUMN_NAME = 'guest_info'`
	err = DB.QueryRow(q, dbName).Scan(&guestInfoCount)
	if err != nil || guestInfoCount == 0 {
		_, err = DB.Exec("ALTER TABLE orders ADD COLUMN guest_info JSON AFTER session_id")
		if err != nil {
			log.Println("Warning: could not add guest_info column to orders:", err)
		} else {
			log.Println("Added guest_info column to orders table")
		}
	}

	// Ensure a default 'Uncategorized' category exists to avoid FK issues
	_, err = DB.Exec("INSERT INTO categories (name, slug) SELECT ?, ? FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM categories WHERE slug = ?)", "Uncategorized", "uncategorized", "uncategorized")
	if err != nil {
		log.Println("Warning: could not ensure default category 'uncategorized':", err)
	} else {
		log.Println("Ensured default category 'uncategorized' exists or was already present")
	}

	// Ensure cart_items.session_id column exists (migration safety)
	// Some older databases may be missing this column; add it only if absent using INFORMATION_SCHEMA.
	var colCount int
	q = `SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = 'cart_items' AND COLUMN_NAME = 'session_id'`
	err = DB.QueryRow(q, dbName).Scan(&colCount)
	if err != nil {
		log.Println("Warning: could not verify cart_items.session_id existence:", err)
	} else if colCount == 0 {
		_, err = DB.Exec("ALTER TABLE cart_items ADD COLUMN session_id VARCHAR(255)")
		if err != nil {
			log.Println("Warning: could not add cart_items.session_id column:", err)
		} else {
			log.Println("Added cart_items.session_id column")
		}
	} else {
		log.Println("cart_items.session_id column already present")
	}

	// Ensure demo products exist
	ensureDemoProducts()

	// Ensure shipping_methods table and default COD method exist
	ensureShippingMethods()
}

// ensureShippingMethods creates the shipping_methods table and default COD method if needed
func ensureShippingMethods() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("PANIC in ensureShippingMethods:", r)
		}
	}()

	// Create table if not exists
	createTableQuery := `CREATE TABLE IF NOT EXISTS shipping_methods (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		description TEXT,
		cost DECIMAL(10, 2) NOT NULL,
		is_active BOOLEAN DEFAULT 1,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`

	_, err := DB.Exec(createTableQuery)
	if err != nil {
		log.Println("Warning: could not create shipping_methods table:", err)
		return
	}
	log.Println("Shipping methods table ready")

	// Ensure shipping_method_id column exists in orders table
	var colCount int
	q := `SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = (SELECT DATABASE()) AND TABLE_NAME = 'orders' AND COLUMN_NAME = 'shipping_method_id'`
	err = DB.QueryRow(q).Scan(&colCount)
	if err != nil {
		log.Println("Warning: could not verify shipping_method_id column:", err)
		return
	}
	if colCount == 0 {
		_, err = DB.Exec("ALTER TABLE orders ADD COLUMN shipping_method_id INT DEFAULT 1")
		if err != nil {
			log.Println("Warning: could not add shipping_method_id column to orders:", err)
			return
		}
		// Also add foreign key
		_, _ = DB.Exec("ALTER TABLE orders ADD FOREIGN KEY (shipping_method_id) REFERENCES shipping_methods(id) ON DELETE SET DEFAULT")
		log.Println("Added shipping_method_id column and FK to orders table")
	} else {
		log.Println("shipping_method_id column already exists")
	}

	// Insert default COD method if no methods exist
	var count int
	err = DB.QueryRow("SELECT COUNT(*) FROM shipping_methods").Scan(&count)
	if err != nil {
		log.Println("Warning: could not count shipping methods:", err)
		return
	}

	if count == 0 {
		_, err = DB.Exec(`INSERT INTO shipping_methods (id, name, description, cost, is_active) 
		                  VALUES (1, 'COD (Bayar di Tempat)', 'Pembayaran saat pesanan tiba', 50000, 1)`)
		if err != nil {
			log.Println("Warning: could not insert default COD method:", err)
		} else {
			log.Println("Default COD shipping method created")
		}
	} else {
		log.Printf("Shipping methods table has %d entries\n", count)
	}
}

// ensureDemoProducts adds demo products if none exist
func ensureDemoProducts() {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if err != nil {
		log.Println("Warning: could not count products:", err)
		return
	}

	if count == 0 {
		log.Println("No products found, adding demo products...")
		_, err = DB.Exec(`INSERT INTO products (name, slug, description, price, stock, image_url, category_id) VALUES
			(?, ?, ?, ?, ?, ?, ?),
			(?, ?, ?, ?, ?, ?, ?),
			(?, ?, ?, ?, ?, ?, ?),
			(?, ?, ?, ?, ?, ?, ?)`,
			"Produk 1", "produk-1", "Kualitas terbaik", 10000, 10, "http://localhost:8080/uploads/product1.jpg", 1,
			"Produk 2", "produk-2", "Sangat bagus", 15000, 8, "http://localhost:8080/uploads/product2.jpg", 1,
			"Produk 3", "produk-3", "Murah meriah", 25000, 5, "http://localhost:8080/uploads/product3.jpg", 1,
			"Produk 4", "produk-4", "Limited Edition", 50000, 3, "http://localhost:8080/uploads/product4.jpg", 1,
		)
		if err != nil {
			log.Println("Warning: could not insert demo products:", err)
		} else {
			log.Println("Demo products added successfully")
		}
	} else {
		log.Printf("Products table already has %d products\n", count)
	}
}

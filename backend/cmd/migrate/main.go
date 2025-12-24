package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Connect to MySQL
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/ecommerce")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Connected to database")

	// Disable foreign key checks to allow dropping tables
	_, err = db.Exec("SET FOREIGN_KEY_CHECKS = 0;")
	if err != nil {
		log.Printf("Warning disabling FK checks: %v", err)
	}

	// Drop dependent tables first (foreign key constraints)
	// Drop in correct order: cart_items, orders, then users
	dropTables := []string{
		"DROP TABLE IF EXISTS cart_items;",
		"DROP TABLE IF EXISTS order_items;",
		"DROP TABLE IF EXISTS orders;",
		"DROP TABLE IF EXISTS customer_info;",
		"DROP TABLE IF EXISTS support_messages;",
		"DROP TABLE IF EXISTS admin_activity_log;",
		"DROP TABLE IF EXISTS comments;",
		"DROP TABLE IF EXISTS users;",
		"DROP TABLE IF EXISTS products;",
		"DROP TABLE IF EXISTS categories;",
	}

	for _, dropSQL := range dropTables {
		_, err = db.Exec(dropSQL)
		if err != nil {
			log.Printf("Warning dropping table: %v", err)
		}
	}
	log.Println("Dropped dependent tables")

	// Create users table with ALL required columns
	createUsersSQL := `
	CREATE TABLE users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		role ENUM('admin', 'customer') DEFAULT 'customer',
		phone VARCHAR(20),
		address TEXT,
		province VARCHAR(100),
		city VARCHAR(100),
		postal_code VARCHAR(20),
		profile_image VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = db.Exec(createUsersSQL)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}
	log.Println("Created users table with proper schema")

	// Create categories table
	createCategoriesSQL := `
	CREATE TABLE categories (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		slug VARCHAR(255) NOT NULL UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = db.Exec(createCategoriesSQL)
	if err != nil {
		log.Fatal("Failed to create categories table:", err)
	}
	log.Println("Created categories table")

	// Create products table
	createProductsSQL := `
	CREATE TABLE products (
		id INT AUTO_INCREMENT PRIMARY KEY,
		category_id INT,
		name VARCHAR(255) NOT NULL,
		slug VARCHAR(255) NOT NULL UNIQUE,
		description TEXT,
		price DECIMAL(10, 2) NOT NULL,
		stock INT DEFAULT 0,
		image_url VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
	);
	`

	_, err = db.Exec(createProductsSQL)
	if err != nil {
		log.Fatal("Failed to create products table:", err)
	}
	log.Println("Created products table")

	// Recreate cart_items table
	createCartItemsSQL := `
	CREATE TABLE cart_items (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT,
		session_id VARCHAR(255),
		product_id INT NOT NULL,
		quantity INT DEFAULT 1,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
		UNIQUE KEY unique_cart_item (user_id, session_id, product_id)
	);
	`

	_, err = db.Exec(createCartItemsSQL)
	if err != nil {
		log.Fatal("Failed to create cart_items table:", err)
	}
	log.Println("Created cart_items table")

	// Recreate orders table
	createOrdersSQL := `
	CREATE TABLE orders (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT NULL,
		session_id VARCHAR(255),
		guest_info JSON,
		total_amount DECIMAL(10, 2) NOT NULL,
		status ENUM('pending', 'paid', 'shipped', 'cancelled') DEFAULT 'pending',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	`

	_, err = db.Exec(createOrdersSQL)
	if err != nil {
		log.Fatal("Failed to create orders table:", err)
	}
	log.Println("Created orders table")

	// Create order_items table
	createOrderItemsSQL := `
	CREATE TABLE order_items (
		id INT AUTO_INCREMENT PRIMARY KEY,
		order_id INT NOT NULL,
		product_id INT NOT NULL,
		quantity INT NOT NULL,
		price DECIMAL(10, 2) NOT NULL,
		FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
		FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
	);
	`

	_, err = db.Exec(createOrderItemsSQL)
	if err != nil {
		log.Fatal("Failed to create order_items table:", err)
	}
	log.Println("Created order_items table")

	// Create customer_info table
	createCustomerInfoSQL := `
	CREATE TABLE customer_info (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT,
		session_id VARCHAR(255),
		full_name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		phone VARCHAR(20),
		address TEXT,
		province VARCHAR(100),
		city VARCHAR(100),
		postal_code VARCHAR(20),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	`

	_, err = db.Exec(createCustomerInfoSQL)
	if err != nil {
		log.Fatal("Failed to create customer_info table:", err)
	}
	log.Println("Created customer_info table")

	// Create comments table
	createCommentsSQL := `
	CREATE TABLE comments (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		text TEXT NOT NULL,
		status ENUM('pending', 'approved', 'rejected') DEFAULT 'pending',
		admin_notes TEXT,
		moderated_by INT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		moderated_at TIMESTAMP NULL,
		FOREIGN KEY (moderated_by) REFERENCES users(id) ON DELETE SET NULL
	);
	`

	_, err = db.Exec(createCommentsSQL)
	if err != nil {
		log.Fatal("Failed to create comments table:", err)
	}
	log.Println("Created comments table")

	// Create support_messages table
	createSupportMessagesSQL := `
	CREATE TABLE support_messages (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT,
		session_id VARCHAR(255),
		message TEXT NOT NULL,
		type ENUM('user', 'admin') DEFAULT 'user',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	`

	_, err = db.Exec(createSupportMessagesSQL)
	if err != nil {
		log.Fatal("Failed to create support_messages table:", err)
	}
	log.Println("Created support_messages table")

	// Create admin_activity_log table
	createActivityLogSQL := `
	CREATE TABLE admin_activity_log (
		id INT AUTO_INCREMENT PRIMARY KEY,
		admin_id INT NOT NULL,
		action VARCHAR(255) NOT NULL,
		entity_type VARCHAR(100),
		entity_id INT,
		details JSON,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (admin_id) REFERENCES users(id) ON DELETE CASCADE
	);
	`

	_, err = db.Exec(createActivityLogSQL)
	if err != nil {
		log.Fatal("Failed to create admin_activity_log table:", err)
	}
	log.Println("Created admin_activity_log table")

	// Re-enable foreign key checks
	_, err = db.Exec("SET FOREIGN_KEY_CHECKS = 1;")
	if err != nil {
		log.Printf("Warning re-enabling FK checks: %v", err)
	}

	// Verify table structure
	rows, err := db.Query("DESCRIBE users;")
	if err != nil {
		log.Fatal("Failed to describe users table:", err)
	}
	defer rows.Close()

	fmt.Println("\n=== Users Table Structure ===")
	for rows.Next() {
		var field, typ, null, key, dflt, extra string
		if err := rows.Scan(&field, &typ, &null, &key, &dflt, &extra); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%-20s | %-30s | %s\n", field, typ, null)
	}

	log.Println("\nMigration complete!")
}

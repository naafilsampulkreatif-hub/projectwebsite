package models // Defines the package name 'models'

import "time" // Import time package

// Product struct represents a product in the store
type Product struct {
	ID          int       `json:"id"`          // Unique identifier
	CategoryID  int       `json:"category_id"` // Foreign key to Category
	Name        string    `json:"name"`        // Product name
	Slug        string    `json:"slug"`        // URL-friendly name
	Description string    `json:"description"` // detailed description
	Price       float64   `json:"price"`       // Product price
	Stock       int       `json:"stock"`       // Available stock quantity
	ImageURL    string    `json:"image_url"`   // URL to product image
	CreatedAt   time.Time `json:"created_at"`  // Timestamp
}

// Category struct represents a product category
type Category struct {
	ID        int       `json:"id"`        // Unique identifier
	Name      string    `json:"name"`      // Category name
	Slug      string    `json:"slug"`      // URL-friendly name
	CreatedAt time.Time `json:"created_at"` // Timestamp
}

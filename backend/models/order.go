package models // Defines the package name 'models'

import "time" // Import time package

// Order struct represents a customer order
type Order struct {
	ID          int         `json:"id"`           // Unique identifier
	UserID      *int        `json:"user_id"`      // ID of the user who placed the order (Nullable for guests)
	SessionID   string      `json:"session_id"`   // Session ID for guest orders
	GuestInfo   string      `json:"guest_info"`   // JSON string for guest details
	TotalAmount float64     `json:"total_amount"` // Total cost of the order
	Status      string      `json:"status"`       // Order status (pending, paid, etc.)
	CreatedAt   time.Time   `json:"created_at"`   // Timestamp
	Items       []OrderItem `json:"items"`        // List of items in the order
}

// OrderItem struct represents a single item within an order
type OrderItem struct {
	ID        int     `json:"id"`         // Unique identifier
	OrderID   int     `json:"order_id"`   // ID of the parent order
	ProductID int     `json:"product_id"` // ID of the product
	Quantity  int     `json:"quantity"`   // Quantity ordered
	Price     float64 `json:"price"`      // Price at the time of order
}

// CartItem struct represents an item in the shopping cart
type CartItem struct {
	ID        int     `json:"id"`         // Unique identifier
	UserID    *int    `json:"user_id"`    // User who owns the cart (Nullable)
	SessionID string  `json:"session_id"` // Session ID for guests
	ProductID int     `json:"product_id"` // Product being added
	Quantity  int     `json:"quantity"`   // Quantity
	Product   Product `json:"product"`    // details of the product (for frontend display)
}

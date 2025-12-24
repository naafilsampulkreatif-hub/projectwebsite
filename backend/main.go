package main // Main package

import (
	"fmt"      // Formatting
	"log"      // Logging
	"net/http" // HTTP
	"os"       // OS

	"ecommerce-backend/db"         // DB
	"ecommerce-backend/handlers"   // Handlers
	"ecommerce-backend/middleware" // Middleware

	"github.com/gorilla/mux"   // Router
	"github.com/joho/godotenv" // Load .env file
	"github.com/rs/cors"       // CORS
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize Database
	db.InitDB()

	// Initialize Router
	r := mux.NewRouter()

	// Static Files
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	// Public Routes (Auth)
	r.HandleFunc("/api/register", handlers.Register).Methods("POST")                               // Register
	r.HandleFunc("/api/login", handlers.Login).Methods("POST")                                     // Login
	r.HandleFunc("/api/products", handlers.GetProducts).Methods("GET")                             // List Products
	r.HandleFunc("/api/products/{id}", handlers.GetProduct).Methods("GET")                         // Get Product
	r.HandleFunc("/api/comments", handlers.GetComments).Methods("GET")                             // Get approved comments
	r.HandleFunc("/api/shipping-methods", handlers.GetShippingMethods).Methods("GET")              // Get active shipping methods
	r.HandleFunc("/api/shipping-methods/{id}/cost", handlers.GetShippingMethodCost).Methods("GET") // Get shipping method cost

	// Hybrid Routes (User or Guest)
	hybridRouter := r.PathPrefix("/api").Subrouter()
	hybridRouter.Use(middleware.OptionalAuthMiddleware)
	hybridRouter.HandleFunc("/cart", handlers.GetCart).Methods("GET")                          // View Cart
	hybridRouter.HandleFunc("/cart", handlers.AddToCart).Methods("POST")                       // Add to Cart
	hybridRouter.HandleFunc("/cart", handlers.UpdateCart).Methods("PUT")                       // Update cart (qty)
	hybridRouter.HandleFunc("/cart/{id}", handlers.RemoveCartItem).Methods("DELETE")           // Remove cart item
	hybridRouter.HandleFunc("/cart/validate-stock", handlers.ValidateStock).Methods("GET")     // Validate Stock
	hybridRouter.HandleFunc("/checkout", handlers.Checkout).Methods("POST")                    // Checkout
	hybridRouter.HandleFunc("/invoice/{id}", handlers.GetOrderInvoice).Methods("GET")          // Get Invoice
	hybridRouter.HandleFunc("/orders", handlers.GetOrders).Methods("GET")                      // My Orders (or Session Orders)
	hybridRouter.HandleFunc("/comments", handlers.CreateComment).Methods("POST")               // Create comment
	hybridRouter.HandleFunc("/support/message", handlers.CreateSupportMessage).Methods("POST") // Send support message
	hybridRouter.HandleFunc("/support/messages", handlers.GetSupportMessages).Methods("GET")   // Get support messages

	// Protected Routes (Authenticated User)
	userRouter := r.PathPrefix("/api").Subrouter()
	userRouter.Use(middleware.AuthMiddleware)
	userRouter.HandleFunc("/profile", handlers.GetUserProfile).Methods("GET")    // Get user profile
	userRouter.HandleFunc("/profile", handlers.UpdateUserProfile).Methods("PUT") // Update user profile

	// Protected Routes (Admin)
	adminRouter := r.PathPrefix("/api/admin").Subrouter()
	adminRouter.Use(middleware.AuthMiddleware)                                         // Must be logged in
	adminRouter.Use(middleware.AdminMiddleware)                                        // Must be admin
	adminRouter.HandleFunc("/products", handlers.CreateProduct).Methods("POST")        // Add Product
	adminRouter.HandleFunc("/products/{id}", handlers.UpdateProduct).Methods("PUT")    // Update Product
	adminRouter.HandleFunc("/products/{id}", handlers.DeleteProduct).Methods("DELETE") // Delete Product
	adminRouter.HandleFunc("/upload", handlers.UploadFile).Methods("POST")             // Upload File
	adminRouter.HandleFunc("/orders", handlers.AdminGetOrders).Methods("GET")          // Admin: list orders
	adminRouter.HandleFunc("/orders/{id}", handlers.AdminUpdateOrder).Methods("PUT")   // Admin: update order status
	// New admin endpoints
	adminRouter.HandleFunc("/users", handlers.GetUsers).Methods("GET")                                 // Get all users
	adminRouter.HandleFunc("/users/{id}/role", handlers.UpdateUserRole).Methods("PUT")                 // Update user role
	adminRouter.HandleFunc("/users/{id}", handlers.DeleteUser).Methods("DELETE")                       // Delete user
	adminRouter.HandleFunc("/customer-info", handlers.GetCustomerInfo).Methods("GET")                  // Get customer checkout data
	adminRouter.HandleFunc("/profile", handlers.GetAdminProfile).Methods("GET")                        // Get admin profile
	adminRouter.HandleFunc("/profile", handlers.UpdateAdminProfile).Methods("PUT")                     // Update admin profile
	adminRouter.HandleFunc("/password", handlers.ChangeAdminPassword).Methods("POST")                  // Change admin password
	adminRouter.HandleFunc("/activity-log", handlers.GetActivityLog).Methods("GET")                    // Get activity log
	adminRouter.HandleFunc("/stats", handlers.GetUserStats).Methods("GET")                             // Get user stats
	adminRouter.HandleFunc("/comments", handlers.GetComments).Methods("GET")                           // Get all comments (admin view)
	adminRouter.HandleFunc("/comments/{id}/moderate", handlers.ModerateComment).Methods("POST")        // Moderate comment
	adminRouter.HandleFunc("/comments/{id}", handlers.DeleteComment).Methods("DELETE")                 // Delete comment
	adminRouter.HandleFunc("/comments/pending-count", handlers.GetPendingCommentsCount).Methods("GET") // Pending comments count
	adminRouter.HandleFunc("/shipping-methods", handlers.AdminGetShippingMethods).Methods("GET")       // List all shipping methods
	adminRouter.HandleFunc("/shipping-methods", handlers.CreateShippingMethod).Methods("POST")         // Create shipping method
	adminRouter.HandleFunc("/shipping-methods/{id}", handlers.UpdateShippingMethod).Methods("PUT")     // Update shipping method
	adminRouter.HandleFunc("/shipping-methods/{id}", handlers.DeleteShippingMethod).Methods("DELETE")  // Delete shipping method

	// CORS Handler
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://127.0.0.1:5173"}, // Vue dev server
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},        // Methods
		AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Session-ID"},  // Headers (Added X-Session-ID)
		AllowCredentials: true,
	})

	handler := c.Handler(r) // Wrap router with CORS

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler)) // Start server
}

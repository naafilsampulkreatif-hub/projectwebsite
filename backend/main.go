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
	r.HandleFunc("/api/register", handlers.Register).Methods("POST")       // Register
	r.HandleFunc("/api/login", handlers.Login).Methods("POST")             // Login
	r.HandleFunc("/api/products", handlers.GetProducts).Methods("GET")     // List Products
	r.HandleFunc("/api/products/{id}", handlers.GetProduct).Methods("GET") // Get Product

	// Hybrid Routes (User or Guest)
	hybridRouter := r.PathPrefix("/api").Subrouter()
	hybridRouter.Use(middleware.OptionalAuthMiddleware)
	hybridRouter.HandleFunc("/cart", handlers.GetCart).Methods("GET")                        // View Cart
	hybridRouter.HandleFunc("/cart", handlers.AddToCart).Methods("POST")                     // Add to Cart
	hybridRouter.HandleFunc("/cart/validate-stock", handlers.ValidateStock).Methods("GET")   // Validate Stock
	hybridRouter.HandleFunc("/checkout", handlers.Checkout).Methods("POST")                  // Checkout
	hybridRouter.HandleFunc("/orders", handlers.GetOrders).Methods("GET")                    // My Orders (or Session Orders)
	hybridRouter.HandleFunc("/orders/{id}/invoice", handlers.GetOrderInvoice).Methods("GET") // Get Invoice

	// Protected Routes (Admin)
	adminRouter := r.PathPrefix("/api/admin").Subrouter()
	adminRouter.Use(middleware.AuthMiddleware)                                         // Must be logged in
	adminRouter.Use(middleware.AdminMiddleware)                                        // Must be admin
	adminRouter.HandleFunc("/products", handlers.CreateProduct).Methods("POST")        // Add Product
	adminRouter.HandleFunc("/products/{id}", handlers.UpdateProduct).Methods("PUT")    // Update Product
	adminRouter.HandleFunc("/products/{id}", handlers.DeleteProduct).Methods("DELETE") // Delete Product
	adminRouter.HandleFunc("/upload", handlers.UploadFile).Methods("POST")             // Upload File
	adminRouter.HandleFunc("/orders", handlers.AdminGetOrders).Methods("GET")          // Admin: list orders

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

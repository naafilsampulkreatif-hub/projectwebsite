package main // Main package

import (
	"fmt" // Formatting
	"log" // Logging
	"net/http" // HTTP
	"os" // OS

	"ecommerce-backend/db" // DB
	"ecommerce-backend/handlers" // Handlers
	"ecommerce-backend/middleware" // Middleware

	"github.com/gorilla/mux" // Router
	"github.com/rs/cors" // CORS
	"github.com/joho/godotenv" // Load .env file
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

	// Public Routes
	r.HandleFunc("/api/register", handlers.Register).Methods("POST") // Register
	r.HandleFunc("/api/login", handlers.Login).Methods("POST") // Login
	r.HandleFunc("/api/products", handlers.GetProducts).Methods("GET") // List Products
	r.HandleFunc("/api/products/{id}", handlers.GetProduct).Methods("GET") // Get Product

	// Protected Routes (User)
	userRouter := r.PathPrefix("/api").Subrouter()
	userRouter.Use(middleware.AuthMiddleware)
	userRouter.HandleFunc("/cart", handlers.GetCart).Methods("GET") // View Cart
	userRouter.HandleFunc("/cart", handlers.AddToCart).Methods("POST") // Add to Cart
	userRouter.HandleFunc("/checkout", handlers.Checkout).Methods("POST") // Checkout
	userRouter.HandleFunc("/orders", handlers.GetOrders).Methods("GET") // My Orders

	// Protected Routes (Admin)
	adminRouter := r.PathPrefix("/api/admin").Subrouter()
	adminRouter.Use(middleware.AuthMiddleware)
	adminRouter.Use(middleware.AdminMiddleware)
	adminRouter.HandleFunc("/products", handlers.CreateProduct).Methods("POST") // Add Product
	adminRouter.HandleFunc("/products/{id}", handlers.UpdateProduct).Methods("PUT") // Update Product
	adminRouter.HandleFunc("/products/{id}", handlers.DeleteProduct).Methods("DELETE") // Delete Product

	// CORS Handler
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://127.0.0.1:5173"}, // Vue dev server
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Methods
		AllowedHeaders:   []string{"Authorization", "Content-Type"}, // Headers
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

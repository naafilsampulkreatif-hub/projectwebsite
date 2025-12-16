# Molla Ecommerce (Go + Vue)

A responsive ecommerce website built with a Golang backend and Vue.js frontend, featuring a multi-level admin dashboard, JWT authentication, and full shopping cart functionality.

## Prerequisites

- Go 1.22+
- Node.js 18+
- MySQL 8.0+

## Setup

### Database

1. Create a MySQL database named `ecommerce`.
2. Run the initialization script:
   ```bash
   mysql -u root -p ecommerce < backend/schema.sql
   ```
3. Insert an admin user manually into the `users` table if you want to access the Admin Dashboard immediately (role='admin').

### Backend

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Run the server:
   ```bash
   # Set environment variables if needed (defaults provided in code)
   export DB_PASS=yourpassword
   go run main.go
   ```
   Server runs on `http://localhost:8080`.

### Frontend

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```
2. Install dependencies:
   ```bash
   npm install
   ```
3. Run the development server:
   ```bash
   npm run dev
   ```
   App runs on `http://localhost:5173`.

## Features

- **Authentication**: Register/Login with JWT.
- **RBAC**: Admin vs Customer roles.
- **Admin Dashboard**: Manage products (CRUD).
- **Shopping**: Browse products, Add to Cart, Checkout.
- **Design**: Responsive UI with TailwindCSS.

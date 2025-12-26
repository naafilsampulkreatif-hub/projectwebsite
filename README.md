# Night Stalkers E-Commerce Platform

A modern e-commerce platform for selling manga and light novels from the Night Stalkers intellectual property. Built with Vue.js 3, Go, and MySQL.

## 🚀 Quick Start

### Prerequisites

Before you begin, ensure you have the following installed on your system:

- **Go** (version 1.22 or higher) - [Download](https://golang.org/dl/)
- **Node.js** (version 18 or higher) - [Download](https://nodejs.org/)
- **npm** (comes with Node.js)
- **MySQL** (version 8.0 or higher) - [Download](https://www.mysql.com/downloads/)
- **Git** - [Download](https://git-scm.com/)

### Installation

#### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/projectwebsite.git
cd projectwebsite
```

#### 2. Backend Setup (Go)

Navigate to the backend directory:

```bash
cd backend
```

Install Go dependencies:

```bash
go mod download
```

Create a `.env` file in the backend directory (if needed for configuration):

```bash
# Copy example if exists
cp .env.example .env

# Or create manually with your database credentials
echo "DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=ecommerce" > .env
```

Build the application (optional):

```bash
go build -o app main.go
```

#### 3. Frontend Setup (Vue.js + npm)

Navigate to the frontend directory:

```bash
cd ../frontend
```

Install dependencies:

```bash
npm install
```

Create environment configuration if needed:

```bash
# Copy example if exists
cp .env.example .env.local
```

---

## 🏃 Running the Project

### Start Backend Server (Go)

From the `backend` directory:

```bash
# Option 1: Run directly
go run main.go

# Option 2: Run compiled binary (if built)
./app
```

The backend API will be available at `http://localhost:8080` (or your configured port).

**Expected output:**
```
Server is running on :8080
```

### Start Frontend Development Server (Vue.js)

From the `frontend` directory:

```bash
npm run dev
```

The frontend will be available at `http://localhost:5173` (Vite default port).

**You should see:**
```
  ➜  Local:   http://localhost:5173/
  ➜  press h to show help
```

---

## 📁 Project Structure

```
projectwebsite/
├── backend/                 # Go API server
│   ├── cmd/                # Command-line tools
│   ├── handlers/           # HTTP request handlers
│   ├── models/             # Data models
│   ├── db/                 # Database initialization
│   ├── middleware/         # Middleware functions
│   ├── main.go             # Entry point
│   ├── go.mod              # Go module file
│   └── schema.sql          # Database schema
│
├── frontend/               # Vue.js application
│   ├── src/
│   │   ├── components/     # Vue components
│   │   ├── views/          # Page views
│   │   ├── stores/         # State management (Pinia)
│   │   ├── services/       # API services
│   │   ├── router/         # Route configuration
│   │   ├── App.vue         # Root component
│   │   └── main.ts         # Entry point
│   ├── public/             # Static assets
│   ├── package.json        # npm dependencies
│   └── vite.config.ts      # Vite configuration
│
└── README.md               # This file
```

---

## 🗄️ Database Setup

### Create Database

Connect to MySQL and create the database:

```bash
mysql -u root -p
```

Inside MySQL shell:

```sql
CREATE DATABASE IF NOT EXISTS ecommerce;
USE ecommerce;
```

### Run Schema

Import the schema file:

```bash
mysql -u root -p ecommerce < backend/schema.sql
```

Or manually run the SQL commands from `backend/schema.sql`.

---

## 🛠️ Development Commands

### Backend (Go)

```bash
# Run with hot reload (using air - install first if needed)
air

# Or install air
go install github.com/cosmtrek/air@latest

# Run tests
go test ./...

# Format code
go fmt ./...

# Vet code for common mistakes
go vet ./...
```

### Frontend (Vue.js)

```bash
# Development server with hot reload
npm run dev

# Build for production
npm run build

# Preview production build locally
npm run preview

# Run tests (if configured)
npm run test

# Lint code
npm run lint

# Format code
npm run format
```

---

## 📝 API Documentation

Key API endpoints:

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login user
- `POST /api/auth/logout` - Logout user

### Products
- `GET /api/products` - Get all products
- `GET /api/products/:id` - Get product by ID

### Cart
- `GET /api/cart` - Get user's cart
- `POST /api/cart/add` - Add item to cart
- `PUT /api/cart/update` - Update cart item quantity
- `DELETE /api/cart/remove/:id` - Remove item from cart

### Orders
- `POST /api/orders` - Create new order
- `GET /api/orders` - Get user's orders
- `GET /api/orders/:id` - Get order details

---

## 🔐 Environment Variables

### Backend (.env)

```
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=ecommerce
PORT=8080
JWT_SECRET=your_jwt_secret_key
```

### Frontend (.env.local)

```
VITE_API_URL=http://localhost:8080/api
VITE_APP_NAME=Night Stalkers
```

---

## 🐛 Troubleshooting

### Port Already in Use

**Backend (Port 8080):**
```bash
# Find process using port 8080
lsof -i :8080

# Kill process (macOS/Linux)
kill -9 <PID>

# Or use different port in .env
PORT=8081
```

**Frontend (Port 5173):**
```bash
# Find process using port 5173
lsof -i :5173

# Kill process (macOS/Linux)
kill -9 <PID>
```

### Database Connection Error

```bash
# Verify MySQL is running
mysql -u root -p

# Check database credentials in .env
# Ensure database exists
mysql -u root -p -e "SHOW DATABASES;"
```

### Module Not Found (Node)

```bash
# Clean install
rm -rf node_modules package-lock.json
npm install
```

### Go Dependencies Issue

```bash
# Clean go modules
go clean -modcache
go mod tidy
go mod download
```

---

## 📚 Tech Stack

| Layer | Technology | Version |
|-------|-----------|---------|
| Frontend | Vue.js | 3.x |
| Frontend Build | Vite | 7.3.0+ |
| Frontend Styling | Tailwind CSS | 3.x |
| Backend | Go | 1.22+ |
| Database | MySQL | 8.0+ |
| State Management | Pinia | - |
| HTTP Client | Axios | - |

---

## 🤝 Contributing

1. Create a feature branch (`git checkout -b feature/AmazingFeature`)
2. Commit your changes (`git commit -m 'Add AmazingFeature'`)
3. Push to the branch (`git push origin feature/AmazingFeature`)
4. Open a Pull Request

---

## 📄 License

This project is proprietary and confidential.

---

## 📞 Support

For issues and questions, please contact the development team or open an issue on the repository.

---

## 🎯 Roadmap

- [ ] Payment gateway integration
- [ ] Order tracking system
- [ ] User reviews and ratings
- [ ] Wishlist feature
- [ ] Admin analytics dashboard
- [ ] Email notifications
- [ ] Mobile app version

---

**Last Updated:** December 26, 2025  
**Maintained By:** Development Team

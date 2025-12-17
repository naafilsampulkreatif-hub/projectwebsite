# 🚀 Quick Start Guide

## 1️⃣ **Start MySQL**

### Option A: Via Windows Services (Recommended)
```powershell
# Open PowerShell as Admin and run:
net start MySQL80

# Verify it's running:
netstat -ano | findstr :3306
```

### Option B: Via Command Line
```bash
cd "C:\Program Files\MySQL\MySQL Server 8.0\bin"
mysqld.exe
```

### Option C: Via MySQL Installer
- Open MySQL Installer
- MySQL Server → Configure now
- Select "Manual MySQL Server Instance Configuration Wizard"
- Ensure "MySQL Server" service starts automatically

---

## 2️⃣ **Setup Database**

```bash
# Connect to MySQL
mysql -u root -p

# Run these commands in MySQL CLI:
CREATE DATABASE IF NOT EXISTS ecommerce;
USE ecommerce;

# Exit MySQL
exit

# Import schema from file
mysql -u root -p ecommerce < backend/schema.sql

# Verify tables created
mysql -u root -p -e "USE ecommerce; SHOW TABLES;"
```

Expected output:
```
+---------------------------+
| Tables_in_ecommerce       |
+---------------------------+
| admin_activity_log        | ⭐ NEW
| cart_items                |
| categories                |
| comments                  | ⭐ NEW
| customer_info             | ⭐ NEW
| order_items               |
| orders                    |
| products                  |
| users                     |
+---------------------------+
```

---

## 3️⃣ **Setup Backend**

```bash
cd backend

# Create .env file (copy from .env.example)
copy .env.example .env

# Edit .env with your MySQL password:
# DB_USER=root
# DB_PASS=YOUR_PASSWORD_HERE
# DB_HOST=127.0.0.1
# DB_PORT=3306
# DB_NAME=ecommerce
# JWT_SECRET=your_secret_key_change_me

# Start the backend server
go run main.go

# Expected output:
# Database connected successfully!
# Server starting on port 8080
```

---

## 4️⃣ **Setup Frontend**

```bash
cd frontend

# Install dependencies
npm install

# Start dev server
npm run dev

# Expected output:
# VITE v4.x.x  ready in XXX ms
# ➜  Local:   http://localhost:5173/
# ➜  press h to show help
```

---

## 5️⃣ **Test Admin Dashboard**

1. Open browser: **http://localhost:5173**
2. Login dengan:
   - Email: `admin@example.com`
   - Password: `admin123`

3. Navigate ke **Admin Dashboard**:
   - ✅ **Products** - CRUD produk
   - ✅ **Orders** - Manage pesanan
   - ✅ **Users** - 3 tabs (Users, Customers, Comments)
   - ✅ **Settings** - 4 tabs (Profile, Password, Stats, Activity Log)

---

## 6️⃣ **Test Guest Checkout** (Optional)

```
1. Open http://localhost:5173
2. Add products ke cart (tanpa login)
3. Go to Checkout
4. Fill checkout form (nama, email, address, etc.)
5. Submit

Data akan disimpan di:
- orders (order details)
- order_items (products)
- customer_info (checkout form data) ⭐
```

---

## 7️⃣ **Verify Everything is Connected**

### Test 1: Backend API
```powershell
# Generate JWT token first (login dulu)
# Atau gunakan token yang valid dari browser localStorage

$token = "YOUR_JWT_TOKEN_HERE"

# Test endpoint
Invoke-WebRequest -Uri "http://localhost:8080/api/admin/stats" `
  -Headers @{"Authorization"="Bearer $token"}

# Should return JSON with stats
```

### Test 2: Check Activity Log
```powershell
$token = "YOUR_JWT_TOKEN_HERE"

Invoke-WebRequest -Uri "http://localhost:8080/api/admin/activity-log" `
  -Headers @{"Authorization"="Bearer $token"} `
  -Method GET
```

### Test 3: Database Direct Query
```bash
mysql -u root -p -e "USE ecommerce; SELECT * FROM admin_activity_log ORDER BY created_at DESC LIMIT 5;"
```

---

## ⚙️ **Environment Variables (.env)**

```dotenv
# Database Configuration
DB_USER=root
DB_PASS=your_mysql_password          # Change this!
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=ecommerce

# JWT Configuration
JWT_SECRET=your_secret_key_change_me # Change this!

# Server Configuration
PORT=8080
```

---

## 🐛 **Troubleshooting**

### Backend starts but crashes
**Problem:** "exit status 0xc000013a"
**Solution:** MySQL is not running. Run `net start MySQL80`

### Cannot connect to database
**Problem:** "Error connecting to the database"
**Solution:** 
- Check MySQL is running: `netstat -ano | findstr :3306`
- Check credentials in .env file
- Check database exists: `mysql -u root -p -e "SHOW DATABASES;"`

### Frontend cannot reach backend
**Problem:** CORS error in browser console
**Solution:**
- Ensure backend is running on port 8080
- Check CORS config in `backend/main.go` allows `localhost:5173`
- Clear browser cache and restart

### Tables not created
**Problem:** Schema import failed
**Solution:**
- Manually run schema.sql in MySQL Workbench
- Or: `mysql -u root -p ecommerce < backend/schema.sql`
- Verify with: `mysql -u root -p -e "USE ecommerce; SHOW TABLES;"`

---

## ✅ **Final Checklist**

- [ ] MySQL running (port 3306 listening)
- [ ] Database `ecommerce` created
- [ ] Tables imported from schema.sql (9 tables)
- [ ] Backend .env file created with credentials
- [ ] Backend running on port 8080
- [ ] Frontend running on port 5173
- [ ] Can login as admin
- [ ] Can navigate admin dashboard
- [ ] Can see Users, Orders, Products
- [ ] Activity log shows admin actions

**All done? Your admin dashboard is fully operational! 🎉**

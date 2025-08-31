# StokQ Backend API

StokQ adalah aplikasi manajemen stok sederhana untuk pemilik usaha kecil. Backend ini dibangun menggunakan Go (Golang) dengan framework Gin dan database PostgreSQL.

## Fitur

- 🔐 **Autentikasi JWT** - Register dan Login pengguna
- 📦 **Manajemen Produk** - CRUD operations untuk produk
- 📊 **Manajemen Stok** - Stock In dan Stock Out
- 🗄️ **Database PostgreSQL** dengan GORM ORM
- 🛡️ **Middleware Authentication** untuk proteksi endpoint

## Teknologi yang Digunakan

- **Go 1.21+** - Programming language
- **Gin** - Web framework
- **PostgreSQL** - Database
- **GORM** - ORM untuk database operations
- **JWT** - JSON Web Token untuk autentikasi
- **bcrypt** - Password hashing

## Struktur Proyek

```
stokq-backend/
├── config/          # Konfigurasi database
├── controllers/     # Business logic handlers
├── dto/            # Data Transfer Objects
├── initializers/   # Inisialisasi aplikasi
├── middleware/     # Middleware functions
├── models/         # Database models
├── routes/         # Route definitions
├── main.go         # Entry point aplikasi
├── go.mod          # Go module dependencies
└── .env            # Environment variables
```

## Instalasi dan Setup

### 1. Prerequisites

- Go 1.21 atau lebih baru
- PostgreSQL database
- Git

### 2. Clone dan Setup Project

```bash
# Clone repository
git clone <repository-url>
cd stokq-backend

# Install dependencies
go mod tidy
```

### 3. Setup Database

```sql
-- Buat database PostgreSQL
CREATE DATABASE stokq_db;
CREATE USER stokq_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE stokq_db TO stokq_user;
```

### 4. Environment Configuration

Buat file `.env` di root directory:

```env
DB_URL="host=localhost user=stokq_user password=your_password dbname=stokq_db port=5432 sslmode=disable"
JWT_SECRET="your_super_secret_jwt_key_here"
PORT="8080"
```

### 5. Jalankan Aplikasi

```bash
# Development mode
go run main.go

# Build dan jalankan
go build -o stokq-backend.exe .
./stokq-backend.exe
```

Server akan berjalan di `http://localhost:8080`

## API Endpoints

### Health Check
- `GET /health` - Check server status

### Authentication
- `POST /api/v1/auth/register` - Register pengguna baru
- `POST /api/v1/auth/login` - Login pengguna

### Products (Protected - Require Authentication)
- `POST /api/v1/products` - Buat produk baru
- `GET /api/v1/products` - Ambil semua produk
- `GET /api/v1/products/:id` - Ambil produk berdasarkan ID
- `PUT /api/v1/products/:id` - Update produk
- `DELETE /api/v1/products/:id` - Hapus produk

### Stock Management (Protected - Require Authentication)
- `POST /api/v1/stock/in` - Tambah stok produk
- `POST /api/v1/stock/out` - Kurangi stok produk

## Contoh Penggunaan API

### 1. Register User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

### 2. Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### 3. Create Product

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "sku": "PROD001",
    "name": "Laptop Dell",
    "stock": 10,
    "price": 15000000
  }'
```

### 4. Stock In

```bash
curl -X POST http://localhost:8080/api/v1/stock/in \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "product_id": 1,
    "quantity": 5
  }'
```

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);
```

### Products Table
```sql
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    sku VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    stock INTEGER DEFAULT 0,
    price DECIMAL(15,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);
```

## Security

- Password di-hash menggunakan bcrypt
- JWT tokens dengan expiry 7 hari
- Protected routes dengan middleware authentication
- CORS enabled untuk cross-origin requests

## Development

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
go build -ldflags="-s -w" -o stokq-backend .
```

## Contributing

1. Fork repository
2. Buat feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push ke branch (`git push origin feature/amazing-feature`)
5. Buat Pull Request

## License

Distributed under the MIT License. See `LICENSE` file for more information.

## Support

Untuk support atau pertanyaan, silakan buat issue di GitHub repository atau hubungi tim development.

---

**StokQ Backend API v1.0.0** - Built with ❤️ using Go & Gin

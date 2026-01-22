# Chirpy Setup Guide

This guide will walk you through setting up the Chirpy API on your local machine.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Database Setup](#database-setup)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [Testing the API](#testing-the-api)
- [Troubleshooting](#troubleshooting)

---

## Prerequisites

Before you begin, ensure you have the following installed on your system:

### Required Software
- **Go** (version 1.25.5 or higher)
  - Download from [golang.org](https://golang.org/dl/)
  - Verify installation: `go version`

- **PostgreSQL** (version 12 or higher recommended)
  - Download from [postgresql.org](https://www.postgresql.org/download/)
  - Verify installation: `psql --version`

- **Git**
  - Download from [git-scm.com](https://git-scm.com/)
  - Verify installation: `git --version`

### Optional Tools
- **curl** or **Postman** - For testing API endpoints
- **goose** - For database migrations (if you want to modify the schema)
  ```bash
  go install github.com/pressly/goose/v3/cmd/goose@latest
  ```

---

## Installation

### 1. Clone the Repository
```bash
git clone <your-repository-url>
cd chirpy
```

### 2. Install Go Dependencies
```bash
go mod download
```

This will download all required packages:
- `github.com/lib/pq` - PostgreSQL driver
- `github.com/golang-jwt/jwt/v5` - JWT authentication
- `github.com/google/uuid` - UUID generation
- `github.com/alexedwards/argon2id` - Password hashing
- `github.com/joho/godotenv` - Environment variable loading

### 3. Verify Installation
```bash
go build -o chirpy
```

If successful, this creates a `chirpy` executable in your current directory.

---

## Database Setup

### 1. Start PostgreSQL
Make sure your PostgreSQL server is running:

**macOS (Homebrew):**
```bash
brew services start postgresql@14
```

**Linux:**
```bash
sudo systemctl start postgresql
```

**Windows:**
Start PostgreSQL from the Services app or pgAdmin

### 2. Create Database
Open PostgreSQL interactive terminal:
```bash
psql -U postgres
```

Create the database:
```sql
CREATE DATABASE chirpy;
```

Exit psql:
```sql
\q
```

### 3. Run Database Migrations

The project uses **goose** for database migrations. The migration files are located in `sql/schema/`.

#### Using Goose (Recommended)
```bash
# Navigate to sql/schema directory
cd sql/schema

# Run migrations
goose postgres "postgres://your_username:your_password@localhost:5432/chirpy?sslmode=disable" up

# Return to project root
cd ../..
```

#### Manual Migration (Alternative)
If you prefer to run migrations manually:
```bash
psql -U your_username -d chirpy -f sql/schema/001_users.sql
psql -U your_username -d chirpy -f sql/schema/002_chirps.sql
psql -U your_username -d chirpy -f sql/schema/003_passwords.sql
psql -U your_username -d chirpy -f sql/schema/004_refresh.sql
psql -U your_username -d chirpy -f sql/schema/005_red.sql
```

### 4. Verify Database Schema
```bash
psql -U your_username -d chirpy
```

Check tables:
```sql
\dt
```

You should see:
- `users`
- `chirps`
- `refresh_tokens`
- `goose_db_version` (if using goose)

---

## Configuration

### 1. Create Environment File
Copy the example environment file:
```bash
cp .env.example .env
```

### 2. Configure Environment Variables
Edit the `.env` file with your settings:

```bash
# Database Configuration
DB_URL="postgres://your_username:your_password@localhost:5432/chirpy?sslmode=disable"

# Application Environment (dev or prod)
PLATFORM="dev"

# JWT Secret (generate a secure random string)
JWT_SECRET="your-secret-jwt-key-here"

# Polka API Key (for webhook authentication)
POLKA_KEY="your-polka-api-key-here"
```

### 3. Generate Secure JWT Secret
For production use, generate a strong random secret:

**Linux/macOS:**
```bash
openssl rand -base64 64
```

**Or using Go:**
```bash
go run -c 'package main; import ("crypto/rand"; "encoding/base64"; "fmt"); func main() { b := make([]byte, 64); rand.Read(b); fmt.Println(base64.StdEncoding.EncodeToString(b)) }'
```

Copy the output and paste it as your `JWT_SECRET` value.

### Environment Variables Explained

| Variable | Description | Example |
|----------|-------------|---------|
| `DB_URL` | PostgreSQL connection string | `postgres://user:pass@localhost:5432/chirpy?sslmode=disable` |
| `PLATFORM` | Environment mode (`dev` or `prod`) | `dev` |
| `JWT_SECRET` | Secret key for signing JWT tokens | (64-character random string) |
| `POLKA_KEY` | API key for Polka webhook authentication | (API key from Polka) |

**Important Notes:**
- **Never commit `.env` to version control** - It contains sensitive credentials
- The `.gitignore` already excludes `.env` files
- Use different secrets for development and production
- `PLATFORM=dev` enables the `/admin/reset` endpoint (dangerous in production!)

---

## Running the Application

### Development Mode
```bash
go run .
```

Or build and run the executable:
```bash
go build -o chirpy
./chirpy
```

### Expected Output
```
2024/01/20 10:30:00 Serving files from . on port: 8080
```

The server is now running on **http://localhost:8080**

### Production Mode
1. Set `PLATFORM="prod"` in your `.env` file
2. Build the optimized binary:
   ```bash
   go build -ldflags="-s -w" -o chirpy
   ```
3. Run the binary:
   ```bash
   ./chirpy
   ```

**Production Considerations:**
- Use a process manager (systemd, supervisor, pm2)
- Set up proper logging
- Configure firewall rules
- Use HTTPS with a reverse proxy (nginx, caddy)
- Implement rate limiting
- Regular database backups

---

## Testing the API

### 1. Health Check
Verify the server is running:
```bash
curl http://localhost:8080/api/healthz
```

Expected response: `OK`

### 2. Create a User
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"testpassword123"}'
```

### 3. Login
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"testpassword123"}'
```

Save the `token` and `refresh_token` from the response.

### 4. Create a Chirp
```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"body":"My first chirp! This is awesome."}'
```

Replace `YOUR_JWT_TOKEN` with the token from step 3.

### 5. Get All Chirps
```bash
curl http://localhost:8080/api/chirps
```

### Complete Testing Workflow
For a full test suite, check out the [API Documentation](./api.md) which includes examples for all endpoints.

---

## Troubleshooting

### Common Issues

#### Database Connection Failed
**Error:** `Error opening database: connection refused`

**Solutions:**
1. Verify PostgreSQL is running:
   ```bash
   # macOS
   brew services list

   # Linux
   sudo systemctl status postgresql
   ```

2. Check your database credentials in `.env`
3. Verify database exists:
   ```bash
   psql -U postgres -l | grep chirpy
   ```

#### Port Already in Use
**Error:** `bind: address already in use`

**Solutions:**
1. Check what's using port 8080:
   ```bash
   # macOS/Linux
   lsof -i :8080
   ```

2. Kill the process:
   ```bash
   kill -9 <PID>
   ```

3. Or change the port in `main.go` (line 40)

#### JWT Token Invalid
**Error:** `401 Unauthorized`

**Solutions:**
1. Verify you're including the token in the header:
   ```
   Authorization: Bearer <token>
   ```
2. Check token hasn't expired (1-hour lifetime)
3. Use the refresh endpoint to get a new token

#### sqlc Code Generation Issues
If you modify SQL queries and need to regenerate Go code:

```bash
# Install sqlc
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Generate code
sqlc generate
```

#### Migration Errors
**Error:** `goose: no migrations to run`

**Solution:**
Make sure you're in the correct directory and the connection string is correct:
```bash
cd sql/schema
goose postgres "YOUR_DB_URL" status
```

---

## Development Tools

### Hot Reload (Optional)
Install **air** for automatic reloading during development:

```bash
go install github.com/cosmtrek/air@latest
```

Create `.air.toml` configuration:
```toml
root = "."
tmp_dir = "tmp"

[build]
cmd = "go build -o ./tmp/main ."
bin = "tmp/main"
include_ext = ["go"]
exclude_dir = ["tmp"]
```

Run with hot reload:
```bash
air
```

### Database GUI (Optional)
Consider using a PostgreSQL GUI for easier database management:
- **pgAdmin** - Full-featured PostgreSQL tool
- **DBeaver** - Universal database tool
- **Postico** (macOS) - Native PostgreSQL client

---

## Project Structure

```
chirpy/
├── main.go                 # Application entry point
├── chirps.go              # Chirp handlers
├── users.go               # User handlers
├── cleanWords.go          # Profanity filter
├── go.mod                 # Go module definition
├── go.sum                 # Dependency checksums
├── .env                   # Environment variables (not in git)
├── .env.example           # Environment template
├── sqlc.yaml              # sqlc configuration
├── docs/
│   ├── api.md            # API documentation
│   └── setup.md          # This file
├── internal/
│   ├── auth/             # Authentication utilities
│   └── database/         # Generated database code (sqlc)
└── sql/
    ├── schema/           # Database migrations
    └── queries/          # SQL queries for sqlc
```

---

## Next Steps

- Read the [API Documentation](./api.md) to learn about all available endpoints
- Experiment with different API calls
- Try building a frontend application
- Explore the code to understand the implementation
- Consider adding new features (see Future Enhancements in API docs)

---

## Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [PostgreSQL Tutorial](https://www.postgresql.org/docs/current/tutorial.html)
- [JWT Introduction](https://jwt.io/introduction)
- [sqlc Documentation](https://docs.sqlc.dev/)
- [goose Migrations](https://github.com/pressly/goose)
- [Boot.dev](https://boot.dev) - Where this project originated

---

## Getting Help

If you encounter issues:

1. Check this troubleshooting guide
2. Review the Boot.dev curriculum materials
3. Check the Go documentation
4. Search for error messages online
5. Review the code comments and structure

Happy coding! 🚀

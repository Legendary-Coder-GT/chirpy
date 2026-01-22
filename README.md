# Chirpy 🐦

> A Twitter-like social media API built with Go - A Boot.dev guided project

Chirpy is a RESTful API backend that enables users to post short messages called "chirps," manage user accounts, and handle authentication. This project was built as part of the [Boot.dev](https://boot.dev) curriculum to learn backend web development with Go.

## About This Project

This project is a hands-on learning experience focused on building REST APIs in Go. Through building Chirpy, I learned about:

- Building RESTful APIs with Go's standard library
- Implementing JWT-based authentication with refresh tokens
- Secure password hashing using Argon2id
- Working with PostgreSQL and type-safe SQL queries (sqlc)
- Database design and migrations
- API middleware and request handling
- Webhook integration (Polka payment system)

## Features

### Core Functionality
- **User Management**: Registration, login, profile updates
- **Chirps (Posts)**: Create, read, and delete short messages (140 character limit)
- **Authentication**: JWT access tokens with refresh token rotation
- **Premium Memberships**: Upgrade users to "Chirpy Red" via webhooks
- **Profanity Filtering**: Automatic censoring of inappropriate words

### Technical Highlights
- **Secure Authentication**: Argon2id password hashing, JWT tokens with 1-hour expiration
- **Database Integration**: PostgreSQL with sqlc for type-safe queries
- **Clean Architecture**: Separation of concerns with internal packages
- **Webhook Support**: Integration with external payment processors

## Tech Stack

- **Language**: Go 1.25.5
- **Database**: PostgreSQL
- **Authentication**: JWT (JSON Web Tokens)
- **Password Hashing**: Argon2id
- **SQL Code Generation**: sqlc
- **Database Migrations**: goose

### Key Dependencies
- `github.com/lib/pq` - PostgreSQL driver
- `github.com/golang-jwt/jwt/v5` - JWT authentication
- `github.com/google/uuid` - UUID generation
- `github.com/alexedwards/argon2id` - Secure password hashing
- `github.com/joho/godotenv` - Environment configuration

## Quick Start

### Prerequisites
- Go 1.25.5+
- PostgreSQL 12+
- Git

### Installation

1. Clone the repository:
```bash
git clone <your-repo-url>
cd chirpy
```

2. Install dependencies:
```bash
go mod download
```

3. Set up PostgreSQL database:
```bash
createdb chirpy
```

4. Configure environment variables:
```bash
cp .env.example .env
# Edit .env with your database credentials and secrets
```

5. Run database migrations:
```bash
cd sql/schema
goose postgres "your-db-url" up
cd ../..
```

6. Start the server:
```bash
go run .
```

The API will be available at `http://localhost:8080`

For detailed setup instructions, see the [Setup Guide](./docs/setup.md).

## Documentation

- **[API Documentation](./docs/api.md)** - Complete API endpoint reference with examples
- **[Setup Guide](./docs/setup.md)** - Detailed installation and configuration instructions

## API Endpoints Overview

### Public Endpoints
- `GET /api/healthz` - Health check
- `POST /api/users` - User registration
- `POST /api/login` - User authentication
- `GET /api/chirps` - Get all chirps (with filtering)
- `GET /api/chirps/{id}` - Get single chirp

### Authenticated Endpoints
- `PUT /api/users` - Update user profile
- `POST /api/chirps` - Create a new chirp
- `DELETE /api/chirps/{id}` - Delete a chirp
- `POST /api/refresh` - Refresh access token
- `POST /api/revoke` - Revoke refresh token (logout)

### Admin & Webhooks
- `GET /admin/metrics` - View server metrics
- `POST /admin/reset` - Reset server (dev only)
- `POST /api/polka/webhooks` - Premium upgrade webhook

See the [API Documentation](./docs/api.md) for complete details and examples.

## Project Structure

```
chirpy/
├── main.go                 # Entry point and route configuration
├── chirps.go              # Chirp CRUD handlers
├── users.go               # User management handlers
├── cleanWords.go          # Profanity filter utility
├── internal/
│   ├── auth/             # JWT & password utilities
│   └── database/         # Generated database code (sqlc)
├── sql/
│   ├── schema/           # Database migrations
│   └── queries/          # SQL queries for sqlc
└── docs/                 # API and setup documentation
```

## Example Usage

### Create a User
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securepass123"}'
```

### Login
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securepass123"}'
```

### Post a Chirp
```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"body":"Learning Go is awesome!"}'
```

## What I Learned

Building Chirpy taught me valuable skills in backend development:

### Go Programming
- Idiomatic Go code structure and patterns
- Working with Go's standard library for HTTP servers
- Error handling and logging best practices
- Using third-party packages effectively

### API Design
- RESTful API principles and conventions
- Request/response handling with JSON
- Middleware implementation for metrics tracking
- URL routing and path parameters

### Authentication & Security
- Implementing JWT-based authentication
- Secure password storage with Argon2id
- Token refresh strategies and expiration
- API key authentication for webhooks
- Protecting endpoints with middleware

### Database Development
- PostgreSQL database design
- Writing efficient SQL queries
- Using sqlc for type-safe database access
- Database migrations with goose
- Foreign key relationships and cascading deletes

### Software Engineering
- Project structure and organization
- Environment-based configuration
- Version control with Git
- API documentation
- Development vs. production environments

## Future Enhancements

Potential features to add for further learning:
- User profiles with avatars
- Chirp likes and favorites
- Follow/unfollow functionality
- User feeds based on follows
- Hashtags and mentions
- Direct messaging
- Email verification
- Password reset via email
- Rate limiting
- Pagination for chirps
- Full-text search
- Media attachments (images/videos)

## Development

### Running Tests
```bash
go test ./...
```

### Hot Reload (Optional)
Install and use [air](https://github.com/cosmtrek/air) for automatic reloading:
```bash
go install github.com/cosmtrek/air@latest
air
```

### Regenerate Database Code
After modifying SQL queries:
```bash
sqlc generate
```

## Acknowledgments

This project was built as part of the **[Boot.dev](https://boot.dev)** backend development curriculum. Boot.dev provides hands-on, project-based learning for aspiring backend developers.

Special thanks to the Boot.dev community and instructors for creating an engaging learning experience that emphasizes practical, real-world development skills.

## License

This is a personal learning project created for educational purposes.

---

**Built with ❤️ as part of my backend development journey**

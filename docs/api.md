# Chirpy API Documentation

## Table of Contents
- [Authentication Overview](#authentication-overview)
- [Error Handling](#error-handling)
- [Endpoints](#endpoints)
  - [Health Check](#health-check)
  - [User Management](#user-management)
  - [Authentication](#authentication)
  - [Chirps](#chirps)
  - [Admin](#admin)
  - [Webhooks](#webhooks)

---

## Authentication Overview

Chirpy uses JWT (JSON Web Tokens) for authentication with a dual-token system:

### Access Tokens
- **Type**: JWT
- **Duration**: 1 hour
- **Purpose**: Authenticate API requests
- **Format**: Include in `Authorization` header as `Bearer <token>`

### Refresh Tokens
- **Type**: Opaque token (stored in database)
- **Duration**: 60 days
- **Purpose**: Obtain new access tokens without re-login
- **Storage**: PostgreSQL database with revocation support

### Authentication Flow
1. User logs in with email/password
2. Server returns both access token and refresh token
3. Client uses access token for API requests
4. When access token expires, use refresh token to get a new access token
5. Refresh tokens can be revoked (logout)

### Using Authenticated Endpoints
Include the JWT access token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

---

## Error Handling

### Standard Error Response Format
```json
{
  "error": "Error message describing what went wrong"
}
```

### Common HTTP Status Codes
- `200 OK` - Request succeeded
- `201 Created` - Resource created successfully
- `204 No Content` - Request succeeded with no response body
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Missing or invalid authentication
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

---

## Endpoints

### Health Check

#### Get API Health Status
Check if the API is running and healthy.

**Endpoint:** `GET /api/healthz`

**Authentication:** Not required

**Response:**
- Status: `200 OK`
- Content-Type: `text/plain`
- Body: `OK`

**Example:**
```bash
curl http://localhost:8080/api/healthz
```

---

### User Management

#### Register New User
Create a new user account.

**Endpoint:** `POST /api/users`

**Authentication:** Not required

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response:** `201 Created`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2024-01-20T10:30:00Z",
  "updated_at": "2024-01-20T10:30:00Z",
  "email": "user@example.com",
  "is_chirpy_red": false
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securepassword123"}'
```

**Notes:**
- Password is hashed using Argon2id before storage
- Email must be unique

---

#### Update User
Update user email and/or password.

**Endpoint:** `PUT /api/users`

**Authentication:** Required (JWT)

**Request Body:**
```json
{
  "email": "newemail@example.com",
  "password": "newpassword123"
}
```

**Response:** `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2024-01-20T10:30:00Z",
  "updated_at": "2024-01-20T10:35:00Z",
  "email": "newemail@example.com",
  "is_chirpy_red": false
}
```

**Example:**
```bash
curl -X PUT http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{"email":"newemail@example.com","password":"newpassword123"}'
```

---

### Authentication

#### Login
Authenticate user and receive access and refresh tokens.

**Endpoint:** `POST /api/login`

**Authentication:** Not required

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response:** `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2024-01-20T10:30:00Z",
  "updated_at": "2024-01-20T10:30:00Z",
  "email": "user@example.com",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "is_chirpy_red": false
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securepassword123"}'
```

**Notes:**
- `token` is the JWT access token (expires in 1 hour)
- `refresh_token` is used to obtain new access tokens (expires in 60 days)

---

#### Refresh Access Token
Get a new access token using a refresh token.

**Endpoint:** `POST /api/refresh`

**Authentication:** Required (Refresh Token)

**Headers:**
```
Authorization: Bearer <your-refresh-token>
```

**Response:** `200 OK`
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/refresh \
  -H "Authorization: Bearer <your-refresh-token>"
```

**Notes:**
- Returns a new JWT access token
- Refresh token must not be expired or revoked

---

#### Revoke Refresh Token (Logout)
Revoke a refresh token to log out.

**Endpoint:** `POST /api/revoke`

**Authentication:** Required (Refresh Token)

**Headers:**
```
Authorization: Bearer <your-refresh-token>
```

**Response:** `204 No Content`

**Example:**
```bash
curl -X POST http://localhost:8080/api/revoke \
  -H "Authorization: Bearer <your-refresh-token>"
```

**Notes:**
- Marks the refresh token as revoked in the database
- Revoked tokens cannot be used to refresh access tokens

---

### Chirps

#### Create Chirp
Post a new chirp (message).

**Endpoint:** `POST /api/chirps`

**Authentication:** Required (JWT)

**Request Body:**
```json
{
  "body": "This is my first chirp! Learning Go is awesome."
}
```

**Response:** `201 Created`
```json
{
  "id": "650e8400-e29b-41d4-a716-446655440001",
  "created_at": "2024-01-20T10:40:00Z",
  "updated_at": "2024-01-20T10:40:00Z",
  "body": "This is my first chirp! Learning Go is awesome.",
  "user_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{"body":"This is my first chirp! Learning Go is awesome."}'
```

**Notes:**
- Maximum length: 140 characters
- Profanity filtering: Words like "kerfuffle", "sharbert", and "fornax" are replaced with "****"
- User must be authenticated

---

#### Get All Chirps
Retrieve all chirps with optional filtering and sorting.

**Endpoint:** `GET /api/chirps`

**Authentication:** Not required

**Query Parameters:**
- `author_id` (optional) - Filter chirps by user ID
- `sort` (optional) - Sort order: `asc` (default) or `desc`

**Response:** `200 OK`
```json
[
  {
    "id": "650e8400-e29b-41d4-a716-446655440001",
    "created_at": "2024-01-20T10:40:00Z",
    "updated_at": "2024-01-20T10:40:00Z",
    "body": "This is my first chirp!",
    "user_id": "550e8400-e29b-41d4-a716-446655440000"
  },
  {
    "id": "750e8400-e29b-41d4-a716-446655440002",
    "created_at": "2024-01-20T10:45:00Z",
    "updated_at": "2024-01-20T10:45:00Z",
    "body": "Another chirp here!",
    "user_id": "550e8400-e29b-41d4-a716-446655440000"
  }
]
```

**Examples:**
```bash
# Get all chirps
curl http://localhost:8080/api/chirps

# Get chirps by specific author
curl "http://localhost:8080/api/chirps?author_id=550e8400-e29b-41d4-a716-446655440000"

# Get all chirps sorted by newest first
curl "http://localhost:8080/api/chirps?sort=desc"

# Combine filters
curl "http://localhost:8080/api/chirps?author_id=550e8400-e29b-41d4-a716-446655440000&sort=desc"
```

---

#### Get Single Chirp
Retrieve a specific chirp by ID.

**Endpoint:** `GET /api/chirps/{chirpID}`

**Authentication:** Not required

**URL Parameters:**
- `chirpID` - UUID of the chirp

**Response:** `200 OK`
```json
{
  "id": "650e8400-e29b-41d4-a716-446655440001",
  "created_at": "2024-01-20T10:40:00Z",
  "updated_at": "2024-01-20T10:40:00Z",
  "body": "This is my first chirp!",
  "user_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Example:**
```bash
curl http://localhost:8080/api/chirps/650e8400-e29b-41d4-a716-446655440001
```

**Error Responses:**
- `404 Not Found` - Chirp does not exist

---

#### Delete Chirp
Delete a chirp (only the author can delete their own chirps).

**Endpoint:** `DELETE /api/chirps/{chirpID}`

**Authentication:** Required (JWT)

**URL Parameters:**
- `chirpID` - UUID of the chirp to delete

**Response:** `204 No Content`

**Example:**
```bash
curl -X DELETE http://localhost:8080/api/chirps/650e8400-e29b-41d4-a716-446655440001 \
  -H "Authorization: Bearer <your-jwt-token>"
```

**Error Responses:**
- `401 Unauthorized` - Missing or invalid token
- `403 Forbidden` - User is not the author of the chirp
- `404 Not Found` - Chirp does not exist

---

### Admin

#### Get Metrics
View admin metrics including page visit count.

**Endpoint:** `GET /admin/metrics`

**Authentication:** Not required (admin only in production)

**Response:** `200 OK`
- Content-Type: `text/html`
- Returns HTML page showing visit count

**Example:**
```bash
curl http://localhost:8080/admin/metrics
```

**Response:**
```html
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited 42 times!</p>
  </body>
</html>
```

---

#### Reset Server (Development Only)
Reset metrics and delete all users. Only available when `PLATFORM=dev`.

**Endpoint:** `POST /admin/reset`

**Authentication:** Not required

**Platform Requirement:** `PLATFORM=dev`

**Response:** `200 OK`
```
Server Hits Reset to 0 - All users cleared
```

**Example:**
```bash
curl -X POST http://localhost:8080/admin/reset
```

**Error Responses:**
- `403 Forbidden` - Not in development mode

**Notes:**
- This endpoint is dangerous and only available in development
- Deletes all users and their associated chirps (cascade delete)
- Resets the fileserver hit counter to 0

---

### Webhooks

#### Polka Premium Upgrade Webhook
Webhook endpoint for handling premium membership upgrades from Polka payment service.

**Endpoint:** `POST /api/polka/webhooks`

**Authentication:** Required (API Key)

**Headers:**
```
Authorization: ApiKey <your-polka-api-key>
```

**Request Body:**
```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

**Response:** `204 No Content`

**Example:**
```bash
curl -X POST http://localhost:8080/api/polka/webhooks \
  -H "Content-Type: application/json" \
  -H "Authorization: ApiKey f271c81ff7084ee5b99a5091b42d486e" \
  -d '{"event":"user.upgraded","data":{"user_id":"550e8400-e29b-41d4-a716-446655440000"}}'
```

**Events:**
- `user.upgraded` - Marks user as Chirpy Red (premium member)

**Error Responses:**
- `401 Unauthorized` - Invalid or missing API key
- `404 Not Found` - User does not exist

**Notes:**
- This endpoint is meant to be called by the Polka payment service
- API key must match the `POLKA_KEY` environment variable
- Sets `is_chirpy_red` to `true` for the specified user

---

## Database Schema

### Users Table
```sql
- id: UUID (Primary Key)
- created_at: TIMESTAMP
- updated_at: TIMESTAMP
- email: TEXT (Unique)
- hashed_password: TEXT
- is_chirpy_red: BOOLEAN (Premium status)
```

### Chirps Table
```sql
- id: UUID (Primary Key)
- created_at: TIMESTAMP
- updated_at: TIMESTAMP
- body: TEXT
- user_id: UUID (Foreign Key -> users.id, CASCADE DELETE)
```

### Refresh Tokens Table
```sql
- token: TEXT (Primary Key)
- created_at: TIMESTAMP
- updated_at: TIMESTAMP
- user_id: UUID (Foreign Key -> users.id, CASCADE DELETE)
- expires_at: TIMESTAMP
- revoked_at: TIMESTAMP (Nullable)
```

---

## Rate Limiting & Security

### Password Security
- Passwords are hashed using **Argon2id** algorithm
- Never store plaintext passwords
- Argon2id is memory-hard and resistant to GPU attacks

### Token Expiration
- **Access tokens**: 1 hour
- **Refresh tokens**: 60 days

### Profanity Filter
The following words are automatically censored (replaced with `****`):
- kerfuffle
- sharbert
- fornax

### CORS
Currently not implemented. Add CORS middleware if needed for frontend applications.

---

## Pagination

Pagination is not currently implemented. All chirps are returned in a single response. For production use, consider implementing:
- Limit/offset pagination
- Cursor-based pagination for better performance

---

## Future Enhancements

Potential improvements for this API:
- User profile pictures
- Chirp likes/favorites
- Follow/unfollow users
- User feeds based on followed accounts
- Media attachments (images, videos)
- Hashtags and mentions
- Direct messages
- Email verification
- Password reset functionality
- Rate limiting
- Pagination
- Search functionality

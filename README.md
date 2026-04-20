# go-payment-api

A RESTful payment API built with Go. It supports user registration and authentication, transferring funds between wallets, querying wallet balances, and listing transactions — backed by a PostgreSQL database.

## Features

- User registration (sign-up) with automatic wallet provisioning (starting balance: 500.00)
- User authentication (sign-in) returning a JWT access token
- JWT-based route protection for wallet, transfer, transaction, and user account endpoints
- Transfer funds between user wallets
- Retrieve a user's wallet balance
- List all recorded transactions
- List transactions for a specific user
- Update authenticated user profile (name and avatar URL)
- Change authenticated user password
- PostgreSQL database with auto-initialised schema
- Built-in database migrations (e.g., `avatar_url` on users)
- HTTP request logging and CORS middleware
- Environment-based configuration via `.env`
- Clean layered architecture: handlers → services → repositories
- Standardised JSON response envelope (`data` / `error`)

## Project Structure

```
go-payment-api/
├── main.go                        # Entry point, connects to DB and starts the HTTP server
├── go.mod
├── go.sum
└── internal/
    ├── database/
    │   ├── database.go            # PostgreSQL connection (reads DATABASE_PUBLIC_URL)
    │   └── schema.go              # Auto-creates users, wallets, transactions tables
    ├── handlers/
    │   ├── user_handler.go        # HTTP handlers for sign-up and sign-in
    │   ├── transaction_handler.go # HTTP handlers for fund transfers and listing transactions
    │   └── wallet_handler.go      # HTTP handler for wallet lookup
    ├── middleware/
    │   ├── auth.go                # JWT auth middleware (Bearer token validation)
    │   ├── logging.go             # HTTP request logging middleware
    │   └── cors.go                # CORS middleware
    ├── models/
    │   ├── user_model.go          # User structs and request/response types
    │   ├── wallet_model.go        # Wallet structs and response types
    │   └── transaction_model.go   # Transaction structs and request type
    ├── repository/
    │   ├── user_repository.go     # PostgreSQL user store
    │   ├── wallet_repository.go   # PostgreSQL wallet store
    │   └── transaction_repository.go # PostgreSQL transaction store
    ├── services/
    │   ├── user_services.go       # Sign-up/sign-in business logic (bcrypt password hashing)
    │   ├── transaction_services.go # Transfer and listing business logic
    │   └── wallet_services.go     # Wallet lookup business logic
    └── utils/
        ├── jwt.go                 # JWT generation and validation helpers
        └── response.go            # Shared JSON response helper
```

## Prerequisites

- [Go](https://golang.org/dl/) 1.21+
- [PostgreSQL](https://www.postgresql.org/) 13+

## Getting Started

1. **Clone the repository**

   ```bash
   git clone https://github.com/Dav16Akin/go-payment-api.git
   cd go-payment-api
   ```

2. **Configure the environment**

   Create a `.env` file in the project root:

   ```env
   DATABASE_PUBLIC_URL=postgres://<user>:<password>@<host>:<port>/<dbname>?sslmode=disable
   PORT=8000
   JWT_SECRET=your-strong-secret
   ```

   > The application will automatically create the `users`, `wallets`, and `transactions` tables on startup if they do not already exist.

3. **Install dependencies**

   ```bash
   go mod tidy
   ```

4. **Run the server**

   ```bash
   go run main.go
   ```

   The server starts on **port 8000** by default (override with the `PORT` env var). All incoming requests are logged to stdout in the format:

   ```
   <METHOD> <PATH> <STATUS_CODE> <DURATION>
   ```

## Response Envelope

All endpoints return a consistent JSON envelope:

```json
{
  "data": <payload or null>,
  "error": <"error message" or null>
}
```

## API Endpoints

### Authentication

- `POST /sign-up` and `POST /sign-in` are public.
- The following endpoints are protected and require `Authorization: Bearer <jwt-token>`:
  - `POST /transfer`
  - `GET /transactions`
  - `GET /transactions/user`
  - `GET /wallet`
  - `PATCH /users/profile`
  - `PATCH /users/password`

### Sign Up

Registers a new user and provisions a wallet with a starting balance of 500.00.

**`POST /sign-up`**

**Request body:**
```json
{
  "name": "Alice",
  "email": "alice@example.com",
  "password": "secret"
}
```

**Success response (`201 Created`):**
```json
{
  "data": {
    "ID": "550e8400-e29b-41d4-a716-446655440000",
    "Name": "Alice",
    "Email": "alice@example.com"
  },
  "error": null
}
```

**Error responses:**
- `400 Bad Request` – missing name, email, or password; email already registered
- `405 Method Not Allowed` – non-POST request

---

### Sign In

Authenticates a user and returns a JWT access token.

**`POST /sign-in`**

**Request body:**
```json
{
  "email": "alice@example.com",
  "password": "secret"
}
```

**Success response (`200 OK`):**
```json
{
  "data": {
    "token": "<jwt-token>",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Alice",
      "email": "alice@example.com",
      "avatar_url": "https://example.com/avatar.png"
    }
  },
  "error": null
}
```

**Error responses:**
- `400 Bad Request` – invalid request body
- `401 Unauthorized` – user not found or incorrect password
- `405 Method Not Allowed` – non-POST request

---

### Transfer Funds

Transfers an amount from one wallet to another.

**`POST /transfer`**

**Request body:**
```json
{
  "sender_id": "<sender-wallet-id>",
  "receiver_id": "<receiver-wallet-id>",
  "amount": 100.00
}
```

> **Note:** `sender_id` and `receiver_id` are wallet IDs (which match the user ID assigned at sign-up).

**Success response (`201 Created`):**
```json
{
  "data": { "message": "transfer successful" },
  "error": null
}
```

**Error responses:**
- `400 Bad Request` – sender/receiver wallet not found, same sender and receiver, amount ≤ 0, or insufficient funds
- `401 Unauthorized` – missing/invalid JWT
- `405 Method Not Allowed` – non-POST request

---

### Get Wallet

Returns the wallet balance for a given user.

**`GET /wallet?user_id=<user_id>`**

**Query parameter:** `user_id` – the ID of the user whose wallet to retrieve.

**Success response (`200 OK`):**
```json
{
  "data": {
    "UserID": "550e8400-e29b-41d4-a716-446655440000",
    "Balance": 500
  },
  "error": null
}
```

**Error responses:**
- `400 Bad Request` – missing `user_id`
- `404 Not Found` – wallet not found for the given user
- `401 Unauthorized` – missing/invalid JWT
- `405 Method Not Allowed` – non-GET request

---

### List All Transactions

Returns all recorded transactions.

**`GET /transactions`**

**Success response (`200 OK`):**
```json
{
  "data": [
    {
      "ID": "550e8400-e29b-41d4-a716-446655440000",
      "SenderID": "wallet1",
      "ReceiverID": "wallet2",
      "Amount": 100,
      "Status": "completed",
      "CreatedAt": "2024-01-15T10:30:00Z"
    }
  ],
  "error": null
}
```

**Error responses:**
- `401 Unauthorized` – missing/invalid JWT
- `405 Method Not Allowed` – non-GET request

---

### List Transactions by User

Returns all transactions where the given user is the sender or receiver.

**`GET /transactions/user?user_id=<user_id>`**

**Query parameter:** `user_id` – the ID of the user.

**Success response (`200 OK`):**
```json
{
  "data": [
    {
      "ID": "550e8400-e29b-41d4-a716-446655440000",
      "SenderID": "wallet1",
      "ReceiverID": "wallet2",
      "Amount": 100,
      "Status": "completed",
      "CreatedAt": "2024-01-15T10:30:00Z"
    }
  ],
  "error": null
}
```

**Error responses:**
- `400 Bad Request` – missing `user_id`
- `404 Not Found` – no transactions found for the given user
- `401 Unauthorized` – missing/invalid JWT
- `405 Method Not Allowed` – non-GET request

> **Note:** The `CreatedAt` field is a UTC timestamp set automatically by the database.

---

### Update Profile

Updates the authenticated user's profile details.

**`PATCH /users/profile`**

**Request body (all fields optional):**
```json
{
  "name": "Alice Updated",
  "avatar_url": "https://example.com/alice.png"
}
```

**Success response (`201 Created`):**
```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Alice Updated",
    "avatar_url": "https://example.com/alice.png"
  },
  "error": null
}
```

**Error responses:**
- `400 Bad Request` – invalid request body or update failed
- `401 Unauthorized` – missing/invalid JWT
- `405 Method Not Allowed` – non-PATCH request

---

### Change Password

Changes the authenticated user's password.

**`PATCH /users/password`**

**Request body:**
```json
{
  "old_password": "secret",
  "new_password": "new-secret"
}
```

**Success response (`200 OK`):**
```json
{
  "data": "password changed succesfully",
  "error": null
}
```

**Error responses:**
- `400 Bad Request` – invalid request body, wrong old password, or invalid new password
- `401 Unauthorized` – missing/invalid JWT
- `405 Method Not Allowed` – non-PATCH request

## Example Usage

```bash
# Register a new user
curl -X POST http://localhost:8000/sign-up \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice", "email": "alice@example.com", "password": "secret"}'

# Sign in
curl -X POST http://localhost:8000/sign-in \
  -H "Content-Type: application/json" \
  -d '{"email": "alice@example.com", "password": "secret"}'

# Use your JWT on protected routes
TOKEN="<jwt-token>"

# Transfer funds (using wallet/user IDs from sign-up)
curl -X POST http://localhost:8000/transfer \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"sender_id": "<alice-id>", "receiver_id": "<bob-id>", "amount": 200}'

# Get a user's wallet balance
curl "http://localhost:8000/wallet?user_id=<alice-id>" \
  -H "Authorization: Bearer $TOKEN"

# List all transactions
curl http://localhost:8000/transactions \
  -H "Authorization: Bearer $TOKEN"

# List transactions for a specific user
curl "http://localhost:8000/transactions/user?user_id=<alice-id>" \
  -H "Authorization: Bearer $TOKEN"

# Update profile
curl -X PATCH http://localhost:8000/users/profile \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Updated","avatar_url":"https://example.com/alice.png"}'

# Change password
curl -X PATCH http://localhost:8000/users/password \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"old_password":"secret","new_password":"new-secret"}'
```

## Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| [github.com/google/uuid](https://github.com/google/uuid) | v1.6.0 | UUID generation for user/transaction IDs |
| [github.com/golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt) | v5.3.1 | JWT creation and validation for authenticated routes |
| [github.com/lib/pq](https://github.com/lib/pq) | v1.12.2 | PostgreSQL driver for `database/sql` |
| [github.com/joho/godotenv](https://github.com/joho/godotenv) | v1.5.1 | Load environment variables from `.env` |
| [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto) | v0.49.0 | bcrypt password hashing |

## License

This project is open source. See [LICENSE](LICENSE) for details.

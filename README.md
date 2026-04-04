# go-payment-api

A simple RESTful payment API built with Go. It supports creating users (with automatic wallet creation), transferring funds between wallets, querying wallet balances, and listing all transactions — backed by a PostgreSQL database.

## Features

- Create users with automatically provisioned wallets
- Transfer funds between user wallets
- Retrieve a user's wallet balance
- List all recorded transactions
- PostgreSQL database with auto-initialised schema
- HTTP request logging middleware
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
    │   ├── database.go            # PostgreSQL connection (reads DATABASE_URL)
    │   └── schema.go              # Auto-creates users, wallets, transactions tables
    ├── handlers/
    │   ├── user_handler.go        # HTTP handler for user creation
    │   ├── transaction_handler.go # HTTP handlers for fund transfers and listing transactions
    │   └── wallet_handler.go      # HTTP handler for wallet lookup
    ├── middleware/
    │   └── logging.go             # HTTP request logging middleware
    ├── models/
    │   ├── user_model.go          # User structs and request/response types
    │   ├── wallet_model.go        # Wallet structs and response types
    │   └── transaction_model.go   # Transaction structs and request type
    ├── repository/
    │   ├── user_repository.go     # PostgreSQL user store
    │   ├── wallet_repository.go   # PostgreSQL wallet store
    │   └── transaction_repository.go # PostgreSQL transaction store
    ├── services/
    │   ├── user_services.go       # User creation business logic
    │   ├── transaction_services.go # Transfer and listing business logic
    │   └── wallet_services.go     # Wallet lookup business logic
    └── utils/
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

2. **Configure the database**

   Create a `.env` file in the project root with your PostgreSQL connection string:

   ```env
   DATABASE_URL=postgres://<user>:<password>@<host>:<port>/<dbname>?sslmode=disable
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

   The server starts on **port 8000**. All incoming requests are logged to stdout in the format:

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

### Create User

Creates a new user and provisions an empty wallet for them.

**`POST /user`**

**Request body:**
```json
{
  "name": "Alice",
  "email": "alice@example.com"
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
- `400 Bad Request` – missing name/email, or email already exists
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

> **Note:** Use the wallet ID (not the user ID). For seed data, the wallet IDs are `wallet1` and `wallet2`.

**Success response (`201 Created`):**
```json
{
  "data": { "message": "transfer successful" },
  "error": null
}
```

**Error responses:**
- `400 Bad Request` – sender/receiver wallet not found, same sender and receiver, amount ≤ 0, or insufficient funds
- `405 Method Not Allowed` – non-POST request

---

### Get Wallet

Returns the wallet balance for a given user.

**`GET /wallet/{user_id}`**

**URL parameter:** `user_id` – the ID of the user whose wallet to retrieve.

**Success response (`200 OK`):**
```json
{
  "data": {
    "UserID": "user1",
    "Balance": 1000
  },
  "error": null
}
```

**Error responses:**
- `400 Bad Request` – missing `user_id`
- `404 Not Found` – wallet not found for the given user
- `405 Method Not Allowed` – non-GET request

---

### List Transactions

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
- `405 Method Not Allowed` – non-GET request

> **Note:** The `created_at` field is a UTC timestamp set automatically by the database.

## Example Usage

```bash
# Create a user
curl -X POST http://localhost:8000/user \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice", "email": "alice@example.com"}'

# Transfer funds between seed wallets (using wallet IDs)
curl -X POST http://localhost:8000/transfer \
  -H "Content-Type: application/json" \
  -d '{"sender_id": "wallet1", "receiver_id": "wallet2", "amount": 200}'

# Get a user's wallet balance
curl http://localhost:8000/wallet/user1

# List all transactions
curl http://localhost:8000/transactions
```

## Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| [github.com/google/uuid](https://github.com/google/uuid) | v1.6.0 | UUID generation for user/transaction IDs |
| [github.com/lib/pq](https://github.com/lib/pq) | v1.12.2 | PostgreSQL driver for `database/sql` |
| [github.com/joho/godotenv](https://github.com/joho/godotenv) | v1.5.1 | Load environment variables from `.env` |

## License

This project is open source. See [LICENSE](LICENSE) for details.

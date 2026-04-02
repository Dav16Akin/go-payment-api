# go-payment-api

A simple RESTful payment API built with Go. It supports creating users (with automatic wallet creation) and transferring funds between wallets using in-memory storage.

## Features

- Create users with automatically provisioned wallets
- Transfer funds between user wallets
- In-memory data store (no external database required)
- Clean layered architecture: handlers → services → repositories

## Project Structure

```
go-payment-api/
├── main.go                        # Entry point, seeds data and starts the HTTP server
├── go.mod
├── go.sum
└── internal/
    ├── handlers/
    │   ├── user_handler.go        # HTTP handler for user creation
    │   └── transaction_handler.go # HTTP handler for fund transfers
    ├── models/
    │   ├── user_model.go          # User structs and request/response types
    │   ├── wallet_model.go        # Wallet struct
    │   └── transaction_model.go   # Transaction structs and request type
    ├── repository/
    │   ├── user_repository.go     # In-memory user store
    │   └── wallet_repository.go   # In-memory wallet store
    └── services/
        ├── user_services.go       # User creation business logic
        └── transaction_services.go # Transfer business logic
```

## Prerequisites

- [Go](https://golang.org/dl/) 1.21+

## Getting Started

1. **Clone the repository**

   ```bash
   git clone https://github.com/Dav16Akin/go-payment-api.git
   cd go-payment-api
   ```

2. **Install dependencies**

   ```bash
   go mod tidy
   ```

3. **Run the server**

   ```bash
   go run main.go
   ```

   The server starts on **port 8000**. On startup, two seed users (`David` and `John`) and their wallets are created automatically.

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
  "ID": "550e8400-e29b-41d4-a716-446655440000",
  "Name": "Alice",
  "Email": "alice@example.com"
}
```

**Error responses:**
- `400 Bad Request` – missing name/email, or email already exists
- `405 Method Not Allowed` – non-POST request

---

### Transfer Funds

Transfers an amount from one user's wallet to another.

**`POST /transfer`**

**Request body:**
```json
{
  "sender_id": "<sender-wallet-id>",
  "receiver_id": "<receiver-wallet-id>",
  "amount": 100.00
}
```

> **Note:** The wallet ID equals the user ID. Use the `ID` returned from `POST /user`, or the seed IDs `user1` / `user2` for testing.

**Success response (`200 OK`):**
```json
{
  "message": "transfer successful"
}
```

**Error responses:**
- `400 Bad Request` – sender/receiver wallet not found, amount ≤ 0, or insufficient funds
- `405 Method Not Allowed` – non-POST request

## Example Usage

```bash
# Create a user
curl -X POST http://localhost:8000/user \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice", "email": "alice@example.com"}'

# Transfer funds between seed wallets
curl -X POST http://localhost:8000/transfer \
  -H "Content-Type: application/json" \
  -d '{"sender_id": "user1", "receiver_id": "user2", "amount": 200}'
```

## Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| [github.com/google/uuid](https://github.com/google/uuid) | v1.6.0 | UUID generation for user/transaction IDs |

## License

This project is open source. See [LICENSE](LICENSE) for details.

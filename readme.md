# go-kanji Backend

A modular Go backend providing authentication, user management, user activity tracking, and rabbitmq wiring for the “go-kanji” dashboard. The project follows clean-architecture principles with feature-first packaging.

## 🧱 Project Structure

```
├── app/                  # Feature registration & dependency wiring
├── core/                 # Core contracts, configuration helpers
├── env/                  # YAML env/config map files
├── infra/                # Infrastructure helpers (DB, MQ, logging)
├── shared/               # Shared utilities (responses, identity, etc.)
├── src/
│   ├── auth/             # OAuth2 Google auth flow
│   ├── bunpo/            # Placeholder Bunpo API domain
│   ├── logs/             # User activity logging
│   └── users/            # User repository, services & delivery
├── go.mod
├── main.go
└── README.md
```

## 🧪 Tech Highlights

- Go 1.24
- Gin HTTP framework
- PostgreSQL (via `database/sql`)
- RabbitMQ connection helper (`github.com/rabbitmq/amqp091-go`)
- Hashid-based user reference encoder to avoid exposing raw IDs
- Clean/hexagonal architecture vibes with modular domains

## 🚀 Getting Started

1. **Install dependencies**
   ```bash
   go mod tidy
   ```
2. **Configure environment**
   - Copy `.env` and update DB / RabbitMQ credentials and salts.
   - Ensure required tables exist (migrations not handled automatically).
3. **Run the API**
   ```bash
   go run main.go
   ```

## 📡 Key Endpoints

| Method | Endpoint                    | Description                               |
|--------|-----------------------------|-------------------------------------------|
| GET    | `/auth/google/login`        | Initiates Google OAuth login flow          |
| GET    | `/auth/google/callback`     | Callback handler for Google OAuth          |
| POST   | `/auth/logout`              | Records logout activity                    |
| GET    | `/api/users`                | List masked user accounts                  |
| GET    | `/api/users/logs`           | Paginated activity logs (optional filter)  |
| GET    | `/api/users/:ref/logs`      | Logs scoped to a specific user reference   |
| GET    | `/bunpo/test`               | Bunpo domain test endpoint                 |

## 🧩 Feature Notes

- **User Directory**: Emails are masked and IDs are encoded to references using hashids to avoid exposing raw database IDs.
- **User Activity**: Activity logs can be filtered globally or per user reference. Schema validation will warn if required tables/indexes are missing.
- **Bunpo Domain**: Stubbed service exposing `/bunpo/test` returning “endpoint success”. Extend here for future feature logic.
- **RabbitMQ**: Connection helper available via `infra/mq`; currently the app ensures connectivity at startup.

## 🛠 Tooling

- Formatting: `gofmt`
- Tests: `go test ./...`
- Task running: use Makefile or scripts as needed (not included).

## 🧾 License

MIT (see LICENSE)

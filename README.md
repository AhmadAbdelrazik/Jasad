# Jasad - Fitness Tracker API

**Jasad** is a Go-based fitness tracker API designed to help users log and
manage workouts, track exercises, and authenticate via Google OAuth. It
leverages PostgreSQL for persistent storage and Redis for session management.

---

## 📁 Project Structure

```

.
├── Makefile                 # Build and run commands
├── cmd                      # Entrypoint to the application
│   └── main.go
├── go.mod / go.sum          # Go module dependencies
├── internal                 # Application and domain logic
│   ├── application          # Handlers, middleware, routes
│   └── model                # Data models, Redis, sessions
├── migrations               # SQL migration files
├── pkg                      # Shared utilities (config, validation)
└── tmp                      # Temporary files (e.g., logs)

````

---

## 🚀 Features

- ✅ Google OAuth login
- ✅ User authentication with session handling via Redis
- ✅ CRUD operations for exercises and workouts
- ✅ Role-based access control
- ✅ PostgreSQL-backed persistence
- ✅ RESTful API with route protection middleware
- ✅ Basic rate limiting

---

## 📦 Tech Stack

- **Language:** Go
- **Database:** PostgreSQL
- **Session Management:** Redis
- **Authentication:** Google OAuth 2.0
- **Routing:** Native `net/http`

---

## 🔌 Environment Variables

Configuration is managed via a `Config` struct with environment variables:

```go
type Config struct {
    DSN                string  `env:"JASAD_DB_DSN"`
    Origin             string  `env:"ORIGIN"`
    Port               int     `env:"PORT"`
    GoogleClientID     string  `env:"GOOGLE_CLIENT_ID"`
    GoogleClientSecret string  `env:"GOOGLE_CLIENT_SECRET"`
    LimiterEnable      bool    `env:"LIMITER_ENABLED" envdefault:"true"`
    LimiterRPS         float64 `env:"LIMITER_RPS" envdefault:"2"`
    LimiterBurst       int     `env:"LIMITER_BURST" envdefault:"4"`
}
````

---

## 🔑 API Endpoints

### 🧠 Authentication

* `GET /google_login` — Initiate OAuth login
* `GET /google_callback` — Handle OAuth callback

### 🏋️ Exercises

* `POST /v1/exercises` — Create an exercise
* `GET /v1/exercises` — Search exercises
* `GET /v1/exercises/{id}` — Get exercise by ID
* `PATCH /v1/exercises/{id}` — Update exercise
* `DELETE /v1/exercises/{id}` — Delete exercise

### 👤 Users

* `GET /v1/users` — Get all users
* `GET /v1/users/{id}` — Get user by ID

### 🏃 Workouts

* `POST /v1/workouts` — Create workout
* `GET /v1/workouts` — List workouts
* `GET /v1/workouts/{id}` — Get workout by ID
* `PATCH /v1/workouts/{id}` — Update workout
* `DELETE /v1/workouts/{id}` — Delete workout

---

## 🛠️ Makefile Commands

```bash
make run/http             # Start the HTTP server
make db/psql              # Connect to the database via psql
make db/migrations/new    # Create new migration files
make db/migrations/up     # Run migrations (requires confirmation)
make db/migrations/down   # Rollback migrations (requires confirmation)
```

Use `make confirm` to acknowledge critical operations interactively.

---

## 🔧 Migrations

Migrations are managed using `migrate` CLI. Files are located in the `migrations/` directory.

Example:

```bash
make db/migrations/new name=create_exercise_table
make db/migrations/up
```

---

## ✅ Getting Started

1. **Clone the repo**

   ```bash
   git clone https://github.com/your-username/jasad.git
   cd jasad
   ```

2. **Set environment variables**

   Create a `.env` file with values for:

   * `JASAD_DB_DSN`
   * `ORIGIN`
   * `PORT`
   * `GOOGLE_CLIENT_ID`
   * `GOOGLE_CLIENT_SECRET`

3. **Run database migrations**

   ```bash
   make db/migrations/up
   ```

4. **Start the server**

   ```bash
   make run/http
   ```

---

## 🧪 Testing

> Test coverage and test files are not included in this version, but the project is structured to accommodate testing at the handler and model layers.

---

## 🙋 Contributions

Issues, pull requests, and feature suggestions are welcome!

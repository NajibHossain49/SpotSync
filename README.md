# 🚗 SpotSync – Smart Parking & EV Charging Reservation API

A centralized backend platform for airports and malls to manage parking zones and the high-demand reservation of limited EV charging spots. Built with **Go + Echo + GORM + PostgreSQL**, following **strict Clean Architecture** with concurrency-safe reservations.

- **Live URL:** `https://spotsync-api.onrender.com` _(replace with your deployed URL)_
- **GitHub:** `https://github.com/yourusername/spotsync-api` _(replace with your repo)_

---

## ✨ Features

- **JWT Authentication** with bcrypt password hashing (cost 12).
- **Role-based access control** — `driver` and `admin` roles enforced via middleware.
- **Parking zone management** (admin) with three types: `general`, `ev_charging`, `covered`.
- **Dynamic availability** — `available_spots` calculated live from active reservations.
- **Concurrency-safe reservations** — row-level locking (`SELECT ... FOR UPDATE`) inside a DB transaction prevents over-capacity (the "EV Spot Bottleneck").
- **Ownership rules** — drivers can only cancel their own reservations.
- **Centralized error handling** — no raw GORM errors leak to clients.
- **Request validation** with `go-playground/validator`.
- **Connection pooling** configured for production.

---

## 🛠️ Tech Stack

| Technology | Purpose |
| --- | --- |
| Go 1.22 | Language |
| Echo v4 | Web framework |
| GORM + PostgreSQL driver | ORM / database access |
| go-playground/validator/v10 | Struct validation |
| golang-jwt/jwt/v5 | JWT generation & verification |
| golang.org/x/crypto/bcrypt | Password hashing |
| NeonDB / Supabase | Managed PostgreSQL |

---

## 🏛️ Architecture

Strict separation of concerns. **Handlers never touch the database directly.** Data flows in one direction, and dependencies are wired manually in `main.go`.

```
HTTP Request
    │
    ▼
┌─────────────┐   Binds & validates DTOs, reads JWT claims, returns JSON
│   Handler   │
└─────────────┘
    │ calls
    ▼
┌─────────────┐   Business logic: hashing, JWT, capacity rules, ownership checks
│   Service   │
└─────────────┘
    │ calls
    ▼
┌─────────────┐   All GORM operations: CRUD, transactions, row locks
│ Repository  │
└─────────────┘
    │ uses
    ▼
┌─────────────┐   GORM structs = database tables
│   Models    │
└─────────────┘
```

| Layer | Directory | Responsibility |
| --- | --- | --- |
| DTO | `dto/` | Request payloads & response shapes. GORM models are never exposed directly. |
| Handler | `handler/` | HTTP layer. Bind/validate, extract JWT claims, call service, return JSON. |
| Service | `service/` | Business logic and rules (capacity, ownership, hashing, JWT). |
| Repository | `repository/` | Data access. All GORM CRUD, transactions, row locks. |
| Models | `models/` | GORM structs for `users`, `parking_zones`, `reservations`. |
| Middleware | `middleware/` | JWT verification + role authorization. |
| Utils | `utils/` | JWT helpers, custom errors, validator. |

### 🔒 Solving the Race Condition

When two drivers grab the last EV spot simultaneously, a naive read-then-write lets both succeed. We prevent this in `repository/reservation_repository.go`:

1. Open a **GORM transaction**.
2. Lock the zone row: `tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&zone, zoneID)` — this issues `SELECT ... FOR UPDATE`, forcing any concurrent request on the same zone to **wait**.
3. Count active reservations safely inside the lock.
4. Reject with `409 Conflict` if full, otherwise create the reservation.
5. Commit — the lock releases and the next waiting request sees the updated count.

---

## 🚀 Local Setup

### Prerequisites
- Go 1.22+
- A PostgreSQL database (NeonDB / Supabase / local)

### Steps

```bash
# 1. Clone
git clone https://github.com/yourusername/spotsync-api.git
cd spotsync-api

# 2. Create your .env from the example
cp .env.example .env
# then edit .env with your real DATABASE_URL and JWT_SECRET

# 3. Install dependencies
go mod tidy

# 4. Run
go run main.go
# or with hot reload:
air
```

The server starts on `http://localhost:8080`. Tables are auto-migrated on startup.

### Required `.env` variables

| Variable | Description |
| --- | --- |
| `DATABASE_URL` | PostgreSQL connection string (`postgres://user:pass@host:5432/db?sslmode=require`) |
| `JWT_SECRET` | Long random secret for signing tokens |
| `JWT_EXPIRY_HOURS` | Token lifetime in hours (default `24`) |
| `PORT` | Server port (default `8080`; auto-set by Render/Railway) |

---

## 🌐 API Endpoints

Base path: `/api/v1`

### Auth
| Method | Path | Access | Description |
| --- | --- | --- | --- |
| POST | `/auth/register` | Public | Register a user |
| POST | `/auth/login` | Public | Log in, returns JWT |

### Parking Zones
| Method | Path | Access | Description |
| --- | --- | --- | --- |
| GET | `/zones` | Public | List zones (with `available_spots`) |
| GET | `/zones/:id` | Public | Get one zone |
| POST | `/zones` | Admin | Create a zone |
| PUT | `/zones/:id` | Admin | Update a zone |
| DELETE | `/zones/:id` | Admin | Delete a zone |

### Reservations
| Method | Path | Access | Description |
| --- | --- | --- | --- |
| POST | `/reservations` | Authenticated | Reserve a spot (concurrency-safe) |
| GET | `/reservations/my-reservations` | Authenticated | List own reservations |
| DELETE | `/reservations/:id` | Owner / Admin | Cancel a reservation |
| GET | `/reservations` | Admin | List all reservations |

Protected routes require the header: `Authorization: Bearer <token>`.

### Response format

```json
// Success
{ "success": true, "message": "...", "data": { } }

// Error
{ "success": false, "message": "...", "errors": { } }
```

---

## ☁️ Deployment (Render example)

1. Push the repo to GitHub.
2. On Render → **New → Web Service** → connect the repo.
3. Environment: **Docker** (uses the included `Dockerfile`) or **Go** with build `go build -o app .` and start `./app`.
4. Add environment variables: `DATABASE_URL`, `JWT_SECRET`, `JWT_EXPIRY_HOURS`.
5. Use NeonDB/Supabase for the PostgreSQL `DATABASE_URL`.

CORS is enabled and env vars are read from the system in production.

---

## 📁 Project Layout

```
spotsync-api/
├── main.go               # Entry point + dependency injection + routes
├── config/               # Env config loader
├── database/             # DB connection, pooling, auto-migrate
├── models/               # GORM models
├── dto/                  # Request/response structs + validation tags
├── repository/           # GORM data access (incl. locking transaction)
├── service/              # Business logic
├── handler/              # HTTP handlers
├── middleware/           # JWT auth + role guard
├── utils/                # JWT, errors, validator
├── Dockerfile
└── .env.example
```

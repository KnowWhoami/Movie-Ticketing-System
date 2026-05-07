# Movie Ticketing — REST API

A Go REST API for booking movie tickets. Users can register as a regular customer or a theatre owner. Theatre owners manage cinemas, screens, and shows. Regular users browse shows and book seats.

## Tech stack

- **Go** — standard library HTTP server
- **MySQL** — primary database (via GORM)
- **JWT** — authentication (HS256, 24-hour tokens)
- **bcrypt** — password hashing
- **In-memory cache** — per-service caching layer (go-cache)
- **Docker / Docker Compose** — local setup

---

## Architecture

```
cmd/
  server/        → routes, HTTP adapter, server entrypoint
  seed/          → seed data CLI command

pkgs/
  api/           → HTTP handlers and auth middleware
  auth/          → JWT generation/validation, bcrypt helpers
  models/        → GORM models (User, Cinema, Screen, Seat, Movie, Show, Booking)
  service/
    cinema/      → cinema and screen business logic
    movie/       → movie and show business logic
    booking/     → seat booking business logic

internal/
  cache/         → in-memory cache interface (go-cache)
  json/          → JSON write helpers
  router/        → HTTP router abstraction (httprouter)
  server/        → server lifecycle abstraction
  workgroup/     → goroutine lifetime management
  testhelpers/   → shared test utilities
```

Each service (`cinema`, `movie`, `booking`) is constructed with a DB handle, a cache instance, and a logger. Handlers call services; services call models.

---

## User roles

| Role | Can do |
|------|--------|
| `REGULAR` | Browse cinemas/shows, book seats |
| `THEATRE_OWNER` | Everything above + create cinemas, screens, movies, shows |

The role is embedded in the JWT. Protected routes check it via `RequireRole` middleware.

---

## API routes

### Public (no token required)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/healthz` | Health check |
| POST | `/register` | Create an account |
| POST | `/login` | Log in and get a JWT |
| GET | `/cinemas` | List all cinemas |
| GET | `/show` | List/query shows |
| GET | `/bookings` | List all bookings |

### Theatre owner only (`THEATRE_OWNER` role)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/cinema` | Add a cinema (hall or multiplex) |
| POST | `/screen` | Add a screen to a cinema |
| POST | `/movie` | Add a movie |
| POST | `/show` | Schedule a show on a screen |

### Regular user only (`REGULAR` role)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/book` | Book one or more seats for a show |

Authenticated requests need `Authorization: Bearer <token>` in the header.

---

## Data model (key entities)

- **City** → has many Cinemas
- **Cinema** → has many CinemaScreens
- **CinemaScreen** → has many CinemaSeats (Recliner, Premium, Front, Balcony)
- **Movie** → has many MovieShows
- **MovieShow** → belongs to a CinemaScreen; on creation, automatically generates a `BookingSeat` row for every seat on that screen
- **Booking** → belongs to a User and a MovieShow; tracks `PENDING → CONFIRMED / FAILED / CANCELLED`
- **BookingSeat** → tracks per-seat status (`AVAILABLE` / `BOOKED`) for a specific show

---

## Notable behaviours

**Race-safe seat booking** — `BookSeats` uses a conditional update (`WHERE status = AVAILABLE`) and checks `RowsAffected`. If another request grabbed a seat in between, the update reverts and returns an error.

**Overlap detection** — before a new show is inserted, `CheckOverlap` runs `SELECT ... FOR UPDATE` to lock competing rows. Combined with an explicit transaction, two shows can't be double-booked on the same screen at the same time.

**Auto seat generation** — the `AfterCreate` GORM hook on `MovieShow` copies all `CinemaSeat` rows for that screen into `BookingSeat` with status `AVAILABLE`. No manual step needed.

---

## Running locally

### With Docker (recommended)

```bash
docker-compose up
```

This starts a MySQL container and the API server. The API is available at `http://localhost:4000`.

### Without Docker

You need a running MySQL instance. Set the connection details via environment variables or a `.env` file, then:

```bash
make compile
./ticketing server
```

### Environment variables

| Variable | Default | Description |
|----------|---------|-------------|
| `AUTH_SECRET` | `default-secret-change-in-production` | JWT signing key — change this in production |

---

## Running tests

```bash
make test
```

Tests use a mock DB (`go-mocket`) and are organised by package. Each service package has a `main_test.go` that sets up shared test state.

---

## Other make targets

```bash
make compile     # download dependencies and build the binary
make fmt         # run go fmt
make clean       # remove build artifacts
make docker      # build a Docker image tagged with the current git ref
```

---

## Documentation

- [API Docs](docs/api/index.md)
- [Database Design](docs/database_design.md)
- [Code Structure](docs/code_structure.md)

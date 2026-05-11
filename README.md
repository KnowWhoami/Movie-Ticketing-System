# Movie Ticketing — REST API

A Go REST API for booking movie tickets. Users can register as a regular customer or a cinema owner. Cinema owners manage cinemas, screens, and shows. Regular users browse shows and book seats. Admins manage users, cities, movies, and cinemas.

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
  models/        → GORM models (User, City, Cinema, Screen, Seat, Movie, Show, Booking)
  service/
    city/        → city business logic
    user/        → user management business logic
    cinema/      → cinema, screen, and seat business logic
    movie/       → movie, show, and show-seat business logic
    booking/     → seat booking business logic

internal/
  cache/         → in-memory cache interface (go-cache)
  json/          → JSON write helpers
  router/        → HTTP router abstraction (httprouter)
  server/        → server lifecycle abstraction
  workgroup/     → goroutine lifetime management
  testhelpers/   → shared test utilities (real MySQL test DB)
```

Each service is constructed with a DB handle, a cache instance, and a logger. Handlers call services; services call models.

---

## User roles

| Role           | Can do                                                             |
|----------------|--------------------------------------------------------------------|
| `REGULAR`      | Browse cinemas/movies/shows, book seats                            |
| `CINEMA_OWNER` | Everything above + manage their own screens and schedule shows     |
| `ADMIN`        | Full access — manage users, cities, cinemas, movies                |

The role is embedded in the JWT. Protected routes check it via `RequireRole` middleware.

Self-registration via `POST /register` always creates a `REGULAR` account. Use `POST /user` (admin-only) to create `CINEMA_OWNER` or `ADMIN` accounts.

---

## API routes

### Public (no token required)

| Method | Path                                  | Description                            |
|--------|---------------------------------------|----------------------------------------|
| GET    | `/healthz`                            | Health check                           |
| POST   | `/register`                           | Self-register (creates REGULAR account)|
| POST   | `/login`                              | Log in and get a JWT                   |
| GET    | `/cities`                             | List all cities                        |
| GET    | `/cinemas`                            | List cinemas (filterable by name/city) |
| GET    | `/cinema/:id`                         | Get a cinema by ID                     |
| GET    | `/screens`                            | List screens (filterable by cinema)    |
| GET    | `/screens/:screen_id/seats`           | List seats in a screen                 |
| GET    | `/screens/:screen_id/seats/:seat_num` | Get a specific seat by number          |
| GET    | `/movies`                             | List movies (filterable by name)       |
| GET    | `/movie/:id`                          | Get a movie by ID                      |
| GET    | `/shows`                              | List shows (filterable by movie/cinema)|
| GET    | `/show/:id`                           | Get a show by ID                       |
| GET    | `/show/:id/seats`                     | List show seats with availability      |
| GET    | `/bookings`                           | List all bookings                      |

### Admin only (`ADMIN` role)

| Method | Path                        | Description                          |
|--------|-----------------------------|--------------------------------------|
| POST   | `/city`                     | Add a city                           |
| POST   | `/user`                     | Create a user with an explicit role  |
| GET    | `/users`                    | List users (filterable)              |
| GET    | `/user/:id`                 | Get a user by ID                     |
| PATCH  | `/user/:id`                 | Update a user                        |
| DELETE | `/user/:id`                 | Soft-delete a user                   |
| POST   | `/cinema`                   | Add a cinema                         |
| PATCH  | `/cinema/:id`               | Update a cinema                      |
| GET    | `/cinemas/:cinema_owner_id` | List cinemas by owner ID             |
| POST   | `/movie`                    | Add a movie                          |
| PATCH  | `/movie/:id`                | Update a movie                       |

### Cinema owner only (`CINEMA_OWNER` role)

| Method | Path                 | Description                                   |
|--------|----------------------|-----------------------------------------------|
| GET    | `/my-cinemas`        | List cinemas owned by the authenticated user  |
| POST   | `/screen`            | Add a screen to one of your cinemas           |
| POST   | `/show`              | Schedule a show on one of your screens        |
| PATCH  | `/show/:id/cancel`   | Cancel one of your shows                      |

### Regular user only (`REGULAR` role)

| Method | Path    | Description                         |
|--------|---------|-------------------------------------|
| POST   | `/book` | Book one or more seats for a show   |

Authenticated requests require `Authorization: Bearer <token>` in the header.

---

## Data model (key entities)

- **City** → has many Cinemas
- **Cinema** → belongs to a City and a cinema owner (`CINEMA_OWNER` user); has many CinemaScreens
- **CinemaScreen** → has many CinemaSeats (types: `RECLINER`, `PREMIUM`, `FRONT`, `BALCONY`)
- **Movie** → has many MovieShows; `duration` is stored in **minutes**
- **MovieShow** → belongs to a CinemaScreen; `end_time` is auto-calculated on creation (`start_time + duration`); on creation the `AfterCreate` hook auto-generates a `MovieShowSeat` row for every seat on the screen
- **MovieShowSeat** → tracks per-seat availability for a specific show via `booking_id` (`NULL` = available, set = booked)
- **Booking** → belongs to a User and a MovieShow; transitions `PENDING → CONFIRMED / FAILED / CANCELLED`

---

## Notable behaviours

**Race-safe seat booking** — `BookSeats` uses a conditional `UPDATE WHERE booking_id IS NULL` and checks `RowsAffected`. If another request claimed a seat concurrently, the affected count won't match and all seats are reverted.

**Overlap detection** — before a new show is inserted, `CheckOverlap` issues `SELECT ... FOR UPDATE` to lock competing rows. Combined with an explicit transaction, two shows on the same screen at the same time cannot both succeed.

**Auto seat generation** — the `AfterCreate` GORM hook on `MovieShow` copies every `CinemaSeat` for that screen into `MovieShowSeat`. No manual step is needed after scheduling a show.

---

## Running locally

### With Docker (recommended)

```bash
docker-compose up
```

This starts a MySQL container and the API server. The API is available at `http://localhost:4000`.

### Without Docker

You need a running MySQL instance. Configure it via environment variables, then:

```bash
make compile
./ticketing server
```

### Environment variables

| Variable      | Default                              | Description                                    |
|---------------|--------------------------------------|------------------------------------------------|
| `AUTH_SECRET` | `default-secret-change-in-production`| JWT signing key — **change this in production**|
| `DB_HOST`     | `localhost`                          | MySQL host                                     |
| `DB_PORT`     | `3306`                               | MySQL port                                     |
| `DB_NAME`     | —                                    | Database name                                  |
| `DB_USER`     | —                                    | Database user                                  |
| `DB_PASSWORD` | —                                    | Database password                              |

---

## Running tests

Tests run against a real MySQL instance (not a mock). The test helper connects to `127.0.0.1:3306` using database `ticketing_test` with credentials `user:user`. Start MySQL locally or via Docker before running tests.

```bash
make test
```

Each service package has a `main_test.go` that runs `AutoMigrate` and truncates all tables between tests to ensure isolation.

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

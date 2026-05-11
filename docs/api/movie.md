## Movie APIs

### `POST` `/movie` — Add a new movie

Requires an **admin** JWT token.

**Request**

```bash
curl -X POST http://localhost:4000/movie \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Inception",
    "description": "A thief who steals corporate secrets through dream-sharing technology.",
    "duration": 148
  }'
```

| Field         | Type   | Required | Description                  |
|---------------|--------|----------|------------------------------|
| `name`        | string | yes      | Movie title                  |
| `description` | string | yes      | Short synopsis               |
| `duration`    | int    | yes      | Runtime in **minutes**       |

**Response** — `201 Created`

```json
{
    "success": true,
    "status_code": 201,
    "data": {
        "movie": {
            "id": 1,
            "name": "Inception",
            "description": "A thief who steals corporate secrets through dream-sharing technology.",
            "duration": 148,
            "shows": null
        }
    },
    "error_message": ""
}
```

**Error — validation failure**

```json
{
    "success": false,
    "status_code": 422,
    "data": null,
    "error_message": "name, description and duration are mandatory fields"
}
```

---

### `GET` `/movies` — List all movies

Public endpoint. Supports an optional `name` query parameter for a case-insensitive substring match.

**Request**

```bash
# all movies
curl -X GET http://localhost:4000/movies

# search by name
curl -X GET "http://localhost:4000/movies?name=inception"
```

**Response** — `200 OK`

```json
{
    "success": true,
    "status_code": 200,
    "data": {
        "movies": [
            {
                "id": 1,
                "name": "Inception",
                "description": "A thief who steals corporate secrets through dream-sharing technology.",
                "duration": 148,
                "shows": null
            },
            {
                "id": 2,
                "name": "Interstellar",
                "description": "A team of explorers travel through a wormhole in space.",
                "duration": 169,
                "shows": null
            }
        ]
    },
    "error_message": ""
}
```

---

### `GET` `/movie/:id` — Get a movie by ID

Public endpoint. Response includes the movie's scheduled shows.

**Request**

```bash
curl -X GET http://localhost:4000/movie/1
```

**Response** — `200 OK`

```json
{
    "success": true,
    "status_code": 200,
    "data": {
        "movie": {
            "id": 1,
            "name": "Inception",
            "description": "A thief who steals corporate secrets through dream-sharing technology.",
            "duration": 148,
            "shows": [
                {
                    "id": 5,
                    "start_time": "2026-06-01T14:00:00Z",
                    "end_time": "2026-06-01T16:28:00Z",
                    "is_cancelled": false,
                    "CinemaScreen": { "id": 0, "name": "", "seats": null },
                    "Movie": { "id": 0, "name": "", "description": "", "duration": 0, "shows": null },
                    "bookings": null
                }
            ]
        }
    },
    "error_message": ""
}
```

**Error — not found**

```json
{
    "success": false,
    "status_code": 404,
    "data": null,
    "error_message": "movie not found"
}
```

---

### `PATCH` `/movie/:id` — Update a movie

Requires an **admin** JWT token. All body fields are optional — only supplied fields are updated.

**Request**

```bash
curl -X PATCH http://localhost:4000/movie/1 \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "duration": 150
  }'
```

**Response** — `200 OK`

```json
{
    "success": true,
    "status_code": 200,
    "data": {
        "movie": {
            "id": 1,
            "name": "Inception",
            "description": "A thief who steals corporate secrets through dream-sharing technology.",
            "duration": 150,
            "shows": null
        }
    },
    "error_message": ""
}
```

**Error — not found**

```json
{
    "success": false,
    "status_code": 422,
    "data": null,
    "error_message": "movie not found"
}
```

---

## Movie Show APIs

### `POST` `/show` — Schedule a show

Requires a **cinema owner** JWT token. The screen must belong to the authenticated owner.

`end_time` is calculated automatically from `start_time + movie.duration`. Show seats are also generated automatically for every seat in the screen. Scheduling a show on a screen that already has an overlapping show is rejected.

**Request**

```bash
curl -X POST http://localhost:4000/show \
  -H "Authorization: Bearer <cinema_owner_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "movie_id": 1,
    "cinema_screen_id": 5,
    "start_time": "2026-06-01T14:00:00Z"
  }'
```

| Field              | Type        | Required | Description                                  |
|--------------------|-------------|----------|----------------------------------------------|
| `movie_id`         | int         | yes      | ID of the movie to screen                    |
| `cinema_screen_id` | int         | yes      | ID of the screen to use                      |
| `start_time`       | RFC 3339    | yes      | Show start time (e.g. `2026-06-01T14:00:00Z`)|

**Response** — `201 Created`

```json
{
    "success": true,
    "status_code": 201,
    "data": {
        "show": {
            "id": 5,
            "start_time": "2026-06-01T14:00:00Z",
            "end_time": "2026-06-01T16:28:00Z",
            "is_cancelled": false,
            "CinemaScreen": {
                "id": 5,
                "name": "Screen 1",
                "seats": null
            },
            "Movie": {
                "id": 1,
                "name": "Inception",
                "description": "A thief who steals corporate secrets through dream-sharing technology.",
                "duration": 148,
                "shows": null
            },
            "bookings": null
        }
    },
    "error_message": ""
}
```

**Error — validation / scheduling failures**

```json
{ "success": false, "status_code": 422, "data": null, "error_message": "movie_id, cinema_screen_id and start_time are mandatory fields" }
{ "success": false, "status_code": 422, "data": null, "error_message": "cinema screen not found" }
{ "success": false, "status_code": 422, "data": null, "error_message": "cinema screen does not belong to your cinema" }
{ "success": false, "status_code": 422, "data": null, "error_message": "movie not found" }
{ "success": false, "status_code": 422, "data": null, "error_message": "the show has an overlap with 1 shows" }
```

---

### `GET` `/shows` — List all shows

Public endpoint. Supports optional query-string filters.

| Query param        | Type | Description                         |
|--------------------|------|-------------------------------------|
| `movie_id`         | int  | Filter by movie                     |
| `cinema_screen_id` | int  | Filter by screen                    |
| `cinema_id`        | int  | Filter by cinema (across all screens)|

**Request**

```bash
# all shows
curl -X GET http://localhost:4000/shows

# shows for a specific movie
curl -X GET "http://localhost:4000/shows?movie_id=1"

# shows in a specific cinema
curl -X GET "http://localhost:4000/shows?cinema_id=1"
```

**Response** — `200 OK`

```json
{
    "success": true,
    "status_code": 200,
    "data": {
        "shows": [
            {
                "id": 5,
                "start_time": "2026-06-01T14:00:00Z",
                "end_time": "2026-06-01T16:28:00Z",
                "is_cancelled": false,
                "CinemaScreen": {
                    "id": 5,
                    "name": "Screen 1",
                    "seats": null
                },
                "Movie": {
                    "id": 1,
                    "name": "Inception",
                    "description": "A thief who steals corporate secrets through dream-sharing technology.",
                    "duration": 148,
                    "shows": null
                },
                "bookings": null
            }
        ]
    },
    "error_message": ""
}
```

---

### `GET` `/show/:id` — Get a show by ID

Public endpoint.

**Request**

```bash
curl -X GET http://localhost:4000/show/5
```

**Response** — `200 OK` — same shape as a single entry from `GET /shows`

**Error — not found**

```json
{
    "success": false,
    "status_code": 404,
    "data": null,
    "error_message": "show not found"
}
```

---

### `PATCH` `/show/:id/cancel` — Cancel a show

Requires a **cinema owner** JWT token. The show's cinema must belong to the authenticated owner.

**Request**

```bash
curl -X PATCH http://localhost:4000/show/5/cancel \
  -H "Authorization: Bearer <cinema_owner_token>"
```

**Response** — `200 OK`

```json
{
    "success": true,
    "status_code": 200,
    "data": {
        "show": {
            "id": 5,
            "start_time": "2026-06-01T14:00:00Z",
            "end_time": "2026-06-01T16:28:00Z",
            "is_cancelled": true,
            "CinemaScreen": { "id": 5, "name": "Screen 1", "seats": null },
            "Movie": {
                "id": 1,
                "name": "Inception",
                "description": "A thief who steals corporate secrets through dream-sharing technology.",
                "duration": 148,
                "shows": null
            },
            "bookings": null
        }
    },
    "error_message": ""
}
```

**Error — cancellation failures**

```json
{ "success": false, "status_code": 422, "data": null, "error_message": "show not found" }
{ "success": false, "status_code": 422, "data": null, "error_message": "you can only cancel shows for your own cinemas" }
{ "success": false, "status_code": 422, "data": null, "error_message": "show is already cancelled" }
```

---

## Movie Show Seat APIs

### `GET` `/show/:id/seats` — List seats for a show

Public endpoint. Returns each seat's availability status. Supports optional query-string filters.

| Query param | Type   | Description                                              |
|-------------|--------|----------------------------------------------------------|
| `type`      | string | Filter by seat type: `RECLINER`, `PREMIUM`, `FRONT`, `BALCONY` |
| `available` | bool   | `true` = available only, `false` = booked only           |

**Request**

```bash
# all seats for show 5
curl -X GET http://localhost:4000/show/5/seats

# only available RECLINER seats
curl -X GET "http://localhost:4000/show/5/seats?type=RECLINER&available=true"
```

**Response** — `200 OK`

```json
{
    "success": true,
    "status_code": 200,
    "data": {
        "seats": [
            {
                "id": 101,
                "seat_number": 1,
                "type": "RECLINER",
                "available": true
            },
            {
                "id": 102,
                "seat_number": 2,
                "type": "RECLINER",
                "available": false
            },
            {
                "id": 103,
                "seat_number": 3,
                "type": "PREMIUM",
                "available": true
            }
        ]
    },
    "error_message": ""
}
```

> `available: false` means the seat has been claimed by a booking. Use `POST /book` to reserve available seats.

**Error — invalid show id**

```json
{
    "success": false,
    "status_code": 400,
    "data": null,
    "error_message": "invalid show id"
}
```

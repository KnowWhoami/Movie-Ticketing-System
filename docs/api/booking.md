## Booking APIs

---

### `POST` `/booking` — Book seats

Requires a valid JWT token (any role). `user_id` is derived from the token — it cannot be supplied in the body.

Use `GET /show/:id/seats?available=true` to discover available `show_seat_ids` before booking.

**Request**

```bash
curl -X POST http://localhost:4000/booking \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "movie_show_id": 3,
    "show_seat_ids": [101, 102]
  }'
```

| Field            | Type     | Required | Description                                     |
|------------------|----------|----------|-------------------------------------------------|
| `movie_show_id`  | int      | yes      | ID of the show to book                          |
| `show_seat_ids`  | []int    | yes      | IDs of the specific show seats to reserve       |

The booking is executed inside a single transaction. If any seat is already taken, all claims are rolled back and the booking is recorded as `FAILED`. Both `CONFIRMED` and `FAILED` bookings are persisted as audit records.

**Response** — `201 Created`

```json
{
    "success": true,
    "status_code": 201,
    "data": {
        "booking": {
            "id": 21,
            "status": "CONFIRMED",
            "user": {
                "id": 1,
                "Name": "Alice",
                "Email": "alice@example.com",
                "user_type": "REGULAR",
                "Bookings": null
            },
            "movie_show": {
                "id": 3,
                "start_time": "2026-06-01T14:00:00Z",
                "end_time": "2026-06-01T16:28:00Z",
                "is_cancelled": false,
                "CinemaScreen": { "id": 0, "name": "", "seats": null },
                "Movie": { "id": 0, "name": "", "description": "", "duration": 0, "shows": null },
                "bookings": null
            },
            "seats": [
                {
                    "id": 101,
                    "cinema_seat": {
                        "id": 5,
                        "seat_number": 3,
                        "type": "PREMIUM"
                    }
                },
                {
                    "id": 102,
                    "cinema_seat": {
                        "id": 6,
                        "seat_number": 4,
                        "type": "PREMIUM"
                    }
                }
            ]
        }
    },
    "error_message": ""
}
```

**Error — seat already taken**

```json
{
    "success": false,
    "status_code": 422,
    "data": {
        "booking": {
            "id": 22,
            "status": "FAILED",
            ...
        }
    },
    "error_message": "some seats are already booked"
}
```

**Error — validation failures**

```json
{ "success": false, "status_code": 400, "data": null, "error_message": "movie_show_id is required" }
{ "success": false, "status_code": 400, "data": null, "error_message": "show_seat_ids must not be empty" }
{ "success": false, "status_code": 400, "data": null, "error_message": "show not found" }
{ "success": false, "status_code": 400, "data": null, "error_message": "show has been cancelled" }
{ "success": false, "status_code": 400, "data": null, "error_message": "one or more seat IDs do not belong to this show" }
```

---

### `GET` `/bookings` — List bookings

Requires an **admin** or **cinema owner** JWT token.

- **ADMIN** — `movie_show_id` and `user_id` filters are optional; returns all matching bookings.
- **CINEMA_OWNER** — `movie_show_id` is **required**; the show must belong to the authenticated owner's cinema.

| Query param     | Type | Description                            |
|-----------------|------|----------------------------------------|
| `movie_show_id` | int  | Filter by show                         |
| `user_id`       | int  | Filter by user (ADMIN only in practice)|

**Request**

```bash
# all bookings for show 3 (cinema owner or admin)
curl -X GET "http://localhost:4000/bookings?movie_show_id=3" \
  -H "Authorization: Bearer <token>"

# all bookings by user 1 (admin)
curl -X GET "http://localhost:4000/bookings?user_id=1" \
  -H "Authorization: Bearer <admin_token>"
```

**Response** — `200 OK`

```json
{
    "success": true,
    "status_code": 200,
    "data": {
        "bookings": [
            {
                "id": 21,
                "status": "CONFIRMED",
                "user": {
                    "id": 1,
                    "Name": "Alice",
                    "Email": "alice@example.com",
                    "user_type": "REGULAR",
                    "Bookings": null
                },
                "movie_show": {
                    "id": 3,
                    "start_time": "2026-06-01T14:00:00Z",
                    "end_time": "2026-06-01T16:28:00Z",
                    "is_cancelled": false,
                    "CinemaScreen": { "id": 0, "name": "", "seats": null },
                    "Movie": { "id": 0, "name": "", "description": "", "duration": 0, "shows": null },
                    "bookings": null
                },
                "seats": [
                    {
                        "id": 101,
                        "cinema_seat": { "id": 5, "seat_number": 3, "type": "PREMIUM" }
                    }
                ]
            }
        ]
    },
    "error_message": ""
}
```

**Error — cinema owner accessing another owner's show**

```json
{
    "success": false,
    "status_code": 422,
    "data": null,
    "error_message": "show does not belong to your cinema"
}
```

**Error — cinema owner omits movie_show_id**

```json
{
    "success": false,
    "status_code": 422,
    "data": null,
    "error_message": "movie_show_id is required for cinema owners"
}
```

---

### `GET` `/my-bookings` — List my bookings

Requires a valid JWT token (any role). Returns all bookings for the authenticated user.

**Request**

```bash
curl -X GET http://localhost:4000/my-bookings \
  -H "Authorization: Bearer <token>"
```

**Response** — `200 OK`

```json
{
    "success": true,
    "status_code": 200,
    "data": {
        "bookings": [
            {
                "id": 21,
                "status": "CONFIRMED",
                "user": { "id": 0, "Name": "", "Email": "", "user_type": "", "Bookings": null },
                "movie_show": {
                    "id": 3,
                    "start_time": "2026-06-01T14:00:00Z",
                    "end_time": "2026-06-01T16:28:00Z",
                    "is_cancelled": false,
                    "CinemaScreen": { "id": 0, "name": "", "seats": null },
                    "Movie": { "id": 0, "name": "", "description": "", "duration": 0, "shows": null },
                    "bookings": null
                },
                "seats": [
                    {
                        "id": 101,
                        "cinema_seat": { "id": 5, "seat_number": 3, "type": "PREMIUM" }
                    }
                ]
            }
        ]
    },
    "error_message": ""
}
```

> `user` is not preloaded on this endpoint — `user_id` is already known from the token.

**Error — missing or invalid token**

```json
{
    "success": false,
    "status_code": 401,
    "data": null,
    "error_message": "missing or invalid authorization header"
}
```

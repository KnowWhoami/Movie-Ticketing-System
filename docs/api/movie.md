## Movie APIs

### `POST` `/movie` — Add a new movie

Requires a **theatre owner** JWT token in the `Authorization` header.

**Request**

```bash
curl -X POST http://localhost:4000/movie \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Inception",
    "description": "A thief who steals corporate secrets through dream-sharing technology.",
    "duration": 7200000000000
  }'
```

> `duration` is in nanoseconds (`time.Duration`). Examples: 1 hour = `3600000000000`, 2 hours = `7200000000000`.

**Response**

```json
{
    "success": true,
    "status_code": 200,
    "data": {
        "movie": {
            "id": 2,
            "name": "Inception",
            "description": "A thief who steals corporate secrets through dream-sharing technology.",
            "duration": 7200000000000,
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

### `POST` `/show` — Add a show to a cinema screen

Requires a **theatre owner** JWT token in the `Authorization` header.

**Request**

```bash
curl -X POST http://localhost:4000/show \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "movie_id": 1,
    "cinema_screen_id": 3,
    "start_time": "2024-06-01T18:00:00+05:30",
    "end_time": "2024-06-01T20:00:00+05:30"
  }'
```

> `start_time` must be before `end_time`. Use RFC3339 format.

**Response**

```json
{
    "success": true,
    "status_code": 200,
    "data": {
        "Show": {
            "id": 3,
            "start_time": "2024-06-01T18:00:00+05:30",
            "end_time": "2024-06-01T20:00:00+05:30",
            "CinemaScreen": {
                "id": 0,
                "name": "",
                "seats": null
            },
            "Movie": {
                "id": 1,
                "name": "Inception",
                "description": "A thief who steals corporate secrets through dream-sharing technology.",
                "duration": 7200000000000,
                "shows": null
            },
            "bookings": [],
            "seats": [
                {
                    "status": "AVAILABLE",
                    "cinema_seat": {
                        "id": 0,
                        "seat_number": 0,
                        "type": ""
                    }
                },
                {
                    "status": "AVAILABLE",
                    "cinema_seat": {
                        "id": 0,
                        "seat_number": 0,
                        "type": ""
                    }
                }
            ]
        }
    },
    "error_message": ""
}
```

**Error — invalid time range**

```json
{
    "success": false,
    "status_code": 422,
    "data": null,
    "error_message": "start time must be before end time"
}
```

---

### `GET` `/show` — Get a show with its seats and bookings

> **Note:** Despite being a GET request, this endpoint reads `show_id` from a JSON request body.

**Request**

```bash
curl -X GET http://localhost:4000/show \
  -H "Content-Type: application/json" \
  -d '{"show_id": 2}'
```

**Response**

```json
{
    "success": true,
    "status_code": 200,
    "data": {
        "show": {
            "id": 2,
            "start_time": "2024-06-01T18:00:00+05:30",
            "end_time": "2024-06-01T20:00:00+05:30",
            "CinemaScreen": {
                "id": 4,
                "name": "Screen 1",
                "seats": null
            },
            "Movie": {
                "id": 1,
                "name": "Inception",
                "description": "A thief who steals corporate secrets through dream-sharing technology.",
                "duration": 7200000000000,
                "shows": null
            },
            "bookings": [
                {
                    "id": 16,
                    "seat_count": 2,
                    "status": "CONFIRMED",
                    "user": {
                        "id": 0,
                        "Name": "",
                        "Email": "",
                        "Bookings": null
                    },
                    "movie_show": {
                        "id": 0,
                        "start_time": "0001-01-01T00:00:00Z",
                        "end_time": "0001-01-01T00:00:00Z",
                        "CinemaScreen": { "id": 0, "name": "", "seats": null },
                        "Movie": { "id": 0, "name": "", "description": "", "duration": 0, "shows": null },
                        "bookings": null,
                        "seats": null
                    },
                    "seats": null
                }
            ],
            "seats": [
                {
                    "status": "BOOKED",
                    "cinema_seat": { "id": 0, "seat_number": 0, "type": "" }
                },
                {
                    "status": "AVAILABLE",
                    "cinema_seat": { "id": 0, "seat_number": 0, "type": "" }
                },
                {
                    "status": "BOOKED",
                    "cinema_seat": { "id": 0, "seat_number": 0, "type": "" }
                }
            ]
        }
    },
    "error_message": ""
}
```

**Error — missing show_id**

```json
{
    "success": false,
    "status_code": 422,
    "data": null,
    "error_message": "show_id is mandatory fields"
}
```

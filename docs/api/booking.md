## Booking APIs

---

### `GET` `/bookings` — List all bookings

**Request**

```bash
curl -X GET http://localhost:4000/bookings
```

**Response**

```json
{
    "success": true,
    "status_code": 200,
    "data": {
        "bookings": [
            {
                "id": 16,
                "seat_count": 2,
                "status": "CONFIRMED",
                "user": {
                    "id": 1,
                    "Name": "Joy",
                    "Email": "joylal4896@gmail.com",
                    "Bookings": null
                },
                "movie_show": {
                    "id": 2,
                    "start_time": "2020-11-02T23:34:05+05:30",
                    "end_time": "2020-11-03T00:34:05+05:30",
                    "CinemaScreen": {
                        "id": 0,
                        "name": "",
                        "seats": null
                    },
                    "Movie": {
                        "id": 0,
                        "name": "",
                        "description": "",
                        "duration": 0,
                        "shows": null
                    },
                    "bookings": null,
                    "seats": null
                },
                "seats": [
                    {
                        "status": "BOOKED",
                        "cinema_seat": {
                            "id": 0,
                            "seat_number": 0,
                            "type": ""
                        }
                    },
                    {
                        "status": "BOOKED",
                        "cinema_seat": {
                            "id": 0,
                            "seat_number": 0,
                            "type": ""
                        }
                    }
                ]
            },
            {
                "id": 14,
                "seat_count": 2,
                "status": "FAILED",
                "user": {
                    "id": 1,
                    "Name": "Joy",
                    "Email": "joylal4896@gmail.com",
                    "Bookings": null
                },
                "movie_show": {
                    "id": 1,
                    "start_time": "2020-11-02T23:34:05+05:30",
                    "end_time": "2020-11-03T00:34:05+05:30",
                    "CinemaScreen": {
                        "id": 0,
                        "name": "",
                        "seats": null
                    },
                    "Movie": {
                        "id": 0,
                        "name": "",
                        "description": "",
                        "duration": 0,
                        "shows": null
                    },
                    "bookings": null,
                    "seats": null
                },
                "seats": []
            }
        ]
    },
    "error_message": ""
}
```

---

### `POST` `/book` — Book seats in a show

Requires a **regular user** JWT token in the `Authorization` header. Obtain a token via [`POST /login`](cinema.md).

**Request**

```bash
curl -X POST http://localhost:4000/book \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "show_id": 3,
    "user_id": 1,
    "seat_numbers": [1],
    "seat_type": "PREMIUM"
  }'
```

> `seat_type` must be one of: `RECLINER`, `PREMIUM`, `FRONT`, `BALCONY`

**Response**

```json
{
    "success": true,
    "status_code": 200,
    "data": {
        "booking": {
            "id": 21,
            "seat_count": 1,
            "status": "CONFIRMED",
            "user": {
                "id": 1,
                "Name": "Joy",
                "Email": "joylal4896@gmail.com",
                "Bookings": null
            },
            "movie_show": {
                "id": 3,
                "start_time": "2020-11-03T01:34:05+05:30",
                "end_time": "2020-11-03T02:34:05+05:30",
                "CinemaScreen": {
                    "id": 0,
                    "name": "",
                    "seats": null
                },
                "Movie": {
                    "id": 0,
                    "name": "",
                    "description": "",
                    "duration": 0,
                    "shows": null
                },
                "bookings": null,
                "seats": null
            },
            "seats": [
                {
                    "status": "BOOKED",
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

**Error — missing or invalid token**

```json
{
    "success": false,
    "status_code": 401,
    "data": null,
    "error_message": "invalid or expired token"
}
```

**Error — validation failure**

```json
{
    "success": false,
    "status_code": 422,
    "data": null,
    "error_message": "show_id, seat_numbers, seat_type and user_id are required parameters"
}
```

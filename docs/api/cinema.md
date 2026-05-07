## Auth APIs

### `POST` `/register` — Register a new user

**Request**

```bash
curl -X POST http://localhost:4000/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Joy",
    "email": "joy@example.com",
    "password": "secret123",
    "user_type": "THEATRE_OWNER"
  }'
```

> `user_type` is optional and defaults to `REGULAR`. Use `THEATRE_OWNER` to manage cinemas and shows.

**Response**

```json
{
    "success": true,
    "status_code": 200,
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "user": {
            "id": 1,
            "name": "Joy",
            "email": "joy@example.com",
            "user_type": "THEATRE_OWNER"
        }
    },
    "error_message": ""
}
```

---

### `POST` `/login` — Login and get a token

**Request**

```bash
curl -X POST http://localhost:4000/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "joy@example.com",
    "password": "secret123"
  }'
```

**Response**

```json
{
    "success": true,
    "status_code": 200,
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "user": {
            "id": 1,
            "name": "Joy",
            "email": "joy@example.com",
            "user_type": "THEATRE_OWNER"
        }
    },
    "error_message": ""
}
```

**Error — wrong credentials**

```json
{
    "success": false,
    "status_code": 401,
    "data": null,
    "error_message": "invalid email or password"
}
```

---

## Cinema APIs

### `GET` `/cinemas` — List all cinemas

**Request**

```bash
curl -X GET http://localhost:4000/cinemas
```

**Response**

```json
{
    "success": true,
    "status_code": 200,
    "data": {
        "cinemas": [
            {
                "id": 1,
                "name": "PVR Cinemas",
                "screens": [
                    {
                        "id": 1,
                        "name": "Screen 1",
                        "seats": [
                            {
                                "id": 1,
                                "seat_number": 1,
                                "type": "RECLINER"
                            },
                            {
                                "id": 2,
                                "seat_number": 2,
                                "type": "PREMIUM"
                            }
                        ]
                    }
                ],
                "city": {
                    "id": 1,
                    "name": "Bengaluru",
                    "zip_code": "560068"
                }
            }
        ]
    },
    "error_message": ""
}
```

---

### `POST` `/cinema` — Add a new cinema

Requires a **theatre owner** JWT token in the `Authorization` header.

**Request**

```bash
curl -X POST http://localhost:4000/cinema \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "cinema_name": "PVR Cinemas",
    "city_id": 1
  }'
```

**Response**

```json
{
    "success": true,
    "status_code": 200,
    "data": {
        "cinema": {
            "id": 2,
            "name": "PVR Cinemas",
            "screens": null,
            "city": {
                "id": 0,
                "name": "",
                "zip_code": ""
            }
        }
    },
}
```

**Error — validation failure**

```json
{
    "success": false,
    "status_code": 422,
    "data": null,
    "error_message": "cinema_name and city_id are required parameters"
}
```

---

### `POST` `/screen` — Add a screen to a cinema

Requires a **theatre owner** JWT token in the `Authorization` header.

**Request**

```bash
curl -X POST http://localhost:4000/screen \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "cinema_id": 1,
    "screen_name": "Screen 1",
    "seats": [
      { "seat_number": 1, "seat_type": "RECLINER" },
      { "seat_number": 2, "seat_type": "RECLINER" },
      { "seat_number": 3, "seat_type": "PREMIUM" },
      { "seat_number": 4, "seat_type": "BALCONY" }
    ]
  }'
```

> `seat_type` must be one of: `RECLINER`, `PREMIUM`, `FRONT`, `BALCONY`

**Response**

```json
{
    "success": true,
    "status_code": 200,
    "data": {
        "cinema_screen": {
            "id": 5,
            "name": "Screen 1",
            "seats": [
                {
                    "id": 14,
                    "seat_number": 1,
                    "type": "RECLINER"
                },
                {
                    "id": 15,
                    "seat_number": 2,
                    "type": "RECLINER"
                },
                {
                    "id": 16,
                    "seat_number": 3,
                    "type": "PREMIUM"
                },
                {
                    "id": 17,
                    "seat_number": 4,
                    "type": "BALCONY"
                }
            ]
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
    "error_message": "cinema_id, screen_name are required and at least one seat must be provided"
}
```

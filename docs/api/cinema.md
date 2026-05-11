## Auth APIs

### `POST` `/register` — Register a new user

**Request**

```bash
curl -X POST http://localhost:4000/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Joy",
    "email": "joy@example.com",
    "password": "secret123"
  }'
```

> Self-registration always creates a `REGULAR` user. Use `POST /user` (admin-only) to create `CINEMA_OWNER` or `ADMIN` accounts.

**Response** — `201 Created`

```json
{
    "success": true,
    "status_code": 201,
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "user": {
            "id": 1,
            "name": "Joy",
            "email": "joy@example.com",
            "user_type": "REGULAR"
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

**Response** — `200 OK`

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
            "user_type": "REGULAR"
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

Public endpoint. Supports optional query-string filters.

| Query param | Type   | Description                        |
|-------------|--------|------------------------------------|
| `name`      | string | Case-insensitive substring match   |
| `city_id`   | int    | Filter cinemas by city             |

**Request**

```bash
# all cinemas
curl -X GET http://localhost:4000/cinemas

# filtered by city
curl -X GET "http://localhost:4000/cinemas?city_id=1"

# filtered by name
curl -X GET "http://localhost:4000/cinemas?name=pvr"
```

**Response** — `200 OK`

```json
{
    "success": true,
    "status_code": 200,
    "data": {
        "cinemas": [
            {
                "id": 1,
                "name": "PVR Cinemas",
                "cinema_owner": {
                    "id": 2,
                    "Name": "Alice",
                    "Email": "alice@example.com",
                    "user_type": "CINEMA_OWNER"
                },
                "screens": [
                    {
                        "id": 1,
                        "name": "Screen 1",
                        "seats": [
                            { "id": 1, "seat_number": 1, "type": "RECLINER" },
                            { "id": 2, "seat_number": 2, "type": "PREMIUM" }
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

### `GET` `/cinema/:id` — Get a cinema by ID

Public endpoint.

**Request**

```bash
curl -X GET http://localhost:4000/cinema/1
```

**Response** — `200 OK`

```json
{
    "success": true,
    "status_code": 200,
    "data": {
        "cinema": {
            "id": 1,
            "name": "PVR Cinemas",
            "cinema_owner": {
                "id": 2,
                "Name": "Alice",
                "Email": "alice@example.com",
                "user_type": "CINEMA_OWNER"
            },
            "screens": [
                {
                    "id": 1,
                    "name": "Screen 1",
                    "seats": [
                        { "id": 1, "seat_number": 1, "type": "RECLINER" },
                        { "id": 2, "seat_number": 2, "type": "PREMIUM" }
                    ]
                }
            ],
            "city": {
                "id": 1,
                "name": "Bengaluru",
                "zip_code": "560068"
            }
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
    "error_message": "cinema not found"
}
```

---

### `GET` `/my-cinemas` — Get cinemas owned by the logged-in user

Requires a **cinema owner** JWT token.

**Request**

```bash
curl -X GET http://localhost:4000/my-cinemas \
  -H "Authorization: Bearer <cinema_owner_token>"
```

**Response** — `200 OK` — same shape as `GET /cinemas`

---

### `GET` `/cinemas/:cinema_owner_id` — Get cinemas by owner ID

Requires an **admin** JWT token.

**Request**

```bash
curl -X GET http://localhost:4000/cinemas/2 \
  -H "Authorization: Bearer <admin_token>"
```

**Response** — `200 OK` — same shape as `GET /cinemas`

---

### `POST` `/cinema` — Add a new cinema

Requires an **admin** JWT token.

**Request**

```bash
curl -X POST http://localhost:4000/cinema \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "cinema_name": "PVR Cinemas",
    "city_id": 1,
    "cinema_owner_id": 2
  }'
```

| Field             | Type | Required | Description                                    |
|-------------------|------|----------|------------------------------------------------|
| `cinema_name`     | string | yes    | Name of the cinema                             |
| `city_id`         | int    | yes    | ID of an existing city                         |
| `cinema_owner_id` | int    | yes    | ID of a user with `CINEMA_OWNER` role          |

**Response** — `201 Created`

```json
{
    "success": true,
    "status_code": 201,
    "data": {
        "cinema": {
            "id": 3,
            "name": "PVR Cinemas",
            "cinema_owner": {
                "id": 2,
                "Name": "Alice",
                "Email": "alice@example.com",
                "user_type": "CINEMA_OWNER"
            },
            "screens": null,
            "city": {
                "id": 1,
                "name": "Bengaluru",
                "zip_code": "560068"
            }
        }
    },
    "error_message": ""
}
```

**Error — validation failures**

```json
{ "success": false, "status_code": 422, "data": null, "error_message": "cinema_name, city_id, and cinema_owner_id are required" }
{ "success": false, "status_code": 422, "data": null, "error_message": "city not found" }
{ "success": false, "status_code": 422, "data": null, "error_message": "cinema owner not found" }
{ "success": false, "status_code": 422, "data": null, "error_message": "user 2 is not a cinema owner" }
{ "success": false, "status_code": 422, "data": null, "error_message": "cinema \"PVR Cinemas\" already exists in this city" }
```

---

### `PATCH` `/cinema/:id` — Update a cinema

Requires an **admin** JWT token. All body fields are optional — only supplied fields are updated.

**Request**

```bash
curl -X PATCH http://localhost:4000/cinema/3 \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "cinema_name": "INOX Cinemas",
    "city_id": 2
  }'
```

**Response** — `200 OK`

```json
{
    "success": true,
    "status_code": 200,
    "data": {
        "cinema": {
            "id": 3,
            "name": "INOX Cinemas",
            "cinema_owner": {
                "id": 2,
                "Name": "Alice",
                "Email": "alice@example.com",
                "user_type": "CINEMA_OWNER"
            },
            "screens": null,
            "city": {
                "id": 2,
                "name": "Mumbai",
                "zip_code": "400001"
            }
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
    "error_message": "cinema not found"
}
```

---

## Cinema Screen APIs

### `POST` `/screen` — Add a screen to a cinema

Requires a **cinema owner** JWT token. The cinema must belong to the authenticated user.

The `seats` field is a map of seat type to count. Seat numbers are auto-assigned sequentially starting from 1.

| Seat type   | Description              |
|-------------|--------------------------|
| `RECLINER`  | Recliner seats           |
| `PREMIUM`   | Premium seats            |
| `FRONT`     | Front-row seats          |
| `BALCONY`   | Balcony seats            |

**Request**

```bash
curl -X POST http://localhost:4000/screen \
  -H "Authorization: Bearer <cinema_owner_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "cinema_id": 1,
    "screen_name": "Screen 1",
    "seats": {
      "RECLINER": 10,
      "PREMIUM": 30,
      "FRONT": 20
    }
  }'
```

| Field         | Type            | Required | Description                                  |
|---------------|-----------------|----------|----------------------------------------------|
| `cinema_id`   | int             | yes      | ID of the cinema this screen belongs to      |
| `screen_name` | string          | yes      | Name of the screen (unique within a cinema)  |
| `seats`       | map[string]int  | yes      | Seat type → count (at least one entry)       |

**Response** — `201 Created`

The response returns a summary with seat counts per type, not individual seat records.

```json
{
    "success": true,
    "status_code": 201,
    "data": {
        "cinema_screen": {
            "id": 5,
            "name": "Screen 1",
            "cinema_id": 1,
            "seats_summary": {
                "RECLINER": 10,
                "PREMIUM": 30,
                "FRONT": 20
            }
        }
    },
    "error_message": ""
}
```

**Error — validation failures**

```json
{ "success": false, "status_code": 422, "data": null, "error_message": "cinema_id and screen_name are required" }
{ "success": false, "status_code": 422, "data": null, "error_message": "at least one seat type with a count must be provided" }
{ "success": false, "status_code": 422, "data": null, "error_message": "invalid seat type \"VIP\"" }
{ "success": false, "status_code": 422, "data": null, "error_message": "cinema not found" }
{ "success": false, "status_code": 422, "data": null, "error_message": "cinema does not belong to your account" }
{ "success": false, "status_code": 422, "data": null, "error_message": "screen \"Screen 1\" already exists in this cinema" }
```

---

### `GET` `/screens` — List all screens

Public endpoint. Supports optional query-string filters.

| Query param | Type   | Description                       |
|-------------|--------|-----------------------------------|
| `cinema_id` | int    | Filter screens by cinema          |
| `name`      | string | Case-insensitive substring match  |

**Request**

```bash
# all screens
curl -X GET http://localhost:4000/screens

# filtered by cinema
curl -X GET "http://localhost:4000/screens?cinema_id=1"
```

**Response** — `200 OK`

```json
{
    "success": true,
    "status_code": 200,
    "data": {
        "cinema_screens": [
            {
                "id": 5,
                "name": "Screen 1",
                "cinema_id": 1,
                "seats_summary": {
                    "RECLINER": 10,
                    "PREMIUM": 30,
                    "FRONT": 20
                }
            },
            {
                "id": 6,
                "name": "Screen 2",
                "cinema_id": 1,
                "seats_summary": {
                    "BALCONY": 50
                }
            }
        ]
    },
    "error_message": ""
}
```

---

## Cinema Seat APIs

### `GET` `/screens/:screen_id/seats` — List all seats in a screen

Public endpoint. Supports an optional `type` query parameter to filter by seat type.

**Request**

```bash
# all seats in screen 5
curl -X GET http://localhost:4000/screens/5/seats

# only RECLINER seats
curl -X GET "http://localhost:4000/screens/5/seats?type=RECLINER"
```

**Response** — `200 OK`

```json
{
    "success": true,
    "status_code": 200,
    "data": {
        "seats": [
            { "id": 1, "seat_number": 1, "type": "RECLINER" },
            { "id": 2, "seat_number": 2, "type": "RECLINER" },
            { "id": 11, "seat_number": 11, "type": "PREMIUM" },
            { "id": 12, "seat_number": 12, "type": "PREMIUM" }
        ]
    },
    "error_message": ""
}
```

**Error — invalid screen_id**

```json
{
    "success": false,
    "status_code": 400,
    "data": null,
    "error_message": "invalid screen_id"
}
```

---

### `GET` `/screens/:screen_id/seats/:seat_number` — Get a specific seat by number

Public endpoint.

**Request**

```bash
curl -X GET http://localhost:4000/screens/5/seats/3
```

**Response** — `200 OK`

```json
{
    "success": true,
    "status_code": 200,
    "data": {
        "seat": {
            "id": 3,
            "seat_number": 3,
            "type": "RECLINER"
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
    "error_message": "seat 3 not found in screen 5"
}
```

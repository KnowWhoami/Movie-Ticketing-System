## User Management APIs

All endpoints in this section require an **admin** JWT token.

> These endpoints are for admin-managed account creation. Users can self-register via `POST /register` (see [Cinema Management APIs](cinema.md#auth-apis)), which always creates a `REGULAR` account.

---

### `POST` `/user` — Create a user with an explicit role

**Request**

```bash
curl -X POST http://localhost:4000/user \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice",
    "email": "alice@example.com",
    "password": "secret123",
    "user_type": "CINEMA_OWNER"
  }'
```

| Field       | Type   | Required | Description                                         |
|-------------|--------|----------|-----------------------------------------------------|
| `name`      | string | yes      | Display name                                        |
| `email`     | string | yes      | Unique email address                                |
| `password`  | string | yes      | Plain-text password (hashed server-side)            |
| `user_type` | string | yes      | One of: `REGULAR`, `CINEMA_OWNER`, `ADMIN`          |

**Response** — `201 Created`

> Note: `Name` and `Email` are capitalized in the JSON response because the model fields have no `json` tags.

```json
{
    "success": true,
    "status_code": 201,
    "data": {
        "user": {
            "id": 3,
            "Name": "Alice",
            "Email": "alice@example.com",
            "user_type": "CINEMA_OWNER"
        }
    },
    "error_message": ""
}
```

**Error — validation failures**

```json
{ "success": false, "status_code": 422, "data": null, "error_message": "name, email, and password are required" }
{ "success": false, "status_code": 422, "data": null, "error_message": "user_type must be one of: REGULAR, CINEMA_OWNER, ADMIN" }
{ "success": false, "status_code": 422, "data": null, "error_message": "email already registered" }
```

---

### `GET` `/users` — List all users

Supports optional query-string filters.

| Query param | Type   | Description                        |
|-------------|--------|------------------------------------|
| `name`      | string | Case-insensitive substring match   |
| `email`     | string | Case-insensitive substring match   |
| `user_type` | string | Exact match: `REGULAR`, `CINEMA_OWNER`, `ADMIN` |

**Request**

```bash
# all users
curl -X GET http://localhost:4000/users \
  -H "Authorization: Bearer <admin_token>"

# filter by role
curl -X GET "http://localhost:4000/users?user_type=CINEMA_OWNER" \
  -H "Authorization: Bearer <admin_token>"
```

**Response** — `200 OK`

```json
{
    "success": true,
    "status_code": 200,
    "data": {
        "users": [
            {
                "id": 1,
                "Name": "Joy",
                "Email": "joy@example.com",
                "user_type": "REGULAR"
            },
            {
                "id": 2,
                "Name": "Alice",
                "Email": "alice@example.com",
                "user_type": "CINEMA_OWNER"
            }
        ]
    },
    "error_message": ""
}
```

---

### `GET` `/user/:id` — Get a user by ID

**Request**

```bash
curl -X GET http://localhost:4000/user/2 \
  -H "Authorization: Bearer <admin_token>"
```

**Response** — `200 OK`

```json
{
    "success": true,
    "status_code": 200,
    "data": {
        "user": {
            "id": 2,
            "Name": "Alice",
            "Email": "alice@example.com",
            "user_type": "CINEMA_OWNER"
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
    "error_message": "user not found"
}
```

---

### `PATCH` `/user/:id` — Update a user

All body fields are optional — only supplied fields are updated.

**Request**

```bash
curl -X PATCH http://localhost:4000/user/2 \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice Smith",
    "user_type": "ADMIN"
  }'
```

**Response** — `200 OK`

```json
{
    "success": true,
    "status_code": 200,
    "data": {
        "user": {
            "id": 2,
            "Name": "Alice Smith",
            "Email": "alice@example.com",
            "user_type": "ADMIN"
        }
    },
    "error_message": ""
}
```

**Error — validation failures**

```json
{ "success": false, "status_code": 422, "data": null, "error_message": "user_type must be one of: REGULAR, CINEMA_OWNER, ADMIN" }
{ "success": false, "status_code": 422, "data": null, "error_message": "email already in use" }
{ "success": false, "status_code": 422, "data": null, "error_message": "user not found" }
```

---

### `DELETE` `/user/:id` — Delete a user

Performs a soft-delete (GORM `DeletedAt`). The user record is retained in the database but excluded from all queries.

**Request**

```bash
curl -X DELETE http://localhost:4000/user/2 \
  -H "Authorization: Bearer <admin_token>"
```

**Response** — `200 OK`

```json
{
    "success": true,
    "status_code": 200,
    "data": null,
    "error_message": ""
}
```

**Error — not found**

```json
{
    "success": false,
    "status_code": 422,
    "data": null,
    "error_message": "user not found"
}
```

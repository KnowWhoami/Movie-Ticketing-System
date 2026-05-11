## City APIs

### `POST` `/city` — Add a new city

Requires an **admin** JWT token in the `Authorization` header.

**Request**

```bash
curl -X POST http://localhost:4000/city \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bengaluru",
    "zip_code": "560068"
  }'
```

| Field      | Type   | Required | Description           |
|------------|--------|----------|-----------------------|
| `name`     | string | yes      | Name of the city      |
| `zip_code` | string | yes      | ZIP / postal code     |

**Response** — `201 Created`

```json
{
    "success": true,
    "status_code": 201,
    "data": {
        "city": {
            "id": 1,
            "name": "Bengaluru",
            "zip_code": "560068"
        }
    },
    "error_message": ""
}
```

**Error — missing fields**

```json
{
    "success": false,
    "status_code": 422,
    "data": null,
    "error_message": "name and zip_code are required"
}
```

**Error — city already exists**

```json
{
    "success": false,
    "status_code": 422,
    "data": null,
    "error_message": "city \"Bengaluru\" already exists"
}
```

**Error — not authenticated / wrong role**

```json
{
    "success": false,
    "status_code": 403,
    "data": null,
    "error_message": "insufficient permissions"
}
```

---

### `GET` `/cities` — List all cities

Public endpoint — no authentication required.

**Request**

```bash
curl -X GET http://localhost:4000/cities
```

**Response** — `200 OK`

```json
{
    "success": true,
    "status_code": 200,
    "data": {
        "cities": [
            {
                "id": 1,
                "name": "Bengaluru",
                "zip_code": "560068"
            },
            {
                "id": 2,
                "name": "Mumbai",
                "zip_code": "400001"
            }
        ]
    },
    "error_message": ""
}
```

**Response — no cities yet**

```json
{
    "success": true,
    "status_code": 200,
    "data": {
        "cities": []
    },
    "error_message": ""
}
```

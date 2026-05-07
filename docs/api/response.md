# API Contracts and Responses

## Standard Response Format

All endpoints return the same JSON envelope:

| Field           | Type    | Description                                  |
|-----------------|---------|----------------------------------------------|
| `success`       | boolean | `true` on success, `false` on error          |
| `status_code`   | integer | HTTP status code                             |
| `data`          | object  | Response payload (null on error)             |
| `error_message` | string  | Error description (empty string on success)  |

**Success**

```json
{
    "success": true,
    "status_code": 200,
    "data": { },
    "error_message": ""
}
```

**Error**

```json
{
    "success": false,
    "status_code": 422,
    "data": null,
    "error_message": "descriptive error message"
}
```

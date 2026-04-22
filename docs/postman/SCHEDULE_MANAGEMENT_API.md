# Schedule Management API - Postman Collection Update

## Overview
This document provides the Postman collection requests for Schedule Management APIs that need to be added to the main CRM Healthcare API collection.

## Folder Structure
```
CRM Healthcare API
└── Schedule Management
    ├── List Schedules
    ├── Get Schedule by ID
    ├── Create Schedule
    ├── Update Schedule
    ├── Delete Schedule
    ├── Assign Schedule
    ├── Bulk Assign Schedules
    ├── Check Conflicts
    ├── Approve Schedule
    └── Reject Schedule
```

---

## 1. List Schedules

**Method:** `GET`  
**URL:** `{{base_url}}/api/v1/schedules`  
**Auth:** Bearer Token

### Query Parameters:
- `page` (integer, optional) - Page number (default: 1)
- `per_page` (integer, optional) - Items per page (default: 10, max: 100)
- `search` (string, optional) - Search in title and description
- `status` (string, optional) - Filter by status: `pending`, `approved`, `rejected`, `completed`, `cancelled`
- `type` (string, optional) - Filter by type: `visit`, `task`, `meeting`, `other`
- `assigned_to` (string, optional) - Filter by assigned user ID
- `start_date_from` (string, optional) - Filter start date from (YYYY-MM-DD)
- `start_date_to` (string, optional) - Filter start date to (YYYY-MM-DD)

### Example Request:
```
GET {{base_url}}/api/v1/schedules?page=1&per_page=20&status=pending&type=visit
```

### Success Response (200):
```json
{
  "success": true,
  "data": [
    {
      "id": "schedule_abc123",
      "title": "Visit to Hospital XYZ",
      "description": "Monthly visit for product presentation",
      "type": "visit",
      "start_time": "2024-01-20T09:00:00+07:00",
      "end_time": "2024-01-20T11:00:00+07:00",
      "assigned_to": "user_123",
      "assigned_user": {
        "id": "user_123",
        "name": "John Doe",
        "email": "john@example.com"
      },
      "assigned_by": "user_456",
      "assigned_by_user": {
        "id": "user_456",
        "name": "Manager Jane",
        "email": "jane@example.com"
      },
      "account_id": "account_789",
      "account": {
        "id": "account_789",
        "name": "Hospital XYZ"
      },
      "contact_id": "contact_321",
      "contact": {
        "id": "contact_321",
        "name": "Dr. Smith",
        "email": "smith@hospitalxyz.com",
        "phone": "+62812345678"
      },
      "location": {
        "lat": -6.2088,
        "lng": 106.8456,
        "address": "Hospital XYZ, Jakarta"
      },
      "status": "pending",
      "is_recurring": false,
      "recurring_pattern": null,
      "created_at": "2024-01-15T10:00:00+07:00",
      "updated_at": "2024-01-15T10:00:00+07:00"
    }
  ],
  "meta": {
    "pagination": {
      "page": 1,
      "per_page": 20,
      "total": 45,
      "total_pages": 3,
      "has_next": true,
      "has_prev": false,
      "next_page": 2
    }
  },
  "timestamp": "2024-01-15T10:30:45+07:00",
  "request_id": "req_abc123xyz"
}
```

---

## 2. Get Schedule by ID

**Method:** `GET`  
**URL:** `{{base_url}}/api/v1/schedules/:id`  
**Auth:** Bearer Token

### Path Parameters:
- `id` (string, required) - Schedule ID

### Example Request:
```
GET {{base_url}}/api/v1/schedules/schedule_abc123
```

### Success Response (200):
```json
{
  "success": true,
  "data": {
    "id": "schedule_abc123",
    "title": "Visit to Hospital XYZ",
    "description": "Monthly visit for product presentation",
    "type": "visit",
    "start_time": "2024-01-20T09:00:00+07:00",
    "end_time": "2024-01-20T11:00:00+07:00",
    "assigned_to": "user_123",
    "assigned_user": {
      "id": "user_123",
      "name": "John Doe",
      "email": "john@example.com"
    },
    "assigned_by": "user_456",
    "assigned_by_user": {
      "id": "user_456",
      "name": "Manager Jane",
      "email": "jane@example.com"
    },
    "account_id": "account_789",
    "account": {
      "id": "account_789",
      "name": "Hospital XYZ"
    },
    "contact_id": "contact_321",
    "contact": {
      "id": "contact_321",
      "name": "Dr. Smith",
      "email": "smith@hospitalxyz.com",
      "phone": "+62812345678"
    },
    "location": {
      "lat": -6.2088,
      "lng": 106.8456,
      "address": "Hospital XYZ, Jakarta"
    },
    "status": "pending",
    "is_recurring": false,
    "recurring_pattern": null,
    "created_at": "2024-01-15T10:00:00+07:00",
    "updated_at": "2024-01-15T10:00:00+07:00"
  },
  "timestamp": "2024-01-15T10:30:45+07:00",
  "request_id": "req_abc123xyz"
}
```

### Error Response (404):
```json
{
  "success": false,
  "error": {
    "code": "SCHEDULE_NOT_FOUND",
    "message": "Schedule not found"
  },
  "timestamp": "2024-01-15T10:30:45+07:00",
  "request_id": "req_abc123xyz"
}
```

---

## 3. Create Schedule

**Method:** `POST`  
**URL:** `{{base_url}}/api/v1/schedules`  
**Auth:** Bearer Token

### Request Body:
```json
{
  "title": "Visit to Hospital XYZ",
  "description": "Monthly visit for product presentation",
  "type": "visit",
  "start_time": "2024-01-20T09:00:00+07:00",
  "end_time": "2024-01-20T11:00:00+07:00",
  "assigned_to": "user_123",
  "account_id": "account_789",
  "contact_id": "contact_321",
  "deal_id": "deal_456",
  "location": {
    "lat": -6.2088,
    "lng": 106.8456,
    "address": "Hospital XYZ, Jakarta"
  },
  "is_recurring": false,
  "recurring_pattern": null
}
```

### Request Body (with Recurring Pattern):
```json
{
  "title": "Weekly Team Meeting",
  "description": "Weekly sales team sync",
  "type": "meeting",
  "start_time": "2024-01-22T09:00:00+07:00",
  "end_time": "2024-01-22T10:00:00+07:00",
  "assigned_to": "user_123",
  "location": {
    "lat": -6.2088,
    "lng": 106.8456,
    "address": "Office - Meeting Room A"
  },
  "is_recurring": true,
  "recurring_pattern": {
    "type": "weekly",
    "interval": 1,
    "days_of_week": [1],
    "end_date": "2024-04-22T00:00:00+07:00"
  }
}
```

### Success Response (201):
```json
{
  "success": true,
  "data": {
    "id": "schedule_new123",
    "title": "Visit to Hospital XYZ",
    "description": "Monthly visit for product presentation",
    "type": "visit",
    "start_time": "2024-01-20T09:00:00+07:00",
    "end_time": "2024-01-20T11:00:00+07:00",
    "assigned_to": "user_123",
    "assigned_by": "user_current",
    "account_id": "account_789",
    "contact_id": "contact_321",
    "deal_id": "deal_456",
    "location": {
      "lat": -6.2088,
      "lng": 106.8456,
      "address": "Hospital XYZ, Jakarta"
    },
    "status": "pending",
    "is_recurring": false,
    "recurring_pattern": null,
    "created_at": "2024-01-15T10:30:45+07:00",
    "updated_at": "2024-01-15T10:30:45+07:00"
  },
  "timestamp": "2024-01-15T10:30:45+07:00",
  "request_id": "req_abc123xyz"
}
```

### Error Response (400):
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "field_errors": [
      {
        "field": "title",
        "message": "Title is required"
      },
      {
        "field": "start_time",
        "message": "Start time must be before end time"
      }
    ]
  },
  "timestamp": "2024-01-15T10:30:45+07:00",
  "request_id": "req_abc123xyz"
}
```

---

## 4. Update Schedule

**Method:** `PUT`  
**URL:** `{{base_url}}/api/v1/schedules/:id`  
**Auth:** Bearer Token

### Path Parameters:
- `id` (string, required) - Schedule ID

### Request Body (all fields optional):
```json
{
  "title": "Updated Visit to Hospital XYZ",
  "description": "Updated description",
  "type": "visit",
  "start_time": "2024-01-20T10:00:00+07:00",
  "end_time": "2024-01-20T12:00:00+07:00",
  "assigned_to": "user_456",
  "status": "approved",
  "location": {
    "lat": -6.2088,
    "lng": 106.8456,
    "address": "Updated Address"
  }
}
```

### Success Response (200):
```json
{
  "success": true,
  "data": {
    "id": "schedule_abc123",
    "title": "Updated Visit to Hospital XYZ",
    "description": "Updated description",
    "type": "visit",
    "start_time": "2024-01-20T10:00:00+07:00",
    "end_time": "2024-01-20T12:00:00+07:00",
    "assigned_to": "user_456",
    "status": "approved",
    "location": {
      "lat": -6.2088,
      "lng": 106.8456,
      "address": "Updated Address"
    },
    "created_at": "2024-01-15T10:00:00+07:00",
    "updated_at": "2024-01-15T11:00:00+07:00"
  },
  "timestamp": "2024-01-15T11:00:00+07:00",
  "request_id": "req_abc123xyz"
}
```

---

## 5. Delete Schedule

**Method:** `DELETE`  
**URL:** `{{base_url}}/api/v1/schedules/:id`  
**Auth:** Bearer Token

### Path Parameters:
- `id` (string, required) - Schedule ID

### Example Request:
```
DELETE {{base_url}}/api/v1/schedules/schedule_abc123
```

### Success Response (200):
```json
{
  "success": true,
  "data": null,
  "timestamp": "2024-01-15T10:30:45+07:00",
  "request_id": "req_abc123xyz"
}
```

---

## 6. Assign Schedule

**Method:** `POST`  
**URL:** `{{base_url}}/api/v1/schedules/:id/assign`  
**Auth:** Bearer Token

### Path Parameters:
- `id` (string, required) - Schedule ID

### Request Body:
```json
{
  "assigned_to": "user_789"
}
```

### Success Response (200):
```json
{
  "success": true,
  "data": {
    "id": "schedule_abc123",
    "title": "Visit to Hospital XYZ",
    "assigned_to": "user_789",
    "assigned_user": {
      "id": "user_789",
      "name": "New Assignee",
      "email": "newassignee@example.com"
    },
    "status": "pending",
    "created_at": "2024-01-15T10:00:00+07:00",
    "updated_at": "2024-01-15T11:00:00+07:00"
  },
  "timestamp": "2024-01-15T11:00:00+07:00",
  "request_id": "req_abc123xyz"
}
```

---

## 7. Bulk Assign Schedules

**Method:** `POST`  
**URL:** `{{base_url}}/api/v1/schedules/bulk-assign`  
**Auth:** Bearer Token

### Request Body:
```json
{
  "schedule_ids": [
    "schedule_abc123",
    "schedule_def456",
    "schedule_ghi789"
  ],
  "user_ids": [
    "user_111",
    "user_222"
  ]
}
```

### Success Response (200):
```json
{
  "success": true,
  "data": {
    "assigned_count": 6,
    "message": "Successfully assigned 3 schedules to 2 users"
  },
  "timestamp": "2024-01-15T11:00:00+07:00",
  "request_id": "req_abc123xyz"
}
```

### Error Response (400):
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "field_errors": [
      {
        "field": "schedule_ids",
        "message": "At least one schedule is required"
      },
      {
        "field": "user_ids",
        "message": "At least one user is required"
      }
    ]
  },
  "timestamp": "2024-01-15T11:00:00+07:00",
  "request_id": "req_abc123xyz"
}
```

---

## 8. Check Conflicts

**Method:** `GET`  
**URL:** `{{base_url}}/api/v1/schedules/conflicts`  
**Auth:** Bearer Token

### Query Parameters:
- `user_id` (string, required) - User ID to check conflicts for
- `start_time` (string, required) - Start time in ISO 8601 format
- `end_time` (string, required) - End time in ISO 8601 format
- `exclude_schedule_id` (string, optional) - Exclude specific schedule from conflict check (useful when updating)

### Example Request:
```
GET {{base_url}}/api/v1/schedules/conflicts?user_id=user_123&start_time=2024-01-20T09:00:00Z&end_time=2024-01-20T11:00:00Z
```

### Success Response (200) - No Conflicts:
```json
{
  "success": true,
  "data": {
    "has_conflict": false,
    "conflicts": [],
    "requested_time": {
      "start": "2024-01-20T09:00:00+07:00",
      "end": "2024-01-20T11:00:00+07:00"
    }
  },
  "timestamp": "2024-01-15T11:00:00+07:00",
  "request_id": "req_abc123xyz"
}
```

### Success Response (200) - With Conflicts:
```json
{
  "success": true,
  "data": {
    "has_conflict": true,
    "conflicts": [
      {
        "schedule_id": "schedule_xyz789",
        "title": "Existing Meeting",
        "start_time": "2024-01-20T08:30:00+07:00",
        "end_time": "2024-01-20T10:00:00+07:00",
        "overlap_duration": 3600
      },
      {
        "schedule_id": "schedule_abc456",
        "title": "Another Visit",
        "start_time": "2024-01-20T10:30:00+07:00",
        "end_time": "2024-01-20T12:00:00+07:00",
        "overlap_duration": 1800
      }
    ],
    "requested_time": {
      "start": "2024-01-20T09:00:00+07:00",
      "end": "2024-01-20T11:00:00+07:00"
    }
  },
  "timestamp": "2024-01-15T11:00:00+07:00",
  "request_id": "req_abc123xyz"
}
```

---

## 9. Approve Schedule

**Method:** `POST`  
**URL:** `{{base_url}}/api/v1/schedules/:id/approve`  
**Auth:** Bearer Token

### Path Parameters:
- `id` (string, required) - Schedule ID

### Example Request:
```
POST {{base_url}}/api/v1/schedules/schedule_abc123/approve
```

### Success Response (200):
```json
{
  "success": true,
  "data": {
    "id": "schedule_abc123",
    "title": "Visit to Hospital XYZ",
    "status": "approved",
    "created_at": "2024-01-15T10:00:00+07:00",
    "updated_at": "2024-01-15T11:00:00+07:00"
  },
  "timestamp": "2024-01-15T11:00:00+07:00",
  "request_id": "req_abc123xyz"
}
```

### Error Response (403):
```json
{
  "success": false,
  "error": {
    "code": "PERMISSION_DENIED",
    "message": "You don't have permission to approve schedules"
  },
  "timestamp": "2024-01-15T11:00:00+07:00",
  "request_id": "req_abc123xyz"
}
```

---

## 10. Reject Schedule

**Method:** `POST`  
**URL:** `{{base_url}}/api/v1/schedules/:id/reject`  
**Auth:** Bearer Token

### Path Parameters:
- `id` (string, required) - Schedule ID

### Example Request:
```
POST {{base_url}}/api/v1/schedules/schedule_abc123/reject
```

### Success Response (200):
```json
{
  "success": true,
  "data": {
    "id": "schedule_abc123",
    "title": "Visit to Hospital XYZ",
    "status": "rejected",
    "created_at": "2024-01-15T10:00:00+07:00",
    "updated_at": "2024-01-15T11:00:00+07:00"
  },
  "timestamp": "2024-01-15T11:00:00+07:00",
  "request_id": "req_abc123xyz"
}
```

---

## Environment Variables

Make sure your Postman environment has these variables set:

```
base_url: http://localhost:8080
token: (auto-set after login)
user_id: (auto-set after login)
```

## Common Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `VALIDATION_ERROR` | 400 | Request validation failed |
| `UNAUTHORIZED` | 401 | Authentication required or token invalid |
| `PERMISSION_DENIED` | 403 | User doesn't have required permission |
| `SCHEDULE_NOT_FOUND` | 404 | Schedule with given ID not found |
| `CONFLICT_DETECTED` | 409 | Schedule conflicts with existing schedules |
| `INTERNAL_ERROR` | 500 | Internal server error |

## Required Permissions

- **View Schedules**: `VIEW_SCHEDULES` or `VIEW_ALL_SCHEDULES`
- **Create Schedule**: `CREATE_SCHEDULES`
- **Update Schedule**: `EDIT_SCHEDULES`
- **Delete Schedule**: `DELETE_SCHEDULES`
- **Assign Schedule**: `ASSIGN_SCHEDULES`
- **Approve Schedule**: `APPROVE_SCHEDULES`
- **Reject Schedule**: `REJECT_SCHEDULES`

---

## Notes

1. All datetime fields should be in ISO 8601 format
2. Recurring pattern is only validated when `is_recurring` is `true`
3. For weekly recurring, `days_of_week` is required (0=Sunday, 1=Monday, ..., 6=Saturday)
4. For monthly recurring, `day_of_month` is required (1-31)
5. Either `end_date` or `occurrences` should be provided for recurring schedules
6. Bulk assign will create schedule_assignments for each combination of schedule and user

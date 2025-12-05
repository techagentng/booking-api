# Staff Management API Documentation

**Base URL:** `http://localhost:8080/api/v1`

**Authentication:** All endpoints require Bearer token in Authorization header.
```
Authorization: Bearer <token>
```

---

## Table of Contents

1. [List Staff](#list-staff)
2. [Get Staff by ID](#get-staff-by-id)
3. [Get Staff by Department](#get-staff-by-department)
4. [Get Staff Statistics](#get-staff-statistics)
5. [Get On-Duty Staff](#get-on-duty-staff)
6. [Get Available Staff](#get-available-staff)
7. [Create Staff](#create-staff)
8. [Update Staff](#update-staff)
9. [Delete Staff](#delete-staff)
10. [Clock In](#clock-in)
11. [Clock Out](#clock-out)
12. [Set Availability](#set-availability)
13. [Auto-Assign Staff to Request](#auto-assign-staff-to-request)
14. [Manual Assign Staff to Request](#manual-assign-staff-to-request)
15. [Complete Service Request](#complete-service-request)

---

## List Staff

**Endpoint:** `GET /staff`

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| page | int | No | Page number (default: 1) |
| page_size | int | No | Items per page (default: 20, max: 100) |
| department | string | No | Filter by department |
| status | string | No | Filter by status: `active`, `inactive`, `on_leave` |
| shift | string | No | Filter by shift: `morning`, `afternoon`, `night` |

**Request:**
```http
GET /api/v1/staff?page=1&page_size=20&department=housekeeping&status=active
```

**Response (200 OK):**
```json
{
  "message": "Staff retrieved",
  "data": {
    "staff": [
      {
        "id": 1,
        "employee_id": "HK-001",
        "first_name": "Maria",
        "last_name": "Garcia",
        "email": "maria.garcia@hotel.com",
        "phone": "+1234567001",
        "department": "housekeeping",
        "position": "Housekeeping Supervisor",
        "status": "active",
        "shift": "morning",
        "hire_date": "2023-12-05T00:00:00Z",
        "profile_image": "",
        "is_on_duty": true,
        "is_available": true,
        "current_task_id": null,
        "last_assigned_at": "2024-12-05T10:30:00Z",
        "clock_in_time": "2024-12-05T08:00:00Z",
        "clock_out_time": null,
        "shift_start_time": "08:00",
        "shift_end_time": "16:00",
        "tasks_today": 3,
        "tasks_completed": 2,
        "created_at": "2024-12-05T00:00:00Z",
        "updated_at": "2024-12-05T10:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 18,
      "total_pages": 1
    }
  },
  "errors": ""
}
```

---

## Get Staff by ID

**Endpoint:** `GET /staff/:id`

**Request:**
```http
GET /api/v1/staff/1
```

**Response (200 OK):**
```json
{
  "message": "Staff retrieved",
  "data": {
    "id": 1,
    "employee_id": "HK-001",
    "first_name": "Maria",
    "last_name": "Garcia",
    "email": "maria.garcia@hotel.com",
    "phone": "+1234567001",
    "department": "housekeeping",
    "position": "Housekeeping Supervisor",
    "status": "active",
    "shift": "morning",
    "hire_date": "2023-12-05T00:00:00Z",
    "is_on_duty": true,
    "is_available": true,
    "current_task_id": null,
    "last_assigned_at": "2024-12-05T10:30:00Z",
    "clock_in_time": "2024-12-05T08:00:00Z",
    "tasks_today": 3,
    "tasks_completed": 2
  },
  "errors": ""
}
```

**Response (404 Not Found):**
```json
{
  "message": "Staff not found",
  "data": null,
  "errors": "record not found"
}
```

---

## Get Staff by Department

**Endpoint:** `GET /staff/department/:department`

**Valid Departments:** `housekeeping`, `maintenance`, `front_desk`, `room_service`, `management`

**Request:**
```http
GET /api/v1/staff/department/housekeeping
```

**Response (200 OK):**
```json
{
  "message": "Staff retrieved",
  "data": [
    {
      "id": 1,
      "employee_id": "HK-001",
      "first_name": "Maria",
      "last_name": "Garcia",
      "department": "housekeeping",
      "position": "Housekeeping Supervisor",
      "status": "active",
      "shift": "morning",
      "is_on_duty": true,
      "is_available": true
    },
    {
      "id": 2,
      "employee_id": "HK-002",
      "first_name": "Ana",
      "last_name": "Martinez",
      "department": "housekeeping",
      "position": "Housekeeper",
      "status": "active",
      "shift": "morning",
      "is_on_duty": false,
      "is_available": true
    }
  ],
  "errors": ""
}
```

---

## Get Staff Statistics

**Endpoint:** `GET /staff/stats`

**Request:**
```http
GET /api/v1/staff/stats
```

**Response (200 OK):**
```json
{
  "message": "Staff stats retrieved",
  "data": [
    {
      "department": "housekeeping",
      "total": 4,
      "active": 3,
      "on_leave": 1,
      "morning_shift": 2,
      "afternoon_shift": 1,
      "night_shift": 0
    },
    {
      "department": "maintenance",
      "total": 4,
      "active": 4,
      "on_leave": 0,
      "morning_shift": 2,
      "afternoon_shift": 1,
      "night_shift": 1
    },
    {
      "department": "front_desk",
      "total": 4,
      "active": 4,
      "on_leave": 0,
      "morning_shift": 2,
      "afternoon_shift": 1,
      "night_shift": 1
    },
    {
      "department": "room_service",
      "total": 5,
      "active": 5,
      "on_leave": 0,
      "morning_shift": 3,
      "afternoon_shift": 1,
      "night_shift": 1
    },
    {
      "department": "management",
      "total": 2,
      "active": 2,
      "on_leave": 0,
      "morning_shift": 2,
      "afternoon_shift": 0,
      "night_shift": 0
    }
  ],
  "errors": ""
}
```

---

## Get On-Duty Staff

**Endpoint:** `GET /staff/on-duty`

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| department | string | No | Filter by department |

**Request:**
```http
GET /api/v1/staff/on-duty?department=housekeeping
```

**Response (200 OK):**
```json
{
  "message": "On-duty staff retrieved",
  "data": [
    {
      "id": 1,
      "employee_id": "HK-001",
      "first_name": "Maria",
      "last_name": "Garcia",
      "department": "housekeeping",
      "position": "Housekeeping Supervisor",
      "is_on_duty": true,
      "is_available": true,
      "clock_in_time": "2024-12-05T08:00:00Z",
      "tasks_today": 3,
      "tasks_completed": 2
    }
  ],
  "errors": ""
}
```

---

## Get Available Staff

**Endpoint:** `GET /staff/available`

Returns staff who are on-duty AND available (not currently assigned to a task).

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| department | string | No | Filter by department |

**Request:**
```http
GET /api/v1/staff/available?department=maintenance
```

**Response (200 OK):**
```json
{
  "message": "Available staff retrieved",
  "data": [
    {
      "id": 5,
      "employee_id": "MN-001",
      "first_name": "John",
      "last_name": "Smith",
      "department": "maintenance",
      "position": "Maintenance Supervisor",
      "is_on_duty": true,
      "is_available": true,
      "current_task_id": null,
      "tasks_today": 1,
      "tasks_completed": 1
    }
  ],
  "errors": ""
}
```

---

## Create Staff

**Endpoint:** `POST /staff`

**Request Body:**
```json
{
  "employee_id": "HK-005",
  "first_name": "Sofia",
  "last_name": "Hernandez",
  "email": "sofia.hernandez@hotel.com",
  "phone": "+1234567005",
  "department": "housekeeping",
  "position": "Housekeeper",
  "shift": "afternoon",
  "hire_date": "2024-12-05"
}
```

**Required Fields:** `employee_id`, `first_name`, `last_name`, `email`, `department`, `position`

**Response (201 Created):**
```json
{
  "message": "Staff created",
  "data": {
    "id": 19,
    "employee_id": "HK-005",
    "first_name": "Sofia",
    "last_name": "Hernandez",
    "email": "sofia.hernandez@hotel.com",
    "phone": "+1234567005",
    "department": "housekeeping",
    "position": "Housekeeper",
    "status": "active",
    "shift": "afternoon",
    "hire_date": "2024-12-05T00:00:00Z",
    "is_on_duty": false,
    "is_available": true,
    "tasks_today": 0,
    "tasks_completed": 0,
    "created_at": "2024-12-05T10:00:00Z"
  },
  "errors": ""
}
```

**Response (409 Conflict):**
```json
{
  "message": "Employee ID already exists",
  "data": null,
  "errors": ""
}
```

---

## Update Staff

**Endpoint:** `PUT /staff/:id`

**Request Body (all fields optional):**
```json
{
  "first_name": "Sofia Maria",
  "last_name": "Hernandez",
  "email": "sofia.m.hernandez@hotel.com",
  "phone": "+1234567099",
  "department": "housekeeping",
  "position": "Senior Housekeeper",
  "status": "active",
  "shift": "morning",
  "profile_image": "/images/staff/sofia.jpg"
}
```

**Response (200 OK):**
```json
{
  "message": "Staff updated",
  "data": {
    "id": 19,
    "employee_id": "HK-005",
    "first_name": "Sofia Maria",
    "last_name": "Hernandez",
    "email": "sofia.m.hernandez@hotel.com",
    "phone": "+1234567099",
    "department": "housekeeping",
    "position": "Senior Housekeeper",
    "status": "active",
    "shift": "morning",
    "profile_image": "/images/staff/sofia.jpg"
  },
  "errors": ""
}
```

---

## Delete Staff

**Endpoint:** `DELETE /staff/:id`

**Request:**
```http
DELETE /api/v1/staff/19
```

**Response (200 OK):**
```json
{
  "message": "Staff deleted",
  "data": null,
  "errors": ""
}
```

---

## Clock In

**Endpoint:** `POST /staff/:id/clock-in`

Marks staff as on-duty and available.

**Request:**
```http
POST /api/v1/staff/1/clock-in
```

**Response (200 OK):**
```json
{
  "message": "Staff clocked in",
  "data": {
    "id": 1,
    "employee_id": "HK-001",
    "first_name": "Maria",
    "last_name": "Garcia",
    "is_on_duty": true,
    "is_available": true,
    "clock_in_time": "2024-12-05T08:00:00Z",
    "clock_out_time": null
  },
  "errors": ""
}
```

---

## Clock Out

**Endpoint:** `POST /staff/:id/clock-out`

Marks staff as off-duty. Staff must not have an active task.

**Request:**
```http
POST /api/v1/staff/1/clock-out
```

**Response (200 OK):**
```json
{
  "message": "Staff clocked out",
  "data": {
    "id": 1,
    "employee_id": "HK-001",
    "first_name": "Maria",
    "last_name": "Garcia",
    "is_on_duty": false,
    "is_available": false,
    "clock_in_time": "2024-12-05T08:00:00Z",
    "clock_out_time": "2024-12-05T16:00:00Z"
  },
  "errors": ""
}
```

**Response (400 Bad Request):**
```json
{
  "message": "Staff has an active task, complete it first",
  "data": null,
  "errors": ""
}
```

---

## Set Availability

**Endpoint:** `PUT /staff/:id/availability`

Toggle staff availability (must be on-duty).

**Request Body:**
```json
{
  "available": false
}
```

**Response (200 OK):**
```json
{
  "message": "Availability updated",
  "data": {
    "id": 1,
    "employee_id": "HK-001",
    "first_name": "Maria",
    "last_name": "Garcia",
    "is_on_duty": true,
    "is_available": false
  },
  "errors": ""
}
```

---

## Auto-Assign Staff to Request

**Endpoint:** `POST /service-requests/:id/auto-assign`

Automatically assigns the best available staff to a service request based on:
1. Department match (housekeeping/maintenance)
2. On-duty and available status
3. Workload (staff with fewer tasks today gets priority)
4. Round-robin (oldest last_assigned_at)

**Request:**
```http
POST /api/v1/service-requests/1/auto-assign
```

**Response (200 OK):**
```json
{
  "message": "Staff assigned successfully",
  "data": {
    "request_id": 1,
    "staff": {
      "id": 1,
      "employee_id": "HK-001",
      "first_name": "Maria",
      "last_name": "Garcia",
      "department": "housekeeping",
      "position": "Housekeeping Supervisor",
      "is_on_duty": true,
      "is_available": false,
      "current_task_id": 1,
      "tasks_today": 4
    }
  },
  "errors": ""
}
```

**Response (400 Bad Request):**
```json
{
  "message": "Request is already assigned",
  "data": null,
  "errors": ""
}
```

**Response (404 Not Found):**
```json
{
  "message": "No available staff for this department",
  "data": null,
  "errors": ""
}
```

---

## Manual Assign Staff to Request

**Endpoint:** `POST /service-requests/:id/assign`

Manually assign a specific staff member to a service request.

**Request Body:**
```json
{
  "staff_id": 2
}
```

**Response (200 OK):**
```json
{
  "message": "Staff assigned successfully",
  "data": {
    "request_id": 1,
    "staff": {
      "id": 2,
      "employee_id": "HK-002",
      "first_name": "Ana",
      "last_name": "Martinez",
      "department": "housekeeping",
      "is_on_duty": true,
      "is_available": false,
      "current_task_id": 1
    }
  },
  "errors": ""
}
```

**Response (400 Bad Request):**
```json
{
  "message": "Staff is not on duty",
  "data": null,
  "errors": ""
}
```

```json
{
  "message": "Staff is not available",
  "data": null,
  "errors": ""
}
```

---

## Complete Service Request

**Endpoint:** `POST /service-requests/:id/complete`

Marks a service request as completed and frees up the assigned staff.

**Request Body (optional):**
```json
{
  "completed_by": "Maria Garcia"
}
```

**Response (200 OK):**
```json
{
  "message": "Request completed",
  "data": null,
  "errors": ""
}
```

---

## Staff Model Reference

```json
{
  "id": 1,
  "employee_id": "HK-001",
  "first_name": "Maria",
  "last_name": "Garcia",
  "email": "maria.garcia@hotel.com",
  "phone": "+1234567001",
  "department": "housekeeping",
  "position": "Housekeeping Supervisor",
  "status": "active",
  "shift": "morning",
  "hire_date": "2023-12-05T00:00:00Z",
  "profile_image": "/images/staff/maria.jpg",
  "is_on_duty": true,
  "is_available": true,
  "current_task_id": null,
  "last_assigned_at": "2024-12-05T10:30:00Z",
  "clock_in_time": "2024-12-05T08:00:00Z",
  "clock_out_time": null,
  "shift_start_time": "08:00",
  "shift_end_time": "16:00",
  "tasks_today": 3,
  "tasks_completed": 2,
  "created_at": "2024-12-05T00:00:00Z",
  "updated_at": "2024-12-05T10:30:00Z"
}
```

### Field Descriptions

| Field | Type | Description |
|-------|------|-------------|
| id | uint | Unique identifier |
| employee_id | string | Employee ID (e.g., "HK-001") |
| first_name | string | First name |
| last_name | string | Last name |
| email | string | Email address (unique) |
| phone | string | Phone number |
| department | string | Department: `housekeeping`, `maintenance`, `front_desk`, `room_service`, `management` |
| position | string | Job position/title |
| status | string | Employment status: `active`, `inactive`, `on_leave` |
| shift | string | Assigned shift: `morning`, `afternoon`, `night` |
| hire_date | datetime | Date of hire |
| profile_image | string | URL to profile image |
| is_on_duty | bool | Currently clocked in |
| is_available | bool | Available for task assignment |
| current_task_id | uint | ID of current service request (null if none) |
| last_assigned_at | datetime | Last time a task was assigned |
| clock_in_time | datetime | Time of last clock-in |
| clock_out_time | datetime | Time of last clock-out |
| shift_start_time | string | Shift start time (24h format) |
| shift_end_time | string | Shift end time (24h format) |
| tasks_today | int | Number of tasks assigned today |
| tasks_completed | int | Number of tasks completed today |

---

## Workflow Example

### 1. Staff Clocks In
```http
POST /api/v1/staff/1/clock-in
```

### 2. Guest Creates Service Request
```http
POST /api/v1/services/housekeeping
{
  "room_id": 5,
  "guest_id": 1,
  "service_type": "cleaning",
  "priority": "normal",
  "notes": "Please clean the bathroom"
}
```

### 3. Auto-Assign Staff
```http
POST /api/v1/service-requests/1/auto-assign
```

### 4. Staff Completes Task
```http
POST /api/v1/service-requests/1/complete
{
  "completed_by": "Maria Garcia"
}
```

### 5. Staff Clocks Out
```http
POST /api/v1/staff/1/clock-out
```

---

## Error Codes

| Status Code | Description |
|-------------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request - Invalid input or business rule violation |
| 401 | Unauthorized - Missing or invalid token |
| 404 | Not Found - Resource doesn't exist |
| 409 | Conflict - Duplicate employee ID or email |
| 500 | Internal Server Error |

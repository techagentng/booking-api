# Hotel Management System

A comprehensive hotel management backend API built with Go, featuring real-time notifications, guest management, reservations, room service, and staff coordination.

---

## Quick Start

### Prerequisites

- **Go** 1.21+
- **PostgreSQL** 14+

### Installation

```bash
# Clone the repository
git clone <repository-url>
cd Hotel

# Install dependencies
go mod download

# Set up environment variables
cp .env.example .env
# Edit .env with your database credentials
```

### Environment Variables

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=hotel
JWT_SECRET=your_jwt_secret
PORT=8080
```

### Run the Application

```bash
# Start the server
go run main.go

# Server runs on http://localhost:8080
```

### Build for Production

```bash
go build -o hotel-server main.go
./hotel-server
```

---

## Features

### 🔐 Authentication
- User registration and login
- JWT-based authentication
- Google OAuth integration
- Token blacklisting for logout

### 👥 Guest Management
- Create, update, delete guests
- Guest preferences tracking (room floors, meal types, special requests)
- AI-generated guest insights
- Guest history and statistics
- Service usage tracking

### 🛏️ Reservations
- Create and manage bookings
- Check-in / Check-out processing
- Reservation status tracking (pending, confirmed, checked-in, checked-out, cancelled)
- Dashboard statistics
- Recent activity feed
- Payment details

### 🚪 Room Management
- Room inventory management
- Room type categorization (Standard, Deluxe, Suite)
- Availability checking
- Room status updates (available, occupied, cleaning, maintenance)
- Room statistics

### 🍽️ Room Service
- Menu item management with categories
- Order creation and tracking
- Order status updates (pending, preparing, ready, delivered, cancelled)
- Active orders monitoring
- Room service statistics

### 🧹 Service Requests
- Housekeeping requests
- Maintenance requests
- Auto-assignment to available staff
- Manual staff assignment
- Request completion tracking

### 👨‍💼 Staff Management
- Staff profiles with departments
- Clock-in / Clock-out tracking
- Availability management
- On-duty staff listing
- Task assignment and completion

### 🔔 Real-Time Notifications (SSE)
- Server-Sent Events for push notifications
- Notification types:
  - `sign_in` / `sign_up` - User authentication
  - `new_guest` - Guest registration
  - `new_booking` - Reservation created
  - `check_in` / `check_out` - Guest movements
  - `service_request` - Housekeeping/maintenance
  - `service_completed` - Request fulfilled

---

## API Endpoints

### Authentication
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/signup` | Register new user |
| POST | `/api/v1/auth/login` | User login |
| POST | `/api/v1/auth/logout` | User logout |
| POST | `/google/user/login` | Google OAuth login |

### Guests
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/guests` | List guests (paginated) |
| GET | `/api/v1/guests/:id` | Get guest details |
| POST | `/api/v1/guests` | Create guest |
| PUT | `/api/v1/guests/:id` | Update guest |
| DELETE | `/api/v1/guests/:id` | Delete guest |
| GET | `/api/v1/guests/:id/history` | Guest history |
| GET | `/api/v1/guests/:id/preferences` | Guest preferences |
| GET | `/api/v1/guests/:id/ai-insights` | AI insights |

### Reservations
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/reservations` | List reservations |
| GET | `/api/v1/reservations/:id` | Get reservation |
| POST | `/api/v1/reservations` | Create reservation |
| PUT | `/api/v1/reservations/:id` | Update reservation |
| DELETE | `/api/v1/reservations/:id` | Delete reservation |
| POST | `/api/v1/reservations/:id/checkin` | Check-in guest |
| POST | `/api/v1/reservations/:id/checkout` | Check-out guest |
| GET | `/api/v1/reservations/dashboard` | Dashboard stats |
| GET | `/api/v1/reservations/recent-activity` | Recent bookings & preferences |

### Rooms
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/rooms` | List rooms |
| GET | `/api/v1/rooms/available` | Available rooms |
| GET | `/api/v1/rooms/:id` | Get room |
| POST | `/api/v1/rooms` | Create room |
| PUT | `/api/v1/rooms/:id` | Update room |
| DELETE | `/api/v1/rooms/:id` | Delete room |
| PUT | `/api/v1/rooms/:id/status` | Update room status |

### Room Service
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/room-service/menu` | Get menu items |
| GET | `/api/v1/room-service/menu/categories` | Menu categories |
| POST | `/api/v1/room-service/menu` | Create menu item |
| GET | `/api/v1/room-service/orders` | List orders |
| POST | `/api/v1/room-service/orders` | Create order |
| PUT | `/api/v1/room-service/orders/:id/status` | Update order status |

### Service Requests
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/service-requests` | List requests |
| POST | `/api/v1/services/housekeeping` | Create housekeeping request |
| POST | `/api/v1/services/maintenance` | Create maintenance request |
| POST | `/api/v1/service-requests/:id/assign` | Assign staff |
| POST | `/api/v1/service-requests/:id/complete` | Complete request |

### Staff
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/staff` | List staff |
| GET | `/api/v1/staff/on-duty` | On-duty staff |
| GET | `/api/v1/staff/available` | Available staff |
| POST | `/api/v1/staff` | Create staff |
| PUT | `/api/v1/staff/:id` | Update staff |
| POST | `/api/v1/staff/:id/clock-in` | Clock in |
| POST | `/api/v1/staff/:id/clock-out` | Clock out |

### Notifications
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/notifications/stream` | SSE notification stream |
| GET | `/api/v1/notifications/stats` | Connected clients count |
| POST | `/api/v1/notifications/test` | Send test notification |

---

## Real-Time Notifications

Connect to the SSE stream for live updates:

```javascript
const eventSource = new EventSource('http://localhost:8080/notifications/stream');

eventSource.addEventListener('notification', (event) => {
  const notification = JSON.parse(event.data);
  console.log('New notification:', notification);
});

eventSource.addEventListener('connected', (event) => {
  console.log('Connected to notification stream');
});
```

---

## Project Structure

```
Hotel/
├── main.go                 # Application entry point
├── db/                     # Database repositories
│   ├── auth_repository.go
│   ├── guest_repository.go
│   ├── reservation_repository.go
│   ├── room_repository.go
│   ├── room_service_repository.go
│   ├── guest_service_repository.go
│   └── staff_repository.go
├── models/                 # Data models
│   ├── user.go
│   ├── guest.go
│   ├── reservation.go
│   ├── room.go
│   ├── room_service.go
│   ├── guest_service.go
│   ├── staff.go
│   └── notification.go
├── server/                 # HTTP handlers
│   ├── server.go
│   ├── router.go
│   ├── auth_handlers.go
│   ├── guest_handlers.go
│   ├── reservation_handlers.go
│   ├── room_handlers.go
│   ├── room_service_handlers.go
│   ├── staff_handlers.go
│   ├── tablet_handlers.go
│   └── notification_handlers.go
├── services/               # Business logic
│   ├── auth_service.go
│   ├── guest_service.go
│   ├── reservation_service.go
│   ├── guest_insights_service.go
│   └── notification_hub.go
└── middleware/             # HTTP middleware
    └── auth.go
```

---

## Tech Stack

- **Framework**: Gin (HTTP router)
- **Database**: PostgreSQL with GORM
- **Authentication**: JWT + Google OAuth
- **Real-time**: Server-Sent Events (SSE)

---

## License

MIT

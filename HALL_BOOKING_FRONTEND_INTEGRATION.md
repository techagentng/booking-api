# Hall Booking Frontend Integration Guide

This guide provides complete instructions for integrating the Hall Booking API with your frontend application.

## 🚀 Quick Start

### Base URL
- **Development**: `http://localhost:8080`
- **API Routes**: All endpoints are available at both `/api/v1/hall-bookings/*` and `/hall-bookings/*` for frontend compatibility

### Authentication
All endpoints require JWT authentication (except health check):
```javascript
headers: {
  'Authorization': `Bearer ${token}`,
  'Content-Type': 'application/json'
}
```

## 📋 API Endpoints

### 1. Create Hall Booking
```http
POST /api/v1/hall-bookings
```

**Request Body:**
```javascript
const bookingData = {
  organizer_name: "John Doe",
  organizer_email: "john.doe@example.com", 
  organizer_phone: "+447123456789",
  event_type: "party", // party, wedding, meeting, christening, funeral, corporate, fete, other
  guest_count: 50,
  special_requests: "Extra tables and decorations needed",
  booking_date: "2024-02-15", // Format: yyyy-MM-dd
  start_time: "18:00", // Format: HH:mm (08:00 - 23:00)
  end_time: "22:00",   // Format: HH:mm (08:00 - 23:00)
  total_price: 150.00,
  deposit_required: 100.00,
  payment_method: "cash" // cash, onsite, online
};
```

**Response:**
```javascript
{
  "message": "Hall booking created successfully",
  "data": {
    "id": 123,
    "booking_id": "HB-20240215-001",
    "organizer_name": "John Doe",
    "organizer_email": "john.doe@example.com",
    "organizer_phone": "+447123456789",
    "event_type": "party",
    "guest_count": 50,
    "special_requests": "Extra tables and decorations needed",
    "booking_date": "2024-02-15",
    "start_time": "18:00",
    "end_time": "22:00",
    "total_price": 150.00,
    "deposit_required": 100.00,
    "payment_method": "cash",
    "status": "pending",
    "created_at": "2024-01-29T18:39:00Z",
    "updated_at": "2024-01-29T18:39:00Z",
    "payments": [
      {
        "id": 456,
        "payment_type": "deposit",
        "payment_method": "cash",
        "amount": 100.00,
        "status": "pending",
        "due_date": "2024-02-08T00:00:00Z",
        "created_at": "2024-01-29T18:39:00Z"
      },
      {
        "id": 457,
        "payment_type": "balance",
        "payment_method": "cash",
        "amount": 50.00,
        "status": "pending",
        "due_date": "2024-02-15T00:00:00Z",
        "created_at": "2024-01-29T18:39:00Z"
      }
    ],
    "invoice": {
      "id": 789,
      "invoice_number": "INV-456789",
      "invoice_date": "2024-01-29T00:00:00Z",
      "due_date": "2024-02-08T00:00:00Z",
      "total_amount": 150.00,
      "status": "sent"
    }
  }
}
```

### 2. Get All Hall Bookings (Paginated)
```http
GET /api/v1/hall-bookings?page=1&page_size=10&search=john
```

**Query Parameters:**
- `page`: Page number (default: 1)
- `page_size`: Items per page (default: 10)
- `search`: Search by organizer name, email, or booking ID

**Response:**
```javascript
{
  "message": "Hall bookings retrieved successfully",
  "data": {
    "data": [
      {
        "id": 123,
        "booking_id": "HB-20240215-001",
        // ... all booking fields
      }
    ],
    "meta": {
      "page": 1,
      "page_size": 10,
      "total": 25,
      "total_pages": 3
    }
  }
}
```

### 3. Get Hall Booking by ID
```http
GET /api/v1/hall-bookings/:id
```

### 4. Update Hall Booking
```http
PUT /api/v1/hall-bookings/:id
```

**Request Body (all fields optional):**
```javascript
{
  "organizer_name": "Updated Name",
  "guest_count": 60,
  "status": "confirmed" // pending, confirmed, cancelled, completed
}
```

### 5. Delete Hall Booking
```http
DELETE /api/v1/hall-bookings/:id
```

### 6. Check Hall Availability
```http
GET /api/v1/hall-bookings/availability?date=2024-02-15&start_time=18:00&end_time=22:00
```

**Response:**
```javascript
{
  "message": "Hall availability checked successfully",
  "data": {
    "available": true,
    "date": "2024-02-15",
    "start_time": "18:00",
    "end_time": "22:00"
  }
}
```

### 7. Get Daily Hall Availability
```http
GET /api/v1/hall-bookings/availability/:date
```

**Response:**
```javascript
{
  "message": "Hall availability retrieved successfully",
  "data": {
    "date": "2024-02-15",
    "slots": [
      {
        "start_time": "08:00",
        "end_time": "09:00",
        "available": true
      },
      {
        "start_time": "18:00",
        "end_time": "19:00",
        "available": false
      }
      // ... all hourly slots from 08:00 to 23:00
    ]
  }
}
```

### 8. Get Bookings by Date
```http
GET /api/v1/hall-bookings/date/:date
```

### 9. Get Bookings by Date Range
```http
GET /api/v1/hall-bookings/range?start_date=2024-02-01&end_date=2024-02-29
```

### 10. Update Booking Status
```http
PUT /api/v1/hall-bookings/:id/status
```

**Request Body:**
```javascript
{
  "status": "confirmed" // pending, confirmed, cancelled, completed
}
```

## 🛠️ Frontend Implementation Examples

### React Component Example

```jsx
import React, { useState, useEffect } from 'react';

const HallBookingModal = ({ isOpen, onClose, selectedDate }) => {
  const [formData, setFormData] = useState({
    organizer_name: '',
    organizer_email: '',
    organizer_phone: '',
    event_type: 'party',
    guest_count: 1,
    special_requests: '',
    booking_date: format(selectedDate, 'yyyy-MM-dd'),
    start_time: '18:00',
    end_time: '22:00',
    total_price: 0,
    deposit_required: 0,
    payment_method: 'cash'
  });
  
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [availability, setAvailability] = useState(null);
  const [showInvoiceModal, setShowInvoiceModal] = useState(false);
  const [createdBooking, setCreatedBooking] = useState(null);

  // Check availability when date/time changes
  useEffect(() => {
    if (formData.booking_date && formData.start_time && formData.end_time) {
      checkAvailability();
    }
  }, [formData.booking_date, formData.start_time, formData.end_time]);

  const checkAvailability = async () => {
    try {
      const response = await fetch(
        `/api/v1/hall-bookings/availability?date=${formData.booking_date}&start_time=${formData.start_time}&end_time=${formData.end_time}`,
        {
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
          }
        }
      );
      const result = await response.json();
      setAvailability(result.data.available);
    } catch (error) {
      console.error('Error checking availability:', error);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsSubmitting(true);

    try {
      const response = await fetch('/api/v1/hall-bookings', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify(formData),
      });

      const result = await response.json();

      if (response.ok) {
        setCreatedBooking(result.data);
        setShowInvoiceModal(true);
        onClose();
      } else {
        throw new Error(result.errors || 'Booking failed');
      }
    } catch (error) {
      console.error('Booking failed:', error);
      alert('Booking failed: ' + error.message);
    } finally {
      setIsSubmitting(false);
    }
  };

  const eventTypes = [
    { value: 'party', label: 'Party' },
    { value: 'wedding', label: 'Wedding' },
    { value: 'meeting', label: 'Meeting' },
    { value: 'christening', label: 'Christening' },
    { value: 'funeral', label: 'Funeral' },
    { value: 'corporate', label: 'Corporate Event' },
    { value: 'fete', label: 'Fete' },
    { value: 'other', label: 'Other' }
  ];

  const timeSlots = Array.from({ length: 16 }, (_, i) => {
    const hour = i + 8;
    return `${hour.toString().padStart(2, '0')}:00`;
  });

  return (
    <>
      <div className={`modal ${isOpen ? 'show' : ''}`}>
        <div className="modal-content">
          <h2>Book Hall</h2>
          
          {availability !== null && (
            <div className={`availability-status ${availability ? 'available' : 'unavailable'}`}>
              {availability ? '✅ Slot Available' : '❌ Slot Not Available'}
            </div>
          )}

          <form onSubmit={handleSubmit}>
            <div className="form-row">
              <input
                type="text"
                placeholder="Organizer Name"
                value={formData.organizer_name}
                onChange={(e) => setFormData({...formData, organizer_name: e.target.value})}
                required
                minLength={2}
                maxLength={100}
              />
              <input
                type="email"
                placeholder="Organizer Email"
                value={formData.organizer_email}
                onChange={(e) => setFormData({...formData, organizer_email: e.target.value})}
                required
              />
            </div>

            <div className="form-row">
              <input
                type="tel"
                placeholder="Phone Number"
                value={formData.organizer_phone}
                onChange={(e) => setFormData({...formData, organizer_phone: e.target.value})}
                required
                minLength={10}
              />
              <select
                value={formData.event_type}
                onChange={(e) => setFormData({...formData, event_type: e.target.value})}
                required
              >
                {eventTypes.map(type => (
                  <option key={type.value} value={type.value}>{type.label}</option>
                ))}
              </select>
            </div>

            <div className="form-row">
              <input
                type="number"
                placeholder="Number of Guests"
                value={formData.guest_count}
                onChange={(e) => setFormData({...formData, guest_count: parseInt(e.target.value)})}
                min={1}
                max={100}
                required
              />
              <input
                type="date"
                value={formData.booking_date}
                onChange={(e) => setFormData({...formData, booking_date: e.target.value})}
                min={format(new Date(), 'yyyy-MM-dd')}
                required
              />
            </div>

            <div className="form-row">
              <select
                value={formData.start_time}
                onChange={(e) => setFormData({...formData, start_time: e.target.value})}
                required
              >
                {timeSlots.map(time => (
                  <option key={time} value={time}>{time}</option>
                ))}
              </select>
              <select
                value={formData.end_time}
                onChange={(e) => setFormData({...formData, end_time: e.target.value})}
                required
              >
                {timeSlots.map(time => (
                  <option key={time} value={time}>{time}</option>
                ))}
              </select>
            </div>

            <div className="form-row">
              <input
                type="number"
                placeholder="Total Price"
                value={formData.total_price}
                onChange={(e) => setFormData({...formData, total_price: parseFloat(e.target.value)})}
                min={0}
                step="0.01"
                required
              />
              <input
                type="number"
                placeholder="Deposit Required"
                value={formData.deposit_required}
                onChange={(e) => setFormData({...formData, deposit_required: parseFloat(e.target.value)})}
                min={0}
                max={formData.total_price}
                step="0.01"
                required
              />
            </div>

            <div className="form-row">
              <select
                value={formData.payment_method}
                onChange={(e) => setFormData({...formData, payment_method: e.target.value})}
                required
              >
                <option value="cash">Cash</option>
                <option value="onsite">On-site Payment</option>
                <option value="online">Online Payment</option>
              </select>
            </div>

            <textarea
              placeholder="Special Requests (optional)"
              value={formData.special_requests}
              onChange={(e) => setFormData({...formData, special_requests: e.target.value})}
              maxLength={500}
              rows={3}
            />

            <div className="form-actions">
              <button type="button" onClick={onClose}>Cancel</button>
              <button 
                type="submit" 
                disabled={isSubmitting || availability === false}
              >
                {isSubmitting ? 'Creating Booking...' : 'Create Booking'}
              </button>
            </div>
          </form>
        </div>
      </div>

      {showInvoiceModal && createdBooking && (
        <InvoiceModal 
          booking={createdBooking} 
          onClose={() => setShowInvoiceModal(false)} 
        />
      )}
    </>
  );
};

// Invoice Modal Component
const InvoiceModal = ({ booking, onClose }) => {
  return (
    <div className="modal show">
      <div className="modal-content">
        <h2>Booking Confirmation</h2>
        
        <div className="booking-details">
          <h3>Booking Details</h3>
          <p><strong>Booking ID:</strong> {booking.booking_id}</p>
          <p><strong>Organizer:</strong> {booking.organizer_name}</p>
          <p><strong>Email:</strong> {booking.organizer_email}</p>
          <p><strong>Phone:</strong> {booking.organizer_phone}</p>
          <p><strong>Event Type:</strong> {booking.event_type}</p>
          <p><strong>Guests:</strong> {booking.guest_count}</p>
          <p><strong>Date:</strong> {booking.booking_date}</p>
          <p><strong>Time:</strong> {booking.start_time} - {booking.end_time}</p>
          <p><strong>Status:</strong> <span className={`status ${booking.status}`}>{booking.status}</span></p>
        </div>

        <div className="payment-details">
          <h3>Payment Schedule</h3>
          {booking.payments.map(payment => (
            <div key={payment.id} className="payment-item">
              <p><strong>Type:</strong> {payment.payment_type}</p>
              <p><strong>Amount:</strong> £{payment.amount.toFixed(2)}</p>
              <p><strong>Due Date:</strong> {new Date(payment.due_date).toLocaleDateString()}</p>
              <p><strong>Status:</strong> <span className={`status ${payment.status}`}>{payment.status}</span></p>
            </div>
          ))}
        </div>

        {booking.invoice && (
          <div className="invoice-details">
            <h3>Invoice</h3>
            <p><strong>Invoice Number:</strong> {booking.invoice.invoice_number}</p>
            <p><strong>Total Amount:</strong> £{booking.invoice.total_amount.toFixed(2)}</p>
            <p><strong>Due Date:</strong> {new Date(booking.invoice.due_date).toLocaleDateString()}</p>
            <p><strong>Status:</strong> <span className={`status ${booking.invoice.status}`}>{booking.invoice.status}</span></p>
          </div>
        )}

        <div className="form-actions">
          <button onClick={onClose}>Close</button>
          <button onClick={() => window.print()}>Print Confirmation</button>
        </div>
      </div>
    </div>
  );
};
```

### Hall Availability Calendar Component

```jsx
const HallAvailabilityCalendar = () => {
  const [selectedDate, setSelectedDate] = useState(new Date());
  const [availability, setAvailability] = useState(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (selectedDate) {
      fetchAvailability();
    }
  }, [selectedDate]);

  const fetchAvailability = async () => {
    setLoading(true);
    try {
      const dateStr = format(selectedDate, 'yyyy-MM-dd');
      const response = await fetch(`/api/v1/hall-bookings/availability/${dateStr}`, {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      });
      const result = await response.json();
      setAvailability(result.data);
    } catch (error) {
      console.error('Error fetching availability:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="availability-calendar">
      <h2>Hall Availability</h2>
      
      <div className="date-selector">
        <input
          type="date"
          value={format(selectedDate, 'yyyy-MM-dd')}
          onChange={(e) => setSelectedDate(new Date(e.target.value))}
          min={format(new Date(), 'yyyy-MM-dd')}
        />
      </div>

      {loading ? (
        <div className="loading">Loading availability...</div>
      ) : availability && (
        <div className="time-slots">
          <h3>Available Time Slots - {format(selectedDate, 'MMMM d, yyyy')}</h3>
          <div className="slots-grid">
            {availability.slots.map((slot, index) => (
              <div 
                key={index} 
                className={`time-slot ${slot.available ? 'available' : 'unavailable'}`}
              >
                <span>{slot.start_time} - {slot.end_time}</span>
                <span>{slot.available ? '✅ Available' : '❌ Booked'}</span>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};
```

### Hall Bookings List Component

```jsx
const HallBookingsList = () => {
  const [bookings, setBookings] = useState([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [search, setSearch] = useState('');

  useEffect(() => {
    fetchBookings();
  }, [page, search]);

  const fetchBookings = async () => {
    setLoading(true);
    try {
      const queryParams = new URLSearchParams({
        page: page.toString(),
        page_size: '10',
        ...(search && { search })
      });

      const response = await fetch(`/api/v1/hall-bookings?${queryParams}`, {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      });
      const result = await response.json();
      setBookings(result.data.data);
      setTotalPages(result.data.meta.total_pages);
    } catch (error) {
      console.error('Error fetching bookings:', error);
    } finally {
      setLoading(false);
    }
  };

  const updateBookingStatus = async (bookingId, newStatus) => {
    try {
      const response = await fetch(`/api/v1/hall-bookings/${bookingId}/status`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify({ status: newStatus }),
      });

      if (response.ok) {
        fetchBookings(); // Refresh the list
      }
    } catch (error) {
      console.error('Error updating booking status:', error);
    }
  };

  return (
    <div className="bookings-list">
      <h2>Hall Bookings</h2>
      
      <div className="list-controls">
        <input
          type="text"
          placeholder="Search by name, email, or booking ID..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
      </div>

      {loading ? (
        <div className="loading">Loading bookings...</div>
      ) : (
        <>
          <div className="bookings-grid">
            {bookings.map(booking => (
              <div key={booking.id} className="booking-card">
                <div className="booking-header">
                  <h3>{booking.booking_id}</h3>
                  <span className={`status ${booking.status}`}>{booking.status}</span>
                </div>
                
                <div className="booking-info">
                  <p><strong>Organizer:</strong> {booking.organizer_name}</p>
                  <p><strong>Email:</strong> {booking.organizer_email}</p>
                  <p><strong>Event:</strong> {booking.event_type}</p>
                  <p><strong>Guests:</strong> {booking.guest_count}</p>
                  <p><strong>Date:</strong> {booking.booking_date}</p>
                  <p><strong>Time:</strong> {booking.start_time} - {booking.end_time}</p>
                  <p><strong>Total:</strong> £{booking.total_price.toFixed(2)}</p>
                </div>

                <div className="booking-actions">
                  <select 
                    value={booking.status}
                    onChange={(e) => updateBookingStatus(booking.id, e.target.value)}
                  >
                    <option value="pending">Pending</option>
                    <option value="confirmed">Confirmed</option>
                    <option value="cancelled">Cancelled</option>
                    <option value="completed">Completed</option>
                  </select>
                </div>
              </div>
            ))}
          </div>

          <div className="pagination">
            <button 
              disabled={page === 1}
              onClick={() => setPage(page - 1)}
            >
              Previous
            </button>
            <span>Page {page} of {totalPages}</span>
            <button 
              disabled={page === totalPages}
              onClick={() => setPage(page + 1)}
            >
              Next
            </button>
          </div>
        </>
      )}
    </div>
  );
};
```

## 🎨 CSS Styling

```css
/* Modal Styles */
.modal {
  display: none;
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  z-index: 1000;
}

.modal.show {
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal-content {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  max-width: 600px;
  width: 90%;
  max-height: 90vh;
  overflow-y: auto;
}

/* Form Styles */
.form-row {
  display: flex;
  gap: 1rem;
  margin-bottom: 1rem;
}

.form-row input,
.form-row select,
.form-row textarea {
  flex: 1;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

.form-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 2rem;
}

.form-actions button {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
}

.form-actions button[type="submit"] {
  background: #007bff;
  color: white;
}

.form-actions button[type="button"] {
  background: #6c757d;
  color: white;
}

/* Availability Status */
.availability-status {
  padding: 1rem;
  border-radius: 4px;
  margin-bottom: 1rem;
  text-align: center;
  font-weight: bold;
}

.availability-status.available {
  background: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

.availability-status.unavailable {
  background: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}

/* Status Badges */
.status {
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.875rem;
  font-weight: bold;
  text-transform: uppercase;
}

.status.pending {
  background: #fff3cd;
  color: #856404;
}

.status.confirmed {
  background: #d4edda;
  color: #155724;
}

.status.cancelled {
  background: #f8d7da;
  color: #721c24;
}

.status.completed {
  background: #d1ecf1;
  color: #0c5460;
}

/* Time Slots */
.slots-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 1rem;
  margin-top: 1rem;
}

.time-slot {
  padding: 1rem;
  border-radius: 4px;
  text-align: center;
  border: 1px solid;
}

.time-slot.available {
  background: #d4edda;
  border-color: #c3e6cb;
  color: #155724;
}

.time-slot.unavailable {
  background: #f8d7da;
  border-color: #f5c6cb;
  color: #721c24;
}

/* Booking Cards */
.bookings-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
  margin-top: 2rem;
}

.booking-card {
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 1.5rem;
  background: white;
}

.booking-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.booking-info p {
  margin-bottom: 0.5rem;
}

.booking-actions {
  margin-top: 1rem;
}

/* Pagination */
.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1rem;
  margin-top: 2rem;
}

.pagination button {
  padding: 0.5rem 1rem;
  border: 1px solid #ddd;
  background: white;
  border-radius: 4px;
  cursor: pointer;
}

.pagination button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Loading */
.loading {
  text-align: center;
  padding: 2rem;
  font-style: italic;
  color: #666;
}
```

## 🔧 Error Handling

### Common Error Responses

```javascript
// Validation Error (400)
{
  "message": "Failed to create hall booking",
  "errors": "booking date must be in the future"
}

// Not Found (404)
{
  "message": "Failed to fetch hall booking",
  "errors": "hall booking not found"
}

// Conflict (409) - Slot not available
{
  "message": "Failed to create hall booking", 
  "errors": "hall is not available for the requested time slot"
}

// Unauthorized (401)
{
  "message": "Unauthorized",
  "errors": "invalid or missing token"
}
```

### Error Handling Example

```javascript
const handleApiCall = async (url, options = {}) => {
  try {
    const response = await fetch(url, {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
        ...options.headers
      },
      ...options
    });

    const result = await response.json();

    if (!response.ok) {
      // Handle different error types
      if (response.status === 401) {
        // Redirect to login
        window.location.href = '/login';
        return;
      }
      
      if (response.status === 409) {
        // Show availability conflict
        alert('This time slot is already booked. Please choose a different time.');
        return;
      }

      throw new Error(result.errors || 'Request failed');
    }

    return result;
  } catch (error) {
    console.error('API Error:', error);
    alert(error.message || 'An unexpected error occurred');
    throw error;
  }
};
```

## 📱 Mobile Responsive Considerations

```css
/* Mobile Styles */
@media (max-width: 768px) {
  .form-row {
    flex-direction: column;
  }
  
  .slots-grid {
    grid-template-columns: 1fr;
  }
  
  .bookings-grid {
    grid-template-columns: 1fr;
  }
  
  .modal-content {
    width: 95%;
    padding: 1rem;
  }
}
```

## 🚀 Integration Checklist

- [ ] Add JWT token to all API calls
- [ ] Implement form validation matching backend rules
- [ ] Add loading states for better UX
- [ ] Handle error responses appropriately
- [ ] Add confirmation dialogs for destructive actions
- [ ] Implement real-time availability checking
- [ ] Add pagination for booking lists
- [ ] Style components to match your design system
- [ ] Test all CRUD operations
- [ ] Add accessibility features

## 📞 Support

If you encounter any issues during integration:

1. Check the browser console for error messages
2. Verify the API is running on `http://localhost:8080`
3. Ensure JWT token is valid and included in headers
4. Check network requests in browser dev tools

The Hall Booking API is now ready for seamless frontend integration!

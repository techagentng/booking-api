# Calendar Availability Backend Plan

## 🎯 **Objective**
Create a calendar availability system that:
1. **Fetches available dates** from the backend
2. **Automatically updates** when bookings are made
3. **Provides real-time availability** for the frontend calendar
4. **Maintains data consistency** across all operations

---

## 🏗️ **System Architecture**

### **Data Flow**
```
Frontend Calendar → Backend API → Database
       ↓                ↓              ↓
  Check Availability → Check Slots → HallBookings Table
       ↓                ↓              ↓
  Display Calendar → Return Dates → Update on Booking
```

### **Core Components**
1. **Calendar Availability API** - New endpoints for calendar data
2. **Real-time Status Updates** - Automatic availability updates
3. **Date Range Queries** - Efficient calendar data fetching
4. **Status Management** - Booking status affects availability

---

## 📋 **Implementation Plan**

### **Phase 1: Calendar API Endpoints**

#### **1. Get Calendar Availability (Date Range)**
```http
GET /api/v1/hall-bookings/calendar?start_date=2024-01-01&end_date=2024-12-31
```

**Response:**
```json
{
  "message": "Calendar availability retrieved successfully",
  "data": {
    "start_date": "2024-01-01",
    "end_date": "2024-12-31",
    "dates": [
      {
        "date": "2024-01-15",
        "status": "available",
        "available_slots": 8,
        "total_slots": 16,
        "bookings_count": 2
      },
      {
        "date": "2024-01-16",
        "status": "limited",
        "available_slots": 3,
        "total_slots": 16,
        "bookings_count": 5
      },
      {
        "date": "2024-01-17",
        "status": "unavailable",
        "available_slots": 0,
        "total_slots": 16,
        "bookings_count": 8
      }
    ]
  }
}
```

#### **2. Get Single Date Availability**
```http
GET /api/v1/hall-bookings/calendar/2024-01-15
```

**Response:**
```json
{
  "message": "Date availability retrieved successfully",
  "data": {
    "date": "2024-01-15",
    "status": "available",
    "available_slots": 8,
    "total_slots": 16,
    "bookings_count": 2,
    "slots": [
      {
        "start_time": "08:00",
        "end_time": "09:00",
        "available": true,
        "booking_id": null
      },
      {
        "start_time": "18:00",
        "end_time": "22:00",
        "available": false,
        "booking_id": "HB-20240115-001"
      }
    ]
  }
}
```

#### **3. Get Month Overview**
```http
GET /api/v1/hall-bookings/calendar/2024/01
```

**Response:**
```json
{
  "message": "Month overview retrieved successfully",
  "data": {
    "month": "2024-01",
    "year": 2024,
    "summary": {
      "total_days": 31,
      "available_days": 25,
      "limited_days": 4,
      "unavailable_days": 2
    },
    "dates": [
      {
        "date": "2024-01-01",
        "status": "available",
        "bookings_count": 0
      }
      // ... all days of the month
    ]
  }
}
```

---

### **Phase 2: Status Management Logic**

#### **Availability Status Rules**
```go
type AvailabilityStatus string

const (
    StatusAvailable   AvailabilityStatus = "available"   // 75%+ slots free
    StatusLimited     AvailabilityStatus = "limited"     // 25-74% slots free
    StatusUnavailable AvailabilityStatus = "unavailable" // <25% slots free
    StatusClosed      AvailabilityStatus = "closed"      // 0 slots free
)

func calculateAvailabilityStatus(availableSlots, totalSlots int) AvailabilityStatus {
    if availableSlots == 0 {
        return StatusClosed
    }
    
    percentage := float64(availableSlots) / float64(totalSlots) * 100
    
    if percentage >= 75 {
        return StatusAvailable
    } else if percentage >= 25 {
        return StatusLimited
    } else {
        return StatusUnavailable
    }
}
```

#### **Booking Status Impact**
```go
// How booking status affects availability
func updateAvailabilityForBooking(booking *HallBooking, oldStatus, newStatus string) {
    switch newStatus {
    case "confirmed":
        // Mark slots as booked
        markSlotsAsBooked(booking.BookingDate, booking.StartTime, booking.EndTime)
    case "cancelled":
        // Release slots back to available
        releaseSlots(booking.BookingDate, booking.StartTime, booking.EndTime)
    case "pending":
        // Tentatively mark slots (optional - depends on business rules)
        markSlotsAsPending(booking.BookingDate, booking.StartTime, booking.EndTime)
    }
}
```

---

### **Phase 3: Database Enhancements**

#### **New Calendar View Table (Optional - for performance)**
```sql
CREATE TABLE hall_calendar_availability (
    id SERIAL PRIMARY KEY,
    date DATE NOT NULL UNIQUE,
    total_slots INTEGER DEFAULT 16, -- 8 slots × 2 hours each
    available_slots INTEGER DEFAULT 16,
    pending_slots INTEGER DEFAULT 0,
    booked_slots INTEGER DEFAULT 0,
    status VARCHAR(20) DEFAULT 'available',
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_hall_calendar_date ON hall_calendar_availability(date);
CREATE INDEX idx_hall_calendar_status ON hall_calendar_availability(status);
```

#### **Trigger for Auto-Updates**
```sql
-- Trigger to update calendar when booking changes
CREATE OR REPLACE FUNCTION update_calendar_availability()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' OR (TG_OP = 'UPDATE' AND OLD.status != NEW.status) THEN
        -- Recalculate availability for the booking date
        INSERT INTO hall_calendar_availability (date, available_slots, booked_slots, status)
        VALUES (
            NEW.booking_date,
            (SELECT 16 - COUNT(*) FROM hall_bookings 
             WHERE booking_date = NEW.booking_date 
             AND status IN ('confirmed', 'pending')
             AND deleted_at IS NULL),
            (SELECT COUNT(*) FROM hall_bookings 
             WHERE booking_date = NEW.booking_date 
             AND status = 'confirmed'
             AND deleted_at IS NULL),
            'calculated'
        )
        ON CONFLICT (date) 
        DO UPDATE SET 
            available_slots = EXCLUDED.available_slots,
            booked_slots = EXCLUDED.booked_slots,
            status = EXCLUDED.status,
            last_updated = CURRENT_TIMESTAMP;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_calendar_availability
    AFTER INSERT OR UPDATE ON hall_bookings
    FOR EACH ROW
    EXECUTE FUNCTION update_calendar_availability();
```

---

### **Phase 4: Repository Layer**

#### **New Repository Methods**
```go
// In hall_booking_repository.go
type HallBookingCalendarRepository interface {
    GetCalendarAvailability(startDate, endDate string) ([]CalendarDate, error)
    GetDateAvailability(date string) (*CalendarDateDetail, error)
    GetMonthOverview(year, month int) (*MonthOverview, error)
    UpdateCalendarAvailability(date string) error
}

type CalendarDate struct {
    Date           string `json:"date"`
    Status         string `json:"status"`
    AvailableSlots int    `json:"available_slots"`
    TotalSlots     int    `json:"total_slots"`
    BookingsCount  int    `json:"bookings_count"`
}

type CalendarDateDetail struct {
    Date           string      `json:"date"`
    Status         string      `json:"status"`
    AvailableSlots int        `json:"available_slots"`
    TotalSlots     int        `json:"total_slots"`
    BookingsCount  int        `json:"bookings_count"`
    Slots          []TimeSlot  `json:"slots"`
}

type MonthOverview struct {
    Month          int           `json:"month"`
    Year           int           `json:"year"`
    Summary        MonthSummary  `json:"summary"`
    Dates          []CalendarDate `json:"dates"`
}

type MonthSummary struct {
    TotalDays      int `json:"total_days"`
    AvailableDays  int `json:"available_days"`
    LimitedDays    int `json:"limited_days"`
    UnavailableDays int `json:"unavailable_days"`
}
```

#### **Implementation**
```go
func (r *hallBookingRepository) GetCalendarAvailability(startDate, endDate string) ([]CalendarDate, error) {
    var results []CalendarDate
    
    // Parse dates
    start, _ := time.Parse("2006-01-02", startDate)
    end, _ := time.Parse("2006-01-02", endDate)
    
    // Generate date range
    for d := start; d.Before(end) || d.Equal(end); d = d.AddDate(0, 0, 1) {
        dateStr := d.Format("2006-01-02")
        
        // Count bookings for this date
        var bookingCount int64
        r.db.Model(&HallBooking{}).
            Where("booking_date = ? AND status IN ? AND deleted_at IS NULL", 
                dateStr, []string{"confirmed", "pending"}).
            Count(&bookingCount)
        
        // Calculate availability
        totalSlots := 16 // 8 slots × 2 hours each
        availableSlots := totalSlots - int(bookingCount)
        status := calculateAvailabilityStatus(availableSlots, totalSlots)
        
        results = append(results, CalendarDate{
            Date:           dateStr,
            Status:         string(status),
            AvailableSlots: availableSlots,
            TotalSlots:     totalSlots,
            BookingsCount:  int(bookingCount),
        })
    }
    
    return results, nil
}

func (r *hallBookingRepository) GetDateAvailability(date string) (*CalendarDateDetail, error) {
    // Get all bookings for the date
    var bookings []HallBooking
    r.db.Where("booking_date = ? AND status IN ? AND deleted_at IS NULL", 
        date, []string{"confirmed", "pending"}).
        Find(&bookings)
    
    // Generate all time slots
    var slots []TimeSlot
    for hour := 8; hour <= 22; hour++ {
        startTime := fmt.Sprintf("%02d:00", hour)
        endTime := fmt.Sprintf("%02d:00", hour+1)
        
        slot := TimeSlot{
            StartTime: startTime,
            EndTime:   endTime,
            Available: true,
        }
        
        // Check if slot is booked
        for _, booking := range bookings {
            if r.timeSlotConflicts(startTime, endTime, booking.StartTime, booking.EndTime) {
                slot.Available = false
                slot.BookingID = booking.BookingID
                break
            }
        }
        
        slots = append(slots, slot)
    }
    
    // Calculate summary
    totalSlots := len(slots)
    availableSlots := 0
    for _, slot := range slots {
        if slot.Available {
            availableSlots++
        }
    }
    
    status := calculateAvailabilityStatus(availableSlots, totalSlots)
    
    return &CalendarDateDetail{
        Date:           date,
        Status:         string(status),
        AvailableSlots: availableSlots,
        TotalSlots:     totalSlots,
        BookingsCount:  len(bookings),
        Slots:          slots,
    }, nil
}
```

---

### **Phase 5: Service Layer**

#### **Calendar Service**
```go
type HallBookingCalendarService interface {
    GetCalendarAvailability(startDate, endDate string) ([]CalendarDate, error)
    GetDateAvailability(date string) (*CalendarDateDetail, error)
    GetMonthOverview(year, month int) (*MonthOverview, error)
    RefreshCalendarAvailability(date string) error
}

type hallBookingCalendarService struct {
    repo HallBookingCalendarRepository
}

func (s *hallBookingCalendarService) GetCalendarAvailability(startDate, endDate string) ([]CalendarDate, error) {
    // Validate date range
    if err := s.validateDateRange(startDate, endDate); err != nil {
        return nil, err
    }
    
    // Get availability from repository
    return s.repo.GetCalendarAvailability(startDate, endDate)
}

func (s *hallBookingCalendarService) RefreshCalendarAvailability(date string) error {
    // This method is called when bookings are created/updated/cancelled
    return s.repo.UpdateCalendarAvailability(date)
}
```

---

### **Phase 6: API Handlers**

#### **Calendar Handlers**
```go
// handleGetCalendarAvailability returns calendar availability for date range
func (s *Server) handleGetCalendarAvailability() gin.HandlerFunc {
    return func(c *gin.Context) {
        startDate := c.Query("start_date")
        endDate := c.Query("end_date")
        
        if startDate == "" || endDate == "" {
            response.JSON(c, "Missing required parameters: start_date, end_date", http.StatusBadRequest, nil, nil)
            return
        }
        
        calendarService := services.NewHallBookingCalendarService(s.HallBookingRepository)
        result, err := calendarService.GetCalendarAvailability(startDate, endDate)
        if err != nil {
            log.Printf("handleGetCalendarAvailability: error: %v", err)
            response.JSON(c, "Failed to get calendar availability", http.StatusInternalServerError, nil, err)
            return
        }
        
        response.JSON(c, "Calendar availability retrieved successfully", http.StatusOK, result, nil)
    }
}

// handleGetDateAvailability returns detailed availability for a specific date
func (s *Server) handleGetDateAvailability() gin.HandlerFunc {
    return func(c *gin.Context) {
        date := c.Param("date")
        if date == "" {
            response.JSON(c, "Missing date parameter", http.StatusBadRequest, nil, nil)
            return
        }
        
        calendarService := services.NewHallBookingCalendarService(s.HallBookingRepository)
        result, err := calendarService.GetDateAvailability(date)
        if err != nil {
            log.Printf("handleGetDateAvailability: error: %v", err)
            response.JSON(c, "Failed to get date availability", http.StatusInternalServerError, nil, err)
            return
        }
        
        response.JSON(c, "Date availability retrieved successfully", http.StatusOK, result, nil)
    }
}

// handleGetMonthOverview returns month overview
func (s *Server) handleGetMonthOverview() gin.HandlerFunc {
    return func(c *gin.Context) {
        year, err1 := strconv.Atoi(c.Param("year"))
        month, err2 := strconv.Atoi(c.Param("month"))
        
        if err1 != nil || err2 != nil {
            response.JSON(c, "Invalid year or month parameter", http.StatusBadRequest, nil, nil)
            return
        }
        
        calendarService := services.NewHallBookingCalendarService(s.HallBookingRepository)
        result, err := calendarService.GetMonthOverview(year, month)
        if err != nil {
            log.Printf("handleGetMonthOverview: error: %v", err)
            response.JSON(c, "Failed to get month overview", http.StatusInternalServerError, nil, err)
            return
        }
        
        response.JSON(c, "Month overview retrieved successfully", http.StatusOK, result, nil)
    }
}
```

---

### **Phase 7: Auto-Update Integration**

#### **Update Hall Booking Service**
```go
// In hall_booking_service.go - modify CreateHallBooking
func (s *hallBookingService) CreateHallBooking(req *models.CreateHallBookingRequest) (*models.HallBookingResponse, error) {
    // ... existing validation and creation logic ...
    
    // Create booking
    createdBooking, err := s.repo.CreateHallBooking(booking)
    if err != nil {
        return nil, err
    }
    
    // 🔄 NEW: Refresh calendar availability for the booking date
    calendarService := services.NewHallBookingCalendarService(s.repo)
    if err := calendarService.RefreshCalendarAvailability(booking.BookingDate.Format("2006-01-02")); err != nil {
        // Log error but don't fail the booking
        log.Printf("Warning: Failed to refresh calendar availability: %v", err)
    }
    
    return s.convertToResponse(createdBooking), nil
}

// Similar updates for UpdateHallBooking and DeleteHallBooking
```

---

### **Phase 8: Routes**

#### **Add Calendar Routes**
```go
// In router.go - add to public routes (no auth required)
// Public Calendar routes (no auth required)
v1.GET("/hall-bookings/calendar", s.handleGetCalendarAvailability())
v1.GET("/hall-bookings/calendar/:date", s.handleGetDateAvailability())
v1.GET("/hall-bookings/calendar/:year/:month", s.handleGetMonthOverview())

// Frontend routes
frontend.GET("/hall-bookings/calendar", func(c *gin.Context) {
    c.Request.URL.Path = "/api/v1/hall-bookings/calendar"
    router.HandleContext(c)
})

frontend.GET("/hall-bookings/calendar/:date", func(c *gin.Context) {
    c.Request.URL.Path = "/api/v1/hall-bookings/calendar/" + c.Param("date")
    router.HandleContext(c)
})

frontend.GET("/hall-bookings/calendar/:year/:month", func(c *gin.Context) {
    c.Request.URL.Path = "/api/v1/hall-bookings/calendar/" + c.Param("year") + "/" + c.Param("month")
    router.HandleContext(c)
})
```

---

## 🚀 **Frontend Integration**

### **Calendar Component**
```jsx
const HallBookingCalendar = () => {
    const [calendarData, setCalendarData] = useState([]);
    const [selectedDate, setSelectedDate] = useState(null);
    const [loading, setLoading] = useState(false);

    // Fetch calendar data for month
    const fetchCalendarData = async (year, month) => {
        setLoading(true);
        try {
            const response = await fetch(
                `http://localhost:8080/api/v1/hall-bookings/calendar/${year}/${month}`,
                {
                    headers: {
                        'Content-Type': 'application/json'
                        // No auth needed!
                    }
                }
            );
            const result = await response.json();
            setCalendarData(result.data.dates);
        } catch (error) {
            console.error('Error fetching calendar data:', error);
        } finally {
            setLoading(false);
        }
    };

    // Get detailed availability for specific date
    const getDateDetails = async (date) => {
        try {
            const response = await fetch(
                `http://localhost:8080/api/v1/hall-bookings/calendar/${date}`,
                {
                    headers: {
                        'Content-Type': 'application/json'
                    }
                }
            );
            const result = await response.json();
            return result.data;
        } catch (error) {
            console.error('Error fetching date details:', error);
            return null;
        }
    };

    // Handle date selection
    const handleDateClick = async (date) => {
        const details = await getDateDetails(date);
        setSelectedDate(details);
        
        // Open booking modal if date is available
        if (details.status !== 'closed') {
            setShowBookingModal(true);
        }
    };

    return (
        <div className="hall-calendar">
            <div className="calendar-grid">
                {calendarData.map(day => (
                    <div 
                        key={day.date}
                        className={`calendar-day ${day.status}`}
                        onClick={() => handleDateClick(day.date)}
                    >
                        <div className="day-number">
                            {new Date(day.date).getDate()}
                        </div>
                        <div className="day-status">
                            {day.bookings_count} bookings
                        </div>
                        <div className="availability">
                            {day.available_slots}/{day.total_slots} slots
                        </div>
                    </div>
                ))}
            </div>
            
            {selectedDate && (
                <BookingModal 
                    date={selectedDate.date}
                    availability={selectedDate}
                    onClose={() => setSelectedDate(null)}
                    onBookingCreated={() => fetchCalendarData(currentYear, currentMonth)} // Refresh calendar
                />
            )}
        </div>
    );
};
```

---

## 📊 **Performance Considerations**

### **Optimization Strategies**
1. **Caching**: Cache calendar data for 5-10 minutes
2. **Database Indexing**: Proper indexes on booking_date and status
3. **Batch Processing**: Process multiple dates in single query
4. **Lazy Loading**: Load detailed slots only when needed

### **Cache Implementation**
```go
// Simple in-memory cache (for production, use Redis)
var calendarCache = sync.Map{}
var cacheTimeout = 10 * time.Minute

func (s *hallBookingCalendarService) GetCalendarAvailability(startDate, endDate string) ([]CalendarDate, error) {
    cacheKey := fmt.Sprintf("calendar_%s_%s", startDate, endDate)
    
    // Check cache
    if cached, ok := calendarCache.Load(cacheKey); ok {
        if cacheEntry, ok := cached.(CacheEntry); ok {
            if time.Since(cacheEntry.Timestamp) < cacheTimeout {
                return cacheEntry.Data.([]CalendarDate), nil
            }
        }
    }
    
    // Get from database
    data, err := s.repo.GetCalendarAvailability(startDate, endDate)
    if err != nil {
        return nil, err
    }
    
    // Cache the result
    calendarCache.Store(cacheKey, CacheEntry{
        Data:      data,
        Timestamp: time.Now(),
    })
    
    return data, nil
}
```

---

## ✅ **Implementation Checklist**

### **Backend Tasks**
- [ ] Create calendar repository methods
- [ ] Implement calendar service layer
- [ ] Add calendar API handlers
- [ ] Update booking service to refresh calendar
- [ ] Add calendar routes (public)
- [ ] Add database triggers for auto-updates
- [ ] Implement caching for performance
- [ ] Add comprehensive error handling

### **Frontend Tasks**
- [ ] Create calendar component
- [ ] Implement date selection
- [ ] Add availability indicators
- [ ] Integrate with booking modal
- [ ] Add real-time updates
- [ ] Handle loading states
- [ ] Add responsive design

### **Testing**
- [ ] Test calendar data fetching
- [ ] Test booking creation updates calendar
- [ ] Test booking cancellation updates calendar
- [ ] Test performance with large date ranges
- [ ] Test concurrent bookings

---

## 🎯 **Expected Outcome**

1. **Real-time Calendar**: Frontend shows live availability
2. **Auto-updates**: Calendar refreshes when bookings change
3. **Performance**: Fast loading with caching
4. **User Experience**: Intuitive date selection with visual feedback
5. **Data Consistency**: Calendar always reflects actual booking status

This plan ensures your frontend calendar stays perfectly synchronized with backend booking data in real-time! 🚀

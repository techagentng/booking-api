package db

import (
	"hotel/models"
	"log"
	"time"

	"gorm.io/gorm"
)

// SeedStaff seeds sample staff data
func SeedStaff(db *gorm.DB) error {
	// Check if staff already exists
	var count int64
	db.Model(&models.Staff{}).Count(&count)
	if count > 0 {
		log.Printf("Staff already seeded (%d staff exist)", count)
		return nil
	}

	log.Println("Seeding staff data...")

	hireDate := time.Now().AddDate(-1, 0, 0) // 1 year ago

	staff := []models.Staff{
		// Housekeeping Department
		{
			EmployeeID: "HK-001",
			FirstName:  "Maria",
			LastName:   "Garcia",
			Email:      "maria.garcia@hotel.com",
			Phone:      "+1234567001",
			Department: "housekeeping",
			Position:   "Housekeeping Supervisor",
			Status:     "active",
			Shift:      "morning",
			HireDate:   hireDate,
		},
		{
			EmployeeID: "HK-002",
			FirstName:  "Ana",
			LastName:   "Martinez",
			Email:      "ana.martinez@hotel.com",
			Phone:      "+1234567002",
			Department: "housekeeping",
			Position:   "Housekeeper",
			Status:     "active",
			Shift:      "morning",
			HireDate:   hireDate.AddDate(0, 3, 0),
		},
		{
			EmployeeID: "HK-003",
			FirstName:  "Rosa",
			LastName:   "Lopez",
			Email:      "rosa.lopez@hotel.com",
			Phone:      "+1234567003",
			Department: "housekeeping",
			Position:   "Housekeeper",
			Status:     "active",
			Shift:      "afternoon",
			HireDate:   hireDate.AddDate(0, 6, 0),
		},
		{
			EmployeeID: "HK-004",
			FirstName:  "Carmen",
			LastName:   "Rodriguez",
			Email:      "carmen.rodriguez@hotel.com",
			Phone:      "+1234567004",
			Department: "housekeeping",
			Position:   "Housekeeper",
			Status:     "on_leave",
			Shift:      "morning",
			HireDate:   hireDate.AddDate(0, 2, 0),
		},

		// Maintenance Department
		{
			EmployeeID: "MN-001",
			FirstName:  "John",
			LastName:   "Smith",
			Email:      "john.smith@hotel.com",
			Phone:      "+1234567010",
			Department: "maintenance",
			Position:   "Maintenance Supervisor",
			Status:     "active",
			Shift:      "morning",
			HireDate:   hireDate,
		},
		{
			EmployeeID: "MN-002",
			FirstName:  "Michael",
			LastName:   "Johnson",
			Email:      "michael.johnson@hotel.com",
			Phone:      "+1234567011",
			Department: "maintenance",
			Position:   "Maintenance Technician",
			Status:     "active",
			Shift:      "morning",
			HireDate:   hireDate.AddDate(0, 4, 0),
		},
		{
			EmployeeID: "MN-003",
			FirstName:  "David",
			LastName:   "Williams",
			Email:      "david.williams@hotel.com",
			Phone:      "+1234567012",
			Department: "maintenance",
			Position:   "Maintenance Technician",
			Status:     "active",
			Shift:      "afternoon",
			HireDate:   hireDate.AddDate(0, 5, 0),
		},
		{
			EmployeeID: "MN-004",
			FirstName:  "Robert",
			LastName:   "Brown",
			Email:      "robert.brown@hotel.com",
			Phone:      "+1234567013",
			Department: "maintenance",
			Position:   "Maintenance Technician",
			Status:     "active",
			Shift:      "night",
			HireDate:   hireDate.AddDate(0, 8, 0),
		},

		// Front Desk Department
		{
			EmployeeID: "FD-001",
			FirstName:  "Sarah",
			LastName:   "Davis",
			Email:      "sarah.davis@hotel.com",
			Phone:      "+1234567020",
			Department: "front_desk",
			Position:   "Front Desk Manager",
			Status:     "active",
			Shift:      "morning",
			HireDate:   hireDate,
		},
		{
			EmployeeID: "FD-002",
			FirstName:  "Emily",
			LastName:   "Wilson",
			Email:      "emily.wilson@hotel.com",
			Phone:      "+1234567021",
			Department: "front_desk",
			Position:   "Front Desk Agent",
			Status:     "active",
			Shift:      "morning",
			HireDate:   hireDate.AddDate(0, 2, 0),
		},
		{
			EmployeeID: "FD-003",
			FirstName:  "Jessica",
			LastName:   "Taylor",
			Email:      "jessica.taylor@hotel.com",
			Phone:      "+1234567022",
			Department: "front_desk",
			Position:   "Front Desk Agent",
			Status:     "active",
			Shift:      "afternoon",
			HireDate:   hireDate.AddDate(0, 3, 0),
		},
		{
			EmployeeID: "FD-004",
			FirstName:  "Amanda",
			LastName:   "Anderson",
			Email:      "amanda.anderson@hotel.com",
			Phone:      "+1234567023",
			Department: "front_desk",
			Position:   "Night Auditor",
			Status:     "active",
			Shift:      "night",
			HireDate:   hireDate.AddDate(0, 6, 0),
		},

		// Room Service Department
		{
			EmployeeID: "RS-001",
			FirstName:  "Chef",
			LastName:   "Antonio",
			Email:      "antonio@hotel.com",
			Phone:      "+1234567030",
			Department: "room_service",
			Position:   "Executive Chef",
			Status:     "active",
			Shift:      "morning",
			HireDate:   hireDate,
		},
		{
			EmployeeID: "RS-002",
			FirstName:  "Carlos",
			LastName:   "Mendez",
			Email:      "carlos.mendez@hotel.com",
			Phone:      "+1234567031",
			Department: "room_service",
			Position:   "Sous Chef",
			Status:     "active",
			Shift:      "morning",
			HireDate:   hireDate.AddDate(0, 1, 0),
		},
		{
			EmployeeID: "RS-003",
			FirstName:  "James",
			LastName:   "Lee",
			Email:      "james.lee@hotel.com",
			Phone:      "+1234567032",
			Department: "room_service",
			Position:   "Room Service Attendant",
			Status:     "active",
			Shift:      "morning",
			HireDate:   hireDate.AddDate(0, 4, 0),
		},
		{
			EmployeeID: "RS-004",
			FirstName:  "Kevin",
			LastName:   "Chen",
			Email:      "kevin.chen@hotel.com",
			Phone:      "+1234567033",
			Department: "room_service",
			Position:   "Room Service Attendant",
			Status:     "active",
			Shift:      "afternoon",
			HireDate:   hireDate.AddDate(0, 5, 0),
		},
		{
			EmployeeID: "RS-005",
			FirstName:  "Daniel",
			LastName:   "Kim",
			Email:      "daniel.kim@hotel.com",
			Phone:      "+1234567034",
			Department: "room_service",
			Position:   "Room Service Attendant",
			Status:     "active",
			Shift:      "night",
			HireDate:   hireDate.AddDate(0, 7, 0),
		},

		// Management
		{
			EmployeeID: "MG-001",
			FirstName:  "William",
			LastName:   "Thompson",
			Email:      "william.thompson@hotel.com",
			Phone:      "+1234567040",
			Department: "management",
			Position:   "General Manager",
			Status:     "active",
			Shift:      "morning",
			HireDate:   hireDate.AddDate(-2, 0, 0),
		},
		{
			EmployeeID: "MG-002",
			FirstName:  "Elizabeth",
			LastName:   "White",
			Email:      "elizabeth.white@hotel.com",
			Phone:      "+1234567041",
			Department: "management",
			Position:   "Assistant Manager",
			Status:     "active",
			Shift:      "morning",
			HireDate:   hireDate.AddDate(-1, 0, 0),
		},
	}

	for _, s := range staff {
		if err := db.Create(&s).Error; err != nil {
			log.Printf("Error creating staff %s %s: %v", s.FirstName, s.LastName, err)
			continue
		}
	}

	log.Printf("Seeded %d staff members", len(staff))
	return nil
}

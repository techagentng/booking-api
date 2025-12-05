package db

import (
	"hotel/models"
	"log"
	"time"

	"gorm.io/gorm"
)

// SeedTabletTestData seeds test data for the tablet interface
func SeedTabletTestData(db *gorm.DB) error {
	log.Println("Seeding tablet test data...")

	// Check if test guest already exists
	var existingGuest models.Guest
	if err := db.Where("email = ?", "john.doe@example.com").First(&existingGuest).Error; err == nil {
		log.Println("Tablet test data already exists, skipping...")
		return nil
	}

	// Create test guests
	guests := []models.Guest{
		{
			Name:        "John Doe",
			Email:       "john.doe@example.com",
			Phone:       "+1234567890",
			IDType:      "passport",
			IDNumber:    "AB123456",
			Nationality: "USA",
		},
		{
			Name:        "Jane Smith",
			Email:       "jane.smith@example.com",
			Phone:       "+1987654321",
			IDType:      "drivers_license",
			IDNumber:    "DL789012",
			Nationality: "USA",
		},
	}

	for i := range guests {
		if err := db.Create(&guests[i]).Error; err != nil {
			log.Printf("Error creating guest %s: %v", guests[i].Name, err)
			continue
		}
	}

	// Get rooms for reservations
	var rooms []models.Room
	if err := db.Limit(5).Find(&rooms).Error; err != nil {
		return err
	}

	if len(rooms) < 2 {
		log.Println("Not enough rooms to create reservations")
		return nil
	}

	// Create reservations
	today := time.Now()
	checkIn := time.Date(today.Year(), today.Month(), today.Day(), 14, 0, 0, 0, today.Location())
	checkOut := checkIn.AddDate(0, 0, 3)

	reservations := []models.Reservation{
		{
			GuestID:            guests[0].ID,
			RoomID:             rooms[0].ID,
			CheckInDate:        checkIn,
			CheckOutDate:       checkOut,
			NumberOfGuests:     2,
			Status:             "confirmed",
			TotalAmount:        450.00,
			ConfirmationNumber: "BK-2024-TEST-001",
			PaymentStatus:      "paid",
			PaymentMethod:      "credit_card",
			SpecialRequests:    "Late check-in requested",
		},
		{
			GuestID:            guests[1].ID,
			RoomID:             rooms[1].ID,
			CheckInDate:        checkIn,
			CheckOutDate:       checkOut.AddDate(0, 0, 2),
			NumberOfGuests:     1,
			Status:             "checked-in",
			TotalAmount:        750.00,
			ConfirmationNumber: "BK-2024-TEST-002",
			PaymentStatus:      "paid",
			PaymentMethod:      "credit_card",
			CheckedInAt:        &today,
		},
	}

	for i := range reservations {
		if err := db.Create(&reservations[i]).Error; err != nil {
			log.Printf("Error creating reservation: %v", err)
			continue
		}
	}

	// Update room status for checked-in reservation
	db.Model(&models.Room{}).Where("id = ?", rooms[1].ID).Update("status", "occupied")

	// Create menu items
	menuItems := []models.MenuItem{
		// Breakfast
		{Name: "Continental Breakfast", Description: "Fresh pastries, fruits, yogurt, and coffee", Category: "Breakfast", Price: 18.00, PreparationTime: 10, Available: true},
		{Name: "American Breakfast", Description: "Eggs, bacon, toast, hash browns, and coffee", Category: "Breakfast", Price: 22.00, PreparationTime: 15, Available: true},
		{Name: "Pancake Stack", Description: "Fluffy pancakes with maple syrup and butter", Category: "Breakfast", Price: 14.00, PreparationTime: 12, Available: true},

		// Main Course
		{Name: "Club Sandwich", Description: "Triple-decker with turkey, bacon, lettuce, tomato", Category: "Main Course", Price: 15.00, PreparationTime: 15, Available: true},
		{Name: "Grilled Salmon", Description: "Atlantic salmon with seasonal vegetables", Category: "Main Course", Price: 32.00, PreparationTime: 25, Available: true},
		{Name: "Beef Burger", Description: "Angus beef patty with fries", Category: "Main Course", Price: 18.00, PreparationTime: 20, Available: true},
		{Name: "Pasta Carbonara", Description: "Creamy pasta with bacon and parmesan", Category: "Main Course", Price: 20.00, PreparationTime: 18, Available: true},
		{Name: "Chicken Caesar Wrap", Description: "Grilled chicken with caesar dressing in a wrap", Category: "Main Course", Price: 16.00, PreparationTime: 12, Available: true},

		// Salads
		{Name: "Caesar Salad", Description: "Romaine lettuce with caesar dressing and croutons", Category: "Salads", Price: 12.00, PreparationTime: 8, Available: true},
		{Name: "Greek Salad", Description: "Fresh vegetables with feta cheese and olives", Category: "Salads", Price: 14.00, PreparationTime: 8, Available: true},
		{Name: "Garden Salad", Description: "Mixed greens with house vinaigrette", Category: "Salads", Price: 10.00, PreparationTime: 5, Available: true},

		// Desserts
		{Name: "Chocolate Lava Cake", Description: "Warm chocolate cake with molten center", Category: "Desserts", Price: 12.00, PreparationTime: 15, Available: true},
		{Name: "New York Cheesecake", Description: "Classic creamy cheesecake with berry compote", Category: "Desserts", Price: 10.00, PreparationTime: 5, Available: true},
		{Name: "Ice Cream Sundae", Description: "Three scoops with chocolate sauce and whipped cream", Category: "Desserts", Price: 8.00, PreparationTime: 5, Available: true},
		{Name: "Tiramisu", Description: "Italian coffee-flavored dessert", Category: "Desserts", Price: 11.00, PreparationTime: 5, Available: true},

		// Beverages
		{Name: "Fresh Orange Juice", Description: "Freshly squeezed orange juice", Category: "Beverages", Price: 6.00, PreparationTime: 3, Available: true},
		{Name: "Coffee", Description: "Freshly brewed premium coffee", Category: "Beverages", Price: 4.00, PreparationTime: 3, Available: true},
		{Name: "Cappuccino", Description: "Espresso with steamed milk foam", Category: "Beverages", Price: 5.00, PreparationTime: 5, Available: true},
		{Name: "Iced Tea", Description: "Refreshing iced tea with lemon", Category: "Beverages", Price: 4.00, PreparationTime: 2, Available: true},
		{Name: "Smoothie", Description: "Mixed berry smoothie", Category: "Beverages", Price: 8.00, PreparationTime: 5, Available: true},
		{Name: "Mineral Water", Description: "Sparkling or still mineral water", Category: "Beverages", Price: 3.00, PreparationTime: 1, Available: true},
	}

	for i := range menuItems {
		// Check if menu item already exists
		var existing models.MenuItem
		if err := db.Where("name = ?", menuItems[i].Name).First(&existing).Error; err == nil {
			continue // Skip if exists
		}
		if err := db.Create(&menuItems[i]).Error; err != nil {
			log.Printf("Error creating menu item %s: %v", menuItems[i].Name, err)
		}
	}

	// Create sample service requests
	var existingRequest models.GuestServiceRequest
	if err := db.Where("request_number = ?", "HK-2024-001").First(&existingRequest).Error; err != nil {
		now := time.Now()
		estimatedTime := now.Add(30 * time.Minute)

		serviceRequests := []models.GuestServiceRequest{
			{
				RequestNumber: "HK-2024-001",
				RoomID:        rooms[0].ID,
				GuestID:       guests[0].ID,
				Type:          "housekeeping",
				ServiceType:   "cleaning",
				Priority:      "normal",
				Status:        "pending",
				Notes:         "Please clean the room and replace towels",
				RequestedAt:   now,
				EstimatedTime: &estimatedTime,
			},
			{
				RequestNumber: "HK-2024-002",
				RoomID:        rooms[1].ID,
				GuestID:       guests[1].ID,
				Type:          "housekeeping",
				ServiceType:   "towels",
				Priority:      "high",
				Status:        "in_progress",
				Notes:         "Need extra towels urgently",
				RequestedAt:   now.Add(-20 * time.Minute),
				EstimatedTime: &estimatedTime,
				AssignedTo:    "Maria",
			},
			{
				RequestNumber: "MN-2024-001",
				RoomID:        rooms[0].ID,
				GuestID:       guests[0].ID,
				Type:          "maintenance",
				IssueType:     "air_conditioning",
				Priority:      "high",
				Status:        "pending",
				Description:   "AC is not cooling properly, making strange noise",
				RequestedAt:   now.Add(-10 * time.Minute),
				EstimatedTime: &estimatedTime,
			},
			{
				RequestNumber: "MN-2024-002",
				RoomID:        rooms[1].ID,
				GuestID:       guests[1].ID,
				Type:          "maintenance",
				IssueType:     "plumbing",
				Priority:      "urgent",
				Status:        "in_progress",
				Description:   "Bathroom sink is clogged",
				RequestedAt:   now.Add(-30 * time.Minute),
				EstimatedTime: &estimatedTime,
				AssignedTo:    "John (Maintenance)",
			},
			{
				RequestNumber: "HK-2024-003",
				RoomID:        rooms[0].ID,
				GuestID:       guests[0].ID,
				Type:          "housekeeping",
				ServiceType:   "turndown",
				Priority:      "normal",
				Status:        "completed",
				Notes:         "Evening turndown service",
				RequestedAt:   now.Add(-2 * time.Hour),
				CompletedBy:   "Sarah",
			},
		}

		for i := range serviceRequests {
			if err := db.Create(&serviceRequests[i]).Error; err != nil {
				log.Printf("Error creating service request: %v", err)
			}
		}
		log.Println("Service requests seeded successfully!")
	}

	log.Println("Tablet test data seeded successfully!")
	log.Printf("Test Guest 1: %s (Room %s - Confirmed)", guests[0].Name, rooms[0].RoomNumber)
	log.Printf("Test Guest 2: %s (Room %s - Checked-In)", guests[1].Name, rooms[1].RoomNumber)

	return nil
}

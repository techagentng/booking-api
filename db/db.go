package db

import (
	"fmt"
	"log"

	"hotel/config"
	"hotel/models"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormDB struct {
	DB *gorm.DB
}

func GetDB(c *config.Config) *GormDB {
	gormDB := &GormDB{}
	gormDB.Init(c)
	return gormDB
}

func (g *GormDB) Init(c *config.Config) {
	g.DB = getPostgresDB(c)

	if err := migrate(g.DB); err != nil {
		log.Fatalf("unable to run migrations: %v", err)
	}
}

func getPostgresDB(c *config.Config) *gorm.DB {
	log.Printf("Connecting to postgres: %+v", c)
	postgresDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Africa/Lagos",
		c.PostgresHost, c.PostgresUser, c.PostgresPassword, c.PostgresDB, c.PostgresPort)

	// Create GORM DB instance
	gormConfig := &gorm.Config{}
	if c.Env != "prod" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN: postgresDSN,
	}), gormConfig)
	if err != nil {
		log.Fatal(err)
	}

	return gormDB
}

func SeedRoles(db *gorm.DB) error {
	roles := []string{
		models.RoleAdmin,
		models.RoleProvider,
		models.RoleCustomer,
		models.RoleModerator,
		models.RoleUser,
	}

	for _, roleName := range roles {
		var existingRole models.Role
		err := db.Where("name = ?", roleName).First(&existingRole).Error

		if err == gorm.ErrRecordNotFound {
			// Role doesn't exist, create it
			newRole := models.Role{
				ID:   uuid.New(),
				Name: roleName,
			}
			if err := db.Create(&newRole).Error; err != nil {
				return fmt.Errorf("failed to create role %s: %w", roleName, err)
			}
			log.Printf("Created role: %s", roleName)
		} else if err != nil {
			// Some other error occurred
			return fmt.Errorf("error checking role %s: %w", roleName, err)
		}
		// Role already exists, skip
	}

	return nil
}

// SeedRooms creates sample rooms if none exist
func SeedRooms(db *gorm.DB) error {
	var count int64
	db.Model(&models.Room{}).Count(&count)
	if count > 0 {
		log.Printf("Rooms already seeded (%d rooms exist)", count)
		return nil
	}

	rooms := []models.Room{
		// Standard rooms (Floor 1)
		{RoomNumber: "101", RoomType: "Standard", Floor: 1, PricePerNight: 100.00, Status: "available", Capacity: 2, BedType: "Queen", Description: "Cozy standard room with city view"},
		{RoomNumber: "102", RoomType: "Standard", Floor: 1, PricePerNight: 100.00, Status: "available", Capacity: 2, BedType: "Queen", Description: "Cozy standard room with garden view"},
		{RoomNumber: "103", RoomType: "Standard", Floor: 1, PricePerNight: 100.00, Status: "available", Capacity: 2, BedType: "Twin", Description: "Standard room with twin beds"},
		{RoomNumber: "104", RoomType: "Standard", Floor: 1, PricePerNight: 100.00, Status: "available", Capacity: 2, BedType: "Queen", Description: "Standard room near elevator"},

		// Deluxe rooms (Floor 2)
		{RoomNumber: "201", RoomType: "Deluxe", Floor: 2, PricePerNight: 150.00, Status: "available", Capacity: 2, BedType: "King", Description: "Spacious deluxe room with balcony"},
		{RoomNumber: "202", RoomType: "Deluxe", Floor: 2, PricePerNight: 150.00, Status: "available", Capacity: 3, BedType: "King", Description: "Deluxe room with ocean view"},
		{RoomNumber: "203", RoomType: "Deluxe", Floor: 2, PricePerNight: 150.00, Status: "available", Capacity: 2, BedType: "King", Description: "Corner deluxe room with extra space"},
		{RoomNumber: "204", RoomType: "Deluxe", Floor: 2, PricePerNight: 150.00, Status: "available", Capacity: 3, BedType: "King", Description: "Deluxe room with work desk"},

		// Suite rooms (Floor 3)
		{RoomNumber: "301", RoomType: "Suite", Floor: 3, PricePerNight: 250.00, Status: "available", Capacity: 4, BedType: "King", Description: "Luxury suite with living area"},
		{RoomNumber: "302", RoomType: "Suite", Floor: 3, PricePerNight: 250.00, Status: "available", Capacity: 4, BedType: "King", Description: "Executive suite with panoramic view"},
		{RoomNumber: "303", RoomType: "Suite", Floor: 3, PricePerNight: 250.00, Status: "available", Capacity: 4, BedType: "King", Description: "Family suite with kitchenette"},
	}

	for _, room := range rooms {
		if err := db.Create(&room).Error; err != nil {
			return fmt.Errorf("failed to create room %s: %w", room.RoomNumber, err)
		}
	}

	log.Printf("Seeded %d rooms", len(rooms))
	return nil
}

func migrate(db *gorm.DB) error {
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	// AutoMigrate all models to create tables
	err := db.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.Blacklist{},
		&models.OAuthState{},
		&models.Guest{},
		&models.GuestPreferences{},
		&models.GuestAIInsights{},
		&models.Room{},
		&models.Reservation{},
		&models.ServiceRequest{},
		&models.MenuItem{},
		&models.RoomServiceOrder{},
		&models.OrderItem{},
		&models.GuestServiceRequest{},
		&models.Staff{},
		&models.HallBooking{},
		&models.Payment{},
		&models.Invoice{},
		&models.CalendarAvailability{},
		&models.TimeSlot{},
	)
	if err != nil {
		return fmt.Errorf("failed to run auto migrations: %w", err)
	}

	// Add indexes for performance
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_calendar_availability_date ON calendar_availability(date)").Error; err != nil {
		return fmt.Errorf("failed to create calendar availability date index: %w", err)
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_time_slots_date ON time_slots(date)").Error; err != nil {
		return fmt.Errorf("failed to create time slots date index: %w", err)
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_time_slots_date_time ON time_slots(date, start_time)").Error; err != nil {
		return fmt.Errorf("failed to create time slots date time index: %w", err)
	}

	return nil
}

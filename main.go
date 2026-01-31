package main

import (
	"log"

	"hotel/config"
	"hotel/db"
	"hotel/server"
	"hotel/services"
)

func main() {
	// Load configuration
	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize database
	gormDB := db.GetDB(conf)

	// Seed roles
	if err := db.SeedRoles(gormDB.DB); err != nil {
		log.Fatalf("error seeding roles: %v", err)
	}

	// Seed rooms
	if err := db.SeedRooms(gormDB.DB); err != nil {
		log.Fatalf("error seeding rooms: %v", err)
	}

	// Seed tablet test data (guests, reservations, menu items)
	if err := db.SeedTabletTestData(gormDB.DB); err != nil {
		log.Fatalf("error seeding tablet test data: %v", err)
	}

	// Seed staff
	if err := db.SeedStaff(gormDB.DB); err != nil {
		log.Fatalf("error seeding staff: %v", err)
	}

	// Repositories
	authRepo := db.NewAuthRepo(gormDB)
	guestRepo := db.NewGuestRepository(gormDB.DB)
	roomRepo := db.NewRoomRepository(gormDB.DB)
	reservationRepo := db.NewReservationRepository(gormDB.DB)
	roomServiceRepo := db.NewRoomServiceRepository(gormDB.DB)
	guestServiceRepo := db.NewGuestServiceRepository(gormDB.DB)
	staffRepo := db.NewStaffRepository(gormDB.DB)
	hallBookingRepo := db.NewHallBookingRepository(gormDB.DB)
	adminHallBookingRepo := db.NewAdminHallBookingRepository(gormDB.DB)
	calendarRepo := db.NewCalendarRepository(gormDB.DB)
	paymentRepo := db.NewPaymentRepository(gormDB.DB)

	// Services
	authService := services.NewAuthService(authRepo, conf)
	notificationHub := services.NewNotificationHub()
	mailService := services.NewMailgunService()

	// Initialize Stripe
	services.InitStripe(conf)

	// Server setup
	s := &server.Server{
		Config:                     conf,
		AuthRepository:             authRepo,
		AuthService:                authService,
		GuestRepository:            guestRepo,
		RoomRepository:             roomRepo,
		ReservationRepository:      reservationRepo,
		RoomServiceRepository:      roomServiceRepo,
		GuestServiceRepository:     guestServiceRepo,
		StaffRepository:            staffRepo,
		HallBookingRepository:      hallBookingRepo,
		AdminHallBookingRepository: adminHallBookingRepo,
		CalendarRepository:         calendarRepo,
		PaymentRepository:          paymentRepo,
		NotificationHub:            notificationHub,
		MailService:                mailService,
		DB:                         gormDB.DB,
	}

	// Start server
	s.Start()
}

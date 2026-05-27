package tripsbook

import (
	"time"
	
	"hotel/tripsbook/models"
)

// MockServiceProviders generates mock service providers with positioning
var MockServiceProviders = []models.ServiceProvider{
	// Hotels Category - Top 20 with admin positioning
	{
		ID: "sp_hotel_001",
		UserID: "user_001",
		BusinessName: "Eko Hotels & Suites",
		DisplayName: "Eko Hotels & Suites",
		Description: "Luxury beachfront hotel with stunning ocean views and world-class amenities",
		CategoryID: "cat_hotels",
		SubCategories: []string{"luxury", "beachfront", "business"},
		
		Phone: "+234-1-2778000",
		Email: "eko.suites@ekohotels.com",
		Website: "https://ekohotels.com",
		Address: "1415 Adetokunbo Ademola Street, Victoria Island",
		City: "Lagos",
		State: "Lagos State",
		
		BusinessType: models.BusinessTypeCompany,
		EstablishedYear: 1977,
		EmployeesCount: 450,
		
		ServiceAreas: []string{"Lagos", "Victoria Island", "Ikoyi", "Lekki"},
		ServiceRadius: 25.0,
		
		Logo: "https://cdn.tripsbook.com/logos/eko-hotels.png",
		BannerImage: "https://cdn.tripsbook.com/banners/eko-hotels.jpg",
		Gallery: []string{
			"https://cdn.tripsbook.com/images/hotels/eko/1.jpg",
			"https://cdn.tripsbook.com/images/hotels/eko/2.jpg",
			"https://cdn.tripsbook.com/images/hotels/eko/3.jpg",
		},
		
		IsVerified: true,
		VerificationStatus: models.VerificationStatusVerified,
		IsActive: true,
		IsFeatured: true,
		
		AdminPosition: 1,
		PositionCategory: "cat_hotels",
		
		AverageRating: 4.6,
		TotalReviews: 1247,
		RatingBreakdown: map[int]int{
			5: 789, 4: 312, 3: 98, 2: 35, 1: 13,
		},
		
		TotalServices: 8,
		ActiveServices: 8,
		CompletedBookings: 8934,
		
		CreatedAt: time.Now().AddDate(-3, 0, 0),
		UpdatedAt: time.Now().Add(-24 * time.Hour),
		LastActiveAt: time.Now().Add(-2 * time.Hour),
	},
	{
		ID: "sp_hotel_002",
		UserID: "user_002",
		BusinessName: "Federal Palace Hotel",
		DisplayName: "Federal Palace Hotel",
		Description: "Historic luxury hotel with modern amenities and rich cultural heritage",
		CategoryID: "cat_hotels",
		SubCategories: []string{"luxury", "historic", "business"},
		
		Phone: "+234-1-2611000",
		Email: "info@federalpalacehotel.com",
		Website: "https://federalpalacehotel.com",
		Address: "1-3 Ahmadu Bello Way, Victoria Island",
		City: "Lagos",
		State: "Lagos State",
		
		BusinessType: models.BusinessTypeCompany,
		EstablishedYear: 1960,
		EmployeesCount: 380,
		
		ServiceAreas: []string{"Lagos", "Victoria Island", "Ikoyi"},
		ServiceRadius: 20.0,
		
		Logo: "https://cdn.tripsbook.com/logos/federal-palace.png",
		BannerImage: "https://cdn.tripsbook.com/banners/federal-palace.jpg",
		Gallery: []string{
			"https://cdn.tripsbook.com/images/hotels/federal/1.jpg",
			"https://cdn.tripsbook.com/images/hotels/federal/2.jpg",
		},
		
		IsVerified: true,
		VerificationStatus: models.VerificationStatusVerified,
		IsActive: true,
		IsFeatured: true,
		
		AdminPosition: 2,
		PositionCategory: "cat_hotels",
		
		AverageRating: 4.5,
		TotalReviews: 987,
		RatingBreakdown: map[int]int{
			5: 623, 4: 245, 3: 87, 2: 22, 1: 10,
		},
		
		TotalServices: 6,
		ActiveServices: 6,
		CompletedBookings: 6234,
		
		CreatedAt: time.Now().AddDate(-5, 0, 0),
		UpdatedAt: time.Now().Add(-48 * time.Hour),
		LastActiveAt: time.Now().Add(-4 * time.Hour),
	},
	{
		ID: "sp_hotel_003",
		UserID: "user_003",
		BusinessName: "Lagos Continental Hotel",
		DisplayName: "Lagos Continental Hotel",
		Description: "Modern luxury hotel with premium facilities and conference center",
		CategoryID: "cat_hotels",
		SubCategories: []string{"luxury", "business", "conference"},
		
		Phone: "+234-1-2777777",
		Email: "lagos@continentalhotels.com",
		Website: "https://continentalhotels.com/lagos",
		Address: "52 Kofo Abayomi Street, Victoria Island",
		City: "Lagos",
		State: "Lagos State",
		
		BusinessType: models.BusinessTypeFranchise,
		EstablishedYear: 2015,
		EmployeesCount: 290,
		
		ServiceAreas: []string{"Lagos", "Victoria Island", "Ikoyi"},
		ServiceRadius: 18.0,
		
		Logo: "https://cdn.tripsbook.com/logos/continental-lagos.png",
		BannerImage: "https://cdn.tripsbook.com/banners/continental-lagos.jpg",
		Gallery: []string{
			"https://cdn.tripsbook.com/images/hotels/continental/1.jpg",
			"https://cdn.tripsbook.com/images/hotels/continental/2.jpg",
		},
		
		IsVerified: true,
		VerificationStatus: models.VerificationStatusVerified,
		IsActive: true,
		IsFeatured: false,
		
		AdminPosition: 3,
		PositionCategory: "cat_hotels",
		
		AverageRating: 4.4,
		TotalReviews: 756,
		RatingBreakdown: map[int]int{
			5: 445, 4: 234, 3: 56, 2: 15, 1: 6,
		},
		
		TotalServices: 7,
		ActiveServices: 7,
		CompletedBookings: 4123,
		
		CreatedAt: time.Now().AddDate(-2, 0, 0),
		UpdatedAt: time.Now().Add(-12 * time.Hour),
		LastActiveAt: time.Now().Add(-1 * time.Hour),
	},
	
	// Restaurants Category - Top 20 with admin positioning
	{
		ID: "sp_restaurant_001",
		UserID: "user_101",
		BusinessName: "Terra Kulture",
		DisplayName: "Terra Kulture",
		Description: "Contemporary Nigerian cuisine with cultural experience and art gallery",
		CategoryID: "cat_restaurants",
		SubCategories: []string{"nigerian", "contemporary", "cultural"},
		
		Phone: "+234-1-2776322",
		Email: "info@terrakulture.com",
		Website: "https://terrakulture.com",
		Address: "1376 Tiamiyu Savage Street, Victoria Island",
		City: "Lagos",
		State: "Lagos State",
		
		BusinessType: models.BusinessTypeCompany,
		EstablishedYear: 2004,
		EmployeesCount: 85,
		
		ServiceAreas: []string{"Lagos", "Victoria Island", "Ikoyi"},
		ServiceRadius: 15.0,
		
		Logo: "https://cdn.tripsbook.com/logos/terra-kulture.png",
		BannerImage: "https://cdn.tripsbook.com/banners/terra-kulture.jpg",
		Gallery: []string{
			"https://cdn.tripsbook.com/images/restaurants/terra/1.jpg",
			"https://cdn.tripsbook.com/images/restaurants/terra/2.jpg",
		},
		
		IsVerified: true,
		VerificationStatus: models.VerificationStatusVerified,
		IsActive: true,
		IsFeatured: true,
		
		AdminPosition: 1,
		PositionCategory: "cat_restaurants",
		
		AverageRating: 4.5,
		TotalReviews: 892,
		RatingBreakdown: map[int]int{
			5: 567, 4: 234, 3: 67, 2: 18, 1: 6,
		},
		
		TotalServices: 4,
		ActiveServices: 4,
		CompletedBookings: 12456,
		
		CreatedAt: time.Now().AddDate(-8, 0, 0),
		UpdatedAt: time.Now().Add(-6 * time.Hour),
		LastActiveAt: time.Now().Add(-30 * time.Minute),
	},
	{
		ID: "sp_restaurant_002",
		UserID: "user_102",
		BusinessName: "Nigerian Kitchen",
		DisplayName: "Nigerian Kitchen",
		Description: "Authentic Nigerian dishes with modern twist and cozy ambiance",
		CategoryID: "cat_restaurants",
		SubCategories: []string{"nigerian", "traditional", "family-friendly"},
		
		Phone: "+234-1-4538921",
		Email: "hello@nigeriankitchen.com",
		Website: "https://nigeriankitchen.com",
		Address: "435 Awolowo Road, Ikoyi",
		City: "Lagos",
		State: "Lagos State",
		
		BusinessType: models.BusinessTypeCompany,
		EstablishedYear: 2012,
		EmployeesCount: 45,
		
		ServiceAreas: []string{"Lagos", "Ikoyi", "Victoria Island"},
		ServiceRadius: 12.0,
		
		Logo: "https://cdn.tripsbook.com/logos/nigerian-kitchen.png",
		BannerImage: "https://cdn.tripsbook.com/banners/nigerian-kitchen.jpg",
		Gallery: []string{
			"https://cdn.tripsbook.com/images/restaurants/nigerian/1.jpg",
		},
		
		IsVerified: true,
		VerificationStatus: models.VerificationStatusVerified,
		IsActive: true,
		IsFeatured: false,
		
		AdminPosition: 2,
		PositionCategory: "cat_restaurants",
		
		AverageRating: 4.6,
		TotalReviews: 723,
		RatingBreakdown: map[int]int{
			5: 489, 4: 189, 3: 34, 2: 8, 1: 3,
		},
		
		TotalServices: 3,
		ActiveServices: 3,
		CompletedBookings: 8934,
		
		CreatedAt: time.Now().AddDate(-4, 0, 0),
		UpdatedAt: time.Now().Add(-3 * time.Hour),
		LastActiveAt: time.Now().Add(-45 * time.Minute),
	},
	
	// Transport Category - Top 20 with admin positioning
	{
		ID: "sp_transport_001",
		UserID: "user_201",
		BusinessName: "Uber Premium",
		DisplayName: "Uber Premium",
		Description: "Premium ride service with professional drivers and luxury vehicles",
		CategoryID: "cat_transport",
		SubCategories: []string{"ride-hailing", "premium", "airport-transfer"},
		
		Phone: "+234-800-000-0000",
		Email: "premium@uber.com",
		Website: "https://uber.com/premium",
		Address: "Tech Hub, Yaba",
		City: "Lagos",
		State: "Lagos State",
		
		BusinessType: models.BusinessTypeFranchise,
		EstablishedYear: 2014,
		EmployeesCount: 1200,
		
		ServiceAreas: []string{"Lagos", "Abuja", "Port Harcourt", "Kano"},
		ServiceRadius: 50.0,
		
		Logo: "https://cdn.tripsbook.com/logos/uber-premium.png",
		BannerImage: "https://cdn.tripsbook.com/banners/uber-premium.jpg",
		Gallery: []string{
			"https://cdn.tripsbook.com/images/transport/uber/1.jpg",
			"https://cdn.tripsbook.com/images/transport/uber/2.jpg",
		},
		
		IsVerified: true,
		VerificationStatus: models.VerificationStatusVerified,
		IsActive: true,
		IsFeatured: true,
		
		AdminPosition: 1,
		PositionCategory: "cat_transport",
		
		AverageRating: 4.7,
		TotalReviews: 2543,
		RatingBreakdown: map[int]int{
			5: 1892, 4: 456, 3: 145, 2: 34, 1: 16,
		},
		
		TotalServices: 5,
		ActiveServices: 5,
		CompletedBookings: 45678,
		
		CreatedAt: time.Now().AddDate(-6, 0, 0),
		UpdatedAt: time.Now().Add(-1 * time.Hour),
		LastActiveAt: time.Now().Add(-15 * time.Minute),
	},
	{
		ID: "sp_transport_002",
		UserID: "user_202",
		BusinessName: "Bolt Professional",
		DisplayName: "Bolt Professional",
		Description: "Professional transport service with experienced drivers and competitive rates",
		CategoryID: "cat_transport",
		SubCategories: []string{"ride-hailing", "professional", "city-transport"},
		
		Phone: "+234-1-888-0000",
		Email: "professional@bolt.ng",
		Website: "https://bolt.ng/professional",
		Address: "Lekki Phase 1, Lagos",
		City: "Lagos",
		State: "Lagos State",
		
		BusinessType: models.BusinessTypeFranchise,
		EstablishedYear: 2015,
		EmployeesCount: 890,
		
		ServiceAreas: []string{"Lagos", "Abuja", "Ibadan", "Benin City"},
		ServiceRadius: 40.0,
		
		Logo: "https://cdn.tripsbook.com/logos/bolt-professional.png",
		BannerImage: "https://cdn.tripsbook.com/banners/bolt-professional.jpg",
		Gallery: []string{
			"https://cdn.tripsbook.com/images/transport/bolt/1.jpg",
		},
		
		IsVerified: true,
		VerificationStatus: models.VerificationStatusVerified,
		IsActive: true,
		IsFeatured: false,
		
		AdminPosition: 2,
		PositionCategory: "cat_transport",
		
		AverageRating: 4.5,
		TotalReviews: 1823,
		RatingBreakdown: map[int]int{
			5: 1234, 4: 389, 3: 145, 2: 42, 1: 13,
		},
		
		TotalServices: 4,
		ActiveServices: 4,
		CompletedBookings: 32145,
		
		CreatedAt: time.Now().AddDate(-5, 0, 0),
		UpdatedAt: time.Now().Add(-2 * time.Hour),
		LastActiveAt: time.Now().Add(-20 * time.Minute),
	},
	
	// Shopping Category - Top 20 with admin positioning
	{
		ID: "sp_shopping_001",
		UserID: "user_301",
		BusinessName: "Palms Shopping Mall",
		DisplayName: "The Palms",
		Description: "Premier shopping destination with international brands and entertainment",
		CategoryID: "cat_shopping",
		SubCategories: []string{"mall", "international-brands", "entertainment"},
		
		Phone: "+234-1-2773000",
		Email: "info@thepalms.com.ng",
		Website: "https://thepalms.com.ng",
		Address: "1, Water Corporation Road, Oniru, Victoria Island",
		City: "Lagos",
		State: "Lagos State",
		
		BusinessType: models.BusinessTypeCompany,
		EstablishedYear: 2005,
		EmployeesCount: 650,
		
		ServiceAreas: []string{"Lagos", "Victoria Island", "Lekki"},
		ServiceRadius: 20.0,
		
		Logo: "https://cdn.tripsbook.com/logos/palms-mall.png",
		BannerImage: "https://cdn.tripsbook.com/banners/palms-mall.jpg",
		Gallery: []string{
			"https://cdn.tripsbook.com/images/shopping/palms/1.jpg",
			"https://cdn.tripsbook.com/images/shopping/palms/2.jpg",
		},
		
		IsVerified: true,
		VerificationStatus: models.VerificationStatusVerified,
		IsActive: true,
		IsFeatured: true,
		
		AdminPosition: 1,
		PositionCategory: "cat_shopping",
		
		AverageRating: 4.4,
		TotalReviews: 1567,
		RatingBreakdown: map[int]int{
			5: 890, 4: 456, 3: 167, 2: 38, 1: 16,
		},
		
		TotalServices: 12,
		ActiveServices: 12,
		CompletedBookings: 234567,
		
		CreatedAt: time.Now().AddDate(-10, 0, 0),
		UpdatedAt: time.Now().Add(-4 * time.Hour),
		LastActiveAt: time.Now().Add(-1 * time.Hour),
	},
	{
		ID: "sp_shopping_002",
		UserID: "user_302",
		BusinessName: "Ikeja City Mall",
		DisplayName: "Ikeja City Mall",
		Description: "Popular shopping mall with diverse retail stores and family entertainment",
		CategoryID: "cat_shopping",
		SubCategories: []string{"mall", "family-friendly", "diverse-retail"},
		
		Phone: "+234-1-2774000",
		Email: "info@ikejacitymall.com",
		Website: "https://ikejacitymall.com",
		Address: "Oba Akran Avenue, Ikeja",
		City: "Lagos",
		State: "Lagos State",
		
		BusinessType: models.BusinessTypeCompany,
		EstablishedYear: 2011,
		EmployeesCount: 480,
		
		ServiceAreas: []string{"Lagos", "Ikeja", "Mainland"},
		ServiceRadius: 25.0,
		
		Logo: "https://cdn.tripsbook.com/logos/ikeja-mall.png",
		BannerImage: "https://cdn.tripsbook.com/banners/ikeja-mall.jpg",
		Gallery: []string{
			"https://cdn.tripsbook.com/images/shopping/ikeja/1.jpg",
		},
		
		IsVerified: true,
		VerificationStatus: models.VerificationStatusVerified,
		IsActive: true,
		IsFeatured: false,
		
		AdminPosition: 2,
		PositionCategory: "cat_shopping",
		
		AverageRating: 4.3,
		TotalReviews: 1234,
		RatingBreakdown: map[int]int{
			5: 678, 4: 345, 3: 145, 2: 45, 1: 21,
		},
		
		TotalServices: 10,
		ActiveServices: 10,
		CompletedBookings: 189234,
		
		CreatedAt: time.Now().AddDate(-8, 0, 0),
		UpdatedAt: time.Now().Add(-5 * time.Hour),
		LastActiveAt: time.Now().Add(-2 * time.Hour),
	},
	
	// Events Category - Top 20 with admin positioning
	{
		ID: "sp_events_001",
		UserID: "user_401",
		BusinessName: "Eko Convention Center",
		DisplayName: "Eko Convention Center",
		Description: "World-class convention center with state-of-the-art facilities",
		CategoryID: "cat_events",
		SubCategories: []string{"convention", "corporate", "large-events"},
		
		Phone: "+234-1-2778001",
		Email: "events@ekohotels.com",
		Website: "https://ekoconventioncenter.com",
		Address: "1415 Adetokunbo Ademola Street, Victoria Island",
		City: "Lagos",
		State: "Lagos State",
		
		BusinessType: models.BusinessTypeCompany,
		EstablishedYear: 2018,
		EmployeesCount: 120,
		
		ServiceAreas: []string{"Lagos", "Victoria Island", "Ikoyi"},
		ServiceRadius: 30.0,
		
		Logo: "https://cdn.tripsbook.com/logos/eko-convention.png",
		BannerImage: "https://cdn.tripsbook.com/banners/eko-convention.jpg",
		Gallery: []string{
			"https://cdn.tripsbook.com/images/events/eko/1.jpg",
			"https://cdn.tripsbook.com/images/events/eko/2.jpg",
		},
		
		IsVerified: true,
		VerificationStatus: models.VerificationStatusVerified,
		IsActive: true,
		IsFeatured: true,
		
		AdminPosition: 1,
		PositionCategory: "cat_events",
		
		AverageRating: 4.8,
		TotalReviews: 345,
		RatingBreakdown: map[int]int{
			5: 289, 4: 45, 3: 8, 2: 2, 1: 1,
		},
		
		TotalServices: 6,
		ActiveServices: 6,
		CompletedBookings: 892,
		
		CreatedAt: time.Now().AddDate(-3, 0, 0),
		UpdatedAt: time.Now().Add(-8 * time.Hour),
		LastActiveAt: time.Now().Add(-3 * time.Hour),
	},
	{
		ID: "sp_events_002",
		UserID: "user_402",
		BusinessName: "Landmark Event Center",
		DisplayName: "Landmark Event Center",
		Description: "Versatile event space for weddings, conferences, and social gatherings",
		CategoryID: "cat_events",
		SubCategories: []string{"events", "weddings", "conferences"},
		
		Phone: "+234-1-4630000",
		Email: "info@landmarkevents.com",
		Website: "https://landmarkevents.com",
		Address: "Water Corporation Road, Oniru, Victoria Island",
		City: "Lagos",
		State: "Lagos State",
		
		BusinessType: models.BusinessTypeCompany,
		EstablishedYear: 2016,
		EmployeesCount: 85,
		
		ServiceAreas: []string{"Lagos", "Victoria Island", "Lekki"},
		ServiceRadius: 25.0,
		
		Logo: "https://cdn.tripsbook.com/logos/landmark-events.png",
		BannerImage: "https://cdn.tripsbook.com/banners/landmark-events.jpg",
		Gallery: []string{
			"https://cdn.tripsbook.com/images/events/landmark/1.jpg",
		},
		
		IsVerified: true,
		VerificationStatus: models.VerificationStatusVerified,
		IsActive: true,
		IsFeatured: false,
		
		AdminPosition: 2,
		PositionCategory: "cat_events",
		
		AverageRating: 4.6,
		TotalReviews: 287,
		RatingBreakdown: map[int]int{
			5: 198, 4: 67, 3: 18, 2: 3, 1: 1,
		},
		
		TotalServices: 5,
		ActiveServices: 5,
		CompletedBookings: 678,
		
		CreatedAt: time.Now().AddDate(-4, 0, 0),
		UpdatedAt: time.Now().Add(-6 * time.Hour),
		LastActiveAt: time.Now().Add(-4 * time.Hour),
	},
}

// MockAdminPositioning represents admin positioning for providers
var MockAdminPositioning = []models.AdminPositioning{
	// Hotels Category Positioning
	{
		ID: "pos_hotel_001",
		CategoryID: "cat_hotels",
		ProviderID: "sp_hotel_001",
		Position: 1,
		IsFeatured: true,
		FeaturedReason: "Premium luxury hotel with excellent reviews",
		PriorityScore: 9.8,
		AdminNotes: "Top performer, maintain position 1",
		AlwaysOnTop: true,
		BoostFactor: 2.0,
		ValidFrom: time.Now().AddDate(-1, 0, 0),
		ValidUntil: nil,
		CreatedBy: "admin_001",
		CreatedAt: time.Now().AddDate(-1, 0, 0),
		UpdatedAt: time.Now().AddDate(-1, 0, 0),
	},
	{
		ID: "pos_hotel_002",
		CategoryID: "cat_hotels",
		ProviderID: "sp_hotel_002",
		Position: 2,
		IsFeatured: true,
		FeaturedReason: "Historic hotel with strong brand recognition",
		PriorityScore: 9.2,
		AdminNotes: "Good performance, featured status",
		BoostFactor: 1.5,
		ValidFrom: time.Now().AddDate(-2, 0, 0),
		ValidUntil: nil,
		CreatedBy: "admin_001",
		CreatedAt: time.Now().AddDate(-2, 0, 0),
		UpdatedAt: time.Now().AddDate(-2, 0, 0),
	},
	{
		ID: "pos_hotel_003",
		CategoryID: "cat_hotels",
		ProviderID: "sp_hotel_003",
		Position: 3,
		IsFeatured: false,
		FeaturedReason: "",
		PriorityScore: 8.7,
		AdminNotes: "Solid performer, maintain position",
		BoostFactor: 1.0,
		ValidFrom: time.Now().AddDate(-3, 0, 0),
		ValidUntil: nil,
		CreatedBy: "admin_001",
		CreatedAt: time.Now().AddDate(-3, 0, 0),
		UpdatedAt: time.Now().AddDate(-3, 0, 0),
	},
	
	// Restaurants Category Positioning
	{
		ID: "pos_restaurant_001",
		CategoryID: "cat_restaurants",
		ProviderID: "sp_restaurant_001",
		Position: 1,
		IsFeatured: true,
		FeaturedReason: "Authentic Nigerian cuisine with cultural experience",
		PriorityScore: 9.5,
		AdminNotes: "Top rated restaurant, maintain featured status",
		AlwaysOnTop: true,
		BoostFactor: 2.0,
		ValidFrom: time.Now().AddDate(-1, 0, 0),
		ValidUntil: nil,
		CreatedBy: "admin_001",
		CreatedAt: time.Now().AddDate(-1, 0, 0),
		UpdatedAt: time.Now().AddDate(-1, 0, 0),
	},
	{
		ID: "pos_restaurant_002",
		CategoryID: "cat_restaurants",
		ProviderID: "sp_restaurant_002",
		Position: 2,
		IsFeatured: false,
		FeaturedReason: "",
		PriorityScore: 9.1,
		AdminNotes: "Excellent ratings, good position",
		BoostFactor: 1.2,
		ValidFrom: time.Now().AddDate(-2, 0, 0),
		ValidUntil: nil,
		CreatedBy: "admin_001",
		CreatedAt: time.Now().AddDate(-2, 0, 0),
		UpdatedAt: time.Now().AddDate(-2, 0, 0),
	},
	
	// Transport Category Positioning
	{
		ID: "pos_transport_001",
		CategoryID: "cat_transport",
		ProviderID: "sp_transport_001",
		Position: 1,
		IsFeatured: true,
		FeaturedReason: "Premium service with highest ratings",
		PriorityScore: 9.7,
		AdminNotes: "Market leader, maintain top position",
		AlwaysOnTop: true,
		BoostFactor: 2.5,
		ValidFrom: time.Now().AddDate(-1, 0, 0),
		ValidUntil: nil,
		CreatedBy: "admin_001",
		CreatedAt: time.Now().AddDate(-1, 0, 0),
		UpdatedAt: time.Now().AddDate(-1, 0, 0),
	},
	{
		ID: "pos_transport_002",
		CategoryID: "cat_transport",
		ProviderID: "sp_transport_002",
		Position: 2,
		IsFeatured: false,
		FeaturedReason: "",
		PriorityScore: 8.9,
		AdminNotes: "Strong performer, good alternative",
		BoostFactor: 1.3,
		ValidFrom: time.Now().AddDate(-2, 0, 0),
		ValidUntil: nil,
		CreatedBy: "admin_001",
		CreatedAt: time.Now().AddDate(-2, 0, 0),
		UpdatedAt: time.Now().AddDate(-2, 0, 0),
	},
	
	// Shopping Category Positioning
	{
		ID: "pos_shopping_001",
		CategoryID: "cat_shopping",
		ProviderID: "sp_shopping_001",
		Position: 1,
		IsFeatured: true,
		FeaturedReason: "Premier shopping destination",
		PriorityScore: 9.3,
		AdminNotes: "Top shopping mall, maintain featured",
		AlwaysOnTop: true,
		BoostFactor: 2.0,
		ValidFrom: time.Now().AddDate(-1, 0, 0),
		ValidUntil: nil,
		CreatedBy: "admin_001",
		CreatedAt: time.Now().AddDate(-1, 0, 0),
		UpdatedAt: time.Now().AddDate(-1, 0, 0),
	},
	{
		ID: "pos_shopping_002",
		CategoryID: "cat_shopping",
		ProviderID: "sp_shopping_002",
		Position: 2,
		IsFeatured: false,
		FeaturedReason: "",
		PriorityScore: 8.8,
		AdminNotes: "Popular mall, good position",
		BoostFactor: 1.1,
		ValidFrom: time.Now().AddDate(-2, 0, 0),
		ValidUntil: nil,
		CreatedBy: "admin_001",
		CreatedAt: time.Now().AddDate(-2, 0, 0),
		UpdatedAt: time.Now().AddDate(-2, 0, 0),
	},
	
	// Events Category Positioning
	{
		ID: "pos_events_001",
		CategoryID: "cat_events",
		ProviderID: "sp_events_001",
		Position: 1,
		IsFeatured: true,
		FeaturedReason: "World-class convention center",
		PriorityScore: 9.6,
		AdminNotes: "Premium venue, maintain top position",
		AlwaysOnTop: true,
		BoostFactor: 2.2,
		ValidFrom: time.Now().AddDate(-1, 0, 0),
		ValidUntil: nil,
		CreatedBy: "admin_001",
		CreatedAt: time.Now().AddDate(-1, 0, 0),
		UpdatedAt: time.Now().AddDate(-1, 0, 0),
	},
	{
		ID: "pos_events_002",
		CategoryID: "cat_events",
		ProviderID: "sp_events_002",
		Position: 2,
		IsFeatured: false,
		FeaturedReason: "",
		PriorityScore: 9.0,
		AdminNotes: "Versatile venue, good positioning",
		BoostFactor: 1.4,
		ValidFrom: time.Now().AddDate(-2, 0, 0),
		ValidUntil: nil,
		CreatedBy: "admin_001",
		CreatedAt: time.Now().AddDate(-2, 0, 0),
		UpdatedAt: time.Now().AddDate(-2, 0, 0),
	},
}

// GetMockServiceProvidersByCategory returns providers for a specific category
func GetMockServiceProvidersByCategory(categoryID string) []models.ServiceProvider {
	var providers []models.ServiceProvider
	for _, provider := range MockServiceProviders {
		if provider.CategoryID == categoryID {
			providers = append(providers, provider)
		}
	}
	return providers
}

// GetMockProviderByID returns a specific provider by ID
func GetMockProviderByID(providerID string) *models.ServiceProvider {
	for _, provider := range MockServiceProviders {
		if provider.ID == providerID {
			return &provider
		}
	}
	return nil
}

// GetMockAdminPositioningByCategory returns admin positioning for a category
func GetMockAdminPositioningByCategory(categoryID string) []models.AdminPositioning {
	var positioning []models.AdminPositioning
	for _, pos := range MockAdminPositioning {
		if pos.CategoryID == categoryID {
			positioning = append(positioning, pos)
		}
	}
	return positioning
}

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/joho/godotenv"
	"github.com/mutsaevz/team-5-ambitious/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	batchSize = 1000

	usersCount    = 5000
	carsCount     = 2000
	tripsCount    = 5000
	bookingsCount = 10000
	reviewsCount  = 8000
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	dsn := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Clearing existing data...")
	db.Exec("TRUNCATE TABLE reviews, bookings, trips, cars, users RESTART IDENTITY CASCADE")
	fmt.Println("Data cleared ✓")

	gofakeit.Seed(0)

	userIDs := seedUsers(db)
	carIDs := seedCars(db, userIDs)
	tripIDs := seedTrips(db, userIDs, carIDs)
	bookingIDs := seedBookings(db, userIDs, tripIDs)
	seedReviews(db, userIDs, tripIDs)

	fmt.Println("\n=== Seeding completed ===")
	fmt.Printf("Users:     %d\n", len(userIDs))
	fmt.Printf("Cars:      %d\n", len(carIDs))
	fmt.Printf("Trips:     %d\n", len(tripIDs))
	fmt.Printf("Bookings:  %d\n", len(bookingIDs))
	fmt.Printf("Reviews:   %d\n", reviewsCount)
}

func seedUsers(db *gorm.DB) []uint {
	const total = usersCount
	users := make([]models.User, 0, batchSize)
	ids := make([]uint, 0, total)

	fmt.Printf("Seeding users... 0/%d", total)
	for i := 0; i < total; i++ {
		u := models.User{
			Name:    gofakeit.Name(),
			Phone:   gofakeit.Phone(),
			Balance: gofakeit.Number(0, 1000),
		}

		// ~5% soft deleted
		if gofakeit.Number(1, 100) <= 5 {
			u.DeletedAt = gorm.DeletedAt{
				Time:  gofakeit.DateRange(time.Now().AddDate(-1, 0, 0), time.Now()),
				Valid: true,
			}
		}

		users = append(users, u)

		if len(users) >= batchSize {
			db.Session(&gorm.Session{SkipHooks: true}).Create(&users)
			for _, u := range users {
				ids = append(ids, u.ID)
			}
			users = users[:0]
			fmt.Printf("\rSeeding users... %d/%d", i+1, total)
		}
	}

	if len(users) > 0 {
		db.Session(&gorm.Session{SkipHooks: true}).Create(&users)
		for _, u := range users {
			ids = append(ids, u.ID)
		}
	}
	fmt.Println(" ✓")
	return ids
}

func seedCars(db *gorm.DB, userIDs []uint) []uint {
	const total = carsCount
	cars := make([]models.Car, 0, batchSize)
	ids := make([]uint, 0, total)

	fmt.Printf("Seeding cars... 0/%d", total)
	brands := []string{"Toyota", "Honda", "BMW", "Mercedes", "Ford"}
	modelsList := []string{"Model A", "Model B", "Model C", "Model D"}

	for i := 0; i < total; i++ {
		c := models.Car{
			OwnerID:  userIDs[gofakeit.Number(0, len(userIDs)-1)],
			Brand:    brands[gofakeit.Number(0, len(brands)-1)],
			CarModel: modelsList[gofakeit.Number(0, len(modelsList)-1)],
			Seats:    gofakeit.Number(2, 7),
		}

		if gofakeit.Number(1, 100) <= 5 {
			c.DeletedAt = gorm.DeletedAt{
				Time:  gofakeit.DateRange(time.Now().AddDate(-1, 0, 0), time.Now()),
				Valid: true,
			}
		}

		cars = append(cars, c)

		if len(cars) >= batchSize {
			db.Session(&gorm.Session{SkipHooks: true}).Create(&cars)
			for _, c := range cars {
				ids = append(ids, c.ID)
			}
			cars = cars[:0]
			fmt.Printf("\rSeeding cars... %d/%d", i+1, total)
		}
	}

	if len(cars) > 0 {
		db.Session(&gorm.Session{SkipHooks: true}).Create(&cars)
		for _, c := range cars {
			ids = append(ids, c.ID)
		}
	}
	fmt.Println(" ✓")
	return ids
}

func seedTrips(db *gorm.DB, userIDs, carIDs []uint) []uint {
	const total = tripsCount
	trips := make([]models.Trip, 0, batchSize)
	ids := make([]uint, 0, total)

	fmt.Printf("Seeding trips... 0/%d", total)
	cities := []string{"Zagreb", "Split", "Dubrovnik", "Rijeka", "Osijek"}
	statuses := []string{"Scheduled", "Ongoing", "Completed", "Cancelled"}

	for i := 0; i < total; i++ {
		start := gofakeit.DateRange(time.Now().AddDate(-1, 0, 0), time.Now())
		totalSeats := gofakeit.Number(2, 7)
		t := models.Trip{
			DriverID:       userIDs[gofakeit.Number(0, len(userIDs)-1)],
			CarID:          carIDs[gofakeit.Number(0, len(carIDs)-1)],
			FromCity:       cities[gofakeit.Number(0, len(cities)-1)],
			ToCity:         cities[gofakeit.Number(0, len(cities)-1)],
			StartTime:      start,
			DurationMin:    gofakeit.Number(30, 300),
			TotalSeats:     totalSeats,
			AvailableSeats: totalSeats,
			Price:          gofakeit.Number(50, 500),
			TripStatus:     statuses[gofakeit.Number(0, len(statuses)-1)],
			AvgRating:      0.0,
		}

		if gofakeit.Number(1, 100) <= 5 {
			t.DeletedAt = gorm.DeletedAt{
				Time:  gofakeit.DateRange(time.Now().AddDate(-1, 0, 0), time.Now()),
				Valid: true,
			}
		}

		trips = append(trips, t)

		if len(trips) >= batchSize {
			db.Session(&gorm.Session{SkipHooks: true}).Create(&trips)
			for _, t := range trips {
				ids = append(ids, t.ID)
			}
			trips = trips[:0]
			fmt.Printf("\rSeeding trips... %d/%d", i+1, total)
		}
	}

	if len(trips) > 0 {
		db.Session(&gorm.Session{SkipHooks: true}).Create(&trips)
		for _, t := range trips {
			ids = append(ids, t.ID)
		}
	}
	fmt.Println(" ✓")
	return ids
}

func seedBookings(db *gorm.DB, userIDs, tripIDs []uint) []uint {
	const total = bookingsCount
	bookings := make([]models.Booking, 0, batchSize)
	ids := make([]uint, 0, total)
	statuses := []string{"Pending", "Approved", "Rejected"}

	fmt.Printf("Seeding bookings... 0/%d", total)
	for i := 0; i < total; i++ {
		b := models.Booking{
			TripID:        tripIDs[gofakeit.Number(0, len(tripIDs)-1)],
			PassengerID:   userIDs[gofakeit.Number(0, len(userIDs)-1)],
			BookingStatus: statuses[gofakeit.Number(0, len(statuses)-1)],
		}

		if gofakeit.Number(1, 100) <= 5 {
			b.DeletedAt = gorm.DeletedAt{
				Time:  gofakeit.DateRange(time.Now().AddDate(-1, 0, 0), time.Now()),
				Valid: true,
			}
		}

		bookings = append(bookings, b)

		if len(bookings) >= batchSize {
			db.Session(&gorm.Session{SkipHooks: true}).Create(&bookings)
			for _, b := range bookings {
				ids = append(ids, b.ID)
			}
			bookings = bookings[:0]
			fmt.Printf("\rSeeding bookings... %d/%d", i+1, total)
		}
	}

	if len(bookings) > 0 {
		db.Session(&gorm.Session{SkipHooks: true}).Create(&bookings)
		for _, b := range bookings {
			ids = append(ids, b.ID)
		}
	}
	fmt.Println(" ✓")
	return ids
}

func seedReviews(db *gorm.DB, userIDs, tripIDs []uint) {
	const total = reviewsCount
	reviews := make([]models.Review, 0, batchSize)

	fmt.Printf("Seeding reviews... 0/%d", total)
	for i := 0; i < total; i++ {
		r := models.Review{
			AuthorID: userIDs[gofakeit.Number(0, len(userIDs)-1)],
			TripID:   tripIDs[gofakeit.Number(0, len(tripIDs)-1)],
			Text:     gofakeit.Sentence(10),
			Rating:   gofakeit.Number(1, 5),
		}

		if gofakeit.Number(1, 100) <= 5 {
			r.DeletedAt = gorm.DeletedAt{
				Time:  gofakeit.DateRange(time.Now().AddDate(-1, 0, 0), time.Now()),
				Valid: true,
			}
		}

		reviews = append(reviews, r)

		if len(reviews) >= batchSize {
			db.Session(&gorm.Session{SkipHooks: true}).Create(&reviews)
			reviews = reviews[:0]
			fmt.Printf("\rSeeding reviews... %d/%d", i+1, total)
		}
	}

	if len(reviews) > 0 {
		db.Session(&gorm.Session{SkipHooks: true}).Create(&reviews)
	}
	fmt.Println(" ✓")
}

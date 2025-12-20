package models

import (
	"time"

	"gorm.io/gorm"
)

type Trip struct {
	gorm.Model

	DriverID       uint
	CarID          uint
	From           string
	To             string
	StartTime      time.Time
	DurationMin    int
	TotalSeats     int
	AvailableSeats int
	Price          int
	TripStatus     string
	AvgRating      float64
}

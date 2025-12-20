package models

import "gorm.io/gorm"

type Booking struct {
	gorm.Model

	TripID        uint
	PassengerID   uint
	BookingStatus string
}

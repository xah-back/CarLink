package models

import (
	"time"

	"github.com/mutsaevz/team-5-ambitious/internal/constants"
)

type Trip struct {
	Base

	DriverID       uint      `json:"driver_id"`
	CarID          uint      `json:"car_id"`
	FromCity       string    `json:"from_city"`
	ToCity         string    `json:"to_city"`
	StartTime      time.Time `json:"start_time"`
	DurationMin    int       `json:"duration_min"`
	TotalSeats     int       `json:"total_seats"`
	AvailableSeats int       `json:"available_seats"`
	Price          int       `json:"price"`
	TripStatus     string    `json:"trip_status"`
	AvgRating      float64   `json:"avg_rating"`
}

type TripCreateRequest struct {
	FromCity       string               `json:"from_city"`
	ToCity         string               `json:"to_city"`
	StartTime      time.Time            `json:"start_time"`
	DurationMin    int                  `json:"duration_min"`
	AvailableSeats int                  `json:"available_seats"`
	Price          int                  `json:"price"`
	TripStatus     constants.TripStatus `json:"trip_status"`
}

type TripFilter struct {
	FromCity       *string               `json:"from_city"`
	ToCity         *string               `json:"to_city"`
	StartTime      *time.Time            `json:"start_time"`
	AvailableSeats *int                  `json:"available_seats"`
	TripStatus     *constants.TripStatus `json:"trip_status"`
}

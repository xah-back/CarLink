package models

type Booking struct {
	Base

	TripID        uint   `json:"trip_id" gorm:"not null;index"`
	PassengerID   uint   `json:"passenger_id" gorm:"not null;index"`
	BookingStatus string `json:"booking_status" gorm:"type:varchar(50);not null;index"`
}

type BookingCreateRequest struct {
	TripID      uint `json:"trip_id" binding:"required"`
	PassengerID uint `json:"passenger_id" binding:"required"`
}

type BookingUpdateRequest struct {
	BookingStatus *string `json:"booking_status" binding:"required"`
}

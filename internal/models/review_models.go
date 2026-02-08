package models

type Review struct {
	Base

	AuthorID uint   `json:"author_id" gorm:"not null;index:"`
	TripID   uint   `json:"trip_id" gorm:"index:idx_trip_id_id,priority:1"`
	Text     string `json:"text" gorm:"type:text;not null"`
	Rating   int    `json:"rating" gorm:"not null;check:rating >= 1 AND rating <= 5"`
}

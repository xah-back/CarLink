package models

import "gorm.io/gorm"

type Review struct {
	gorm.Model

	AuthorID uint
	TripID   uint
	Text     string
	Rating   int
}

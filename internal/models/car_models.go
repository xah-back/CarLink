package models

import "gorm.io/gorm"

type Car struct {
	gorm.Model

	OwnerID  uint
	Brand    string
	CarModel string
	Seats    int
}

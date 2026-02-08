package models

import (
	"time"

	"gorm.io/gorm"
)

type Base struct {
	ID        uint           `json:"id" gorm:"index:idx_trip_id_id,priority:2"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

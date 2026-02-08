package models

type Page struct {
	Page     int
	PageSize int
	TripID   *uint
	AuthorID *uint
	LastID   *uint
}

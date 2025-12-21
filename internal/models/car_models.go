package models

type Car struct {
	Base

	OwnerID  uint   `json:"owner_id"`
	Brand    string `json:"brand"`
	CarModel string `json:"car_model"`
	Seats    int    `json:"seats"`
}

type CarCreateRequest struct {
	Brand    string `json:"brand"`
	CarModel string `json:"car_model"`
	Seats    int    `json:"seats"`
}

package dto

import "time"

type CreateReservationRequest struct {
	ZoneID       uint   `json:"zone_id" validate:"required"`
	LicensePlate string `json:"license_plate" validate:"required,max=15"`
}

type ReservationResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	ZoneID       uint      `json:"zone_id"`
	LicensePlate string    `json:"license_plate"`
	Status       string    `json:"status"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	TotalCost    float64   `json:"total_cost"`
}

type ReservationWithZoneResponse struct {
	ID           uint         `json:"id"`
	UserID       uint         `json:"user_id"`
	ZoneID       uint         `json:"zone_id"`
	Zone         ZoneResponse `json:"zone"`
	LicensePlate string       `json:"license_plate"`
	Status       string       `json:"status"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

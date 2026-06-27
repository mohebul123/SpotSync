package models

import (
	"time"
)

type Reservation struct {
	ID           uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint        `gorm:"not null" json:"user_id"`
	User         User        `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"user,omitempty"`
	ZoneID       uint        `gorm:"not null" json:"zone_id"`
	Zone         ParkingZone `gorm:"foreignKey:ZoneID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"zone,omitempty"`
	LicensePlate string      `gorm:"type:varchar(15);not null" json:"license_plate"`
	Status       string      `gorm:"type:varchar(50);default:'active';not null" json:"status"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

package models

import (
	"time"

	"gorm.io/gorm"
)

type Day struct {
	gorm.Model
	Code         string        `gorm:"unique;not null"`
	Description  string        `gorm:"size:255;not null"`
	StartAt      time.Time     `gorm:"not null"`
	EndAt        time.Time     `gorm:"not null"`
	Status       bool          `gorm:"default:true"`
	Appointments []Appointment `gorm:"foreignKey:DayID"`
}

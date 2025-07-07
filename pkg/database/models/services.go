package models

import (
	"time"

	"gorm.io/gorm"
)

type Service struct {
	gorm.Model
	Code                string               `gorm:"unique;not null"`
	Name                string               `gorm:"size:255;not null"`
	EstimatedTime       time.Time            `gorm:"not null"`
	Status              bool                 `gorm:"default:true"`
	AppointmentServices []AppointmentService `gorm:"foreignKey:ServiceID"`
}

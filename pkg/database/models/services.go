package models

import (
	"gorm.io/gorm"
)

type Service struct {
	gorm.Model
	Code                string               `gorm:"unique;not null"`
	Name                string               `gorm:"size:255;not null"`
	EstimatedTime       uint                 `gorm:"not null"`
	Status              bool                 `gorm:"default:true"`
	AppointmentServices []AppointmentService `gorm:"foreignKey:ServiceID"`
}

func ActiveService(db *gorm.DB) *gorm.DB {
	return db.Where("status = ?", true)
}

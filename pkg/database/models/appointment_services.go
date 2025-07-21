package models

import (
	"gorm.io/gorm"
)

type AppointmentService struct {
	gorm.Model
	ServiceID     uint
	Service       Service `gorm:"foreignKey:ServiceID"`
	AppointmentID uint
	Appointment   Appointment `gorm:"foreignKey:AppointmentID"`
}

package models

import (
	"gorm.io/gorm"
)

type AppointmentService struct {
	gorm.Model
	ServiceID     int
	Service       Service `gorm:"foreignKey:ServiceID"`
	AppointmentID int
	Appointment   Appointment `gorm:"foreignKey:AppointmentID"`
}

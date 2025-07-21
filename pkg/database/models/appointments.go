package models

import (
	"time"

	"gorm.io/gorm"
)

type Appointment struct {
	gorm.Model
	StartAt             time.Time `gorm:"not null"`
	EndAt               time.Time `gorm:"not null"`
	DayID               uint
	Day                 Day `gorm:"foreignKey:DayID"`
	UserID              uint
	User                User `gorm:"foreignKey:UserID"`
	EmployeeID          uint
	Employee            Employee             `gorm:"foreignKey:EmployeeID"`
	AppointmentServices []AppointmentService `gorm:"foreignKey:AppointmentID"`
}

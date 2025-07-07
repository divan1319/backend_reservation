package models

import (
	"time"

	"gorm.io/gorm"
)

type Appointment struct {
	gorm.Model
	StartAt             time.Time `gorm:"not null"`
	EndAt               time.Time `gorm:"not null"`
	DayID               int
	Day                 Day `gorm:"foreignKey:DayID"`
	UserID              int
	User                User `gorm:"foreignKey:UserID"`
	EmployeeID          int
	Employee            Employee             `gorm:"foreignKey:EmployeeID"`
	AppointmentServices []AppointmentService `gorm:"foreignKey:AppointmentID"`
}

package models

import (
	"gorm.io/gorm"
)

type Employee struct {
	gorm.Model
	Name         string `gorm:"unique;not null"`
	RoleID       int
	Role         Role          `gorm:"foreignKey:RoleID"`
	Status       bool          `gorm:"default:true"`
	Appointments []Appointment `gorm:"foreignKey:EmployeeID"`
}

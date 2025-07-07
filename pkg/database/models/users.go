package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name         string `gorm:"size:255;not null"`
	Password     string `gorm:"size:255;not null"`
	Phone        int    `gorm:"unique;not null"`
	Email        string `gorm:"unique;not null"`
	RoleID       int
	Role         Role          `gorm:"foreignKey:RoleID"`
	Appointments []Appointment `gorm:"foreignKey:UserID"`
}

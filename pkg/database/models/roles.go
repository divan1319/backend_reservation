package models

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Code        string `gorm:"unique;not null"`
	Description string `gorm:"size:255;not null"`
	UserID      int
}

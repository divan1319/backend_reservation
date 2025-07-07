package migrations

import (
	"backend_reservation/pkg/database/models"
	"fmt"
	"log"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	log.Println("Ejecutando migraciones...")

	modelsToMigrate := []interface{}{
		&models.User{},
		&models.Employee{},
		&models.Role{},
		&models.Service{},
		&models.Day{},
		&models.Appointment{},
		&models.AppointmentService{},
	}

	err := db.AutoMigrate(modelsToMigrate...)
	if err != nil {
		return fmt.Errorf("error al ejecutar las migraciones: %v", err)
	}

	log.Println("Migraciones ejecutadas correctamente")
	return nil
}

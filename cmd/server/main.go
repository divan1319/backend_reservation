package main

import (
	"backend_reservation/pkg/database/connection"
	"backend_reservation/pkg/database/models"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno desde el archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Printf("advertencia: no se pudo cargar el archivo .env: %v", err)
	}

	database, gormDB, err := connection.InitDB()
	if err != nil {
		log.Fatalf("error al inicializar la base de datos: %v", err)
	}

	gormDB.AutoMigrate(models.User{}, models.Employee{}, models.Role{}, models.Service{}, models.Day{}, models.Appointment{}, models.AppointmentService{})

	defer database.Close()

}

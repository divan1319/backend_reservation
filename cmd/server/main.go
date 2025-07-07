package main

import (
	"backend_reservation/pkg/database/connection"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno desde el archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Printf("advertencia: no se pudo cargar el archivo .env: %v", err)
	}

	database, _, err := connection.InitDB()
	if err != nil {
		log.Fatalf("error al inicializar la base de datos: %v", err)
	}

	defer database.Close()

}

package main

import (
	"backend_reservation/pkg/database/connection"
	"backend_reservation/pkg/database/migrations"
	"flag"
	"log"

	"github.com/joho/godotenv"
)

func main() {

	var migrate = flag.Bool("migrate", false, "Ejecutar las migraciones")

	flag.Parse()

	err := godotenv.Load()

	if err != nil {
		log.Printf("advertencia: no se pudo cargar el archivo .env: %v", err)
	}

	_, gormDB, err := connection.ConnectDB()

	if err != nil {
		log.Fatalf("error al inicializar la base de datos: %v", err)
	}

	sqlDB, err := gormDB.DB()

	if err != nil {
		log.Fatalf("error al obtener la conexión a la base de datos: %v", err)
	}

	defer sqlDB.Close()

	if *migrate {
		if err := migrations.RunMigrations(gormDB); err != nil {
			log.Fatalf("error al ejecutar las migraciones: %v", err)
		}
	}
}

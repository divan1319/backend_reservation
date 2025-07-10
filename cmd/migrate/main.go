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

	_, gormDB, err := connection.GetDB()
	if err != nil {
		log.Fatalf("error al inicializar la base de datos: %v", err)
	}

	// Configurar cierre de conexión al terminar
	defer func() {
		if err := connection.CloseDB(); err != nil {
			log.Printf("error al cerrar la conexión a la base de datos: %v", err)
		}
	}()

	if *migrate {
		if err := migrations.RunMigrations(gormDB); err != nil {
			log.Fatalf("error al ejecutar las migraciones: %v", err)
		}
		log.Println("Migraciones ejecutadas correctamente")
	}
}

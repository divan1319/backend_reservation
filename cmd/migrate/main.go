package main

import (
	"backend_reservation/pkg/database/connection"
	"backend_reservation/pkg/database/migrations"
	"flag"
	"log"

	"github.com/joho/godotenv"
)

// main es la función principal del programa de migración de la base de datos.
// Proporciona la funcionalidad para ejecutar migraciones de la base de datos a través de un flag.
//
// El proceso es el siguiente:
// 1. Define y parsea el flag -migrate para determinar si se deben ejecutar las migraciones
// 2. Carga las variables de entorno desde el archivo .env
// 3. Inicializa la conexión a la base de datos
// 4. Configura el cierre de la conexión al terminar la ejecución
// 5. Si el flag -migrate está activo, ejecuta las migraciones
//
// Los posibles errores que maneja son:
// - Error al cargar el archivo .env (advertencia)
// - Error al inicializar la base de datos (fatal)
// - Error al ejecutar las migraciones (fatal)
// - Error al cerrar la conexión (log)
func main() {
	// Definir y parsear el flag -migrate para determinar si se deben ejecutar las migraciones
	var migrate = flag.Bool("migrate", false, "Ejecutar las migraciones")

	flag.Parse()

	// Cargar variables de entorno desde el archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Printf("advertencia: no se pudo cargar el archivo .env: %v", err)
	}

	// Inicializar conexión a la base de datos
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

	// Si el flag -migrate está activo, ejecutar las migraciones
	if *migrate {
		if err := migrations.RunMigrations(gormDB); err != nil {
			log.Fatalf("error al ejecutar las migraciones: %v", err)
		}
		log.Println("Migraciones ejecutadas correctamente")
	}
}

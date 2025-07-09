package main

import (
	"backend_reservation/internal/infrastructure/web/routes"
	"backend_reservation/pkg/database/connection"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno desde el archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Printf("advertencia: no se pudo cargar el archivo .env: %v", err)
	}

	database, _, err := connection.ConnectDB()
	if err != nil {
		log.Fatalf("error al inicializar la base de datos: %v", err)
	}

	defer database.Close()

	port := os.Getenv("PORT")
	router := routes.MainRouter()

	server := &http.Server{
		Addr:    port,
		Handler: router,
	}

	fmt.Printf("Servidor corriendo en el puerto %s", port)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("error al iniciar el servidor: %v", err)
	}

}

package main

import (
	"backend_reservation/internal/infrastructure/web/routes"
	"backend_reservation/pkg/database/connection"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno desde el archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Printf("advertencia: no se pudo cargar el archivo .env: %v", err)
	}

	// Inicializar conexión a base de datos
	_, _, err = connection.GetDB()
	if err != nil {
		log.Fatalf("error al inicializar la base de datos: %v", err)
	}

	// Configurar cierre graceful
	defer func() {
		if err := connection.CloseDB(); err != nil {
			log.Printf("error al cerrar la conexión a la base de datos: %v", err)
		}
	}()

	port := os.Getenv("PORT")
	router := routes.MainRouter()

	server := &http.Server{
		Addr:    port,
		Handler: router,
	}

	// Canal para manejar señales del sistema
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Ejecutar servidor en una goroutine
	go func() {
		fmt.Printf("Servidor corriendo en el puerto %s\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error al iniciar el servidor: %v", err)
		}
	}()

	// Esperar señal de cierre
	<-quit
	log.Println("Cerrando servidor...")

	// Crear contexto con timeout para el cierre graceful
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Cerrar servidor de manera graceful
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("error al cerrar el servidor: %v", err)
	}

	log.Println("Servidor cerrado correctamente")
}

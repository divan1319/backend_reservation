package main

import (
	"backend_reservation/internal/infrastructure/web/middleware"
	"backend_reservation/internal/infrastructure/web/routes"
	"backend_reservation/pkg/database/connection"
	"backend_reservation/pkg/firmador"
	"backend_reservation/pkg/logger"
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

// main es el punto de entrada de la aplicación del servidor HTTP.
// Se encarga de inicializar las dependencias principales (variables de entorno, base de datos, logger, firmador de tokens),
// configurar los middlewares de seguridad (CORS, rate limiting), y arrancar el servidor HTTP.
// Además, implementa un cierre graceful para asegurar que los recursos se liberen correctamente al finalizar.
func main() {
	// Cargar variables de entorno desde el archivo .env (si existe).
	// Esto permite configurar la aplicación sin hardcodear valores sensibles.
	err := godotenv.Load()
	if err != nil {
		log.Printf("advertencia: no se pudo cargar el archivo .env: %v", err)
	}

	// Inicializar la conexión a la base de datos.
	// Si ocurre un error crítico, se detiene la ejecución.
	_, _, err = connection.GetDB()
	if err != nil {
		log.Fatalf("error al inicializar la base de datos: %v", err)
	}

	// Configurar el cierre graceful de la base de datos.
	// Esta función diferida se ejecutará al finalizar main, cerrando la conexión de forma segura.
	defer func() {
		if err := connection.CloseDB(); err != nil {
			log.Printf("error al cerrar la conexión a la base de datos: %v", err)
		}
	}()

	// Inicializar el firmador de tokens PASETO.
	// Esto prepara la infraestructura para la autenticación basada en tokens.
	firmador.InitPaseto()

	// Configuración del logger a partir de variables de entorno.
	// Permite controlar el nivel de log, el archivo de salida y la rotación de logs.
	config := logger.Config{
		Environment: os.Getenv("APP_ENV"),
		Level:       os.Getenv("LOG_LEVEL"),
		Rotation: logger.RotationConfig{
			Filename:   os.Getenv("LOG_FILE"),
			MaxSize:    10,   // Tamaño máximo del archivo de log en MB antes de rotar
			MaxBackups: 3,    // Número máximo de archivos de backup de log
			MaxAge:     30,   // Días máximos a conservar los logs
			Compress:   true, // Comprimir logs antiguos
		},
	}
	// Inicializar el logger global de la aplicación.
	logger.InitLogger(config)

	// Obtener el puerto de escucha del servidor desde las variables de entorno.
	port := os.Getenv("PORT")

	// Inicializar el middleware de rate limiting (limitador de solicitudes por IP).
	// Permite 15 solicitudes por IP cada 120 segundos.
	rateLimiter := middleware.NewRateLimiter(15, 120*time.Second)
	defer rateLimiter.Stop() // Asegura que el goroutine de reseteo se detenga al cerrar el servidor.

	// Inicializar el router principal de la aplicación (todas las rutas y handlers).
	router := routes.MainRouter()

	// Se encadenan los middlewares de seguridad para el servidor HTTP siguiendo el patrón de decorador.
	// El orden de ejecución de los middlewares es de adentro hacia afuera en la declaración:
	//
	// 1. router: El router principal que maneja todas las rutas de la aplicación
	// 2. rateLimiter.Throttle(): Middleware de limitación de tasa que envuelve al router
	//    - Controla la cantidad de solicitudes por IP (15 solicitudes cada 120 segundos)
	//    - Si se excede el límite, retorna HTTP 429 (Too Many Requests) sin procesar la solicitud
	// 3. middleware.Cors(): Middleware de CORS que envuelve al rate limiter
	//    - Valida que el origen de la solicitud esté en la lista de orígenes permitidos
	//    - Configura los headers CORS necesarios para el intercambio de recursos
	//    - Maneja las solicitudes preflight (OPTIONS)
	//
	// Flujo de ejecución para cada solicitud HTTP:
	// Solicitud → CORS → Rate Limiter → Router → Handler específico → Respuesta
	//
	// El orden es crítico: CORS debe ejecutarse primero para rechazar solicitudes de orígenes no autorizados
	// antes de que consuman recursos del rate limiter, optimizando el rendimiento y la seguridad.
	secureMux := middleware.Cors(rateLimiter.Throttle(router))

	// Configurar el servidor HTTP con el handler seguro y el puerto especificado.
	server := &http.Server{
		Addr:    port,
		Handler: secureMux,
	}

	// Crear un canal para recibir señales del sistema (SIGINT, SIGTERM) y permitir un cierre graceful.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Ejecutar el servidor en una goroutine para no bloquear el hilo principal.
	go func() {
		fmt.Printf("Servidor corriendo en el puerto %s\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error al iniciar el servidor: %v", err)
		}
	}()

	// Esperar a recibir una señal de cierre (Ctrl+C o kill).
	<-quit
	log.Println("Cerrando servidor...")

	// Crear un contexto con timeout de 30 segundos para el cierre graceful del servidor.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Intentar cerrar el servidor de forma ordenada, permitiendo finalizar las conexiones activas.
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("error al cerrar el servidor: %v", err)
	}

	log.Println("Servidor cerrado correctamente")
}

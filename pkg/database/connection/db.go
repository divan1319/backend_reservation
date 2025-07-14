package connection

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbInstance   *sql.DB
	gormInstance *gorm.DB
	once         sync.Once
	initError    error
)

// GetDB devuelve la instancia singleton de la base de datos
// GetDB devuelve las instancias singleton de la base de datos SQL y GORM, inicializándolas solo una vez.
// Utiliza sync.Once para asegurar que la conexión se establezca una única vez durante el ciclo de vida de la aplicación.
// Retorna:
//   - *sql.DB: instancia de la base de datos SQL estándar
//   - *gorm.DB: instancia de la base de datos usando GORM
//   - error: error de inicialización, si ocurrió alguno
func GetDB() (*sql.DB, *gorm.DB, error) {
	once.Do(func() {
		dbInstance, gormInstance, initError = connectDB()
	})
	return dbInstance, gormInstance, initError
}

// connectDB establece la conexión inicial (función privada)
func connectDB() (*sql.DB, *gorm.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, nil, fmt.Errorf("error al inicializar el manejador de la base de datos: %v", err)
	}

	// Configurar pool de conexiones
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	if err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("error al inicializar el manejador de la base de datos con gorm: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pingErr := db.PingContext(ctx)
	if pingErr != nil {
		db.Close()
		return nil, nil, fmt.Errorf("error al conectar a la base de datos: %v", pingErr)
	}

	fmt.Println("Conexión exitosa a la base de datos")
	return db, gormDB, nil
}

// CloseDB cierra la conexión (llamar solo al finalizar la aplicación)
func CloseDB() error {
	if dbInstance != nil {
		return dbInstance.Close()
	}
	return nil
}

// ConnectDB - mantener por compatibilidad pero marcar como deprecated
// Deprecated: Usa GetDB() en su lugar
func ConnectDB() (*sql.DB, *gorm.DB, error) {
	return GetDB()
}

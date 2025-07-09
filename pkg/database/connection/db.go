package connection

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*sql.DB, *gorm.DB, error) {

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

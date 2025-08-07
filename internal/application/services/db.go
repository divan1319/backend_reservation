package services

import (
	"backend_reservation/pkg/database/connection"

	"gorm.io/gorm"
)

// connectDB conecta a la base de datos.
// Retorna una conexión a la base de datos o un error si ocurre algún problema.
//
// El proceso es el siguiente:
// 1. Obtiene una conexión a la base de datos.
// 2. Retorna la conexión a la base de datos o un error si ocurre algún problema.
func ConnectDB() (*gorm.DB, error) {
	_, gormDB, err := connection.GetDB()
	if err != nil {
		return nil, err
	}
	return gormDB, nil
}

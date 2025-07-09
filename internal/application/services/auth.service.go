package services

import (
	"backend_reservation/internal/application/dto"
	"backend_reservation/pkg/database/connection"
	"backend_reservation/pkg/database/models"
	"errors"
)

func Login(loginDto *dto.LoginDTO) (*models.User, error) {
	_, gormDB, _ := connection.ConnectDB()

	var user models.User

	result := gormDB.Find(&user, models.User{Email: loginDto.Email})

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("usuario no encontrado")
	}

	return &user, nil
}

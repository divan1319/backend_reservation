package services

import (
	"backend_reservation/internal/application/dto"
	"backend_reservation/pkg/database/connection"
	"backend_reservation/pkg/database/models"
	"backend_reservation/pkg/utils"
	"errors"
)

func Login(loginDto *dto.LoginDTO) (*models.User, error) {
	_, gormDB, err := connection.GetDB()
	if err != nil {
		return nil, err
	}

	var user models.User

	result := gormDB.Where("email = ?", loginDto.Email).First(&user)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, result.Error
	}

	if !utils.ComparePassword(user.Password, loginDto.Password) {
		return nil, errors.New("contraseña incorrecta")
	}

	return &user, nil
}

func Register(registerDto *dto.RegisterDTO) (*models.User, error) {
	_, gormDB, err := connection.GetDB()
	if err != nil {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(registerDto.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Name:     registerDto.Name,
		Email:    registerDto.Email,
		Password: hashedPassword,
		RoleID:   2,
		Phone:    registerDto.Phone,
	}

	result := gormDB.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("no se pudo crear el usuario")
	}

	return &user, nil
}

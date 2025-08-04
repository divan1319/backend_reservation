package middleware

import (
	"backend_reservation/pkg/database/connection"
	"backend_reservation/pkg/database/models"
	"strconv"
)

func HasPermission(userID string, code string) (bool, error) {
	_, gormDB, err := connection.GetDB()
	if err != nil {
		return false, err
	}

	// Convertir userID a int
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return false, err
	}

	// Buscar el rol de administrador
	var role models.Role
	if err := gormDB.Where("code = ?", code).First(&role).Error; err != nil {
		return false, err
	}

	// Buscar el usuario
	var user models.User
	if err := gormDB.Where("id = ?", userIDInt).First(&user).Error; err != nil {
		return false, err
	}

	// Verificar si el usuario tiene el rol correspondiente
	return role.ID == user.RoleID, nil
}

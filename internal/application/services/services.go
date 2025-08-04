package services

import (
	"backend_reservation/internal/application/dto"
	"backend_reservation/pkg/database/models"
	"errors"

	"gorm.io/gorm"
)

func CrearServicio(servicio dto.Service) (*models.Service, error) {
	gormDB, err := ConnectDB()

	if err != nil {
		return nil, err
	}

	defer func() {
		if sqlDB, err := gormDB.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	service := models.Service{
		Name:          servicio.Name,
		Code:          servicio.Code,
		EstimatedTime: servicio.EstimatedTime,
	}

	result := gormDB.Create(&service)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("No se pudo crear el servicio")
	}

	return &service, nil
}

func ActualizarServicio(id uint, servicio dto.Service) (*models.Service, error) {
	gormDB, err := ConnectDB()

	if err != nil {
		return nil, err
	}

	defer func() {
		if sqlDB, err := gormDB.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	service := models.Service{
		Model:         gorm.Model{ID: id},
		Name:          servicio.Name,
		Code:          servicio.Code,
		EstimatedTime: servicio.EstimatedTime,
		Status:        servicio.Status,
	}

	result := gormDB.Save(&service)

	if result.Error != nil {
		return nil, errors.New("Hubo un error al actualizar")
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("No se pudo actualizar el servicio")
	}

	return &service, nil

}

func ActivarDesactivarServicio(id uint) (*models.Service, error) {
	gormDB, err := ConnectDB()

	if err != nil {
		return nil, err
	}

	defer func() {
		if sqlDB, err := gormDB.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	var servicioModel models.Service

	if err := gormDB.First(&servicioModel, id).Error; err != nil {
		return nil, errors.New("No se pudo encontrar el servicio")
	}

	servicioModel.Status = !servicioModel.Status

	if err := gormDB.Save(&servicioModel).Error; err != nil {
		return nil, errors.New("No se pudo actualizar el status del servicio")
	}

	return &servicioModel, nil

}

func EliminarServicio(id uint) (bool, error) {
	gormDB, err := ConnectDB()

	if err != nil {
		return false, err
	}

	defer func() {
		if sqlDB, err := gormDB.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	if err := gormDB.Delete(&models.Service{}, id).Error; err != nil {
		return false, errors.New("No se pudo eliminar el servicio")
	}

	return true, nil
}

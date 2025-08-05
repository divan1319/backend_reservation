package services

import (
	"backend_reservation/internal/application/dto"
	"backend_reservation/pkg/database/models"
	"errors"

	"gorm.io/gorm"
)

func ObtenerServicios() ([]models.Service, error) {
	gormDB, err := ConnectDB()

	if err != nil {
		return nil, err
	}
	var servicios []models.Service

	if err := gormDB.Find(&servicios).Error; err != nil {
		return nil, err
	}

	return servicios, nil
}

func CrearServicio(servicio *dto.Service) (*models.Service, error) {
	gormDB, err := ConnectDB()

	if err != nil {
		return nil, err
	}

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

func ObtenerServicio(id uint) (*models.Service, error) {
	gormDB, err := ConnectDB()

	if err != nil {
		return nil, err
	}

	var service models.Service
	result := gormDB.First(&service, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("No se encontró el servicio")
		}
		return nil, result.Error
	}

	return &service, nil
}

func ActualizarServicio(id uint, servicio *dto.Service) (*models.Service, error) {
	gormDB, err := ConnectDB()

	if err != nil {
		return nil, err
	}

	// First, get the existing service
	var service models.Service
	result := gormDB.First(&service, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("No se encontró el servicio")
		}
		return nil, result.Error
	}

	// Only update fields that are not empty in the DTO
	if servicio.Name != "" {
		service.Name = servicio.Name
	}

	if servicio.Code != "" {
		service.Code = servicio.Code
	}

	if servicio.EstimatedTime != 0 {
		service.EstimatedTime = servicio.EstimatedTime
	}

	result = gormDB.Save(&service)

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

	if err := gormDB.Delete(&models.Service{}, id).Error; err != nil {
		return false, errors.New("No se pudo eliminar el servicio")
	}

	return true, nil
}

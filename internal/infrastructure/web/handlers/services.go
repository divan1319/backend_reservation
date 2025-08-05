package handlers

import (
	"backend_reservation/internal/application/dto"
	"backend_reservation/internal/application/services"
	"backend_reservation/pkg/handler"
	"net/http"
	"strconv"
)

func ObtenerServiciosHandler(w http.ResponseWriter, r *http.Request) {
	servicios, err := services.ObtenerServicios()

	if err != nil {
		handler.Error(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	dataServicios := make([]map[string]interface{}, len(servicios))

	for i, servicio := range servicios {
		dataServicios[i] = map[string]interface{}{
			"id":             servicio.ID,
			"name":           servicio.Name,
			"code":           servicio.Code,
			"estimated_time": servicio.EstimatedTime,
			"status":         servicio.Status,
		}
	}

	handler.Success(w, r, "", dataServicios)
}

func CrearServicioHandler(w http.ResponseWriter, r *http.Request) {
	estimatedTime, err := strconv.Atoi(r.FormValue("estimated_time"))

	if err != nil {
		handler.Error(w, r, http.StatusBadRequest, "Tiempo estimado no válido")
		return
	}

	registerService := dto.Service{
		Name:          r.FormValue("name"),
		Code:          r.FormValue("code"),
		EstimatedTime: uint(estimatedTime),
	}

	register, err := services.CrearServicio(&registerService)

	if err != nil {
		handler.Error(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	dataRegister := map[string]interface{}{
		"id":             register.ID,
		"name":           register.Name,
		"code":           register.Code,
		"estimated_time": register.EstimatedTime,
		"status":         register.Status,
	}
	handler.Success(w, r, "", dataRegister)
}

func ObtenerServicioHandler(w http.ResponseWriter, r *http.Request) {
	// Implement the logic to retrieve a service by ID
	serviceId := r.PathValue("id")

	if serviceId == "" {
		handler.Error(w, r, http.StatusBadRequest, "ID de servicio no proporcionado")
		return
	}

	parseServiceId, err := strconv.Atoi(serviceId)

	if err != nil {
		handler.Error(w, r, http.StatusBadRequest, "ID de servicio no válido")
		return
	}

	servicio, err := services.ObtenerServicio(uint(parseServiceId))

	if err != nil {
		handler.Error(w, r, http.StatusNotFound, err.Error())
		return
	}

	dataServicio := map[string]interface{}{
		"name":           servicio.Name,
		"code":           servicio.Code,
		"estimated_time": servicio.EstimatedTime,
		"status":         servicio.Status,
	}
	handler.Success(w, r, "", dataServicio)
}

func ActualizarServicioHandler(w http.ResponseWriter, r *http.Request) {

}

package handlers

import (
	"backend_reservation/internal/application/dto"
	"backend_reservation/internal/application/services"
	"backend_reservation/internal/domain"
	"backend_reservation/pkg/firmador"
	"backend_reservation/pkg/handler"
	"net/http"
	"strconv"
	"time"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	loginDto := dto.LoginDTO{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	user, err := services.Login(&loginDto)

	if err != nil {
		handler.Error(w, r, http.StatusNotFound, "Login failed")
		return
	}
	data := map[string]string{
		"user_id": strconv.Itoa(int(user.ID)),
		"email":   user.Email,
		"name":    user.Name,
		"role_id": strconv.Itoa(int(user.RoleID)),
	}

	token, err := firmador.FirmarToken(data, 1440*time.Minute) //token valido por 24 horas

	if err != nil {
		handler.Error(w, r, http.StatusInternalServerError, "No se pudo firmar el token")
		return
	}
	dataUser := domain.User{
		Name:  user.Name,
		Email: user.Email,
	}

	returnData := map[string]interface{}{
		"token": token,
		"user":  dataUser,
	}

	handler.Success(w, r, "Login successful", returnData)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	registerDto := dto.RegisterDTO{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
		Name:     r.FormValue("name"),
		Phone:    r.FormValue("phone"),
	}

	user, err := services.Register(&registerDto)

	if err != nil {
		handler.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	dataUser := domain.User{
		Name:  user.Name,
		Phone: user.Phone,
		Email: user.Email,
	}

	handler.Success(w, r, "Register successful", dataUser)
}

package handlers

import (
	"backend_reservation/internal/application/dto"
	"backend_reservation/internal/application/services"
	"backend_reservation/internal/domain"
	"backend_reservation/internal/infrastructure/web/middleware"
	"backend_reservation/pkg/firmador"
	"backend_reservation/pkg/handler"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// parseLoginData parsea los datos de login desde form-data o JSON
func parseLoginData(r *http.Request) (*dto.LoginDTO, error) {
	contentType := r.Header.Get("Content-Type")

	if strings.Contains(contentType, "application/json") {
		var loginDto dto.LoginDTO
		if err := json.NewDecoder(r.Body).Decode(&loginDto); err != nil {
			return nil, err
		}
		return &loginDto, nil
	}

	// Default: form-data
	loginDto := &dto.LoginDTO{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	return loginDto, nil
}

// parseRegisterData parsea los datos de registro desde form-data o JSON
func parseRegisterData(r *http.Request) (*dto.RegisterDTO, error) {
	contentType := r.Header.Get("Content-Type")

	if strings.Contains(contentType, "application/json") {
		var registerDto dto.RegisterDTO
		if err := json.NewDecoder(r.Body).Decode(&registerDto); err != nil {
			return nil, err
		}
		return &registerDto, nil
	}

	// Default: form-data
	registerDto := &dto.RegisterDTO{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
		Name:     r.FormValue("name"),
		Phone:    r.FormValue("phone"),
	}
	return registerDto, nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	loginDto, err := parseLoginData(r)
	if err != nil {
		handler.Error(w, r, http.StatusBadRequest, "Invalid request data")
		return
	}

	user, err := services.Login(loginDto)

	if err != nil {
		handler.Error(w, r, http.StatusNotFound, "Login failed")
		return
	}
	data := map[string]string{
		"user_id": strconv.Itoa(int(user.ID)),
		"email":   user.Email,
		"name":    user.Name,
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

	registerDto, err := parseRegisterData(r)
	if err != nil {
		handler.Error(w, r, http.StatusBadRequest, "Invalid request data")
		return
	}

	user, err := services.Register(registerDto)

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

func GetUserDataHandler(w http.ResponseWriter, r *http.Request) {
	userId, _ := middleware.GetUserIDFromContext(r.Context())

	parseUserId, err := strconv.Atoi(userId)

	if err != nil {
		handler.Error(w, r, http.StatusBadRequest, "Invalid user id")
		return
	}

	user, err, code := services.CheckUser(uint(parseUserId))

	if err != nil {
		handler.Error(w, r, code, err.Error())
		return
	}

	userData := map[string]any{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role.Code,
	}

	handler.Success(w, r, "User data retrieved successfully", userData)
}

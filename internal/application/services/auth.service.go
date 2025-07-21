package services

import (
	"backend_reservation/internal/application/dto"
	"backend_reservation/pkg/database/connection"
	"backend_reservation/pkg/database/models"
	"backend_reservation/pkg/utils"
	"errors"

	"gorm.io/gorm"
)

// connectDB conecta a la base de datos.
// Retorna una conexión a la base de datos o un error si ocurre algún problema.
//
// El proceso es el siguiente:
// 1. Obtiene una conexión a la base de datos.
// 2. Retorna la conexión a la base de datos o un error si ocurre algún problema.
func connectDB() (*gorm.DB, error) {
	_, gormDB, err := connection.GetDB()
	if err != nil {
		return nil, err
	}
	return gormDB, nil
}

// Login autentica a un usuario en la base de datos.
// Recibe un puntero a LoginDTO con los datos del usuario a autenticar.
// Retorna un puntero al modelo User autenticado o un error si ocurre algún problema.
//
// El proceso es el siguiente:
// 1. Conecta a la base de datos.
// 2. Busca un usuario con el email proporcionado.
func Login(loginDto *dto.LoginDTO) (*models.User, error) {
	gormDB, err := connectDB()
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

// Register registra un nuevo usuario en la base de datos.
// Recibe un puntero a RegisterDTO con los datos del usuario a registrar.
// Retorna un puntero al modelo User creado o un error si ocurre algún problema.
//
// El proceso es el siguiente:
// 1. Conecta a la base de datos.
// 2. Verifica si ya existe un usuario con el email proporcionado.
// 3. Si el usuario ya existe, retorna un error.
// 4. Hashea la contraseña proporcionada.
// 5. Obtiene el rol "user" desde la base de datos.
// 6. Crea el usuario con los datos proporcionados y el rol obtenido.
// 7. Retorna el usuario creado o un error si ocurre algún problema.
func Register(registerDto *dto.RegisterDTO) (*models.User, error) {
	gormDB, err := connectDB()
	if err != nil {
		return nil, err
	}

	// Verificar si el usuario ya existe
	userExists, err := VerifyUser(registerDto.Email, registerDto.Phone, gormDB)
	if err != nil {
		return nil, err
	}

	if userExists {
		return nil, errors.New("el usuario ya existe")
	}

	// Hashear la contraseña
	hashedPassword, err := utils.HashPassword(registerDto.Password)
	if err != nil {
		return nil, err
	}

	// Obtener el rol del usuario con código "user"
	role := models.Role{}
	resultRole := gormDB.Where("code = ?", "user").First(&role)
	if resultRole.Error != nil {
		return nil, errors.New("error al obtener el rol del usuario")
	}

	// Crear el usuario
	user := models.User{
		Name:     registerDto.Name,
		Email:    registerDto.Email,
		Password: hashedPassword,
		RoleID:   role.ID,
		Phone:    registerDto.Phone,
	}

	result := gormDB.Create(&user)
	if result.Error != nil {
		return nil, errors.New("error al crear el usuario")
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("no se pudo crear el usuario")
	}

	return &user, nil
}

// VerifyUser verifica si un usuario ya existe en la base de datos.
// Recibe el email del usuario y una conexión a la base de datos.
// Retorna true si el usuario existe, false en caso contrario o un error si ocurre algún problema.
//
// El proceso es el siguiente:
// 1. Busca un usuario con el email proporcionado.
// 2. Retorna true si se encuentra un usuario, false en caso contrario.
func VerifyUser(email string, phone string, gormDB *gorm.DB) (bool, error) {

	var user models.User
	result := gormDB.Where("email = ? OR phone = ?", email, phone).First(&user)

	if result.Error != nil && result.Error.Error() == "record not found" {
		return false, nil
	}

	return true, nil
}

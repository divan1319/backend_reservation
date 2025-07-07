package models

import (
	"database/sql"
	"time"
)

// Role representa la tabla 'roles' en la base de datos.
type Role struct {
	ID          int64     `json:"id,omitempty" db:"id"`
	Code        string    `json:"code,omitempty" db:"code"`
	Description string    `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// User representa la tabla 'users' en la base de datos.
type User struct {
	ID        int64     `json:"id,omitempty" db:"id"`
	Name      string    `json:"name,omitempty" db:"name"`
	Password  string    `json:"-" db:"password"` // Omitir siempre en JSON por seguridad
	Phone     int       `json:"phone,omitempty" db:"phone"`
	Email     string    `json:"email,omitempty" db:"email"`
	RoleID    int64     `json:"role_id,omitempty" db:"role_id"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`

	// Relación para JOINs (no es una columna de la tabla 'users')
	Role *Role `json:"role,omitempty" db:"-"`
}

// Employee representa la tabla 'employees' en la base de datos.
type Employee struct {
	ID        int64     `json:"id,omitempty" db:"id"`
	Name      string    `json:"name,omitempty" db:"name"`
	RoleID    int64     `json:"role_id,omitempty" db:"role_id"`
	Status    bool      `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`

	// Relación para JOINs
	Role *Role `json:"role,omitempty" db:"-"`
}

// Day representa la tabla 'days' de la semana laboral.
type Day struct {
	ID          int64     `json:"id,omitempty" db:"id"`
	Code        string    `json:"code,omitempty" db:"code"`
	Description string    `json:"description,omitempty" db:"description"`
	StartAt     string    `json:"start_at,omitempty" db:"start_at"` // Formato "HH:MM:SS"
	EndAt       string    `json:"end_at,omitempty" db:"end_at"`     // Formato "HH:MM:SS"
	Status      bool      `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// Service representa la tabla 'services'.
type Service struct {
	ID            int64     `json:"id,omitempty" db:"id"`
	Code          string    `json:"code,omitempty" db:"code"`
	Name          string    `json:"name,omitempty" db:"name"`
	EstimatedTime string    `json:"estimated_time,omitempty" db:"estimated_time"` // Formato "HH:MM:SS"
	Status        bool      `json:"status" db:"status"`
	CreatedAt     time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// Appointment representa la tabla 'appointments'.
type Appointment struct {
	ID         int64         `json:"id,omitempty" db:"id"`
	StartAt    time.Time     `json:"start_at,omitempty" db:"start_at"`
	EndAt      time.Time     `json:"end_at,omitempty" db:"end_at"`
	DayID      int64         `json:"day_id,omitempty" db:"day_id"`
	UserID     int64         `json:"user_id,omitempty" db:"user_id"`
	EmployeeID sql.NullInt64 `json:"employee_id,omitempty" db:"employee_id"`
	CreatedAt  time.Time     `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at,omitempty" db:"updated_at"`

	// Relaciones para JOINs
	Day      *Day      `json:"day,omitempty" db:"-"`
	User     *User     `json:"user,omitempty" db:"-"`
	Employee *Employee `json:"employee,omitempty" db:"-"`
	Services []Service `json:"services,omitempty" db:"-"`
}

// AppointmentService representa la tabla de unión 'appointment_services'.
type AppointmentService struct {
	ID            int64 `json:"id,omitempty" db:"id"`
	ServiceID     int64 `json:"service_id,omitempty" db:"service_id"`
	AppointmentID int64 `json:"appointment_id,omitempty" db:"appointment_id"`
}

// --- DTOs (Data Transfer Objects) ---

// UserCreateRequest define el cuerpo esperado para crear un nuevo usuario.
type UserCreateRequest struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
	Phone    int    `json:"phone" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	RoleID   int64  `json:"role_id" validate:"required"`
}

// UserResponse define la estructura de un usuario para respuestas de API.
type UserResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Phone     int       `json:"phone"`
	Email     string    `json:"email"`
	Role      *Role     `json:"role,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// AppointmentCreateRequest define el cuerpo para crear una nueva cita.
type AppointmentCreateRequest struct {
	StartAt    time.Time `json:"start_at" validate:"required"`
	DayID      int64     `json:"day_id" validate:"required"`
	UserID     int64     `json:"user_id" validate:"required"`
	EmployeeID *int64    `json:"employee_id,omitempty"` // Usar puntero para que sea opcional
	ServiceIDs []int64   `json:"service_ids" validate:"required,min=1"`
}

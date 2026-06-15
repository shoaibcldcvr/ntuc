package domain

import (
	"time"

	"github.com/google/uuid"
)

type Student struct {
	ID        uuid.UUID `db:"id"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	Phone     string    `db:"phone"`
	RollNo    string    `db:"roll_no"`
	Class     string    `db:"class"`
	GPA       float64   `db:"gpa"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// DTOs for requests
type CreateStudentRequest struct {
	FirstName string  `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string  `json:"last_name" validate:"required,min=2,max=50"`
	Email     string  `json:"email" validate:"required,email"`
	Phone     string  `json:"phone" validate:"required,len=10"`
	RollNo    string  `json:"roll_no" validate:"required,alphanum"`
	Class     string  `json:"class" validate:"required,oneof=10A 10B 11A 11B 12A 12B"`
	GPA       float64 `json:"gpa" validate:"required,min=0,max=4"`
}

type UpdateStudentRequest struct {
	FirstName string  `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string  `json:"last_name" validate:"required,min=2,max=50"`
	Email     string  `json:"email" validate:"required,email"`
	Phone     string  `json:"phone" validate:"required,len=10"`
	Class     string  `json:"class" validate:"required,oneof=10A 10B 11A 11B 12A 12B"`
	GPA       float64 `json:"gpa" validate:"required,min=0,max=4"`
}

// DTO for responses
type StudentResponse struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	RollNo    string    `json:"roll_no"`
	Class     string    `json:"class"`
	GPA       float64   `json:"gpa"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Student) ToResponse() *StudentResponse {
	return &StudentResponse{
		ID:        s.ID.String(),
		FirstName: s.FirstName,
		LastName:  s.LastName,
		Email:     s.Email,
		Phone:     s.Phone,
		RollNo:    s.RollNo,
		Class:     s.Class,
		GPA:       s.GPA,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

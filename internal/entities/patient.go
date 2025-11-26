package entities

import (
	"time"

	"gorm.io/gorm"
)

type Patient struct {
	gorm.Model
	FirstNameTh  string    `json:"first_name_th"`
	MiddleNameTh string    `json:"middle_name_th"`
	LastNameTh   string    `json:"last_name_th"`
	FirstNameEn  string    `json:"first_name_en"`
	MiddleNameEn string    `json:"middle_name_en"`
	LastNameEn   string    `json:"last_name_en"`
	DateBirth    time.Time `json:"date_of_birth"`
	PatientHn    string    `json:"patient_hn"`
	NationalId   string    `json:"national_id" validate:"required"`
	PassportId   string    `json:"passport_id"`
	PhoneNumber  string    `json:"phone_number"`
	Email        string    `json:"email"`
	Gender       string    `json:"gender" binding:"required,oneof=male female"`
	Hospital     string    `json:"hospital" validate:"required"`
}

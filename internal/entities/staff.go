package entities

import (
	"gorm.io/gorm"
)

type Staff struct {
	gorm.Model
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Hospital string `json:"hospital" validate:"required"`
}

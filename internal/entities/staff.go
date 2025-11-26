package entities

import (
	"gorm.io/gorm"
)

type Staff struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	Hospital string `json:"hospital"`
}

package adapters

import (
	"agnos/internal/adapters/staff/dto"
	"agnos/internal/entities"
	"agnos/internal/usecases/staff"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type GormStaffRepository struct {
	db *gorm.DB
}

func NewGormStaffRepository(db *gorm.DB) staff.StaffRepository {
	return &GormStaffRepository{db: db}
}

func (r *GormStaffRepository) Save(staff *entities.Staff) (*entities.Staff, error) {
	if err := r.db.Where("username = ?", staff.Username).First(staff).Error; err == nil {
		return nil, fmt.Errorf("username already exist")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if err := r.db.Save(staff).Error; err != nil {
		return nil, err
	}
	return staff, nil
}

func (r *GormStaffRepository) Login(dto *dto.LoginStaffDto) (*entities.Staff, error) {
	var staff entities.Staff
	if err := r.db.Where("username = ?", dto.Username).First(&staff).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with username %s not found", dto.Username)
		}
		return nil, err
	}
	return &staff, nil
}

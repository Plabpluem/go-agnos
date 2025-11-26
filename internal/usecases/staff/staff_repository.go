package staff

import (
	"agnos/internal/adapters/staff/dto"
	"agnos/internal/entities"
)

type StaffRepository interface {
	Save(staff *entities.Staff) (*entities.Staff, error)
	Login(staff *dto.LoginStaffDto) (*entities.Staff, error)
}

package staff

import (
	"agnos/internal/adapters/staff/dto"
	"agnos/internal/entities"
)

type StaffUseCase interface {
	CreateStaff(staff *entities.Staff) (*entities.Staff, error)
	Login(staff *dto.LoginStaffDto) (*entities.Staff, error)
}

type StaffService struct {
	repo StaffRepository
}

func NewStaffService(repo StaffRepository) StaffUseCase {
	return &StaffService{repo: repo}
}

func (s *StaffService) CreateStaff(staff *entities.Staff) (*entities.Staff, error) {
	return s.repo.Save(staff)
}

func (s *StaffService) Login(staff *dto.LoginStaffDto) (*entities.Staff, error) {
	return s.repo.Login(staff)
}

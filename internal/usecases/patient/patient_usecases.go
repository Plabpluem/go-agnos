package patient

import (
	"agnos/internal/adapters/patient/dto"
	"agnos/internal/entities"
)

type PatientUseCase interface {
	CreatePatient(patient *entities.Patient) (*entities.Patient, error)
	SearchPatient(query *dto.SearchPatientDto) ([]*entities.Patient, error)
}

type PatientService struct {
	repo PatientRepository
}

func NewPatientService(repo PatientRepository) PatientUseCase {
	return &PatientService{repo: repo}
}

func (s *PatientService) CreatePatient(patient *entities.Patient) (*entities.Patient, error) {
	return s.repo.Save(patient)
}

func (s *PatientService) SearchPatient(query *dto.SearchPatientDto) ([]*entities.Patient, error) {
	return s.repo.Findone(query)
}

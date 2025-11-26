package patient

import (
	"agnos/internal/adapters/patient/dto"
	"agnos/internal/entities"
)

type PatientRepository interface {
	Save(patient *entities.Patient) (*entities.Patient, error)
	Findone(query *dto.SearchPatientDto) ([]*entities.Patient, error)
	FindoneId(id string) (*entities.Patient, error)
}

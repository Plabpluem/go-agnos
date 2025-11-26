package adapters

import (
	"agnos/internal/adapters/patient/dto"
	"agnos/internal/entities"
	"agnos/internal/usecases/patient"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type GormPatientRepository struct {
	db *gorm.DB
}

func NewGormPatientRepository(db *gorm.DB) patient.PatientRepository {
	return &GormPatientRepository{db: db}
}

func (r *GormPatientRepository) Save(patient *entities.Patient) (*entities.Patient, error) {
	if err := r.db.Where("national_id = ?", patient.NationalId).First(patient).Error; err == nil {
		return nil, fmt.Errorf("national_id already exist")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if err := r.db.Save(patient).Error; err != nil {
		return nil, err
	}
	return patient, nil
}

func (r *GormPatientRepository) Findone(query *dto.SearchPatientDto) ([]*entities.Patient, error) {
	var patient []*entities.Patient

	db := r.db.Model(patient)

	if query.FirstName != "" {
		db = db.Where("LOWER(first_name_th) LIKE ?", "%"+strings.ToLower(query.FirstName)+"%").Or("LOWER(first_name_en) LIKE ?", "%"+strings.ToLower(query.FirstName))
	}

	if query.LastName != "" {
		db = db.Where("LOWER(last_name_th) LIKE ?", "%"+strings.ToLower(query.FirstName)+"%").Or("LOWER(last_name_en) LIKE ?", "%"+strings.ToLower(query.FirstName))
	}

	if query.MiddleName != "" {
		db = db.Where("LOWER(middle_name_th) LIKE ?", "%"+strings.ToLower(query.FirstName)+"%").Or("LOWER(middle_name_en) LIKE ?", "%"+strings.ToLower(query.FirstName))
	}

	if query.PassportId != "" {
		db = db.Where("LOWER(passport_id) LIKE ?", "%"+strings.ToLower(query.PassportId)+"%")
	}

	if query.Email != "" {
		db = db.Where("LOWER(email) LIKE ?", "%"+strings.ToLower(query.Email)+"%")
	}

	if query.PhoneNumber != "" {
		db = db.Where("LOWER(phone_number) LIKE ?", "%"+strings.ToLower(query.PhoneNumber)+"%")
	}

	if query.NationalId != "" {
		db = db.Where("LOWER(national_id) LIKE ?", "%"+strings.ToLower(query.NationalId)+"%")
	}

	if query.Hospital != "" {
		lower_query := strings.ToLower(query.Hospital)
		db = db.Where("LOWER(hospital) = ?", lower_query)
	}
	err := db.Find(&patient).Error

	if err != nil {
		return nil, err
	}
	return patient, nil
}

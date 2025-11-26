package dto

import "time"

type SearchPatientDto struct {
	NationalId  string
	PassportId  string
	FirstName   string
	MiddleName  string
	LastName    string
	DateofBirth time.Time
	PhoneNumber string
	Email       string
	Hospital    string
}

type SearchPatientId struct {
	NationalId string
	PassportId string
}

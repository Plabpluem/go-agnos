package adapters

import (
	"agnos/internal/adapters/patient/dto"
	"agnos/internal/entities"
	"agnos/internal/usecases/patient"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type HttpPatientHandler struct {
	patientUseCase patient.PatientUseCase
}

func NewHttpPatientRepository(usecase patient.PatientUseCase) *HttpPatientHandler {
	return &HttpPatientHandler{patientUseCase: usecase}
}

func (h *HttpPatientHandler) CreatePatient(c *gin.Context) {
	var data entities.Patient

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	patient, err := h.patientUseCase.CreatePatient(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "create success", "statusCode": 201, "data": patient})
}

func (h *HttpPatientHandler) SearchPatient(c *gin.Context) {
	payload, exist := c.Get("payload")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	claims := payload.(jwt.MapClaims)
	params := dto.SearchPatientDto{
		FirstName:   c.Query("first_name"),
		LastName:    c.Query("last_name"),
		MiddleName:  c.Query("middle_name"),
		PassportId:  c.Query("passport_id"),
		Email:       c.Query("email"),
		PhoneNumber: c.Query("phone_number"),
		NationalId:  c.Query("national_id"),
		Hospital:    claims["hospital"].(string),
	}

	patient, err := h.patientUseCase.SearchPatient(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "search success", "statusCode": 200, "data": patient})
}

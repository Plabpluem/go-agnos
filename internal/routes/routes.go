package routes

import (
	adaptersStaff "agnos/internal/adapters/staff"
	usecasesStaff "agnos/internal/usecases/staff"
	"agnos/pkg/middleware"

	adaptersPatient "agnos/internal/adapters/patient"
	usecasesPatient "agnos/internal/usecases/patient"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func StaffRoutes(router *gin.RouterGroup, db *gorm.DB) {
	staffRepo := adaptersStaff.NewGormStaffRepository(db)
	staffService := usecasesStaff.NewStaffService(staffRepo)
	staffHttp := adaptersStaff.NewHttpStaffRepository(staffService)

	router.POST("/staff/create", staffHttp.CreateStaff)
	router.POST("/staff/login", staffHttp.Login)
}

func PatientRoutes(router *gin.RouterGroup, db *gorm.DB) {
	patientRepo := adaptersPatient.NewGormPatientRepository(db)
	patientService := usecasesPatient.NewPatientService(patientRepo)
	patientHttp := adaptersPatient.NewHttpPatientRepository(patientService)

	patientGroup := router.Group("/patient")
	patientGroup.Use(middleware.AuthRequired)

	patientGroup.POST("/create", patientHttp.CreatePatient)
	patientGroup.GET("/search", patientHttp.SearchPatient)
	patientGroup.GET("/search/:id", patientHttp.SearchPatientId)
}

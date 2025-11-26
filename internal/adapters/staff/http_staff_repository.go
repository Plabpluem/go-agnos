package adapters

import (
	"agnos/internal/adapters/staff/dto"
	"agnos/internal/entities"
	"agnos/internal/usecases/staff"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type HttpStaffHandler struct {
	staffUseCase staff.StaffUseCase
}

func NewHttpStaffRepository(usecase staff.StaffUseCase) *HttpStaffHandler {
	return &HttpStaffHandler{staffUseCase: usecase}
}

func (h *HttpStaffHandler) CreateStaff(c *gin.Context) {
	var data entities.Staff

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.New().Struct(&data); err != nil {
		errs := err.(validator.ValidationErrors)

		messages := make([]string, 0)
		for _, e := range errs {
			messages = append(messages, fmt.Sprintf("%s is %s", strings.ToLower(e.Field()), strings.ToLower(e.Tag())))
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": messages,
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	data.Password = string(hashedPassword)

	hospital, err := h.staffUseCase.CreateStaff(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "create success", "statusCode": 201, "data": hospital})
}

func (h *HttpStaffHandler) Login(c *gin.Context) {
	var data dto.LoginStaffDto

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.New().Struct(&data); err != nil {
		errs := err.(validator.ValidationErrors)

		messages := make([]string, 0)
		for _, e := range errs {
			messages = append(messages, fmt.Sprintf("%s is %s", strings.ToLower(e.Field()), strings.ToLower(e.Tag())))
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": messages,
		})
		return
	}

	staff, err := h.staffUseCase.Login(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(staff.Password), []byte(data.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password not matched"})
		return
	}

	claims := jwt.MapClaims{
		"username": staff.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"user_id":  staff.ID,
		"hospital": staff.Hospital,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("supersecret"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"data": staff, "token": t, "statusCode": 200})
}

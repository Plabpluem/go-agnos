package routes_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"agnos/internal/entities"
	"agnos/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	user     = "myuser"
	password = "mypassword"
	dbname   = "mydatabase"
	port     = 5433
)

func clearDatabase(db *gorm.DB) {
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&entities.Staff{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&entities.Patient{})
}

func setupTestRouter() (*gin.Engine, *gorm.DB) {

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	group := r.Group("/")

	dbs := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := gorm.Open(postgres.Open(dbs), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	db.AutoMigrate(&entities.Staff{}, &entities.Patient{})

	routes.StaffRoutes(group, db)
	routes.PatientRoutes(group, db)

	return r, db
}

func createStaffViaApi(t *testing.T, r *gin.Engine, staffData map[string]string) {
	jsonBody, _ := json.Marshal(staffData)

	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer(jsonBody))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)
}

func createLoginStaffViaApi(t *testing.T, r *gin.Engine, staffData entities.Staff) string {
	createDto := map[string]string{
		"username": staffData.Username,
		"password": staffData.Password,
		"hospital": staffData.Hospital,
	}
	createStaffViaApi(t, r, createDto)

	jsonBody, _ := json.Marshal(staffData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/staff/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusAccepted, w.Code)
	assert.Equal(t, float64(200), response["statusCode"])
	return response["token"].(string)
}

func TestStaffRoutes_CreateStaff_Success(t *testing.T) {
	r, db := setupTestRouter()
	defer clearDatabase(db)

	inputData := map[string]string{
		"username": "shinepp",
		"password": "89058905",
		"hospital": "Bangkok Hospital",
	}
	jsonBody, _ := json.Marshal(inputData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, float64(201), response["statusCode"])

	var createdStaff entities.Staff
	err := db.First(&createdStaff, "username = ?", "shinepp").Error
	assert.NoError(t, err)
	assert.Equal(t, "Bangkok Hospital", createdStaff.Hospital)
	assert.Equal(t, "shinepp", createdStaff.Username)
}

func TestStaffRoutes_CreateStaff_FailUsernameAvailable(t *testing.T) {
	r, db := setupTestRouter()
	defer clearDatabase(db)

	firstData := map[string]string{
		"username": "walawala",
		"password": "89058905",
		"hospital": "Bangkok Hospital",
	}
	jsonBody, _ := json.Marshal(firstData)

	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer(jsonBody))
	req1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w1, req1)

	var staff entities.Staff
	err := db.First(&staff, "username = ?", "walawala").Error
	assert.NoError(t, err)
	assert.Equal(t, "walawala", staff.Username)

	inputData := map[string]string{
		"username": "walawala",
		"password": "89058905",
		"hospital": "Bangkok Hospital",
	}
	jsonBody2, _ := json.Marshal(inputData)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/staff/create", bytes.NewBuffer(jsonBody2))
	req2.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusBadRequest, w2.Code)
	var response map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &response)

	assert.Equal(t, "username already exist", response["error"])
}

func TestStaffRoutes_LoginStaff_Success(t *testing.T) {
	r, db := setupTestRouter()
	defer clearDatabase(db)

	dto := map[string]string{
		"username": "walawala",
		"password": "89058905",
		"hospital": "Bangkok Hospital",
	}
	createStaffViaApi(t, r, dto)

	jsonBody, _ := json.Marshal(dto)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/staff/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, float64(200), response["statusCode"])
	assert.NotEmpty(t, response["token"], "ควรได้รับ JWT token")
}

func TestStaffRoutes_LoginStaff_Fail(t *testing.T) {
	r, db := setupTestRouter()
	defer clearDatabase(db)

	createDto := map[string]string{
		"username": "walawala12",
		"password": "89058905",
		"hospital": "Bangkok Hospital",
	}
	loginDto := map[string]string{
		"username": "walawala",
		"password": "89058905",
		"hospital": "Bangkok Hospital",
	}
	createStaffViaApi(t, r, createDto)
	jsonBody, _ := json.Marshal(loginDto)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/staff/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	fmt.Println(response)

	assert.Equal(t, fmt.Sprintf("user with username %s not found", loginDto["username"]), response["error"])

}

func TestStaffRoutes_LoginStaff_FailPasswordWrong(t *testing.T) {
	r, db := setupTestRouter()
	defer clearDatabase(db)
	createDto := map[string]string{
		"username": "walawala",
		"password": "123456",
		"hospital": "Bangkok Hospital",
	}
	loginDto := map[string]string{
		"username": "walawala",
		"password": "89058905",
		"hospital": "Bangkok Hospital",
	}
	createStaffViaApi(t, r, createDto)
	jsonBody, _ := json.Marshal(loginDto)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/staff/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "password not matched", response["error"])
}

func TestPatient_CreatePatient_Success(t *testing.T) {
	r, db := setupTestRouter()
	defer clearDatabase(db)

	createStaffDto := entities.Staff{
		Username: "walawala12",
		Password: "89058905",
		Hospital: "Bangkok Hospital",
	}
	token := createLoginStaffViaApi(t, r, createStaffDto)

	createPatientDto := map[string]string{
		"first_name_th":  "ปลาบปลื้ม",
		"middle_name_th": "-",
		"last_name_th":   "ยอดจันทร์",
		"first_name_en":  "Plabpluem",
		"middle_name_en": "D",
		"last_name_en":   "Yodchan",
		"date_of_birth":  "1995-07-21T00:00:00Z",
		"patient_hn":     "HN00001",
		"national_id":    "890589058905",
		"passport_id":    "",
		"phone_number":   "0812345678",
		"email":          "plabpluem@example.com",
		"gender":         "male",
		"hospital":       "hua-hin hospital",
	}
	jsonBody, _ := json.Marshal(createPatientDto)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/patient/create", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	r.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, float64(201), response["statusCode"])
	assert.Equal(t, "create success", response["message"])
}

func createPatientViaApi(t *testing.T, r *gin.Engine, patient map[string]string, token string) {
	jsonBody, _ := json.Marshal(patient)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/patient/create", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPatient_CreatePatient_FailAvailable(t *testing.T) {
	r, db := setupTestRouter()
	defer clearDatabase(db)

	createStaffDto := entities.Staff{
		Username: "walawala12",
		Password: "89058905",
		Hospital: "Bangkok Hospital",
	}
	token := createLoginStaffViaApi(t, r, createStaffDto)

	patientDto := map[string]string{
		"first_name_th":  "ปลาบปลื้ม",
		"middle_name_th": "-",
		"last_name_th":   "ยอดจันทร์",
		"first_name_en":  "Plabpluem",
		"middle_name_en": "D",
		"last_name_en":   "Yodchan",
		"date_of_birth":  "1995-07-21T00:00:00Z",
		"patient_hn":     "HN00001",
		"national_id":    "890589058905",
		"passport_id":    "",
		"phone_number":   "0812345678",
		"email":          "plabpluem@example.com",
		"gender":         "male",
		"hospital":       "hua-hin hospital",
	}
	createPatientViaApi(t, r, patientDto, token)

	jsonBody, _ := json.Marshal(patientDto)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/patient/create", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	r.ServeHTTP(w, req)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "national_id already exist", response["error"])
}

func TestPatient_SearchPatient_Success(t *testing.T) {
	r, db := setupTestRouter()
	defer clearDatabase(db)

	createStaffDto := entities.Staff{
		Username: "walawala12",
		Password: "89058905",
		Hospital: "Bangkok Hospital",
	}
	token := createLoginStaffViaApi(t, r, createStaffDto)

	firstDto := map[string]string{
		"first_name_th":  "ปลาบปลื้ม",
		"middle_name_th": "-",
		"last_name_th":   "ยอดจันทร์",
		"first_name_en":  "Plabpluem",
		"middle_name_en": "D",
		"last_name_en":   "Yodchan",
		"date_of_birth":  "1995-07-21T00:00:00Z",
		"patient_hn":     "HN00001",
		"national_id":    "890589058905",
		"passport_id":    "",
		"phone_number":   "0812345678",
		"email":          "plabpluem@example.com",
		"gender":         "male",
		"hospital":       "hua-hin hospital",
	}
	secondDto := map[string]string{
		"first_name_th":  "สมศักดิ์",
		"middle_name_th": "-",
		"last_name_th":   "ชวนศรี",
		"first_name_en":  "Somsak",
		"middle_name_en": "D",
		"last_name_en":   "Chunsri",
		"date_of_birth":  "1995-07-21T00:00:00Z",
		"patient_hn":     "HN00002",
		"national_id":    "123456",
		"passport_id":    "",
		"phone_number":   "0812345678",
		"email":          "plabpluem@example.com",
		"gender":         "male",
		"hospital":       "Bangkok Hospital",
	}
	thirdDto := map[string]string{
		"first_name_th":  "กัญชนก",
		"middle_name_th": "-",
		"last_name_th":   "ชวนศรี",
		"first_name_en":  "Kanchanok",
		"middle_name_en": "D",
		"last_name_en":   "Chunsri",
		"date_of_birth":  "1995-07-21T00:00:00Z",
		"patient_hn":     "HN00003",
		"national_id":    "12345678",
		"passport_id":    "",
		"phone_number":   "0812345678",
		"email":          "plabpluem@example.com",
		"gender":         "male",
		"hospital":       "Bangkok Hospital",
	}
	createPatientViaApi(t, r, firstDto, token)
	createPatientViaApi(t, r, secondDto, token)
	createPatientViaApi(t, r, thirdDto, token)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/patient/search", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	dataArray := response["data"].([]interface{})
	firstObject := dataArray[0].(map[string]interface{})
	secondObject := dataArray[1].(map[string]interface{})

	assert.Equal(t, float64(200), response["statusCode"])
	assert.Equal(t, int(2), len(response["data"].([]interface{})))
	assert.Equal(t, "สมศักดิ์", firstObject["first_name_th"])
	assert.Equal(t, "กัญชนก", secondObject["first_name_th"])
}
func TestPatient_SearchPatientWithQuery_Success(t *testing.T) {
	r, db := setupTestRouter()
	defer clearDatabase(db)

	createStaffDto := entities.Staff{
		Username: "walawala12",
		Password: "89058905",
		Hospital: "Bangkok Hospital",
	}
	token := createLoginStaffViaApi(t, r, createStaffDto)

	firstDto := map[string]string{
		"first_name_th":  "ปลาบปลื้ม",
		"middle_name_th": "-",
		"last_name_th":   "ยอดจันทร์",
		"first_name_en":  "Plabpluem",
		"middle_name_en": "D",
		"last_name_en":   "Yodchan",
		"date_of_birth":  "1995-07-21T00:00:00Z",
		"patient_hn":     "HN00001",
		"national_id":    "890589058905",
		"passport_id":    "",
		"phone_number":   "0812345678",
		"email":          "plabpluem@example.com",
		"gender":         "male",
		"hospital":       "hua-hin hospital",
	}
	secondDto := map[string]string{
		"first_name_th":  "สมศักดิ์",
		"middle_name_th": "-",
		"last_name_th":   "ชวนศรี",
		"first_name_en":  "Somsak",
		"middle_name_en": "D",
		"last_name_en":   "Chunsri",
		"date_of_birth":  "1995-07-21T00:00:00Z",
		"patient_hn":     "HN00002",
		"national_id":    "123456",
		"passport_id":    "",
		"phone_number":   "0812345678",
		"email":          "plabpluem@example.com",
		"gender":         "male",
		"hospital":       "Bangkok Hospital",
	}
	thirdDto := map[string]string{
		"first_name_th":  "กัญชนก",
		"middle_name_th": "-",
		"last_name_th":   "ชวนศรี",
		"first_name_en":  "Kanchanok",
		"middle_name_en": "D",
		"last_name_en":   "Chunsri",
		"date_of_birth":  "1995-07-21T00:00:00Z",
		"patient_hn":     "HN00003",
		"national_id":    "12345678",
		"passport_id":    "",
		"phone_number":   "0812345678",
		"email":          "plabpluem@example.com",
		"gender":         "male",
		"hospital":       "Bangkok Hospital",
	}
	createPatientViaApi(t, r, firstDto, token)
	createPatientViaApi(t, r, secondDto, token)
	createPatientViaApi(t, r, thirdDto, token)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/patient/search?first_name=กัญชนก", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	r.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	dataArray := response["data"].([]interface{})
	firstObject := dataArray[0].(map[string]interface{})

	assert.Equal(t, float64(200), response["statusCode"])
	assert.Equal(t, int(1), len(response["data"].([]interface{})))
	assert.Equal(t, "กัญชนก", firstObject["first_name_th"])
}

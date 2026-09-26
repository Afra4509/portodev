package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"rest-api-ipk-mahasiswa-its/internal/handler"
	"rest-api-ipk-mahasiswa-its/internal/model"
	"rest-api-ipk-mahasiswa-its/internal/repository"
	"rest-api-ipk-mahasiswa-its/internal/routes"
	"rest-api-ipk-mahasiswa-its/internal/service"
	"rest-api-ipk-mahasiswa-its/internal/testutil"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	db := testutil.SetupTestDB()

	studentRepo := repository.NewStudentRepository(db)
	courseRepo := repository.NewCourseRepository(db)
	gradeRepo := repository.NewGradeRepository(db)

	studentService := service.NewStudentService(studentRepo)
	courseService := service.NewCourseService(courseRepo)
	gradeService := service.NewGradeService(gradeRepo, studentRepo, courseRepo)
	academicService := service.NewAcademicService(studentRepo, gradeRepo)

	studentHandler := handler.NewStudentHandler(studentService)
	courseHandler := handler.NewCourseHandler(courseService)
	gradeHandler := handler.NewGradeHandler(gradeService)
	academicHandler := handler.NewAcademicHandler(academicService)

	return routes.SetupRouter(routes.RouterDependencies{
		StudentHandler:  studentHandler,
		CourseHandler:   courseHandler,
		GradeHandler:    gradeHandler,
		AcademicHandler: academicHandler,
	})
}

func TestAPI_EndToEndFlow(t *testing.T) {
	router := setupTestRouter()

	// 1. POST /api/students
	studentPayload := map[string]any{
		"nrp":           "5025211001",
		"nama":          "Budi Santoso",
		"program_studi": "Teknik Informatika",
		"fakultas":      "FTEIC",
		"angkatan":      2021,
	}
	body, _ := json.Marshal(studentPayload)
	req, _ := http.NewRequest(http.MethodPost, "/api/students", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var studentResp model.WebResponse
	err := json.Unmarshal(w.Body.Bytes(), &studentResp)
	require.NoError(t, err)
	assert.True(t, studentResp.Success)

	// Extract student ID
	studentData := studentResp.Data.(map[string]any)
	studentID := uint(studentData["id"].(float64))

	// 2. POST /api/courses
	course1Payload := map[string]any{
		"kode_mata_kuliah": "IF184101",
		"nama_mata_kuliah": "Dasar Pemrograman",
		"sks":              3,
	}
	body, _ = json.Marshal(course1Payload)
	req, _ = http.NewRequest(http.MethodPost, "/api/courses", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var courseResp model.WebResponse
	_ = json.Unmarshal(w.Body.Bytes(), &courseResp)
	course1Data := courseResp.Data.(map[string]any)
	course1ID := uint(course1Data["id"].(float64))

	course2Payload := map[string]any{
		"kode_mata_kuliah": "KM184101",
		"nama_mata_kuliah": "Kalkulus 1",
		"sks":              3,
	}
	body, _ = json.Marshal(course2Payload)
	req, _ = http.NewRequest(http.MethodPost, "/api/courses", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	_ = json.Unmarshal(w.Body.Bytes(), &courseResp)
	course2Data := courseResp.Data.(map[string]any)
	course2ID := uint(course2Data["id"].(float64))

	// 3. POST /api/grades
	// Assign Course 1: Nilai A (4.00) in Semester 1
	grade1Payload := map[string]any{
		"student_id":  studentID,
		"course_id":   course1ID,
		"semester":    1,
		"nilai_huruf": "A",
	}
	body, _ = json.Marshal(grade1Payload)
	req, _ = http.NewRequest(http.MethodPost, "/api/grades", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Assign Course 2: Nilai AB (3.50) in Semester 1
	grade2Payload := map[string]any{
		"student_id":  studentID,
		"course_id":   course2ID,
		"semester":    1,
		"nilai_huruf": "AB",
	}
	body, _ = json.Marshal(grade2Payload)
	req, _ = http.NewRequest(http.MethodPost, "/api/grades", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// 4. Test Duplicate Grade validation: Course 1 in Semester 1 again -> 409 Conflict
	req, _ = http.NewRequest(http.MethodPost, "/api/grades", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusConflict, w.Code)

	// 5. GET /api/students/:id/ip/:semester
	// Total SKS = 6, Total Mutu = (4.00*3) + (3.50*3) = 12 + 10.5 = 22.5
	// IP = 22.5 / 6 = 3.75
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/api/students/%d/ip/1", studentID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var ipResp model.WebResponse
	err = json.Unmarshal(w.Body.Bytes(), &ipResp)
	require.NoError(t, err)
	ipData := ipResp.Data.(map[string]any)
	assert.Equal(t, float64(6), ipData["total_sks"])
	assert.Equal(t, 22.5, ipData["total_mutu"])
	assert.Equal(t, 3.75, ipData["ip"])

	// 6. GET /api/students/:id/ipk
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/api/students/%d/ipk", studentID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var ipkResp model.WebResponse
	err = json.Unmarshal(w.Body.Bytes(), &ipkResp)
	require.NoError(t, err)
	ipkData := ipkResp.Data.(map[string]any)
	assert.Equal(t, float64(6), ipkData["total_sks"])
	assert.Equal(t, 3.75, ipkData["ipk"])
	assert.Equal(t, "5025211001", ipkData["nrp"])

	// 7. GET /api/students/:id/grades
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/api/students/%d/grades", studentID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 8. PUT /api/students/:id
	updatePayload := map[string]any{
		"nrp":           "5025211001",
		"nama":          "Budi Santoso, S.Kom",
		"program_studi": "Teknik Informatika",
		"fakultas":      "FTEIC",
		"angkatan":      2021,
	}
	body, _ = json.Marshal(updatePayload)
	req, _ = http.NewRequest(http.MethodPut, fmt.Sprintf("/api/students/%d", studentID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 9. DELETE /api/students/:id
	req, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/students/%d", studentID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"rest-api-ipk-mahasiswa-its/internal/model"
	"rest-api-ipk-mahasiswa-its/internal/service"

	"github.com/gin-gonic/gin"
)

type StudentHandler struct {
	studentService service.StudentService
}

func NewStudentHandler(studentService service.StudentService) *StudentHandler {
	return &StudentHandler{studentService: studentService}
}

// CreateStudent handles POST /api/students
func (h *StudentHandler) CreateStudent(c *gin.Context) {
	var req model.CreateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendBadRequest(c, "Invalid request body or missing required fields", err.Error())
		return
	}

	student, err := h.studentService.CreateStudent(&req)
	if err != nil {
		SendError(c, err)
		return
	}

	SendSuccess(c, http.StatusCreated, "Student created successfully", student)
}

// GetAllStudents handles GET /api/students
func (h *StudentHandler) GetAllStudents(c *gin.Context) {
	search := c.Query("search")
	students, err := h.studentService.GetAllStudents(search)
	if err != nil {
		SendError(c, err)
		return
	}

	SendSuccess(c, http.StatusOK, "Students retrieved successfully", students)
}

// GetStudentByID handles GET /api/students/:id
func (h *StudentHandler) GetStudentByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		SendBadRequest(c, "Invalid student ID parameter", fmt.Sprintf("ID '%s' is not a valid unsigned integer", idParam))
		return
	}

	student, err := h.studentService.GetStudentByID(uint(id))
	if err != nil {
		SendError(c, err)
		return
	}

	SendSuccess(c, http.StatusOK, "Student retrieved successfully", student)
}

// UpdateStudent handles PUT /api/students/:id
func (h *StudentHandler) UpdateStudent(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		SendBadRequest(c, "Invalid student ID parameter", fmt.Sprintf("ID '%s' is not a valid unsigned integer", idParam))
		return
	}

	var req model.UpdateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendBadRequest(c, "Invalid request body or missing required fields", err.Error())
		return
	}

	student, err := h.studentService.UpdateStudent(uint(id), &req)
	if err != nil {
		SendError(c, err)
		return
	}

	SendSuccess(c, http.StatusOK, "Student updated successfully", student)
}

// DeleteStudent handles DELETE /api/students/:id
func (h *StudentHandler) DeleteStudent(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		SendBadRequest(c, "Invalid student ID parameter", fmt.Sprintf("ID '%s' is not a valid unsigned integer", idParam))
		return
	}

	if err := h.studentService.DeleteStudent(uint(id)); err != nil {
		SendError(c, err)
		return
	}

	SendSuccess(c, http.StatusOK, "Student deleted successfully", nil)
}

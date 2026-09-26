package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"rest-api-ipk-mahasiswa-its/internal/model"
	"rest-api-ipk-mahasiswa-its/internal/service"

	"github.com/gin-gonic/gin"
)

type GradeHandler struct {
	gradeService service.GradeService
}

func NewGradeHandler(gradeService service.GradeService) *GradeHandler {
	return &GradeHandler{gradeService: gradeService}
}

// CreateGrade handles POST /api/grades
func (h *GradeHandler) CreateGrade(c *gin.Context) {
	var req model.CreateGradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendBadRequest(c, "Invalid request body or missing required fields", err.Error())
		return
	}

	grade, err := h.gradeService.CreateGrade(&req)
	if err != nil {
		SendError(c, err)
		return
	}

	SendSuccess(c, http.StatusCreated, "Grade created successfully", grade)
}

// GetGradesByStudentID handles GET /api/students/:id/grades
func (h *GradeHandler) GetGradesByStudentID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		SendBadRequest(c, "Invalid student ID parameter", fmt.Sprintf("ID '%s' is not a valid unsigned integer", idParam))
		return
	}

	grades, err := h.gradeService.GetGradesByStudentID(uint(id))
	if err != nil {
		SendError(c, err)
		return
	}

	SendSuccess(c, http.StatusOK, "Student grades retrieved successfully", grades)
}

// UpdateGrade handles PUT /api/grades/:id
func (h *GradeHandler) UpdateGrade(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		SendBadRequest(c, "Invalid grade ID parameter", fmt.Sprintf("ID '%s' is not a valid unsigned integer", idParam))
		return
	}

	var req model.UpdateGradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendBadRequest(c, "Invalid request body or missing required fields", err.Error())
		return
	}

	grade, err := h.gradeService.UpdateGrade(uint(id), &req)
	if err != nil {
		SendError(c, err)
		return
	}

	SendSuccess(c, http.StatusOK, "Grade updated successfully", grade)
}

// DeleteGrade handles DELETE /api/grades/:id
func (h *GradeHandler) DeleteGrade(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		SendBadRequest(c, "Invalid grade ID parameter", fmt.Sprintf("ID '%s' is not a valid unsigned integer", idParam))
		return
	}

	if err := h.gradeService.DeleteGrade(uint(id)); err != nil {
		SendError(c, err)
		return
	}

	SendSuccess(c, http.StatusOK, "Grade deleted successfully", nil)
}

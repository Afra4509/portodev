package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"rest-api-ipk-mahasiswa-its/internal/model"
	"rest-api-ipk-mahasiswa-its/internal/service"

	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	courseService service.CourseService
}

func NewCourseHandler(courseService service.CourseService) *CourseHandler {
	return &CourseHandler{courseService: courseService}
}

// CreateCourse handles POST /api/courses
func (h *CourseHandler) CreateCourse(c *gin.Context) {
	var req model.CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendBadRequest(c, "Invalid request body or missing required fields", err.Error())
		return
	}

	course, err := h.courseService.CreateCourse(&req)
	if err != nil {
		SendError(c, err)
		return
	}

	SendSuccess(c, http.StatusCreated, "Course created successfully", course)
}

// GetAllCourses handles GET /api/courses
func (h *CourseHandler) GetAllCourses(c *gin.Context) {
	search := c.Query("search")
	courses, err := h.courseService.GetAllCourses(search)
	if err != nil {
		SendError(c, err)
		return
	}

	SendSuccess(c, http.StatusOK, "Courses retrieved successfully", courses)
}

// GetCourseByID handles GET /api/courses/:id
func (h *CourseHandler) GetCourseByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		SendBadRequest(c, "Invalid course ID parameter", fmt.Sprintf("ID '%s' is not a valid unsigned integer", idParam))
		return
	}

	course, err := h.courseService.GetCourseByID(uint(id))
	if err != nil {
		SendError(c, err)
		return
	}

	SendSuccess(c, http.StatusOK, "Course retrieved successfully", course)
}

// UpdateCourse handles PUT /api/courses/:id
func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		SendBadRequest(c, "Invalid course ID parameter", fmt.Sprintf("ID '%s' is not a valid unsigned integer", idParam))
		return
	}

	var req model.UpdateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendBadRequest(c, "Invalid request body or missing required fields", err.Error())
		return
	}

	course, err := h.courseService.UpdateCourse(uint(id), &req)
	if err != nil {
		SendError(c, err)
		return
	}

	SendSuccess(c, http.StatusOK, "Course updated successfully", course)
}

// DeleteCourse handles DELETE /api/courses/:id
func (h *CourseHandler) DeleteCourse(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		SendBadRequest(c, "Invalid course ID parameter", fmt.Sprintf("ID '%s' is not a valid unsigned integer", idParam))
		return
	}

	if err := h.courseService.DeleteCourse(uint(id)); err != nil {
		SendError(c, err)
		return
	}

	SendSuccess(c, http.StatusOK, "Course deleted successfully", nil)
}

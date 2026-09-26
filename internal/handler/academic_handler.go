package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"rest-api-ipk-mahasiswa-its/internal/service"

	"github.com/gin-gonic/gin"
)

type AcademicHandler struct {
	academicService service.AcademicService
}

func NewAcademicHandler(academicService service.AcademicService) *AcademicHandler {
	return &AcademicHandler{academicService: academicService}
}

// GetIPSemester handles GET /api/students/:id/ip/:semester
func (h *AcademicHandler) GetIPSemester(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		SendBadRequest(c, "Invalid student ID parameter", fmt.Sprintf("ID '%s' is not a valid unsigned integer", idParam))
		return
	}

	semParam := c.Param("semester")
	semester, err := strconv.Atoi(semParam)
	if err != nil {
		SendBadRequest(c, "Invalid semester parameter", fmt.Sprintf("Semester '%s' is not a valid integer", semParam))
		return
	}

	res, err := h.academicService.CalculateIPSemester(uint(id), semester)
	if err != nil {
		SendError(c, err)
		return
	}

	SendSuccess(c, http.StatusOK, fmt.Sprintf("IP Semester %d calculated successfully", semester), res)
}

// GetIPK handles GET /api/students/:id/ipk
func (h *AcademicHandler) GetIPK(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		SendBadRequest(c, "Invalid student ID parameter", fmt.Sprintf("ID '%s' is not a valid unsigned integer", idParam))
		return
	}

	res, err := h.academicService.CalculateIPK(uint(id))
	if err != nil {
		SendError(c, err)
		return
	}

	SendSuccess(c, http.StatusOK, "IPK calculated successfully", res)
}

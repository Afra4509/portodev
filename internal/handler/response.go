package handler

import (
	"errors"
	"net/http"

	"rest-api-ipk-mahasiswa-its/internal/model"
	"rest-api-ipk-mahasiswa-its/internal/service"

	"github.com/gin-gonic/gin"
)

// SendSuccess sends a standardized successful JSON response
func SendSuccess(c *gin.Context, statusCode int, message string, data any) {
	c.JSON(statusCode, model.WebResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SendError formats and sends a standardized error JSON response based on error type
func SendError(c *gin.Context, err error) {
	var appErr *service.AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.Code, model.WebResponse{
			Success: false,
			Message: appErr.Message,
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, model.WebResponse{
		Success: false,
		Message: err.Error(),
		Data:    nil,
	})
}

// SendBadRequest formats validation / binding errors
func SendBadRequest(c *gin.Context, message string, errDetails any) {
	c.JSON(http.StatusBadRequest, model.WebResponse{
		Success: false,
		Message: message,
		Data:    nil,
		Errors:  errDetails,
	})
}

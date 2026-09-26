package routes

import (
	"net/http"

	"rest-api-ipk-mahasiswa-its/internal/handler"
	"rest-api-ipk-mahasiswa-its/internal/model"

	"github.com/gin-gonic/gin"
)

type RouterDependencies struct {
	StudentHandler  *handler.StudentHandler
	CourseHandler   *handler.CourseHandler
	GradeHandler    *handler.GradeHandler
	AcademicHandler *handler.AcademicHandler
}

// corsMiddleware handles Cross-Origin Resource Sharing
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// SetupRouter initializes and configures all API routes
func SetupRouter(deps RouterDependencies) *gin.Engine {
	r := gin.New()

	// Middlewares
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())

	// Custom 404 Handler for unmatched routes
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, model.WebResponse{
			Success: false,
			Message: "Endpoint not found",
			Data:    nil,
		})
	})

	// Health Check
	r.GET("/", func(c *gin.Context) {
		handler.SendSuccess(c, http.StatusOK, "REST API Manajemen IPK Mahasiswa ITS is running", gin.H{
			"version": "1.0.0",
			"docs":    "/api/health",
		})
	})

	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			handler.SendSuccess(c, http.StatusOK, "System is healthy", gin.H{"status": "UP"})
		})

		// 1. Student Management
		api.POST("/students", deps.StudentHandler.CreateStudent)
		api.GET("/students", deps.StudentHandler.GetAllStudents)
		api.GET("/students/:id", deps.StudentHandler.GetStudentByID)
		api.PUT("/students/:id", deps.StudentHandler.UpdateStudent)
		api.DELETE("/students/:id", deps.StudentHandler.DeleteStudent)

		// 2. Course Management
		api.POST("/courses", deps.CourseHandler.CreateCourse)
		api.GET("/courses", deps.CourseHandler.GetAllCourses)
		api.GET("/courses/:id", deps.CourseHandler.GetCourseByID)
		api.PUT("/courses/:id", deps.CourseHandler.UpdateCourse)
		api.DELETE("/courses/:id", deps.CourseHandler.DeleteCourse)

		// 3. Grade Management
		api.POST("/grades", deps.GradeHandler.CreateGrade)
		api.GET("/students/:id/grades", deps.GradeHandler.GetGradesByStudentID)
		api.PUT("/grades/:id", deps.GradeHandler.UpdateGrade)
		api.DELETE("/grades/:id", deps.GradeHandler.DeleteGrade)

		// 4. IP Semester Calculation
		api.GET("/students/:id/ip/:semester", deps.AcademicHandler.GetIPSemester)

		// 5. IPK Cumulative Calculation
		api.GET("/students/:id/ipk", deps.AcademicHandler.GetIPK)
	}

	return r
}

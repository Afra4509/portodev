package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"rest-api-ipk-mahasiswa-its/internal/config"
	"rest-api-ipk-mahasiswa-its/internal/handler"
	"rest-api-ipk-mahasiswa-its/internal/repository"
	"rest-api-ipk-mahasiswa-its/internal/routes"
	"rest-api-ipk-mahasiswa-its/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load Environment Configuration
	cfg := config.LoadConfig()

	if cfg.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 2. Initialize Database Connection & Auto Migration
	db, err := config.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("[FATAL] Could not initialize database: %v\n", err)
	}

	// 3. Initialize Repositories (Data Access Layer)
	studentRepo := repository.NewStudentRepository(db)
	courseRepo := repository.NewCourseRepository(db)
	gradeRepo := repository.NewGradeRepository(db)

	// 4. Initialize Services (Business Logic Layer)
	studentService := service.NewStudentService(studentRepo)
	courseService := service.NewCourseService(courseRepo)
	gradeService := service.NewGradeService(gradeRepo, studentRepo, courseRepo)
	academicService := service.NewAcademicService(studentRepo, gradeRepo)

	// 5. Initialize Handlers (HTTP Presentation Layer)
	studentHandler := handler.NewStudentHandler(studentService)
	courseHandler := handler.NewCourseHandler(courseService)
	gradeHandler := handler.NewGradeHandler(gradeService)
	academicHandler := handler.NewAcademicHandler(academicService)

	// 6. Setup Router & Endpoints
	router := routes.SetupRouter(routes.RouterDependencies{
		StudentHandler:  studentHandler,
		CourseHandler:   courseHandler,
		GradeHandler:    gradeHandler,
		AcademicHandler: academicHandler,
	})

	// 7. Start HTTP Server with Graceful Shutdown
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("====================================================\n")
		log.Printf(" REST API Manajemen IPK Mahasiswa ITS is running on port %s\n", cfg.ServerPort)
		log.Printf(" Health Check URL: http://localhost:%s/api/health\n", cfg.ServerPort)
		log.Printf("====================================================\n")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[FATAL] Server error: %v\n", err)
		}
	}()

	// Wait for termination signal
	<-quit
	log.Println("[INFO] Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("[FATAL] Server forced to shutdown: %v\n", err)
	}

	log.Println("[INFO] Server exited successfully.")
}

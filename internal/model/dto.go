package model

import (
	"fmt"
	"strings"
)

// Standard ITS Grading Scale
var ValidGradePoints = map[string]float64{
	"A":  4.00,
	"AB": 3.50,
	"B":  3.00,
	"BC": 2.50,
	"C":  2.00,
	"D":  1.00,
	"E":  0.00,
}

// ConvertGrade converts letter grade to numerical grade point
func ConvertGrade(huruf string) (float64, error) {
	upper := strings.ToUpper(strings.TrimSpace(huruf))
	val, ok := ValidGradePoints[upper]
	if !ok {
		return 0, fmt.Errorf("invalid letter grade: '%s'. Allowed: A, AB, B, BC, C, D, E", huruf)
	}
	return val, nil
}

// Student DTOs
type CreateStudentRequest struct {
	NRP          string `json:"nrp" binding:"required"`
	Nama         string `json:"nama" binding:"required"`
	ProgramStudi string `json:"program_studi" binding:"required"`
	Fakultas     string `json:"fakultas" binding:"required"`
	Angkatan     int    `json:"angkatan" binding:"required,min=2000,max=2100"`
}

type UpdateStudentRequest struct {
	NRP          string `json:"nrp" binding:"required"`
	Nama         string `json:"nama" binding:"required"`
	ProgramStudi string `json:"program_studi" binding:"required"`
	Fakultas     string `json:"fakultas" binding:"required"`
	Angkatan     int    `json:"angkatan" binding:"required,min=2000,max=2100"`
}

// Course DTOs
type CreateCourseRequest struct {
	KodeMataKuliah string `json:"kode_mata_kuliah" binding:"required"`
	NamaMataKuliah string `json:"nama_mata_kuliah" binding:"required"`
	SKS            int    `json:"sks" binding:"required,min=1,max=6"`
}

type UpdateCourseRequest struct {
	KodeMataKuliah string `json:"kode_mata_kuliah" binding:"required"`
	NamaMataKuliah string `json:"nama_mata_kuliah" binding:"required"`
	SKS            int    `json:"sks" binding:"required,min=1,max=6"`
}

// Grade DTOs
type CreateGradeRequest struct {
	StudentID  uint   `json:"student_id" binding:"required"`
	CourseID   uint   `json:"course_id" binding:"required"`
	Semester   int    `json:"semester" binding:"required,min=1,max=14"`
	NilaiHuruf string `json:"nilai_huruf" binding:"required"`
}

type UpdateGradeRequest struct {
	Semester   *int    `json:"semester,omitempty" binding:"omitempty,min=1,max=14"`
	NilaiHuruf *string `json:"nilai_huruf,omitempty"`
}

// Academic Calculation Responses
type GradeDetailResponse struct {
	ID             uint    `json:"id"`
	CourseID       uint    `json:"course_id"`
	KodeMataKuliah string  `json:"kode_mata_kuliah"`
	NamaMataKuliah string  `json:"nama_mata_kuliah"`
	SKS            int     `json:"sks"`
	Semester       int     `json:"semester"`
	NilaiHuruf     string  `json:"nilai_huruf"`
	NilaiAngka     float64 `json:"nilai_angka"`
	NilaiMutu      float64 `json:"nilai_mutu"` // NilaiAngka * SKS
}

type IPSemesterResponse struct {
	StudentID uint                  `json:"student_id"`
	NRP       string                `json:"nrp"`
	Nama      string                `json:"nama"`
	Semester  int                   `json:"semester"`
	TotalSKS  int                   `json:"total_sks"`
	TotalMutu float64               `json:"total_mutu"`
	IP        float64               `json:"ip"`
	Courses   []GradeDetailResponse `json:"courses,omitempty"`
}

type IPKResponse struct {
	StudentID       uint                  `json:"student_id"`
	Student         string                `json:"student"` // formatted as "NRP - Nama" for convenience
	NRP             string                `json:"nrp"`
	Nama            string                `json:"nama"`
	ProgramStudi    string                `json:"program_studi"`
	Fakultas        string                `json:"fakultas"`
	Angkatan        int                   `json:"angkatan"`
	TotalSKS        int                   `json:"total_sks"`
	TotalMutu       float64               `json:"total_mutu"`
	IPK             float64               `json:"ipk"`
	TotalMataKuliah int                   `json:"total_mata_kuliah"`
	Courses         []GradeDetailResponse `json:"courses,omitempty"`
}

// Consistent JSON API Responses
type WebResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Errors  any    `json:"errors,omitempty"`
}

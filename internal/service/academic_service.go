package service

import (
	"fmt"
	"math"

	"rest-api-ipk-mahasiswa-its/internal/model"
	"rest-api-ipk-mahasiswa-its/internal/repository"
)

type AcademicService interface {
	CalculateIPSemester(studentID uint, semester int) (*model.IPSemesterResponse, error)
	CalculateIPK(studentID uint) (*model.IPKResponse, error)
}

type academicService struct {
	studentRepo repository.StudentRepository
	gradeRepo   repository.GradeRepository
}

func NewAcademicService(
	studentRepo repository.StudentRepository,
	gradeRepo repository.GradeRepository,
) AcademicService {
	return &academicService{
		studentRepo: studentRepo,
		gradeRepo:   gradeRepo,
	}
}

// round2 rounds float64 to 2 decimal places
func round2(val float64) float64 {
	return math.Round(val*100.0) / 100.0
}

func (s *academicService) CalculateIPSemester(studentID uint, semester int) (*model.IPSemesterResponse, error) {
	if semester < 1 || semester > 14 {
		return nil, NewBadRequestError("Semester must be between 1 and 14", ErrInvalidSemester)
	}

	// 1. Verify Student Exists
	student, err := s.studentRepo.FindByID(studentID)
	if err != nil {
		return nil, NewInternalError("failed to query student", err)
	}
	if student == nil {
		return nil, NewNotFoundError(fmt.Sprintf("Student with ID %d not found", studentID), ErrStudentNotFound)
	}

	// 2. Fetch Grades for that Semester
	grades, err := s.gradeRepo.FindByStudentIDAndSemester(studentID, semester)
	if err != nil {
		return nil, NewInternalError("failed to query grades", err)
	}

	var totalSKS int
	var totalMutu float64
	courseDetails := make([]model.GradeDetailResponse, 0, len(grades))

	for _, g := range grades {
		sks := 0
		kodeMK := ""
		namaMK := ""
		if g.Course != nil {
			sks = g.Course.SKS
			kodeMK = g.Course.KodeMataKuliah
			namaMK = g.Course.NamaMataKuliah
		}

		mutu := round2(g.NilaiAngka * float64(sks))
		totalSKS += sks
		totalMutu += mutu

		courseDetails = append(courseDetails, model.GradeDetailResponse{
			ID:             g.ID,
			CourseID:       g.CourseID,
			KodeMataKuliah: kodeMK,
			NamaMataKuliah: namaMK,
			SKS:            sks,
			Semester:       g.Semester,
			NilaiHuruf:     g.NilaiHuruf,
			NilaiAngka:     g.NilaiAngka,
			NilaiMutu:      mutu,
		})
	}

	var ip float64
	if totalSKS > 0 {
		ip = round2(totalMutu / float64(totalSKS))
	}

	return &model.IPSemesterResponse{
		StudentID: student.ID,
		NRP:       student.NRP,
		Nama:      student.Nama,
		Semester:  semester,
		TotalSKS:  totalSKS,
		TotalMutu: round2(totalMutu),
		IP:        ip,
		Courses:   courseDetails,
	}, nil
}

func (s *academicService) CalculateIPK(studentID uint) (*model.IPKResponse, error) {
	// 1. Verify Student Exists
	student, err := s.studentRepo.FindByID(studentID)
	if err != nil {
		return nil, NewInternalError("failed to query student", err)
	}
	if student == nil {
		return nil, NewNotFoundError(fmt.Sprintf("Student with ID %d not found", studentID), ErrStudentNotFound)
	}

	// 2. Fetch all grades of this student
	grades, err := s.gradeRepo.FindByStudentID(studentID)
	if err != nil {
		return nil, NewInternalError("failed to query grades", err)
	}

	// In academic systems (ITS regulation):
	// If a student repeats a course in different semesters, the best/highest grade is counted for IPK.
	bestGradeByCourse := make(map[uint]model.Grade)
	for _, g := range grades {
		existing, found := bestGradeByCourse[g.CourseID]
		if !found || g.NilaiAngka > existing.NilaiAngka {
			bestGradeByCourse[g.CourseID] = g
		}
	}

	var totalSKS int
	var totalMutu float64
	courseDetails := make([]model.GradeDetailResponse, 0, len(bestGradeByCourse))

	for _, g := range bestGradeByCourse {
		sks := 0
		kodeMK := ""
		namaMK := ""
		if g.Course != nil {
			sks = g.Course.SKS
			kodeMK = g.Course.KodeMataKuliah
			namaMK = g.Course.NamaMataKuliah
		}

		mutu := round2(g.NilaiAngka * float64(sks))
		totalSKS += sks
		totalMutu += mutu

		courseDetails = append(courseDetails, model.GradeDetailResponse{
			ID:             g.ID,
			CourseID:       g.CourseID,
			KodeMataKuliah: kodeMK,
			NamaMataKuliah: namaMK,
			SKS:            sks,
			Semester:       g.Semester,
			NilaiHuruf:     g.NilaiHuruf,
			NilaiAngka:     g.NilaiAngka,
			NilaiMutu:      mutu,
		})
	}

	var ipk float64
	if totalSKS > 0 {
		ipk = round2(totalMutu / float64(totalSKS))
	}

	return &model.IPKResponse{
		StudentID:       student.ID,
		Student:         fmt.Sprintf("%s - %s", student.NRP, student.Nama),
		NRP:             student.NRP,
		Nama:            student.Nama,
		ProgramStudi:    student.ProgramStudi,
		Fakultas:        student.Fakultas,
		Angkatan:        student.Angkatan,
		TotalSKS:        totalSKS,
		TotalMutu:       round2(totalMutu),
		IPK:             ipk,
		TotalMataKuliah: len(bestGradeByCourse),
		Courses:         courseDetails,
	}, nil
}

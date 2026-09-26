package service

import (
	"fmt"
	"strings"

	"rest-api-ipk-mahasiswa-its/internal/model"
	"rest-api-ipk-mahasiswa-its/internal/repository"
)

type GradeService interface {
	CreateGrade(req *model.CreateGradeRequest) (*model.Grade, error)
	GetGradesByStudentID(studentID uint) ([]model.Grade, error)
	GetGradeByID(id uint) (*model.Grade, error)
	UpdateGrade(id uint, req *model.UpdateGradeRequest) (*model.Grade, error)
	DeleteGrade(id uint) error
}

type gradeService struct {
	gradeRepo   repository.GradeRepository
	studentRepo repository.StudentRepository
	courseRepo  repository.CourseRepository
}

func NewGradeService(
	gradeRepo repository.GradeRepository,
	studentRepo repository.StudentRepository,
	courseRepo repository.CourseRepository,
) GradeService {
	return &gradeService{
		gradeRepo:   gradeRepo,
		studentRepo: studentRepo,
		courseRepo:  courseRepo,
	}
}

func (s *gradeService) CreateGrade(req *model.CreateGradeRequest) (*model.Grade, error) {
	// 1. Validate Semester
	if req.Semester < 1 || req.Semester > 14 {
		return nil, NewBadRequestError("Semester must be between 1 and 14", ErrInvalidSemester)
	}

	// 2. Validate Letter Grade and Convert to Numerical Point
	cleanHuruf := strings.ToUpper(strings.TrimSpace(req.NilaiHuruf))
	nilaiAngka, err := model.ConvertGrade(cleanHuruf)
	if err != nil {
		return nil, NewBadRequestError(err.Error(), ErrInvalidNilaiHuruf)
	}

	// 3. Verify Student Exists
	student, err := s.studentRepo.FindByID(req.StudentID)
	if err != nil {
		return nil, NewInternalError("failed to verify student", err)
	}
	if student == nil {
		return nil, NewNotFoundError(fmt.Sprintf("Student with ID %d not found", req.StudentID), ErrStudentNotFound)
	}

	// 4. Verify Course Exists
	course, err := s.courseRepo.FindByID(req.CourseID)
	if err != nil {
		return nil, NewInternalError("failed to verify course", err)
	}
	if course == nil {
		return nil, NewNotFoundError(fmt.Sprintf("Course with ID %d not found", req.CourseID), ErrCourseNotFound)
	}

	// 5. Check for Duplicate Grade in the same semester
	existing, err := s.gradeRepo.FindByStudentAndCourseAndSemester(req.StudentID, req.CourseID, req.Semester)
	if err != nil {
		return nil, NewInternalError("failed to check existing grade", err)
	}
	if existing != nil {
		return nil, NewConflictError(
			fmt.Sprintf("Grade for student ID %d, course ID %d, in semester %d already exists", req.StudentID, req.CourseID, req.Semester),
			ErrDuplicateGrade,
		)
	}

	grade := &model.Grade{
		StudentID:  req.StudentID,
		CourseID:   req.CourseID,
		Semester:   req.Semester,
		NilaiHuruf: cleanHuruf,
		NilaiAngka: nilaiAngka,
	}

	if err := s.gradeRepo.Create(grade); err != nil {
		return nil, NewInternalError("failed to save grade", err)
	}

	grade.Student = student
	grade.Course = course
	return grade, nil
}

func (s *gradeService) GetGradesByStudentID(studentID uint) ([]model.Grade, error) {
	// Verify student exists
	student, err := s.studentRepo.FindByID(studentID)
	if err != nil {
		return nil, NewInternalError("failed to verify student", err)
	}
	if student == nil {
		return nil, NewNotFoundError(fmt.Sprintf("Student with ID %d not found", studentID), ErrStudentNotFound)
	}

	grades, err := s.gradeRepo.FindByStudentID(studentID)
	if err != nil {
		return nil, NewInternalError("failed to retrieve grades", err)
	}
	if grades == nil {
		grades = []model.Grade{}
	}
	return grades, nil
}

func (s *gradeService) GetGradeByID(id uint) (*model.Grade, error) {
	grade, err := s.gradeRepo.FindByID(id)
	if err != nil {
		return nil, NewInternalError("failed to query database", err)
	}
	if grade == nil {
		return nil, NewNotFoundError(fmt.Sprintf("Grade with ID %d not found", id), ErrGradeNotFound)
	}
	return grade, nil
}

func (s *gradeService) UpdateGrade(id uint, req *model.UpdateGradeRequest) (*model.Grade, error) {
	grade, err := s.gradeRepo.FindByID(id)
	if err != nil {
		return nil, NewInternalError("failed to query database", err)
	}
	if grade == nil {
		return nil, NewNotFoundError(fmt.Sprintf("Grade with ID %d not found", id), ErrGradeNotFound)
	}

	targetSemester := grade.Semester
	if req.Semester != nil {
		if *req.Semester < 1 || *req.Semester > 14 {
			return nil, NewBadRequestError("Semester must be between 1 and 14", ErrInvalidSemester)
		}
		targetSemester = *req.Semester
	}

	targetHuruf := grade.NilaiHuruf
	targetAngka := grade.NilaiAngka
	if req.NilaiHuruf != nil {
		cleanHuruf := strings.ToUpper(strings.TrimSpace(*req.NilaiHuruf))
		angka, err := model.ConvertGrade(cleanHuruf)
		if err != nil {
			return nil, NewBadRequestError(err.Error(), ErrInvalidNilaiHuruf)
		}
		targetHuruf = cleanHuruf
		targetAngka = angka
	}

	// Check duplicate if semester changed
	if targetSemester != grade.Semester {
		existing, err := s.gradeRepo.FindByStudentAndCourseAndSemester(grade.StudentID, grade.CourseID, targetSemester)
		if err != nil {
			return nil, NewInternalError("failed to check existing grade", err)
		}
		if existing != nil && existing.ID != grade.ID {
			return nil, NewConflictError(
				fmt.Sprintf("Grade for student ID %d, course ID %d, in semester %d already exists", grade.StudentID, grade.CourseID, targetSemester),
				ErrDuplicateGrade,
			)
		}
	}

	grade.Semester = targetSemester
	grade.NilaiHuruf = targetHuruf
	grade.NilaiAngka = targetAngka

	if err := s.gradeRepo.Update(grade); err != nil {
		return nil, NewInternalError("failed to update grade", err)
	}

	return grade, nil
}

func (s *gradeService) DeleteGrade(id uint) error {
	grade, err := s.gradeRepo.FindByID(id)
	if err != nil {
		return NewInternalError("failed to query database", err)
	}
	if grade == nil {
		return NewNotFoundError(fmt.Sprintf("Grade with ID %d not found", id), ErrGradeNotFound)
	}

	if err := s.gradeRepo.Delete(id); err != nil {
		return NewInternalError("failed to delete grade", err)
	}
	return nil
}

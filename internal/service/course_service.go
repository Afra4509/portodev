package service

import (
	"fmt"
	"strings"

	"rest-api-ipk-mahasiswa-its/internal/model"
	"rest-api-ipk-mahasiswa-its/internal/repository"
)

type CourseService interface {
	CreateCourse(req *model.CreateCourseRequest) (*model.Course, error)
	GetAllCourses(search string) ([]model.Course, error)
	GetCourseByID(id uint) (*model.Course, error)
	UpdateCourse(id uint, req *model.UpdateCourseRequest) (*model.Course, error)
	DeleteCourse(id uint) error
}

type courseService struct {
	courseRepo repository.CourseRepository
}

func NewCourseService(courseRepo repository.CourseRepository) CourseService {
	return &courseService{courseRepo: courseRepo}
}

func (s *courseService) validateCourseData(kode, nama string, sks int) error {
	if strings.TrimSpace(kode) == "" {
		return NewBadRequestError("Kode mata kuliah cannot be empty", ErrRequiredFieldEmpty)
	}
	if strings.TrimSpace(nama) == "" {
		return NewBadRequestError("Nama mata kuliah cannot be empty", ErrRequiredFieldEmpty)
	}
	if sks < 1 || sks > 6 {
		return NewBadRequestError("SKS must be between 1 and 6", ErrInvalidSKS)
	}
	return nil
}

func (s *courseService) CreateCourse(req *model.CreateCourseRequest) (*model.Course, error) {
	if err := s.validateCourseData(req.KodeMataKuliah, req.NamaMataKuliah, req.SKS); err != nil {
		return nil, err
	}

	cleanKode := strings.ToUpper(strings.TrimSpace(req.KodeMataKuliah))
	existing, err := s.courseRepo.FindByKode(cleanKode)
	if err != nil {
		return nil, NewInternalError("failed to query database", err)
	}
	if existing != nil {
		return nil, NewConflictError(fmt.Sprintf("Course with code '%s' already exists", cleanKode), ErrDuplicateKodeMK)
	}

	course := &model.Course{
		KodeMataKuliah: cleanKode,
		NamaMataKuliah: strings.TrimSpace(req.NamaMataKuliah),
		SKS:            req.SKS,
	}

	if err := s.courseRepo.Create(course); err != nil {
		return nil, NewInternalError("failed to save course", err)
	}

	return course, nil
}

func (s *courseService) GetAllCourses(search string) ([]model.Course, error) {
	courses, err := s.courseRepo.FindAll(strings.TrimSpace(search))
	if err != nil {
		return nil, NewInternalError("failed to retrieve courses", err)
	}
	if courses == nil {
		courses = []model.Course{}
	}
	return courses, nil
}

func (s *courseService) GetCourseByID(id uint) (*model.Course, error) {
	course, err := s.courseRepo.FindByID(id)
	if err != nil {
		return nil, NewInternalError("failed to query database", err)
	}
	if course == nil {
		return nil, NewNotFoundError(fmt.Sprintf("Course with ID %d not found", id), ErrCourseNotFound)
	}
	return course, nil
}

func (s *courseService) UpdateCourse(id uint, req *model.UpdateCourseRequest) (*model.Course, error) {
	course, err := s.courseRepo.FindByID(id)
	if err != nil {
		return nil, NewInternalError("failed to query database", err)
	}
	if course == nil {
		return nil, NewNotFoundError(fmt.Sprintf("Course with ID %d not found", id), ErrCourseNotFound)
	}

	if err := s.validateCourseData(req.KodeMataKuliah, req.NamaMataKuliah, req.SKS); err != nil {
		return nil, err
	}

	cleanKode := strings.ToUpper(strings.TrimSpace(req.KodeMataKuliah))
	if cleanKode != course.KodeMataKuliah {
		existing, err := s.courseRepo.FindByKode(cleanKode)
		if err != nil {
			return nil, NewInternalError("failed to verify duplicate course code", err)
		}
		if existing != nil && existing.ID != course.ID {
			return nil, NewConflictError(fmt.Sprintf("Course with code '%s' already exists", cleanKode), ErrDuplicateKodeMK)
		}
	}

	course.KodeMataKuliah = cleanKode
	course.NamaMataKuliah = strings.TrimSpace(req.NamaMataKuliah)
	course.SKS = req.SKS

	if err := s.courseRepo.Update(course); err != nil {
		return nil, NewInternalError("failed to update course", err)
	}

	return course, nil
}

func (s *courseService) DeleteCourse(id uint) error {
	course, err := s.courseRepo.FindByID(id)
	if err != nil {
		return NewInternalError("failed to query database", err)
	}
	if course == nil {
		return NewNotFoundError(fmt.Sprintf("Course with ID %d not found", id), ErrCourseNotFound)
	}

	if err := s.courseRepo.Delete(id); err != nil {
		return NewInternalError("failed to delete course", err)
	}
	return nil
}

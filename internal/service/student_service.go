package service

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"rest-api-ipk-mahasiswa-its/internal/model"
	"rest-api-ipk-mahasiswa-its/internal/repository"
)

var nrpRegex = regexp.MustCompile(`^[0-9]{10,14}$`)

type StudentService interface {
	CreateStudent(req *model.CreateStudentRequest) (*model.Student, error)
	GetAllStudents(search string) ([]model.Student, error)
	GetStudentByID(id uint) (*model.Student, error)
	UpdateStudent(id uint, req *model.UpdateStudentRequest) (*model.Student, error)
	DeleteStudent(id uint) error
}

type studentService struct {
	studentRepo repository.StudentRepository
}

func NewStudentService(studentRepo repository.StudentRepository) StudentService {
	return &studentService{studentRepo: studentRepo}
}

func (s *studentService) validateStudentData(nrp, nama, prodi, fakultas string, angkatan int) error {
	nrp = strings.TrimSpace(nrp)
	if nrp == "" {
		return NewBadRequestError("NRP cannot be empty", ErrRequiredFieldEmpty)
	}
	if !nrpRegex.MatchString(nrp) {
		return NewBadRequestError("NRP must consist of 10 to 14 numeric digits (ITS standard format)", ErrInvalidNRP)
	}
	if strings.TrimSpace(nama) == "" {
		return NewBadRequestError("Nama cannot be empty", ErrRequiredFieldEmpty)
	}
	if strings.TrimSpace(prodi) == "" {
		return NewBadRequestError("Program Studi cannot be empty", ErrRequiredFieldEmpty)
	}
	if strings.TrimSpace(fakultas) == "" {
		return NewBadRequestError("Fakultas cannot be empty", ErrRequiredFieldEmpty)
	}
	currentYear := time.Now().Year()
	if angkatan < 2000 || angkatan > currentYear+1 {
		return NewBadRequestError(fmt.Sprintf("Angkatan must be between 2000 and %d", currentYear+1), errors.New("invalid angkatan"))
	}
	return nil
}

func (s *studentService) CreateStudent(req *model.CreateStudentRequest) (*model.Student, error) {
	if err := s.validateStudentData(req.NRP, req.Nama, req.ProgramStudi, req.Fakultas, req.Angkatan); err != nil {
		return nil, err
	}

	cleanNRP := strings.TrimSpace(req.NRP)
	existing, err := s.studentRepo.FindByNRP(cleanNRP)
	if err != nil {
		return nil, NewInternalError("failed to query database", err)
	}
	if existing != nil {
		return nil, NewConflictError(fmt.Sprintf("Student with NRP '%s' already exists", cleanNRP), ErrDuplicateNRP)
	}

	student := &model.Student{
		NRP:          cleanNRP,
		Nama:         strings.TrimSpace(req.Nama),
		ProgramStudi: strings.TrimSpace(req.ProgramStudi),
		Fakultas:     strings.TrimSpace(req.Fakultas),
		Angkatan:     req.Angkatan,
	}

	if err := s.studentRepo.Create(student); err != nil {
		return nil, NewInternalError("failed to save student to database", err)
	}

	return student, nil
}

func (s *studentService) GetAllStudents(search string) ([]model.Student, error) {
	students, err := s.studentRepo.FindAll(strings.TrimSpace(search))
	if err != nil {
		return nil, NewInternalError("failed to retrieve students", err)
	}
	if students == nil {
		students = []model.Student{}
	}
	return students, nil
}

func (s *studentService) GetStudentByID(id uint) (*model.Student, error) {
	student, err := s.studentRepo.FindByID(id)
	if err != nil {
		return nil, NewInternalError("failed to retrieve student", err)
	}
	if student == nil {
		return nil, NewNotFoundError(fmt.Sprintf("Student with ID %d not found", id), ErrStudentNotFound)
	}
	return student, nil
}

func (s *studentService) UpdateStudent(id uint, req *model.UpdateStudentRequest) (*model.Student, error) {
	student, err := s.studentRepo.FindByID(id)
	if err != nil {
		return nil, NewInternalError("failed to query database", err)
	}
	if student == nil {
		return nil, NewNotFoundError(fmt.Sprintf("Student with ID %d not found", id), ErrStudentNotFound)
	}

	if err := s.validateStudentData(req.NRP, req.Nama, req.ProgramStudi, req.Fakultas, req.Angkatan); err != nil {
		return nil, err
	}

	cleanNRP := strings.TrimSpace(req.NRP)
	if cleanNRP != student.NRP {
		existing, err := s.studentRepo.FindByNRP(cleanNRP)
		if err != nil {
			return nil, NewInternalError("failed to verify duplicate NRP", err)
		}
		if existing != nil && existing.ID != student.ID {
			return nil, NewConflictError(fmt.Sprintf("Student with NRP '%s' already exists", cleanNRP), ErrDuplicateNRP)
		}
	}

	student.NRP = cleanNRP
	student.Nama = strings.TrimSpace(req.Nama)
	student.ProgramStudi = strings.TrimSpace(req.ProgramStudi)
	student.Fakultas = strings.TrimSpace(req.Fakultas)
	student.Angkatan = req.Angkatan

	if err := s.studentRepo.Update(student); err != nil {
		return nil, NewInternalError("failed to update student", err)
	}

	return student, nil
}

func (s *studentService) DeleteStudent(id uint) error {
	student, err := s.studentRepo.FindByID(id)
	if err != nil {
		return NewInternalError("failed to query database", err)
	}
	if student == nil {
		return NewNotFoundError(fmt.Sprintf("Student with ID %d not found", id), ErrStudentNotFound)
	}

	if err := s.studentRepo.Delete(id); err != nil {
		return NewInternalError("failed to delete student", err)
	}
	return nil
}

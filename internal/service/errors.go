package service

import "errors"

var (
	ErrStudentNotFound = errors.New("student not found")
	ErrCourseNotFound  = errors.New("course not found")
	ErrGradeNotFound   = errors.New("grade not found")

	ErrDuplicateNRP    = errors.New("student with this NRP already exists")
	ErrDuplicateKodeMK = errors.New("course with this code already exists")
	ErrDuplicateGrade  = errors.New("grade for this student and course in this semester already exists")

	ErrInvalidNRP         = errors.New("invalid NRP format: must be 10 to 14 numeric digits (ITS standard)")
	ErrInvalidSKS         = errors.New("invalid SKS: must be between 1 and 6")
	ErrInvalidSemester    = errors.New("invalid semester: must be between 1 and 14")
	ErrInvalidNilaiHuruf  = errors.New("invalid letter grade: allowed values are A, AB, B, BC, C, D, E")
	ErrRequiredFieldEmpty = errors.New("required field cannot be empty")
)

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewNotFoundError(message string, err error) *AppError {
	return &AppError{Code: 404, Message: message, Err: err}
}

func NewConflictError(message string, err error) *AppError {
	return &AppError{Code: 409, Message: message, Err: err}
}

func NewBadRequestError(message string, err error) *AppError {
	return &AppError{Code: 400, Message: message, Err: err}
}

func NewInternalError(message string, err error) *AppError {
	return &AppError{Code: 500, Message: message, Err: err}
}

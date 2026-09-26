package repository

import (
	"errors"
	"rest-api-ipk-mahasiswa-its/internal/model"

	"gorm.io/gorm"
)

type GradeRepository interface {
	Create(grade *model.Grade) error
	FindByID(id uint) (*model.Grade, error)
	FindByStudentAndCourseAndSemester(studentID, courseID uint, semester int) (*model.Grade, error)
	FindByStudentID(studentID uint) ([]model.Grade, error)
	FindByStudentIDAndSemester(studentID uint, semester int) ([]model.Grade, error)
	Update(grade *model.Grade) error
	Delete(id uint) error
}

type gradeRepository struct {
	db *gorm.DB
}

func NewGradeRepository(db *gorm.DB) GradeRepository {
	return &gradeRepository{db: db}
}

func (r *gradeRepository) Create(grade *model.Grade) error {
	return r.db.Create(grade).Error
}

func (r *gradeRepository) FindByID(id uint) (*model.Grade, error) {
	var grade model.Grade
	err := r.db.Preload("Course").Preload("Student").First(&grade, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &grade, nil
}

func (r *gradeRepository) FindByStudentAndCourseAndSemester(studentID, courseID uint, semester int) (*model.Grade, error) {
	var grade model.Grade
	err := r.db.Where("student_id = ? AND course_id = ? AND semester = ?", studentID, courseID, semester).First(&grade).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &grade, nil
}

func (r *gradeRepository) FindByStudentID(studentID uint) ([]model.Grade, error) {
	var grades []model.Grade
	err := r.db.Preload("Course").
		Where("student_id = ?", studentID).
		Order("semester ASC, id ASC").
		Find(&grades).Error
	return grades, err
}

func (r *gradeRepository) FindByStudentIDAndSemester(studentID uint, semester int) ([]model.Grade, error) {
	var grades []model.Grade
	err := r.db.Preload("Course").
		Where("student_id = ? AND semester = ?", studentID, semester).
		Order("id ASC").
		Find(&grades).Error
	return grades, err
}

func (r *gradeRepository) Update(grade *model.Grade) error {
	return r.db.Save(grade).Error
}

func (r *gradeRepository) Delete(id uint) error {
	return r.db.Delete(&model.Grade{}, id).Error
}

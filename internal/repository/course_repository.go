package repository

import (
	"errors"
	"rest-api-ipk-mahasiswa-its/internal/model"

	"gorm.io/gorm"
)

type CourseRepository interface {
	Create(course *model.Course) error
	FindAll(search string) ([]model.Course, error)
	FindByID(id uint) (*model.Course, error)
	FindByKode(kode string) (*model.Course, error)
	Update(course *model.Course) error
	Delete(id uint) error
}

type courseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) CourseRepository {
	return &courseRepository{db: db}
}

func (r *courseRepository) Create(course *model.Course) error {
	return r.db.Create(course).Error
}

func (r *courseRepository) FindAll(search string) ([]model.Course, error) {
	var courses []model.Course
	query := r.db.Order("id ASC")
	if search != "" {
		likeSearch := "%" + search + "%"
		query = query.Where("nama_mata_kuliah ILIKE ? OR kode_mata_kuliah ILIKE ?", likeSearch, likeSearch)
	}
	err := query.Find(&courses).Error
	return courses, err
}

func (r *courseRepository) FindByID(id uint) (*model.Course, error) {
	var course model.Course
	err := r.db.First(&course, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &course, nil
}

func (r *courseRepository) FindByKode(kode string) (*model.Course, error) {
	var course model.Course
	err := r.db.Where("kode_mata_kuliah = ?", kode).First(&course).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &course, nil
}

func (r *courseRepository) Update(course *model.Course) error {
	return r.db.Save(course).Error
}

func (r *courseRepository) Delete(id uint) error {
	return r.db.Delete(&model.Course{}, id).Error
}

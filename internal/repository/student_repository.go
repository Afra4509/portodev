package repository

import (
	"errors"
	"rest-api-ipk-mahasiswa-its/internal/model"

	"gorm.io/gorm"
)

type StudentRepository interface {
	Create(student *model.Student) error
	FindAll(search string) ([]model.Student, error)
	FindByID(id uint) (*model.Student, error)
	FindByNRP(nrp string) (*model.Student, error)
	Update(student *model.Student) error
	Delete(id uint) error
}

type studentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) StudentRepository {
	return &studentRepository{db: db}
}

func (r *studentRepository) Create(student *model.Student) error {
	return r.db.Create(student).Error
}

func (r *studentRepository) FindAll(search string) ([]model.Student, error) {
	var students []model.Student
	query := r.db.Order("id ASC")
	if search != "" {
		likeSearch := "%" + search + "%"
		query = query.Where("nama ILIKE ? OR nrp LIKE ? OR program_studi ILIKE ?", likeSearch, likeSearch, likeSearch)
	}
	err := query.Find(&students).Error
	return students, err
}

func (r *studentRepository) FindByID(id uint) (*model.Student, error) {
	var student model.Student
	err := r.db.First(&student, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &student, nil
}

func (r *studentRepository) FindByNRP(nrp string) (*model.Student, error) {
	var student model.Student
	err := r.db.Where("nrp = ?", nrp).First(&student).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &student, nil
}

func (r *studentRepository) Update(student *model.Student) error {
	return r.db.Save(student).Error
}

func (r *studentRepository) Delete(id uint) error {
	return r.db.Delete(&model.Student{}, id).Error
}

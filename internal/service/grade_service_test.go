package service_test

import (
	"testing"

	"rest-api-ipk-mahasiswa-its/internal/model"
	"rest-api-ipk-mahasiswa-its/internal/repository"
	"rest-api-ipk-mahasiswa-its/internal/service"
	"rest-api-ipk-mahasiswa-its/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGradeService_CRUD(t *testing.T) {
	db := testutil.SetupTestDB()
	studentRepo := repository.NewStudentRepository(db)
	courseRepo := repository.NewCourseRepository(db)
	gradeRepo := repository.NewGradeRepository(db)

	gradeService := service.NewGradeService(gradeRepo, studentRepo, courseRepo)

	// Seed prerequisite student and course
	student := &model.Student{
		NRP:          "5025211001",
		Nama:         "Budi Santoso",
		ProgramStudi: "Teknik Informatika",
		Fakultas:     "FTEIC",
		Angkatan:     2021,
	}
	require.NoError(t, studentRepo.Create(student))

	course := &model.Course{
		KodeMataKuliah: "IF184101",
		NamaMataKuliah: "Dasar Pemrograman",
		SKS:            3,
	}
	require.NoError(t, courseRepo.Create(course))

	// 1. Create Grade Success
	createReq := &model.CreateGradeRequest{
		StudentID:  student.ID,
		CourseID:   course.ID,
		Semester:   1,
		NilaiHuruf: "A",
	}
	grade, err := gradeService.CreateGrade(createReq)
	require.NoError(t, err)
	require.NotNil(t, grade)
	assert.Equal(t, "A", grade.NilaiHuruf)
	assert.Equal(t, 4.00, grade.NilaiAngka)

	// 2. Duplicate Grade in Same Semester -> 409 Conflict
	_, err = gradeService.CreateGrade(createReq)
	assert.Error(t, err)
	var appErr *service.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, 409, appErr.Code)

	// 3. Create Grade with Invalid Letter Grade -> 400 Bad Request
	invalidGradeReq := &model.CreateGradeRequest{
		StudentID:  student.ID,
		CourseID:   course.ID,
		Semester:   2,
		NilaiHuruf: "X",
	}
	_, err = gradeService.CreateGrade(invalidGradeReq)
	assert.Error(t, err)
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, 400, appErr.Code)

	// 4. Create Grade for Non-Existent Student -> 404 Not Found
	nonExistentStudentReq := &model.CreateGradeRequest{
		StudentID:  9999,
		CourseID:   course.ID,
		Semester:   1,
		NilaiHuruf: "A",
	}
	_, err = gradeService.CreateGrade(nonExistentStudentReq)
	assert.Error(t, err)
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, 404, appErr.Code)

	// 5. Get Grades by Student ID
	grades, err := gradeService.GetGradesByStudentID(student.ID)
	require.NoError(t, err)
	assert.Len(t, grades, 1)

	// 6. Update Grade
	newHuruf := "AB"
	updateReq := &model.UpdateGradeRequest{
		NilaiHuruf: &newHuruf,
	}
	updated, err := gradeService.UpdateGrade(grade.ID, updateReq)
	require.NoError(t, err)
	assert.Equal(t, "AB", updated.NilaiHuruf)
	assert.Equal(t, 3.50, updated.NilaiAngka)

	// 7. Delete Grade
	err = gradeService.DeleteGrade(grade.ID)
	require.NoError(t, err)

	// Verify Deleted
	_, err = gradeService.GetGradeByID(grade.ID)
	assert.Error(t, err)
}

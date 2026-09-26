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

func TestStudentService_CRUD(t *testing.T) {
	db := testutil.SetupTestDB()
	studentRepo := repository.NewStudentRepository(db)
	studentService := service.NewStudentService(studentRepo)

	// 1. Create Student Success
	req := &model.CreateStudentRequest{
		NRP:          "5025211001",
		Nama:         "Budi Santoso",
		ProgramStudi: "Teknik Informatika",
		Fakultas:     "FTEIC",
		Angkatan:     2021,
	}
	student, err := studentService.CreateStudent(req)
	require.NoError(t, err)
	require.NotNil(t, student)
	assert.Equal(t, "5025211001", student.NRP)
	assert.Equal(t, "Budi Santoso", student.Nama)

	// 2. Create Duplicate NRP -> Conflict Error
	_, err = studentService.CreateStudent(req)
	assert.Error(t, err)
	var appErr *service.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, 409, appErr.Code)

	// 3. Create Invalid NRP (less than 10 digits or non-numeric)
	invalidReq := &model.CreateStudentRequest{
		NRP:          "12345", // too short
		Nama:         "Invalid Student",
		ProgramStudi: "Teknik Informatika",
		Fakultas:     "FTEIC",
		Angkatan:     2021,
	}
	_, err = studentService.CreateStudent(invalidReq)
	assert.Error(t, err)
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, 400, appErr.Code)

	// 4. Get Student By ID
	fetched, err := studentService.GetStudentByID(student.ID)
	require.NoError(t, err)
	assert.Equal(t, student.Nama, fetched.Nama)

	// 5. Get Student By Non-Existent ID -> 404 Not Found
	_, err = studentService.GetStudentByID(9999)
	assert.Error(t, err)
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, 404, appErr.Code)

	// 6. GetAllStudents
	list, err := studentService.GetAllStudents("")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(list), 1)

	// 7. Update Student
	updateReq := &model.UpdateStudentRequest{
		NRP:          "5025211001",
		Nama:         "Budi Santoso S.Kom",
		ProgramStudi: "Teknik Informatika",
		Fakultas:     "FTEIC",
		Angkatan:     2021,
	}
	updated, err := studentService.UpdateStudent(student.ID, updateReq)
	require.NoError(t, err)
	assert.Equal(t, "Budi Santoso S.Kom", updated.Nama)

	// 8. Delete Student
	err = studentService.DeleteStudent(student.ID)
	require.NoError(t, err)

	// Verify Deleted
	_, err = studentService.GetStudentByID(student.ID)
	assert.Error(t, err)
}

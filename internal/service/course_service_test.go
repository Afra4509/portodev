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

func TestCourseService_CRUD(t *testing.T) {
	db := testutil.SetupTestDB()
	courseRepo := repository.NewCourseRepository(db)
	courseService := service.NewCourseService(courseRepo)

	// 1. Create Course Success
	createReq := &model.CreateCourseRequest{
		KodeMataKuliah: "IF184101",
		NamaMataKuliah: "Dasar Pemrograman",
		SKS:            3,
	}
	course, err := courseService.CreateCourse(createReq)
	require.NoError(t, err)
	require.NotNil(t, course)
	assert.Equal(t, "IF184101", course.KodeMataKuliah)
	assert.Equal(t, 3, course.SKS)

	// 2. Duplicate Course Code -> 409 Conflict
	_, err = courseService.CreateCourse(createReq)
	assert.Error(t, err)
	var appErr *service.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, 409, appErr.Code)

	// 3. Invalid SKS (< 1 or > 6) -> 400 Bad Request
	invalidReq := &model.CreateCourseRequest{
		KodeMataKuliah: "IF999999",
		NamaMataKuliah: "Invalid SKS Course",
		SKS:            10,
	}
	_, err = courseService.CreateCourse(invalidReq)
	assert.Error(t, err)
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, 400, appErr.Code)

	// 4. Get Course By ID
	fetched, err := courseService.GetCourseByID(course.ID)
	require.NoError(t, err)
	assert.Equal(t, course.NamaMataKuliah, fetched.NamaMataKuliah)

	// 5. Get All Courses
	list, err := courseService.GetAllCourses("")
	require.NoError(t, err)
	assert.Len(t, list, 1)

	// 6. Update Course
	updateReq := &model.UpdateCourseRequest{
		KodeMataKuliah: "IF184101",
		NamaMataKuliah: "Dasar Pemrograman Lanjut",
		SKS:            4,
	}
	updated, err := courseService.UpdateCourse(course.ID, updateReq)
	require.NoError(t, err)
	assert.Equal(t, "Dasar Pemrograman Lanjut", updated.NamaMataKuliah)
	assert.Equal(t, 4, updated.SKS)

	// 7. Delete Course
	err = courseService.DeleteCourse(course.ID)
	require.NoError(t, err)

	// Verify Deleted
	_, err = courseService.GetCourseByID(course.ID)
	assert.Error(t, err)
}

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

func TestGradeConversion(t *testing.T) {
	testCases := []struct {
		input       string
		expected    float64
		expectError bool
	}{
		{"A", 4.00, false},
		{"a", 4.00, false},
		{"AB", 3.50, false},
		{"ab", 3.50, false},
		{"B", 3.00, false},
		{"BC", 2.50, false},
		{"C", 2.00, false},
		{"D", 1.00, false},
		{"E", 0.00, false},
		{"F", 0.00, true},
		{"A+", 0.00, true},
		{"", 0.00, true},
	}

	for _, tc := range testCases {
		val, err := model.ConvertGrade(tc.input)
		if tc.expectError {
			assert.Error(t, err, "expected error for grade %s", tc.input)
		} else {
			assert.NoError(t, err, "unexpected error for grade %s", tc.input)
			assert.Equal(t, tc.expected, val, "grade point mismatch for %s", tc.input)
		}
	}
}

func TestAcademicService_CalculateIPSemester(t *testing.T) {
	db := testutil.SetupTestDB()
	studentRepo := repository.NewStudentRepository(db)
	courseRepo := repository.NewCourseRepository(db)
	gradeRepo := repository.NewGradeRepository(db)

	academicService := service.NewAcademicService(studentRepo, gradeRepo)

	// 1. Create Student
	student := &model.Student{
		NRP:          "5025211001",
		Nama:         "Budi Santoso",
		ProgramStudi: "Teknik Informatika",
		Fakultas:     "FTEIC",
		Angkatan:     2021,
	}
	require.NoError(t, studentRepo.Create(student))

	// 2. Create Courses
	course1 := &model.Course{KodeMataKuliah: "IF184101", NamaMataKuliah: "Dasar Pemrograman", SKS: 3}
	course2 := &model.Course{KodeMataKuliah: "KM184101", NamaMataKuliah: "Kalkulus 1", SKS: 3}
	course3 := &model.Course{KodeMataKuliah: "UG184914", NamaMataKuliah: "Bahasa Inggris", SKS: 2}
	require.NoError(t, courseRepo.Create(course1))
	require.NoError(t, courseRepo.Create(course2))
	require.NoError(t, courseRepo.Create(course3))

	// 3. Assign Grades for Semester 1:
	// Course1: 3 SKS, A (4.00)  -> Mutu = 12.0
	// Course2: 3 SKS, AB (3.50) -> Mutu = 10.5
	// Course3: 2 SKS, B (3.00)  -> Mutu = 6.0
	// Total SKS = 8, Total Mutu = 28.5, IP = 28.5 / 8 = 3.5625 -> rounded to 3.56
	require.NoError(t, gradeRepo.Create(&model.Grade{
		StudentID: student.ID, CourseID: course1.ID, Semester: 1, NilaiHuruf: "A", NilaiAngka: 4.00,
	}))
	require.NoError(t, gradeRepo.Create(&model.Grade{
		StudentID: student.ID, CourseID: course2.ID, Semester: 1, NilaiHuruf: "AB", NilaiAngka: 3.50,
	}))
	require.NoError(t, gradeRepo.Create(&model.Grade{
		StudentID: student.ID, CourseID: course3.ID, Semester: 1, NilaiHuruf: "B", NilaiAngka: 3.00,
	}))

	// Execute CalculateIPSemester
	res, err := academicService.CalculateIPSemester(student.ID, 1)
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, student.ID, res.StudentID)
	assert.Equal(t, "5025211001", res.NRP)
	assert.Equal(t, 1, res.Semester)
	assert.Equal(t, 8, res.TotalSKS)
	assert.Equal(t, 28.5, res.TotalMutu)
	assert.Equal(t, 3.56, res.IP)
	assert.Equal(t, 3, len(res.Courses))

	// Test Semester with no courses
	emptyRes, err := academicService.CalculateIPSemester(student.ID, 2)
	require.NoError(t, err)
	assert.Equal(t, 0, emptyRes.TotalSKS)
	assert.Equal(t, 0.0, emptyRes.TotalMutu)
	assert.Equal(t, 0.0, emptyRes.IP)

	// Test Non-existent student
	_, err = academicService.CalculateIPSemester(9999, 1)
	assert.Error(t, err)

	// Test Invalid semester
	_, err = academicService.CalculateIPSemester(student.ID, 15)
	assert.Error(t, err)
}

func TestAcademicService_CalculateIPK(t *testing.T) {
	db := testutil.SetupTestDB()
	studentRepo := repository.NewStudentRepository(db)
	courseRepo := repository.NewCourseRepository(db)
	gradeRepo := repository.NewGradeRepository(db)

	academicService := service.NewAcademicService(studentRepo, gradeRepo)

	student := &model.Student{
		NRP:          "5025211045",
		Nama:         "Siti Aminah",
		ProgramStudi: "Teknik Informatika",
		Fakultas:     "FTEIC",
		Angkatan:     2021,
	}
	require.NoError(t, studentRepo.Create(student))

	course1 := &model.Course{KodeMataKuliah: "IF184101", NamaMataKuliah: "Dasar Pemrograman", SKS: 3}
	course2 := &model.Course{KodeMataKuliah: "KM184101", NamaMataKuliah: "Kalkulus 1", SKS: 3}
	course3 := &model.Course{KodeMataKuliah: "IF184201", NamaMataKuliah: "Struktur Data", SKS: 3}
	require.NoError(t, courseRepo.Create(course1))
	require.NoError(t, courseRepo.Create(course2))
	require.NoError(t, courseRepo.Create(course3))

	// Semester 1:
	// Course1: A (4.0) -> Mutu = 12.0
	// Course2: C (2.0) -> Mutu = 6.0
	require.NoError(t, gradeRepo.Create(&model.Grade{
		StudentID: student.ID, CourseID: course1.ID, Semester: 1, NilaiHuruf: "A", NilaiAngka: 4.00,
	}))
	require.NoError(t, gradeRepo.Create(&model.Grade{
		StudentID: student.ID, CourseID: course2.ID, Semester: 1, NilaiHuruf: "C", NilaiAngka: 2.00,
	}))

	// Semester 2:
	// Course3: A (4.0) -> Mutu = 12.0
	// Course2 (Retake Kalkulus 1!): AB (3.50) -> Best grade should be AB (3.50), Mutu = 10.5
	require.NoError(t, gradeRepo.Create(&model.Grade{
		StudentID: student.ID, CourseID: course3.ID, Semester: 2, NilaiHuruf: "A", NilaiAngka: 4.00,
	}))
	require.NoError(t, gradeRepo.Create(&model.Grade{
		StudentID: student.ID, CourseID: course2.ID, Semester: 2, NilaiHuruf: "AB", NilaiAngka: 3.50,
	}))

	// Calculate IPK:
	// Unique courses:
	// - Course1: 3 SKS, A (4.00)  -> Mutu = 12.0
	// - Course2: 3 SKS, AB (3.50) -> Mutu = 10.5 (highest grade chosen over C)
	// - Course3: 3 SKS, A (4.00)  -> Mutu = 12.0
	// Total SKS = 9, Total Mutu = 34.5, IPK = 34.5 / 9 = 3.8333 -> 3.83
	res, err := academicService.CalculateIPK(student.ID)
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, student.ID, res.StudentID)
	assert.Equal(t, "5025211045", res.NRP)
	assert.Equal(t, "Siti Aminah", res.Nama)
	assert.Equal(t, 9, res.TotalSKS)
	assert.Equal(t, 34.5, res.TotalMutu)
	assert.Equal(t, 3.83, res.IPK)
	assert.Equal(t, 3, res.TotalMataKuliah)
}

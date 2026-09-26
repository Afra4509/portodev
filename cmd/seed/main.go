package main

import (
	"log"

	"rest-api-ipk-mahasiswa-its/internal/config"
	"rest-api-ipk-mahasiswa-its/internal/model"
)

func main() {
	log.Println("[INFO] Starting database seeder...")

	cfg := config.LoadConfig()
	db, err := config.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("[FATAL] Could not connect to database for seeding: %v\n", err)
	}

	// 1. Seed Students
	students := []model.Student{
		{
			NRP:          "5025211001",
			Nama:         "Budi Santoso",
			ProgramStudi: "Teknik Informatika",
			Fakultas:     "FTEIC",
			Angkatan:     2021,
		},
		{
			NRP:          "5025211045",
			Nama:         "Siti Aminah",
			ProgramStudi: "Teknik Informatika",
			Fakultas:     "FTEIC",
			Angkatan:     2021,
		},
		{
			NRP:          "5026211012",
			Nama:         "Ahmad Fauzi",
			ProgramStudi: "Sistem Informasi",
			Fakultas:     "FTEIC",
			Angkatan:     2021,
		},
	}

	for _, s := range students {
		var existing model.Student
		if err := db.Where("nrp = ?", s.NRP).First(&existing).Error; err != nil {
			if err := db.Create(&s).Error; err != nil {
				log.Printf("[ERROR] Failed to seed student %s: %v\n", s.NRP, err)
			} else {
				log.Printf("[SUCCESS] Seeded student: %s (%s)\n", s.Nama, s.NRP)
			}
		} else {
			log.Printf("[SKIP] Student %s already exists\n", s.NRP)
		}
	}

	// 2. Seed Courses
	courses := []model.Course{
		{KodeMataKuliah: "IF184101", NamaMataKuliah: "Dasar Pemrograman", SKS: 3},
		{KodeMataKuliah: "IF184201", NamaMataKuliah: "Struktur Data", SKS: 3},
		{KodeMataKuliah: "IF184301", NamaMataKuliah: "Pemrograman Berorientasi Objek", SKS: 3},
		{KodeMataKuliah: "IF184401", NamaMataKuliah: "Basis Data", SKS: 4},
		{KodeMataKuliah: "KM184101", NamaMataKuliah: "Kalkulus 1", SKS: 3},
		{KodeMataKuliah: "UG184914", NamaMataKuliah: "Bahasa Inggris", SKS: 2},
	}

	for _, c := range courses {
		var existing model.Course
		if err := db.Where("kode_mata_kuliah = ?", c.KodeMataKuliah).First(&existing).Error; err != nil {
			if err := db.Create(&c).Error; err != nil {
				log.Printf("[ERROR] Failed to seed course %s: %v\n", c.KodeMataKuliah, err)
			} else {
				log.Printf("[SUCCESS] Seeded course: %s (%s, %d SKS)\n", c.NamaMataKuliah, c.KodeMataKuliah, c.SKS)
			}
		} else {
			log.Printf("[SKIP] Course %s already exists\n", c.KodeMataKuliah)
		}
	}

	// 3. Seed Sample Grades for Student Budi Santoso (NRP 5025211001)
	var budi model.Student
	if err := db.Where("nrp = ?", "5025211001").First(&budi).Error; err == nil {
		var dasprog, kalkulus, bing, strukdat, pbo, basdat model.Course
		db.Where("kode_mata_kuliah = ?", "IF184101").First(&dasprog)
		db.Where("kode_mata_kuliah = ?", "KM184101").First(&kalkulus)
		db.Where("kode_mata_kuliah = ?", "UG184914").First(&bing)
		db.Where("kode_mata_kuliah = ?", "IF184201").First(&strukdat)
		db.Where("kode_mata_kuliah = ?", "IF184301").First(&pbo)
		db.Where("kode_mata_kuliah = ?", "IF184401").First(&basdat)

		sampleGrades := []model.Grade{
			// Semester 1
			{StudentID: budi.ID, CourseID: dasprog.ID, Semester: 1, NilaiHuruf: "A", NilaiAngka: 4.00},
			{StudentID: budi.ID, CourseID: kalkulus.ID, Semester: 1, NilaiHuruf: "AB", NilaiAngka: 3.50},
			{StudentID: budi.ID, CourseID: bing.ID, Semester: 1, NilaiHuruf: "A", NilaiAngka: 4.00},

			// Semester 2
			{StudentID: budi.ID, CourseID: strukdat.ID, Semester: 2, NilaiHuruf: "A", NilaiAngka: 4.00},
			{StudentID: budi.ID, CourseID: pbo.ID, Semester: 2, NilaiHuruf: "B", NilaiAngka: 3.00},
			{StudentID: budi.ID, CourseID: basdat.ID, Semester: 2, NilaiHuruf: "AB", NilaiAngka: 3.50},
		}

		for _, g := range sampleGrades {
			var existing model.Grade
			if err := db.Where("student_id = ? AND course_id = ? AND semester = ?", g.StudentID, g.CourseID, g.Semester).First(&existing).Error; err != nil {
				if err := db.Create(&g).Error; err != nil {
					log.Printf("[ERROR] Failed to seed grade: %v\n", err)
				} else {
					log.Printf("[SUCCESS] Seeded grade: Student %d, Course %d, Semester %d: %s (%.2f)\n", g.StudentID, g.CourseID, g.Semester, g.NilaiHuruf, g.NilaiAngka)
				}
			}
		}
	}

	log.Println("[INFO] Database seeding completed successfully!")
}

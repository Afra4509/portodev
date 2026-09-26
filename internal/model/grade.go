package model

import "time"

// Grade represents student course grade in a specific semester (Nilai Mahasiswa)
type Grade struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	StudentID  uint      `gorm:"not null;uniqueIndex:idx_student_course_semester" json:"student_id"`
	CourseID   uint      `gorm:"not null;uniqueIndex:idx_student_course_semester" json:"course_id"`
	Semester   int       `gorm:"not null;uniqueIndex:idx_student_course_semester" json:"semester"`
	NilaiHuruf string    `gorm:"type:varchar(5);not null" json:"nilai_huruf"`
	NilaiAngka float64   `gorm:"type:decimal(4,2);not null" json:"nilai_angka"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	Student *Student `gorm:"foreignKey:StudentID;constraint:OnDelete:CASCADE;" json:"student,omitempty"`
	Course  *Course  `gorm:"foreignKey:CourseID;constraint:OnDelete:RESTRICT;" json:"course,omitempty"`
}

func (Grade) TableName() string {
	return "grades"
}

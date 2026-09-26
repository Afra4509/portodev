package model

import "time"

// Student represents the student entity (Mahasiswa ITS)
type Student struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	NRP          string    `gorm:"type:varchar(20);uniqueIndex;not null" json:"nrp"`
	Nama         string    `gorm:"type:varchar(100);not null" json:"nama"`
	ProgramStudi string    `gorm:"type:varchar(100);not null" json:"program_studi"`
	Fakultas     string    `gorm:"type:varchar(100);not null" json:"fakultas"`
	Angkatan     int       `gorm:"not null" json:"angkatan"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	Grades []Grade `gorm:"foreignKey:StudentID;constraint:OnDelete:CASCADE;" json:"grades,omitempty"`
}

func (Student) TableName() string {
	return "students"
}

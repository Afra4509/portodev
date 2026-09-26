package model

import "time"

// Course represents academic course entity (Mata Kuliah)
type Course struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	KodeMataKuliah string    `gorm:"type:varchar(20);uniqueIndex;not null" json:"kode_mata_kuliah"`
	NamaMataKuliah string    `gorm:"type:varchar(100);not null" json:"nama_mata_kuliah"`
	SKS            int       `gorm:"not null" json:"sks"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	Grades []Grade `gorm:"foreignKey:CourseID;constraint:OnDelete:RESTRICT;" json:"grades,omitempty"`
}

func (Course) TableName() string {
	return "courses"
}

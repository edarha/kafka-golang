package models

// the relationship between student and class table. many-many
type ClassStudent struct {
	ID        string `gorm:"id;not null;primaryKey"`
	ClassID   string `gorm:"class_id;not null"`
	StudentID string `gorm:"class_id;not null"`
}

func (*ClassStudent) TableName() string {
	return "class_student"
}

package domain

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleStudent    Role = "STUDENT"
	RoleInstructor Role = "INSTRUCTOR"
	RoleAdmin      Role = "ADMIN"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	FullName  string         `gorm:"not null" json:"full_name"`
	Role      Role           `gorm:"type:varchar(20);default:'STUDENT'" json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Course struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	InstructorID uint          `gorm:"not null" json:"instructor_id"`
	Instructor  User           `gorm:"foreignKey:InstructorID" json:"instructor,omitempty"`
	Lessons     []Lesson       `gorm:"foreignKey:CourseID" json:"lessons,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Lesson struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CourseID    uint           `gorm:"not null" json:"course_id"`
	Title       string         `gorm:"not null" json:"title"`
	Content     string         `gorm:"type:text" json:"content"`
	VideoURL    string         `json:"video_url"`
	Sequence    int            `gorm:"default:0" json:"sequence"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Enrollment struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"uniqueIndex:idx_user_course;not null" json:"user_id"`
	CourseID   uint      `gorm:"uniqueIndex:idx_user_course;not null" json:"course_id"`
	User       User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Course     Course    `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	EnrolledAt time.Time `json:"enrolled_at"`
}

type LessonProgress struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"uniqueIndex:idx_user_lesson;not null" json:"user_id"`
	LessonID    uint      `gorm:"uniqueIndex:idx_user_lesson;not null" json:"lesson_id"`
	IsCompleted bool      `gorm:"default:false" json:"is_completed"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
}


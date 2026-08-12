package http

import (
	"net/http"

	"mini-lms/internal/domain"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CourseHandler struct {
	db *gorm.DB
}

func NewCourseHandler(db *gorm.DB) *CourseHandler {
	return &CourseHandler{db: db}
}

type CreateCourseRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

func (h *CourseHandler) CreateCourse(c *gin.Context) {
	var req CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	instructorID := c.MustGet("userID").(uint)

	course := domain.Course{
		Title:        req.Title,
		Description:  req.Description,
		InstructorID: instructorID,
	}

	if err := h.db.Create(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Course created successfully",
		"course":  course,
	})
}

func (h *CourseHandler) ListCourses(c *gin.Context) {
	var courses []domain.Course
	if err := h.db.Preload("Instructor").Preload("Lessons").Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch courses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"courses": courses,
	})
}

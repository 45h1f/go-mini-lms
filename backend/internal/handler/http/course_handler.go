package http

import (
	"log"
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
	if h.db == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database connection is not initialized"})
		return
	}

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
		log.Printf("Error creating course: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create course: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Course created successfully",
		"course":  course,
	})
}

func (h *CourseHandler) ListCourses(c *gin.Context) {
	courses := make([]domain.Course, 0)

	if h.db == nil {
		c.JSON(http.StatusOK, gin.H{
			"courses": courses,
		})
		return
	}

	if err := h.db.Preload("Instructor").Preload("Lessons").Find(&courses).Error; err != nil {
		log.Printf("Error querying courses: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"courses": make([]domain.Course, 0),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"courses": courses,
	})
}

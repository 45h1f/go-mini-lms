package http

import (
	"net/http"
	"strconv"
	"time"

	"mini-lms/internal/domain"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EnrollmentHandler struct {
	db *gorm.DB
}

func NewEnrollmentHandler(db *gorm.DB) *EnrollmentHandler {
	return &EnrollmentHandler{db: db}
}

func (h *EnrollmentHandler) EnrollCourse(c *gin.Context) {
	courseIDParam := c.Param("id")
	courseID, err := strconv.ParseUint(courseIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	userID := c.MustGet("userID").(uint)

	enrollment := domain.Enrollment{
		UserID:     userID,
		CourseID:   uint(courseID),
		EnrolledAt: time.Now(),
	}

	if err := h.db.Create(&enrollment).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Already enrolled in this course or invalid course ID"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Successfully enrolled in course",
		"enrollment": enrollment,
	})
}

func (h *EnrollmentHandler) GetMyEnrollments(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var enrollments []domain.Enrollment
	if err := h.db.Preload("Course").Preload("Course.Instructor").Where("user_id = ?", userID).Find(&enrollments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch enrollments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"enrollments": enrollments,
	})
}

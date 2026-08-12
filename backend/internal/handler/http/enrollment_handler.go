package http

import (
	"log"
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
	if h.db == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database connection is not available"})
		return
	}

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
		log.Printf("Enrollment error: %v", err)
		c.JSON(http.StatusConflict, gin.H{"error": "Already enrolled in this course or invalid course ID"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Successfully enrolled in course",
		"enrollment": enrollment,
	})
}

func (h *EnrollmentHandler) GetMyEnrollments(c *gin.Context) {
	var enrollments []domain.Enrollment

	if h.db == nil {
		c.JSON(http.StatusOK, gin.H{
			"enrollments": []domain.Enrollment{},
		})
		return
	}

	userID := c.MustGet("userID").(uint)

	if err := h.db.Preload("Course").Preload("Course.Instructor").Where("user_id = ?", userID).Find(&enrollments).Error; err != nil {
		log.Printf("Error fetching enrollments: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"enrollments": []domain.Enrollment{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"enrollments": enrollments,
	})
}

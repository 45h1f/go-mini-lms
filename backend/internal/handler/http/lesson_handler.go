package http

import (
	"log"
	"net/http"
	"strconv"

	"mini-lms/internal/domain"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LessonHandler struct {
	db *gorm.DB
}

func NewLessonHandler(db *gorm.DB) *LessonHandler {
	return &LessonHandler{db: db}
}

type CreateLessonRequest struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content"`
	VideoURL string `json:"video_url"`
	Sequence int    `json:"sequence"`
}

func (h *LessonHandler) CreateLesson(c *gin.Context) {
	if h.db == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database connection unavailable"})
		return
	}

	courseIDParam := c.Param("id")
	courseID, err := strconv.ParseUint(courseIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	var course domain.Course
	if err := h.db.First(&course, courseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	userID := c.MustGet("userID").(uint)
	userRole := c.MustGet("userRole").(domain.Role)

	if userRole != domain.RoleAdmin && course.InstructorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only the course instructor or an admin can add lessons"})
		return
	}

	var req CreateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lesson := domain.Lesson{
		CourseID: uint(courseID),
		Title:    req.Title,
		Content:  req.Content,
		VideoURL: req.VideoURL,
		Sequence: req.Sequence,
	}

	if err := h.db.Create(&lesson).Error; err != nil {
		log.Printf("Create lesson error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create lesson"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Lesson created successfully",
		"lesson":  lesson,
	})
}

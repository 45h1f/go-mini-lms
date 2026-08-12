package http

import (
	"net/http"
	"strconv"
	"time"

	"mini-lms/internal/domain"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProgressHandler struct {
	db *gorm.DB
}

func NewProgressHandler(db *gorm.DB) *ProgressHandler {
	return &ProgressHandler{db: db}
}

type TrackProgressRequest struct {
	LessonID    uint `json:"lesson_id" binding:"required"`
	IsCompleted bool `json:"is_completed"`
}

func (h *ProgressHandler) MarkLessonProgress(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var req TrackProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var lesson domain.Lesson
	if err := h.db.First(&lesson, req.LessonID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Lesson not found"})
		return
	}

	// Verify user enrollment in course
	var enrollment domain.Enrollment
	if err := h.db.Where("user_id = ? AND course_id = ?", userID, lesson.CourseID).First(&enrollment).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "User is not enrolled in this course"})
		return
	}

	var progress domain.LessonProgress
	err := h.db.Where("user_id = ? AND lesson_id = ?", userID, req.LessonID).First(&progress).Error
	if err == gorm.ErrRecordNotFound {
		progress = domain.LessonProgress{
			UserID:      userID,
			LessonID:    req.LessonID,
			IsCompleted: req.IsCompleted,
			CompletedAt: time.Now(),
		}
		if err := h.db.Create(&progress).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record lesson progress"})
			return
		}
	} else {
		progress.IsCompleted = req.IsCompleted
		if req.IsCompleted {
			progress.CompletedAt = time.Now()
		}
		if err := h.db.Save(&progress).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update lesson progress"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Lesson progress updated",
		"progress": progress,
	})
}

func (h *ProgressHandler) GetCourseProgress(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	courseIDParam := c.Param("id")
	courseID, err := strconv.ParseUint(courseIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	var totalLessons int64
	h.db.Model(&domain.Lesson{}).Where("course_id = ?", courseID).Count(&totalLessons)

	var completedLessons int64
	h.db.Table("lesson_progresses").
		Joins("JOIN lessons ON lessons.id = lesson_progresses.lesson_id").
		Where("lessons.course_id = ? AND lesson_progresses.user_id = ? AND lesson_progresses.is_completed = ?", courseID, userID, true).
		Count(&completedLessons)

	percentage := 0.0
	if totalLessons > 0 {
		percentage = (float64(completedLessons) / float64(totalLessons)) * 100.0
	}

	c.JSON(http.StatusOK, gin.H{
		"course_id":         courseID,
		"total_lessons":     totalLessons,
		"completed_lessons": completedLessons,
		"progress_percent":  percentage,
	})
}

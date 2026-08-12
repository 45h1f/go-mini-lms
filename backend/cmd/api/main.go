package main

import (
	"log"

	"mini-lms/internal/domain"
	handler "mini-lms/internal/handler/http"
	"mini-lms/internal/middleware"
	"mini-lms/internal/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Printf("Warning: DB init failed (ensure Postgres is running): %v", err)
	}

	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK", "service": "mini-lms-api"})
	})

	authHandler := handler.NewAuthHandler(db)
	courseHandler := handler.NewCourseHandler(db)
	lessonHandler := handler.NewLessonHandler(db)
	enrollmentHandler := handler.NewEnrollmentHandler(db)
	progressHandler := handler.NewProgressHandler(db)

	v1 := r.Group("/api/v1")
	{
		// Public Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Public Course routes
		v1.GET("/courses", courseHandler.ListCourses)

		// Protected Student/Instructor/Admin routes
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// Student Enrollment & Progress Tracking
			protected.POST("/courses/:id/enroll", enrollmentHandler.EnrollCourse)
			protected.GET("/my-enrollments", enrollmentHandler.GetMyEnrollments)
			protected.POST("/progress", progressHandler.MarkLessonProgress)
			protected.GET("/courses/:id/progress", progressHandler.GetCourseProgress)

			// Instructor/Admin Lessons & Courses
			protected.POST("/courses", middleware.RequireRole(domain.RoleInstructor, domain.RoleAdmin), courseHandler.CreateCourse)
			protected.POST("/courses/:id/lessons", middleware.RequireRole(domain.RoleInstructor, domain.RoleAdmin), lessonHandler.CreateLesson)
		}
	}

	log.Println("Enterprise Mini LMS API running on :8080...")
	r.Run(":8080")
}

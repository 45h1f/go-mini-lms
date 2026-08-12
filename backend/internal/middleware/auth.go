package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"mini-lms/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("SUPER_SECRET_ENTERPRISE_KEY_CHANGE_IN_PROD")

type Claims struct {
	UserID uint        `json:"user_id"`
	Role   domain.Role `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(user *domain.User) (string, error) {
	claims := Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(parts[1], claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		c.Next()
	}
}

func RequireRole(roles ...domain.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "User role not found in context"})
			c.Abort()
			return
		}

		userRole := roleVal.(domain.Role)
		for _, r := range roles {
			if userRole == r {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied for this role"})
		c.Abort()
	}
}

package utils

import (
	"Blog/models"
	"Blog/storage"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func LogAudit(c *gin.Context, action, object string, objectID uint, metadata string) {

	userID := c.GetUint("user_id")
	ip := c.ClientIP()
	ua := c.GetHeader("User-Agent")

	logEntry := models.AuditLog{
		UserID:    userID,
		Action:    action,
		Object:    object,
		ObjectID:  objectID,
		Timestamp: time.Now(),
		IP:        ip,
		UserAgent: ua,
		Metadata:  metadata,
	}

	if err := storage.DB.Create(&logEntry).Error; err != nil {
		log.Println("Ошибка при записи в аудит лог:", err)
	}
}

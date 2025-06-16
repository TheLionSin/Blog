package utils

import (
	"Blog/models"
	"Blog/storage"
	"log"
	"time"
)

func LogAudit(userID uint, action, object string, objectID uint) {
	logEntry := models.AuditLog{
		UserID:    userID,
		Action:    action,
		Object:    object,
		ObjectID:  objectID,
		Timestamp: time.Now(),
	}

	if err := storage.DB.Create(&logEntry).Error; err != nil {
		log.Println("Ошибка при записи в аудит лог:", err)
	}
}

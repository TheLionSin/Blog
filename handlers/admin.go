package handlers

import (
	"Blog/dto"
	"Blog/models"
	"Blog/storage"
	"Blog/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func GetAuditLogs(c *gin.Context) {
	var logs []models.AuditLog
	if err := storage.DB.Order("timestamp desc").Limit(100).Find(&logs).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Не удалось загрузить журнал")
		return
	}

	logResponses := dto.ToAuditLogList(logs)

	utils.RespondOK(c, gin.H{
		"logs": logResponses,
	})
}

func RestoreUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Некорректный ID")
		return
	}

	var user models.User
	if err := storage.DB.Unscoped().First(&user, userID).Error; err != nil {
		utils.RespondError(c, http.StatusNotFound, "Пользователь не найден")
		return
	}

	if user.DeletedAt.Valid == false {
		utils.RespondError(c, http.StatusBadRequest, "Пользователь уже активен")
		return
	}

	user.DeletedAt = gorm.DeletedAt{}
	if err := storage.DB.Unscoped().Save(&user).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Не удалось восстановить пользователя")
		return
	}

	utils.LogAudit(c, "undelete_user", "user", user.ID, "")

	utils.RespondOK(c, gin.H{
		"message": "Пользователь восстановлен",
	})
}

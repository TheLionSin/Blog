package handlers

import (
	"Blog/dto"
	"Blog/models"
	"Blog/storage"
	"Blog/utils"
	"github.com/gin-gonic/gin"
	"net/http"
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

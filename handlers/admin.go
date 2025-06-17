package handlers

import (
	"Blog/dto"
	"Blog/models"
	"Blog/storage"
	"Blog/utils"
	"encoding/csv"
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

func ExportUsersCSV(c *gin.Context) {
	var users []models.User
	query := storage.DB.Model(&models.User{})
	if search := c.Query("search"); search != "" {
		term := "%" + search + "%"
		query = query.Where("nickname ILIKE ? OR email ILIKE ?", term, term)
	}
	if role := c.Query("role"); role != "" {
		query = query.Where("role = ?", role)
	}
	if email := c.Query("email"); email != "" {
		query = query.Where("email ILIKE ?", "%"+email+"%")
	}
	if nickname := c.Query("nickname"); nickname != "" {
		query = query.Where("nickname ILIKE ?", "%"+nickname+"%")
	}

	if err := query.Find(&users).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Не удалось получить пользователя")
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", `attachment; filename="users_export.csv"`)

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	writer.Write([]string{"ID", "Email", "Nickname", "Role", "Created At"})

	for _, u := range users {
		record := []string{
			strconv.FormatUint(uint64(u.ID), 10),
			u.Email,
			u.Nickname,
			u.Role,
			u.RegisteredAt.Format("2006-01-02 15:04:05"),
		}
		writer.Write(record)
	}
}

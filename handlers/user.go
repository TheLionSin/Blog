package handlers

import (
	"Blog/dto"
	"Blog/models"
	"Blog/storage"
	"Blog/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var validate = validator.New()

func GetCurrentUser(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "Пользователь не найден")
		return
	}

	var user models.User

	if err := storage.DB.First(&user, userID).Error; err != nil {
		utils.RespondError(c, http.StatusNotFound, "Пользователь не найден")
		return
	}

	utils.RespondOK(c, gin.H{
		"user": dto.ToUserResponse(user),
	})
}

func CreateUser(c *gin.Context) {

	var input dto.CreateUserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		errors := utils.FormatValidationError(err)
		utils.RespondError(c, http.StatusBadRequest, errors)
		return
	}

	var existing models.User
	if err := storage.DB.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		utils.RespondError(c, http.StatusBadRequest, "Email уже зарегистрирован")
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Не удалось зашифровать пароль")
		return
	}

	user := models.User{
		Nickname: input.Nickname,
		Email:    input.Email,
		Password: hashedPassword,
	}

	if err := storage.DB.Create(&user).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Ошибка при создании пользователя")
		return
	}

	utils.RespondOK(c, gin.H{
		"user": dto.ToUserResponse(user),
	})

}

func UpdateUser(c *gin.Context) {

	idParam := c.Param("id")
	targetID, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Некорректный ID")
		return
	}

	var user models.User
	if err := storage.DB.First(&user, targetID).Error; err != nil {
		utils.RespondError(c, http.StatusNotFound, "Пользователь не найден")
		return
	}

	var input dto.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		errors := utils.FormatValidationError(err)
		utils.RespondError(c, http.StatusBadRequest, errors)
		return
	}

	if input.Nickname != "" {
		user.Nickname = input.Nickname
	}
	if input.Email != "" {
		user.Email = input.Email
	}
	if input.Password != "" {
		hashedPassword, err := utils.HashPassword(input.Password)
		if err != nil {
			utils.RespondError(c, http.StatusInternalServerError, "Не удалось зашифровать пароль")
			return
		}
		user.Password = hashedPassword
	}

	if err := storage.DB.Save(&user).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Ошибка при обновлении пользователя")
		return
	}

	utils.RespondOK(c, gin.H{
		"user": dto.ToUserResponse(user),
	})

	inputForLog := input
	inputForLog.Password = ""
	metadata := fmt.Sprintf("input: %+v", inputForLog)
	utils.LogAudit(c, "update_user", "user", user.ID, metadata)

}

func DeleteUser(c *gin.Context) {

	idParam := c.Param("id")
	targetID, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Некорректный ID")
		return
	}

	var user models.User
	if err := storage.DB.First(&user, targetID).Error; err != nil {
		utils.RespondError(c, http.StatusNotFound, "Пользователь не найден")
		return
	}

	if err := storage.DB.Delete(&user).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Ошибка при удалении")
		return
	}

	utils.RespondOK(c, gin.H{
		"message": "Пользователь удален",
	})

	utils.LogAudit(c, "delete_user", "user", user.ID, "")

}

func GetUsers(c *gin.Context) {

	var users []models.User

	limitStr := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "1")
	search := c.Query("search")

	limit, _ := strconv.Atoi(limitStr)
	page, _ := strconv.Atoi(pageStr)
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	query := storage.DB.Limit(limit).Offset(offset)
	if search != "" {
		term := "%" + search + "%"
		query = query.Where("name ILIKE ? or email ILIKE ?", term, term)
	}

	if err := query.Find(&users).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Ошибка при получении пользователей")
		return
	}

	var response []dto.UserResponse
	for _, u := range users {
		response = append(response, dto.ToUserResponse(u))
	}

	utils.RespondOK(c, gin.H{
		"users": response,
	})

}

func GetUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Некорректный ID")
		return
	}

	var user models.User
	if err := storage.DB.First(&user, userID).Error; err != nil {
		utils.RespondError(c, http.StatusNotFound, "Пользователь не найден")
		return
	}

	utils.RespondOK(c, gin.H{
		"user": dto.ToUserResponse(user),
	})
}

func UploadAvatar(c *gin.Context) {
	userID := c.GetUint("user_id")

	file, err := c.FormFile("avatar")
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Файл не найден")
		return
	}

	var user models.User
	if err := storage.DB.First(&user, userID).Error; err != nil {
		utils.RespondError(c, http.StatusNotFound, "Пользователь не найден")
		return
	}

	if user.AvatarURL != "" {
		oldPath := strings.TrimPrefix(user.AvatarURL, "/")
		if err := os.Remove(oldPath); err != nil && !os.IsNotExist(err) {
			fmt.Println("Ошибка при удалении старого аватара:", err)
		}
	}

	filename := fmt.Sprintf("avatar_%d%s", userID, filepath.Ext(file.Filename))
	path := filepath.Join("uploads", filename)

	if err := c.SaveUploadedFile(file, path); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Не удалось сохранить файл")
		return
	}

	avatarURL := "/uploads/" + filename
	if err := storage.DB.Model(&models.User{}).Where("id = ?", userID).Update("avatar_url", avatarURL).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Не удалось обновить профиль")
		return
	}

	utils.RespondOK(c, gin.H{
		"avatar_url": avatarURL,
	})

	utils.LogAudit(c, "upload_avatar", "user", user.ID, file.Filename)

}

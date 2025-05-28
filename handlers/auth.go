package handlers

import (
	"Blog/dto"
	"Blog/models"
	"Blog/storage"
	"Blog/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Register(c *gin.Context) {
	var input dto.RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Неверный JSON")
		return
	}

	if err := validate.Struct(input); err != nil {
		errors := utils.FormatValidationError(err)
		utils.RespondError(c, http.StatusBadRequest, errors)
		return
	}

	var existing models.User
	if err := storage.DB.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		utils.RespondError(c, http.StatusConflict, "Пользователь с таким email уже существует")
		return
	}

	hashed, err := utils.HashPassword(input.Password)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Ошибка хеширования пароля")
		return
	}

	user := models.User{
		Nickname: input.Nickname,
		Email:    input.Email,
		Password: hashed,
	}

	if err := storage.DB.Create(&user).Error; err != nil {
		// Проверка на дубликат по email
		if strings.Contains(err.Error(), "duplicate key value") &&
			strings.Contains(err.Error(), "users_email_key") {
			utils.RespondError(c, http.StatusConflict, "Email уже используется")
			return
		}

		utils.RespondError(c, http.StatusInternalServerError, "Ошибка при создании пользователя")
		return
	}

	accessToken, err := utils.GenerateJWT(user.ID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Ошибка генерации токена")
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Ошибка генерации токена")
		return
	}

	c.SetCookie("refresh_token", refreshToken, 7*24*60*60, "/", "", true, true)

	utils.RespondCreated(c, gin.H{
		"user":         dto.ToUserResponse(user),
		"access_token": accessToken,
	})
}

func Login(c *gin.Context) {
	var input dto.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadGateway, "Неверный JSON")
		return
	}

	if err := validate.Struct(input); err != nil {
		errors := utils.FormatValidationError(err)
		utils.RespondError(c, http.StatusBadRequest, errors)
		return
	}

	var user models.User

	if err := storage.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		utils.RespondError(c, http.StatusUnauthorized, "Пользователь не найден")
		return
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		utils.RespondError(c, http.StatusUnauthorized, "Неверный пароль")
		return
	}

	accessToken, err := utils.GenerateJWT(user.ID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Ошибка генерации токена")
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Ошибка генерации токена")
		return
	}

	c.SetCookie("refresh_token", refreshToken, 7*24*60*60, "/", "", true, true)

	user.Password = ""

	utils.RespondOK(c, gin.H{
		"user":         dto.ToUserResponse(user),
		"access_token": accessToken,
	})
}

func RefreshToken(c *gin.Context) {

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, "Нет refresh токена")
		return
	}

	userID, err := utils.ParseRefreshToken(refreshToken)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, err.Error())
		return
	}

	newAccessToken, err := utils.GenerateJWT(userID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Не удалось выдать новый токен")
		return
	}

	utils.RespondOK(c, gin.H{
		"access_token": newAccessToken,
	})
}

func Logout(c *gin.Context) {
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)

	utils.RespondOK(c, gin.H{
		"message": "Вы вышли из системы",
	})
}

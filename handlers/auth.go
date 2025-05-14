package handlers

import (
	"Blog/models"
	"Blog/storage"
	"Blog/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func Register(c *gin.Context) {
	var input models.User

	if err := c.BindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Неверный JSON")
		return
	}

	if err := validate.Struct(input); err != nil {
		errors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			errors[e.Field()] = fmt.Sprintf("Не проходит '%s'", e.Tag())
		}
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

	input.Password = hashed

	if err := storage.DB.Create(&input).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Ошибка при создании")
		return
	}

	token, err := utils.GenerateJWT(input.ID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Ошибка генерации токена")
		return
	}

	utils.RespondCreated(c, gin.H{
		"user":  input,
		"token": token,
	})
}

func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadGateway, "Неверный JSON")
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

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Ошибка генерации токена")
		return
	}

	user.Password = ""

	utils.RespondOK(c, gin.H{
		"user":  user,
		"token": token,
	})
}

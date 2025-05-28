package handlers

import (
	"Blog/dto"
	"Blog/models"
	"Blog/storage"
	"Blog/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
	"strconv"
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
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid JSON")
		return
	}
	if err := validate.Struct(user); err != nil {
		errors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			errors[e.Field()] = fmt.Sprintf("Ошибка поля '%s'", e.Tag())
		}
		utils.RespondError(c, http.StatusBadRequest, errors)
		return
	}
	if err := storage.DB.Create(&user).Error; err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "uni_users_nickname":
				utils.RespondError(c, http.StatusConflict, fmt.Sprintf("Никнейм '%s' уже используется", user.Nickname))
				return
			case "uni_users_email":
				utils.RespondError(c, http.StatusConflict, fmt.Sprintf("Email '%s' уже используется", user.Email))
				return
			default:
				utils.RespondError(c, http.StatusConflict, "Пользователь с таким email или никнеймом уже существует")
				return
			}
		}
		utils.RespondError(c, http.StatusBadRequest, "Ошибка при сохранении пользователя")
		return
	}
	utils.RespondCreated(c, user)
}

func UpdateUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	if err := storage.DB.First(&user, id).Error; err != nil {
		utils.RespondError(c, http.StatusNotFound, "Пользователь не найден")
		return
	}

	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Неверный JSON")
		return
	}
	if err := validate.Struct(input); err != nil {
		errors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			errors[e.Field()] = fmt.Sprintf("Ошибка поля '%s'", e.Tag())
		}
		utils.RespondError(c, http.StatusBadRequest, errors)
		return
	}

	user.Nickname = input.Nickname
	user.Email = input.Email

	if err := storage.DB.Save(&user).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Ошибка при сохранении пользователя")
		return
	}
	utils.RespondOK(c, gin.H{"Пользователь обновлен": user})

}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := storage.DB.First(&user, id).Error; err != nil {
		utils.RespondError(c, http.StatusNotFound, "Пользователь не существует")
		return
	}
	if err := storage.DB.Delete(&user).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Ошибка при удалении")
	}
	utils.RespondOK(c, "Пользователь удален")
}

func GetUsers(c *gin.Context) {

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	var users []models.User
	var total int64

	storage.DB.Model(&models.User{}).Count(&total)

	err := storage.DB.Limit(limit).Offset(offset).Find(&users).Error

	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Ошибка при получении пользователей")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"total":   total,
		"page":    page,
		"limit":   limit,
		"data":    users,
	})

}

func GetUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	if err := storage.DB.First(&user, id).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Не удалось получить пользователя")
		return
	}
	utils.RespondOK(c, user)
}

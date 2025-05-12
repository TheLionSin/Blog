package handlers

import (
	"Blog/models"
	"Blog/storage"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
)

var validate = validator.New()

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка JSON"})
		return
	}
	if err := validate.Struct(user); err != nil {
		errors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			errors[e.Field()] = fmt.Sprintf("Ошибка поля '%s'", e.Tag())
		}
		c.JSON(http.StatusBadRequest, gin.H{"validation_errors": errors})
		return
	}
	if err := storage.DB.Create(&user).Error; err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "uni_users_nickname":
				c.JSON(http.StatusConflict, gin.H{"error": fmt.Sprintf("Никнейм '%s' уже используется", user.Nickname)})
				return
			case "uni_users_email":
				c.JSON(http.StatusConflict, gin.H{"error": fmt.Sprintf("Email '%s' уже используется", user.Email)})
				return
			default:
				c.JSON(http.StatusConflict, gin.H{"error": "Пользователь с таким email или никнеймом уже существует"})
				return
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при сохранении пользователя"})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func UpdateUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	if err := storage.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON"})
		return
	}
	if err := validate.Struct(input); err != nil {
		errors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			errors[e.Field()] = fmt.Sprintf("Ошибка поля '%s'", e.Tag())
		}
		c.JSON(http.StatusBadRequest, gin.H{"validation_errors": errors})
		return
	}

	user.Nickname = input.Nickname
	user.Email = input.Email

	if err := storage.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении пользователя"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Пользователь обновлен:": user})

}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := storage.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не существует"})
		return
	}
	if err := storage.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении"})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Пользователь удален"})
}

func GetUsers(c *gin.Context) {
	var users []models.User
	if err := storage.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить пользователей"})
		return
	} else if len(users) == 0 {
		c.JSON(http.StatusOK, gin.H{"status": "В данный момент нету записей."})
	}

	c.JSON(http.StatusOK, users)

}

func GetUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	if err := storage.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить пользователя"})
		return
	}
	c.JSON(http.StatusOK, user)
}

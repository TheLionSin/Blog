package handlers

import (
	"Blog/dto"
	"Blog/models"
	"Blog/storage"
	"Blog/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Некорректный ID")
		return
	}

	var input dto.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		errors := utils.FormatValidationError(err)
		utils.RespondError(c, http.StatusBadRequest, errors)
		return
	}

	var user models.User
	if err := storage.DB.First(&user, userID).Error; err != nil {
		utils.RespondError(c, http.StatusNotFound, "Пользователь не найден")
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

}

func DeleteUser(c *gin.Context) {

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

	if err := storage.DB.Delete(&user).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Ошибка при удалении")
		return
	}

	utils.RespondOK(c, gin.H{
		"message": "Пользователь удален",
	})

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

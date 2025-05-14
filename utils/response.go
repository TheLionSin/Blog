package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RespondOK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

func RespondCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    data,
	})
}

func RespondError(c *gin.Context, status int, err interface{}) {
	c.JSON(status, gin.H{
		"success": false,
		"error":   err,
	})
}

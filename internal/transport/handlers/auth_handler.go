package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type DummyLoginRequest struct {
	Role string `json:"role" binding:"required,oneof=moderator employee"`
}

func DummyLogin(c *gin.Context) {
	var req DummyLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": req.Role})
}

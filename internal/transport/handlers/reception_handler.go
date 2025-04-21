package handlers

import (
	"avito-internship/internal/database"
	"avito-internship/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ReceptionRequest struct {
	PVZID string `json:"pvzId" binding:"required"`
}

func CreateReception(c *gin.Context) {
	role := c.GetString("role")
	if role != "employee" {
		c.JSON(http.StatusForbidden, gin.H{"message": "Only employees can create receptions"})
		return
	}

	var req ReceptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	reception, err := services.CreateReception(database.DB, req.PVZID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, reception)
}

func CloseReception(c *gin.Context) {
	role := c.GetString("role")
	if role != "employee" {
		c.JSON(http.StatusForbidden, gin.H{"message": "Only employees can close receptions"})
		return
	}

	pvzID := c.Param("pvzId")
	err := services.CloseLastReception(database.DB, pvzID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "reception has been closed"})
}

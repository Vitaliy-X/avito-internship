package handlers

import (
	"avito-internship/internal/database"
	"avito-internship/internal/models"
	"avito-internship/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type CreatePVZRequest struct {
	ID               string    `json:"id" binding:"required,uuid"`
	RegistrationDate time.Time `json:"registrationDate" binding:"required"`
	City             string    `json:"city" binding:"required"`
}

func CreatePVZ(c *gin.Context) {
	role := c.GetString("role")
	if role != "moderator" {
		c.JSON(http.StatusForbidden, gin.H{"message": "Only moderators can create PVZ"})
		return
	}

	var req CreatePVZRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	pvz := models.PVZ{
		ID:               req.ID,
		City:             req.City,
		RegistrationDate: req.RegistrationDate,
	}

	result, err := services.CreatePVZ(database.DB, pvz)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func GetPVZList(c *gin.Context) {
	role := c.GetString("role")
	if role != "employee" && role != "moderator" {
		c.JSON(http.StatusForbidden, gin.H{"message": "Access denied"})
		return
	}

	layout := time.RFC3339
	var startDate, endDate *time.Time

	if s := c.Query("startDate"); s != "" {
		t, err := time.Parse(layout, s)
		if err == nil {
			startDate = &t
		}
	}

	if e := c.Query("endDate"); e != "" {
		t, err := time.Parse(layout, e)
		if err == nil {
			endDate = &t
		}
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit > 30 {
		limit = 30
	}

	result, err := services.GetPVZList(database.DB, startDate, endDate, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

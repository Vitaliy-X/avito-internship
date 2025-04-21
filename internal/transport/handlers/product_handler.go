package handlers

import (
	"avito-internship/internal/database"
	"avito-internship/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AddProductRequest struct {
	Type  string `json:"type" binding:"required,oneof=электроника одежда обувь"`
	PVZID string `json:"pvzId" binding:"required"`
}

func AddProduct(c *gin.Context) {
	role := c.GetString("role")
	if role != "employee" {
		c.JSON(http.StatusForbidden, gin.H{"message": "Only employees can add products"})
		return
	}

	var req AddProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	product, err := services.AddProduct(database.DB, req.PVZID, req.Type)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func DeleteLastProduct(c *gin.Context) {
	role := c.GetString("role")
	if role != "employee" {
		c.JSON(http.StatusForbidden, gin.H{"message": "Only employees can delete products"})
		return
	}

	pvzID := c.Param("pvzId")
	err := services.DeleteLastProduct(database.DB, pvzID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Last product deleted successfully"})
}

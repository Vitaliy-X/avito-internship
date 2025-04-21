package transport

import (
	"avito-internship/internal/transport/handlers"
	"avito-internship/internal/transport/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Публичная ручка
	r.POST("/dummyLogin", handlers.DummyLogin)

	// Защищённые ручки
	auth := r.Group("/", middleware.AuthMiddleware())
	{
		auth.POST("/pvz", handlers.CreatePVZ)
		auth.GET("/pvz", handlers.GetPVZList)

		auth.POST("/receptions", handlers.CreateReception)
		auth.POST("/pvz/:pvzId/close_last_reception", handlers.CloseReception)

		auth.POST("/products", handlers.AddProduct)
		auth.POST("/pvz/:pvzId/delete_last_product", handlers.DeleteLastProduct)
	}

	return r
}

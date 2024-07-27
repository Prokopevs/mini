package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "github.com/Prokopevs/mini/server/docs"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (h *HTTP) setRoutes(r *gin.Engine) {
	r.Use(CORSMiddleware())

	api := r.Group("/api/v1")
	{
		api.POST("/create", h.CreateMessage)
		api.GET("/getMessages", h.GetMessages)
		api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}

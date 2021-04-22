package handler

import (
	"net/http"

	_ "github.com/gakkiyomi/compare-excel-column/docs"
	"github.com/gin-gonic/gin"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// SetupRouter create gin router and return
func SetupRouter() *gin.Engine {
	r := gin.Default()
	//cors
	r.Use(Cors())
	r.GET("/api/compare-excel-column/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/api/compare", Compare)
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong!")
	})
	return r
}

// Cors deal with cors
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "false")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

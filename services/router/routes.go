package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func genericCrudRoutes(engine *gin.Engine) {
	engine.POST("/:table/",)
	engine.GET("/:table/:id",)
	engine.PUT("/:table/:id",)
	engine.DELETE("/:table/:id",)
}

func RegisterRoutes(engine *gin.Engine) {
	engine.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	genericCrudRoutes(engine)
}
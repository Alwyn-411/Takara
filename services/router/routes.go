package router

import (
	"takara/services/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RegisterRoutes(engine *gin.Engine, db *sqlx.DB) {
	// User Routes
	engine.POST("/v1/user/", func(ctx *gin.Context) { handlers.CreateUser(ctx, db) })
	engine.GET("/v1/user/:id", func(ctx *gin.Context) { handlers.GetUserById(ctx, db) })
	engine.PUT("/v1/user/:id", func(ctx *gin.Context) { handlers.UpdateUserById(ctx, db) })
	engine.DELETE("/v1/user/:id", func(ctx *gin.Context) { handlers.DeleteUserById(ctx, db) })

	// Account Routes
	engine.POST("/v1/account/:userid")
	engine.GET("/v1/account/:account_id")
	engine.PUT("/v1/account/:account_id")
	engine.DELETE("/v1/account/:account_id")

	// Trasactions Routes
	engine.POST("/v1/transactions/:account_id")
	engine.GET("/v1/transactions/:transaction_id")
	engine.PUT("/v1/transactions/:transaction_id")
	engine.DELETE("/v1/transactions/:transaction_id")
}

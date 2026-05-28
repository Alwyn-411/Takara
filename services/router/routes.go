package router

import (
	"net/http"
	"takara/services/forex"
	"takara/services/handlers"
	"takara/services/handlers/transactions"
	"takara/services/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RegisterRoutes(engine *gin.Engine, db *sqlx.DB) {
	forexSvc := forex.NewAccessor()
	tokenSvc := middleware.NewTokenService("")

	authHandler := handlers.NewAuthHandler(db, tokenSvc)
	userHandler := handlers.NewUserHandler(db)
	tagsHandler := handlers.NewTagsHandler(db)
	categoryHandler := handlers.NewCategoryHandler(db)
	merchantHandler := handlers.NewMerchantHandler(db)
	accountHandler := handlers.NewAccountHandler(db)
	transactionsHandler := transactions.NewTransactionsHandler(db, forexSvc)

	engine.GET("/ping", func(ctx *gin.Context) { ctx.JSON(http.StatusAccepted, gin.H{"message": "pong"}) })
	engine.POST("/v1/auth", authHandler.Login)
	engine.POST("/v1/user/", userHandler.CreateUser)

	api := engine.Group("/v1")
	api.Use(tokenSvc.RequireAuth())
	{
		// User
		api.GET("/user/:id", userHandler.GetUserById)
		api.PUT("/user/:id", userHandler.UpdateUserById)
		api.DELETE("/user/:id", userHandler.DeleteUserById)

		// Account
		api.POST("/account/", accountHandler.CreateAccount)
		api.GET("/account/:accountId", accountHandler.GetAccountById)
		api.PUT("/account/:accountId", accountHandler.UpdateAccountById)
		api.DELETE("/account/:accountId", accountHandler.DeleteAccountById)
		api.GET("/account/list", accountHandler.ListAccounts)

		// Category
		api.GET("/category/list", categoryHandler.ListCategories)

		// Tags
		api.GET("/tag/list", tagsHandler.ListTags)

		// Merchants
		api.GET("/merchants/list", merchantHandler.ListMerchants)

		// Transactions
		api.POST("/transaction/", transactionsHandler.CreateTransaction)
		api.GET("/transaction/:transactionId", transactionsHandler.GetTransactionById)
		api.PUT("/transaction/:transactionId", transactionsHandler.UpdateTransaction)
		api.DELETE("/transaction/:transactionId", transactionsHandler.DeleteTransaction)
		api.GET("/transaction/account/:accountId", transactionsHandler.GetTransactionsByAccountId)

		// Transaction X Tags Routes
		// api.POST("/transaction/:transactionId/tags", tagsHandler.AddTagsToTransaction)
		// api.DELETE("/transaction/:transactionId/tags/:tagId", tagsHandler.RemoveTagFromTransaction)
		// api.GET("/transaction/:transactionId/tags", tagsHandler.GetTagsForTransaction)
	}
}

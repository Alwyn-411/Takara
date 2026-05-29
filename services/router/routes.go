package router

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"takara/services/forex"
	"takara/services/handlers"
	"takara/services/handlers/transactions"
	"takara/services/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func Env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func RegisterRoutes(engine *gin.Engine, db *sqlx.DB) {
	forexSvc := forex.NewAccessor()
	tokenSvc := middleware.NewTokenService(Env("TOKEN_SECRET", ""))
	authHandler := handlers.NewAuthHandler(db, tokenSvc)
	userHandler := handlers.NewUserHandler(db)
	tagsHandler := handlers.NewTagsHandler(db)
	categoryHandler := handlers.NewCategoryHandler(db)
	merchantHandler := handlers.NewMerchantHandler(db)
	accountHandler := handlers.NewAccountHandler(db)
	transactionsHandler := transactions.NewTransactionsHandler(db, forexSvc)

	// --- API routes (unchanged) ---
	engine.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusAccepted, gin.H{"message": "pong"})
	})
	engine.POST("/v1/auth", authHandler.Login)
	engine.POST("/v1/user/", userHandler.CreateUser)

	api := engine.Group("/v1")
	api.Use(tokenSvc.RequireAuth())
	{
		api.GET("/user/:id", userHandler.GetUserById)
		api.PUT("/user/:id", userHandler.UpdateUserById)
		api.DELETE("/user/:id", userHandler.DeleteUserById)

		api.POST("/account/", accountHandler.CreateAccount)
		api.GET("/account/:accountId", accountHandler.GetAccountById)
		api.PUT("/account/:accountId", accountHandler.UpdateAccountById)
		api.DELETE("/account/:accountId", accountHandler.DeleteAccountById)
		api.GET("/account/list", accountHandler.ListAccounts)

		api.GET("/category/list", categoryHandler.ListCategories)
		api.GET("/tag/list", tagsHandler.ListTags)
		api.GET("/merchants/list", merchantHandler.ListMerchants)

		api.POST("/transaction/", transactionsHandler.CreateTransaction)
		api.GET("/transaction/:transactionId", transactionsHandler.GetTransactionById)
		api.PUT("/transaction/:transactionId", transactionsHandler.UpdateTransaction)
		api.DELETE("/transaction/:transactionId", transactionsHandler.DeleteTransaction)
		api.GET("/transaction/account/:accountId", transactionsHandler.GetTransactionsByAccountId)
	}

	staticDir := os.Getenv("STATIC_DIR")
	if staticDir != "" {
		engine.Static("/assets", filepath.Join(staticDir, "assets"))
		engine.StaticFile("/favicon.ico", filepath.Join(staticDir, "favicon.ico"))
		engine.StaticFile("/vite.svg", filepath.Join(staticDir, "vite.svg"))

		engine.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			if strings.HasPrefix(path, "/v1") || path == "/ping" {
				c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
				return
			}
			c.File(filepath.Join(staticDir, "index.html"))
		})
	}
}

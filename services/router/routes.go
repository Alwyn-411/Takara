package router

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"takara/services/forex"
	"takara/services/handlers"
	"takara/services/handlers/holdings"
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

func StaticAssetsRoutes(engine *gin.Engine) {
	staticDir := Env("STATIC_DIR", "")

	if staticDir != "" {
		engine.Static("/assets", filepath.Join(staticDir, "assets"))
		engine.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			if strings.HasPrefix(path, "/v1") || path == "/ping" {
				c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
				return
			}

			requested := filepath.Join(staticDir, path)
			if info, err := os.Stat(requested); err == nil && !info.IsDir() {
				c.File(requested)
				return
			}

			c.File(filepath.Join(staticDir, "index.html"))
		})
	}
}

func RegisterRoutes(engine *gin.Engine, db *sqlx.DB) {
	forexSvc := forex.NewAccessor()
	tokenSvc := middleware.NewTokenService(Env("TOKEN_SECRET", ""))
	authHandler := handlers.NewAuthHandler(db, tokenSvc)
	userHandler := handlers.NewUserHandler(db)
	userPrefHandler := handlers.NewUserPrefHandler(db)
	tagsHandler := handlers.NewTagsHandler(db)
	categoryHandler := handlers.NewCategoryHandler(db)
	merchantHandler := handlers.NewMerchantHandler(db)
	accountHandler := handlers.NewAccountHandler(db)
	transactionsHandler := transactions.NewTransactionsHandler(db, forexSvc)
	holdingsHandler := holdings.NewHoldingsHandler(db)

	engine.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusAccepted, gin.H{"message": "pong"})
	})

	api := engine.Group("/v1")

	api.POST("/auth", authHandler.Login)
	api.POST("/user", userHandler.CreateUser)
	api.POST("/user/pref", userPrefHandler.CreatePref)
	api.GET("/user/avatar/:id", userHandler.GetAvatar)

	api.Use(tokenSvc.RequireAuth())
	{
		api.GET("/user/:id", userHandler.GetUserById)
		api.PUT("/user", userHandler.UpdateUserById)
		api.DELETE("/user", userHandler.DeleteUserById)
		api.PUT("/user/avatar", userHandler.UpdateAvatar)

		api.GET("/user/pref", userPrefHandler.GetPref)
		api.PUT("/user/pref", userPrefHandler.UpdateUserPref)

		api.POST("/account", accountHandler.CreateAccount)
		api.GET("/account/:accountId", accountHandler.GetAccountById)
		api.PUT("/account/:accountId", accountHandler.UpdateAccountById)
		api.DELETE("/account/:accountId", accountHandler.DeleteAccountById)
		api.GET("/account/list", accountHandler.ListAccounts)

		api.GET("/category/list", categoryHandler.ListCategories)
		api.GET("/tag/list", tagsHandler.ListTags)
		api.GET("/merchants/list", merchantHandler.ListMerchants)

		api.POST("/transaction", transactionsHandler.CreateTransaction)
		api.GET("/transaction/:transactionId", transactionsHandler.GetTransactionById)
		api.PUT("/transaction/:transactionId", transactionsHandler.UpdateTransaction)
		api.DELETE("/transaction/:transactionId", transactionsHandler.DeleteTransaction)
		api.GET("/transaction/account/:accountId", transactionsHandler.GetTransactionsByAccountId)

		api.POST("/holdings", holdingsHandler.CreateHolding)
		api.GET("/holdings", holdingsHandler.ListHoldings)
		api.GET("/holdings/:id", holdingsHandler.GetHoldingById)
		api.PUT("/holdings/:id", holdingsHandler.UpdateHoldingById)
		api.DELETE("/holdings/:id", holdingsHandler.DeleteHoldingById)
		api.POST("/holdings/:id/valuations", holdingsHandler.CreateValuation)
		api.GET("/holdings/:id/valuations", holdingsHandler.ListValuations)
		api.DELETE("/valuations/:id", holdingsHandler.DeleteValuation)
	}

	StaticAssetsRoutes(engine)
}

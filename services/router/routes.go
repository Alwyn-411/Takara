package router

import (
	"net/http"
	"takara/services/forex"
	"takara/services/handlers"
	"takara/services/handlers/transactions"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RegisterRoutes(engine *gin.Engine, db *sqlx.DB) {
	forexSvc := forex.NewAccessor()

	userHandler := handlers.NewUserHandler(db)
	authHandler := handlers.NewAuthHandler(db)
	accountHandler := handlers.NewAccountHandler(db)
	transactionsHandler := transactions.NewTransactionsHandler(db, forexSvc)

	// v1 routes
	engine.GET("/ping", func(ctx *gin.Context) { ctx.JSON(http.StatusAccepted, gin.H{"message": "pong"}) })
	engine.GET("/v1/auth", func(ctx *gin.Context) { authHandler.AuthorizeUserWithUserNameAndPassword(ctx) })

	// User Routes
	engine.POST("/v1/user/", func(ctx *gin.Context) { userHandler.CreateUser(ctx) })
	engine.GET("/v1/user/:id", func(ctx *gin.Context) { userHandler.GetUserById(ctx) })
	engine.PUT("/v1/user/:id", func(ctx *gin.Context) { userHandler.UpdateUserById(ctx) })
	engine.DELETE("/v1/user/:id", func(ctx *gin.Context) { userHandler.DeleteUserById(ctx) })

	// Account Routes
	engine.POST("/v1/account/", func(ctx *gin.Context) { accountHandler.CreateAccount(ctx) })
	engine.GET("/v1/account/:accountId", func(ctx *gin.Context) { accountHandler.GetAccountById(ctx) })
	engine.PUT("/v1/account/:accountId", func(ctx *gin.Context) { accountHandler.UpdateAccountById(ctx) })
	engine.DELETE("/v1/account/:accountId", func(ctx *gin.Context) { accountHandler.DeleteAccountById(ctx) })
	engine.GET("/v1/account/user/:userId", func(ctx *gin.Context) { accountHandler.GetAccountsByUserId(ctx) })

	// Trasactions Routes
	engine.POST("/v1/transaction/", func(ctx *gin.Context) { transactionsHandler.CreateTransaction(ctx) })
	engine.GET("/v1/transaction/:transactionId", func(ctx *gin.Context) { transactionsHandler.GetTransactionById(ctx) })
	engine.PUT("/v1/transaction/:transactionId", func(ctx *gin.Context) { transactionsHandler.UpdateTransaction(ctx) })
	engine.DELETE("/v1/transaction/:transactionId", func(ctx *gin.Context) { transactionsHandler.DeleteTransaction(ctx) })
	engine.GET("/v1/transaction/account/:accountId", func(ctx *gin.Context) { transactionsHandler.GetTransactionsByAccountId(ctx) })
}

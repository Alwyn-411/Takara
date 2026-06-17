package main

import (
	"log"
	"os"
	"takara/services/middleware"
	"takara/services/router"
	"takara/services/schema"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

var db *sqlx.DB

func initDatabase() (*sqlx.DB, error) {
	var err error
	dbPath := router.Env("DB_PATH", "takara.db")
	db, err = sqlx.Connect("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	initTables(db)
	return db, nil
}

func initTables(db *sqlx.DB) {
	schema.ForeignKeysEnabled(db)
	schema.InitUsers(db)
	schema.InitAccounts(db)
	schema.InitCategories(db)
	schema.InitMerchant(db)
	schema.InitTags(db)
	schema.InitTransactions(db)
	schema.InitTransactionTags(db)
	schema.InitIndexUsers(db)
	schema.InitIndexTransactions(db)
}

func main() {
	_ = godotenv.Load("../.env")

	if os.Getenv("TOKEN_SECRET") == "" {
		log.Fatal("TOKEN_SECRET is required")
	}

	if os.Getenv("FOREX_ENDPOINT") == "" {
		log.Fatal("FOREX_ENDPOINT is required")
	}

	engine := gin.Default()
	engine.Use(middleware.CreateCorsMiddleware())

	db, err := initDatabase()
	if err != nil {
		log.Fatalf("init database: %v", err)
	}
	defer db.Close()

	router.RegisterRoutes(engine, db)

	port := router.Env("PORT", "8080")
	if err := engine.Run(":" + port); err != nil {
		log.Fatalf("server: %v", err)
	}
}

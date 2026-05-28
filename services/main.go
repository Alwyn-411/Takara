package main

import (
	"takara/services/middleware"
	"takara/services/router"
	"takara/services/schema"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	_ "modernc.org/sqlite"
)

var db *sqlx.DB

func initDatabase() (*sqlx.DB, error) {
	var err error
	db, err = sqlx.Connect("sqlite", "takara.db")
	if err != nil {
		return nil, err
	}

	initTables(db)

	return db, nil
}

func initTables(db *sqlx.DB) {
	// Enable Foreign Keys if not panic
	schema.ForeignKeysEnabled(db)

	// Set Up Tables if not panic
	schema.InitUsers(db)
	schema.InitAccounts(db)
	schema.InitCategories(db)
	schema.InitMerchant(db)
	schema.InitTags(db)
	schema.InitTransactions(db)
	schema.InitTransactionTags(db)

	// SetUp Indexes if not panic
	schema.InitIndexUsers(db)
	schema.InitIndexTransactions(db)
}

func main() {
	engine := gin.Default()
	engine.Use(middleware.CreateCorsMiddleware())

	db, err := initDatabase()
	if err != nil {
		return
	}
	defer db.Close()

	router.RegisterRoutes(engine, db)

	engine.Run()
}

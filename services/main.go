package main

import (
	"database/sql"
	"takara/services/router"

	"github.com/gin-gonic/gin"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func initDatabase() error {
	var err error
	db, err = sql.Open("sqlite", "takara.db")
	if err != nil {
		return err
	}
	
	return nil
}

func main()  {
	engine := gin.Default()
	err := initDatabase()
	if err != nil {
		return
	}

	defer db.Close()

	router.RegisterRoutes(engine)

	engine.Run()
}

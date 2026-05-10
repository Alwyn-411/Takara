package handlers

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"

	"github.com/gin-gonic/gin"
)

func Get(ctx *gin.Context, dbInstance *sql.DB) error {
	tableName := ctx.Param("table")

	query := fmt.Sprintf("SELECT * FROM %s", tableName)

	rows, err := dbInstance.Query(query)
	if err != nil {
		return err
	}

	rows.Scan()

	return nil
}

func Create(ctx *gin.Context){

}

func Update(ctx *gin.Context){

}

func Delete(ctx *gin.Context){

}
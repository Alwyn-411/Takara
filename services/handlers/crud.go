package handlers

// import (
// 	"database/sql"
// 	"fmt"
// 	"net/http"
// 	"reflect"
// 	"strings"
// 	"takara/services/schema"

// 	_ "modernc.org/sqlite"

// 	"github.com/gin-gonic/gin"
// 	"github.com/jmoiron/sqlx"
// )

// type crudModel struct {
// 	Table   string
// 	Columns []string
// 	Model   any
// }

// var models = map[string]crudModel{
// 	"user": {
// 		Table: "user",
// 		Columns: []string{
// 			"name",
// 			"email",
// 			"password",
// 			"active",
// 		},
// 		Model: schema.User{},
// 	},
// }

// func GetRowbyId(ctx *gin.Context, dbInstance *sqlx.DB) (int, any, error) {
// 	tableName := ctx.Param("table")
// 	id := ctx.Param("id")

// 	meta, ok := models[tableName]
// 	if !ok {
// 		return http.StatusBadRequest, nil, fmt.Errorf("invalid table")
// 	}

// 	data := reflect.New(reflect.TypeOf(meta.Model)).Interface()
// 	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ? AND active = ?", meta.Table)

// 	err := dbInstance.Get(data, query, id, true)
// 	if err == sql.ErrNoRows {
// 		return http.StatusNotFound, nil, err
// 	}

// 	if err != nil {
// 		return http.StatusInternalServerError, nil, err
// 	}

// 	return http.StatusOK, data, nil
// }

// func Create(ctx *gin.Context, dbInstance *sqlx.DB) (int, any, error) {
// 	tableName := ctx.Param("table")

// 	meta, ok := models[tableName]
// 	if !ok {
// 		return http.StatusBadRequest, nil, fmt.Errorf("invalid table")
// 	}

// 	data := reflect.New(reflect.TypeOf(meta.Model)).Interface()

// 	if err := ctx.BindJSON(data); err != nil {
// 		return http.StatusBadRequest, nil, err
// 	}

// 	query := fmt.Sprintf(
// 		"INSERT INTO %s (%s) VALUES (%s)",
// 		meta.Table,
// 		strings.Join(meta.Columns, ", "),
// 		":"+strings.Join(meta.Columns, ", :"),
// 	)

// 	result, err := dbInstance.NamedExec(query, data)
// 	if err != nil {
// 		return http.StatusInternalServerError, nil, err
// 	}

// 	id, _ := result.LastInsertId()

// 	return http.StatusCreated, gin.H{
// 		"id": id,
// 	}, nil
// }

// func UpdateRowbyId(ctx *gin.Context, dbInstance *sqlx.DB) (int, any, error) {
// 	tableName := ctx.Param("table")
// 	id := ctx.Param("id")

// 	meta, ok := models[tableName]
// 	if !ok {
// 		return http.StatusBadRequest, nil, fmt.Errorf("invalid table")
// 	}

// 	data := map[string]interface{}{}

// 	if err := ctx.BindJSON(&data); err != nil {
// 		return http.StatusBadRequest, nil, err
// 	}

// 	data["id"] = id

// 	sets := []string{}

// 	for _, col := range meta.Columns {
// 		sets = append(sets, fmt.Sprintf("%s = :%s", col, col))
// 	}

// 	sets = append(sets, "updatedTimeStamp = CURRENT_TIMESTAMP")

// 	query := fmt.Sprintf(
// 		"UPDATE %s SET %s WHERE id = :id",
// 		meta.Table,
// 		strings.Join(sets, ", "),
// 	)

// 	_, err := dbInstance.NamedExec(query, data)
// 	if err != nil {
// 		return http.StatusInternalServerError, nil, err
// 	}

// 	return http.StatusOK, nil, nil
// }

// func DeleteRowbyId(ctx *gin.Context, dbInstance *sqlx.DB) (int, any, error) {
// 	tableName := ctx.Param("table")
// 	id := ctx.Param(("id"))

// 	meta, ok := models[tableName]
// 	if !ok {
// 		return http.StatusBadRequest, nil, fmt.Errorf("invalid table")
// 	}

// 	update := fmt.Sprintf(
// 		"UPDATE %s SET active = ? WHERE id = ?", meta.Table,
// 	)

// 	_, err := dbInstance.Exec(update, false, id)
// 	if err != nil {
// 		return http.StatusInternalServerError, nil, err
// 	}

// 	return http.StatusOK, nil, nil

// }

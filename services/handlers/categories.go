package handlers

import (
	"database/sql"
	"net/http"
	"takara/services/schema"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// type Categories struct {
// 	Base

// 	CategoryId   string `db:"categoryId" json:"categoryId"`
// 	CategoryName string `db:"categoryName" json:"categoryName"`

// 	Timestamp
// }

func CreateCategory(ctx *gin.Context, dbInstance *sqlx.DB) (string, error) {
	categoryId := uuid.New().String()

	var data schema.Categories
	err := ctx.ShouldBindJSON(&data)

	if err != nil {
		return "", err
	}

	data.CategoryId = categoryId

	query := `INSERT INTO categories (userId, categoryId, categoryName) VALUES (:userId, :categoryId, :categoryName)`

	_, err = dbInstance.NamedExec(query, data)
	if err != nil {
		return "", err
	}

	return categoryId, nil
}

func UpdateCategory(categoryId string, ctx *gin.Context, dbInstance *sqlx.DB) (string, error) {
	data := schema.Categories{}
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		return "", err
	}

	namedArgs := gin.H{
		"categoryId":   data.CategoryId,
		"categoryName": data.CategoryName,
		"updatedAt":    time.Now(),
	}

	query := `UPDATE categories SET (categoryName = :categoryName, updatedAt = :updatedAt) WHERE categoryId = :categoryId`

	_, err = dbInstance.NamedExec(query, namedArgs)
	if err != nil {
		return "", err
	}

	return categoryId, nil
}

func GetCategoryById(categoryId string, ctx *gin.Context, dbInstance *sqlx.DB) (schema.Categories, error) {
	query := `SELECT categoryId, categoryName, active FROM categories WHERE categoryId = ? AND active = ?`

	category := schema.Categories{}
	err := dbInstance.Get(&category, query, categoryId, true)

	if err == sql.ErrNoRows {
		return category, err
	}

	if err != nil {
		return category, err
	}

	return category, nil
}

func GetCategoriesById(userId string, ctx *gin.Context, dbInstance *sqlx.DB) {
	query := `SELECT userId, categoryId, categoryName, active FROM categories WHERE userId = ? AND active = ? LIMIT 20 `

	categories := []schema.Categories{}
	err := dbInstance.Select(&categories, query, userId, true)

	if len(categories) == 0 {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var response schema.ListResponse[schema.Categories]

	response.Count = len(categories)
	response.Records = categories

	ctx.JSON(http.StatusOK, response)
}

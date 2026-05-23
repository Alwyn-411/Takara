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

	data := schema.Category{}
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
	data := schema.Category{}

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

func GetCategoryById(categoryId string, ctx *gin.Context, dbInstance *sqlx.DB) (schema.Category, error) {
	query := `SELECT categoryId, categoryName, active FROM categories WHERE categoryId = ? AND active = ?`

	category := schema.Category{}
	err := dbInstance.Get(&category, query, categoryId, 1)

	if err == sql.ErrNoRows {
		return category, err
	}

	if err != nil {
		return category, err
	}

	return category, nil
}

func GetCategoriesByUserId(userId string, ctx *gin.Context, dbInstance *sqlx.DB) {
	query := `SELECT userId, categoryId, categoryName, active FROM categories WHERE userId = ? AND active = ? LIMIT 20 `

	categories := []schema.Category{}
	err := dbInstance.Select(&categories, query, userId, 1)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if len(categories) == 0 {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}

	response := schema.ListResponse[schema.Category]{}

	response.Count = len(categories)
	response.Records = categories

	ctx.JSON(http.StatusOK, response)
}

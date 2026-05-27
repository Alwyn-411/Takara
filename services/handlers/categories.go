package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"takara/services/middleware"
	"takara/services/schema"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CategoryHandler struct {
	dbInstance *sqlx.DB
}

func NewCategoryHandler(db *sqlx.DB) *CategoryHandler {
	return &CategoryHandler{dbInstance: db}
}

type CreateCategoryRequest struct {
	UserId       string `json:"userId" binding:"required"`
	CategoryName string `json:"categoryName" binding:"required"`
}

type UpdateCategoryRequest struct {
	CategoryName *string `json:"categoryName"`
}

func (handler *CategoryHandler) CreateCategory(ctx *gin.Context) {
	var req CreateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	categoryId := uuid.New().String()
	// TODO: req.UserId from the body for now; replace with the authenticated caller's id once auth middleware exists.
	if _, err := handler.dbInstance.NamedExec(
		`INSERT INTO categories (userId, categoryId, categoryName) VALUES (:userId, :categoryId, :categoryName)`,
		map[string]interface{}{
			"userId":       req.UserId,
			"categoryId":   categoryId,
			"categoryName": req.CategoryName,
		},
	); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": categoryId})
}

func (handler *CategoryHandler) UpdateCategory(ctx *gin.Context) {
	categoryId := ctx.Param("categoryId")

	var req UpdateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	if req.CategoryName == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "nothing to update"})
		return
	}

	// TODO: add `AND userId = ?` once auth provides the caller's id.
	res, err := handler.dbInstance.NamedExec(
		`UPDATE categories SET categoryName = :categoryName, updatedAt = :updatedAt
		 WHERE categoryId = :categoryId AND active = 1`,
		map[string]interface{}{
			"categoryId":   categoryId,
			"categoryName": *req.CategoryName,
			"updatedAt":    time.Now().Unix(),
		},
	)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (handler *CategoryHandler) DeleteCategory(ctx *gin.Context) {
	categoryId := ctx.Param("categoryId")

	res, err := handler.dbInstance.NamedExec(
		`UPDATE categories SET active = 0, updatedAt = :updatedAt WHERE categoryId = :categoryId AND active = 1`,
		map[string]interface{}{
			"categoryId": categoryId,
			"updatedAt":  time.Now().Unix(),
		},
	)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (handler *CategoryHandler) GetCategoryById(ctx *gin.Context) {
	categoryId := ctx.Param("categoryId")

	category := schema.Category{}
	err := handler.dbInstance.Get(&category,
		`SELECT categoryId, categoryName, userId, active, createdAt, updatedAt
		 FROM categories WHERE categoryId = ? AND active = 1`, categoryId)
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, category)
}

func (handler *CategoryHandler) ListCategories(ctx *gin.Context) {
	userId, ok := middleware.CurrentUserId(ctx)
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	categories := []schema.Category{}
	err := handler.dbInstance.Select(&categories,
		`SELECT userId, categoryId, categoryName, active, createdAt, updatedAt
		 FROM categories WHERE userId = ? AND active = 1 LIMIT 20`, userId)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, schema.ListResponse[schema.Category]{
		Count:   len(categories),
		Records: categories,
	})
}

func ResolveOrCreateCategory(dbTx *sqlx.Tx, userId, name string) (string, error) {
	if name == "" {
		return "", nil
	}

	var id string
	err := dbTx.Get(&id,
		`SELECT categoryId FROM categories WHERE userId = ? AND categoryName = ? AND active = 1`,
		userId, name,
	)
	if err == nil {
		return id, nil
	}
	if err != sql.ErrNoRows {
		return "", fmt.Errorf("lookup category: %w", err)
	}

	id = uuid.New().String()
	_, err = dbTx.Exec(
		`INSERT INTO categories (categoryId, userId, categoryName) VALUES (?, ?, ?)`,
		id, userId, name,
	)
	if err != nil {
		return "", fmt.Errorf("create category: %w", err)
	}
	return id, nil
}

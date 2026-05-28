package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"takara/services/middleware"
	"takara/services/schema"

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

func (handler *CategoryHandler) ListCategories(ctx *gin.Context) {
	userId, ok := middleware.CurrentUserId(ctx)
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	term := ctx.Query("q")

	categories := []schema.Category{}
	var err error

	if term != "" {
		err = handler.dbInstance.Select(&categories,
			`SELECT userId, categoryId, categoryName, active, createdAt, updatedAt
		 FROM categories WHERE userId = ? AND active = 1 AND categoryName LIKE ? LIMIT 20`, userId, "%"+term+"%")
	} else {
		err = handler.dbInstance.Select(&categories,
			`SELECT userId, categoryId, categoryName, active, createdAt, updatedAt
		 FROM categories WHERE userId = ? AND active = 1 LIMIT 20`, userId)
	}

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
	name = strings.TrimSpace(name)
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

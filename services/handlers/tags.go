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

type TagsHandler struct {
	dbInstance *sqlx.DB
}

func NewTagsHandler(db *sqlx.DB) *TagsHandler {
	return &TagsHandler{dbInstance: db}
}

type CreateTagRequest struct {
	UserId  string `json:"userId" binding:"required"`
	TagName string `json:"tagName" binding:"required"`
}

func (handler *TagsHandler) CreateTag(ctx *gin.Context) {
	var req CreateTagRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	tagId := uuid.New().String()
	// TODO: req.UserId comes from the body for now. Once auth middleware is in
	// place, take the userId from the authenticated context instead and ignore
	// any body-supplied userId.
	if _, err := handler.dbInstance.NamedExec(
		`INSERT INTO tags (userId, tagId, tagName) VALUES (:userId, :tagId, :tagName)`,
		map[string]interface{}{
			"userId":  req.UserId,
			"tagId":   tagId,
			"tagName": req.TagName,
		},
	); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": tagId})
}

type UpdateTagRequest struct {
	TagName *string `json:"tagName"`
}

func (handler *TagsHandler) UpdateTagById(ctx *gin.Context) {
	tagId := ctx.Param("tagId")

	var req UpdateTagRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	if req.TagName == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "nothing to update"})
		return
	}

	res, err := handler.dbInstance.NamedExec(
		`UPDATE tags SET tagName = :tagName, updatedAt = :updatedAt WHERE tagId = :tagId AND active = 1`,
		gin.H{
			"tagId":     tagId,
			"tagName":   *req.TagName,
			"updatedAt": time.Now().Unix(),
		})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": tagId})
}

func (handler *TagsHandler) GetTagById(ctx *gin.Context) {
	tagId := ctx.Param("tagId")

	tag := schema.Tag{}
	err := handler.dbInstance.Get(&tag,
		`SELECT tagId, tagName, userId, active, createdAt, updatedAt
		 FROM tags WHERE tagId = ? AND active = 1`, tagId)
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, tag)
}

func (handler *TagsHandler) ListTags(ctx *gin.Context) {
	userId, ok := middleware.CurrentUserId(ctx)
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	tags := []schema.Tag{}
	err := handler.dbInstance.Select(&tags,
		`SELECT userId, tagId, tagName, active, createdAt, updatedAt
		 FROM tags WHERE userId = ? AND active = 1 LIMIT 20`, userId)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, schema.ListResponse[schema.Tag]{
		Count:   len(tags),
		Records: tags,
	})
}

type AddTagsRequest struct {
	TagIds []string `json:"tagIds" binding:"required"`
}

func (handler *TagsHandler) AddTagsToTransaction(ctx *gin.Context) {
	transactionId := ctx.Param("transactionId")

	var req AddTagsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dbTx, err := handler.dbInstance.Beginx()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start db transaction"})
		return
	}

	for _, tagId := range req.TagIds {
		if _, err := dbTx.Exec(
			`INSERT OR IGNORE INTO transaction_tags (transactionId, tagId) VALUES (?, ?)`,
			transactionId, tagId,
		); err != nil {
			dbTx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add tag: " + tagId})
			return
		}
	}

	if err := dbTx.Commit(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (handler *TagsHandler) RemoveTagFromTransaction(ctx *gin.Context) {
	transactionId := ctx.Param("transactionId")
	tagId := ctx.Param("tagId")

	if _, err := handler.dbInstance.Exec(
		`DELETE FROM transaction_tags WHERE transactionId = ? AND tagId = ?`,
		transactionId, tagId,
	); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove tag"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (handler *TagsHandler) GetTagsForTransaction(ctx *gin.Context) {
	transactionId := ctx.Param("transactionId")

	tags := []schema.Tag{}
	if err := handler.dbInstance.Select(&tags, `
		SELECT t.tagId, t.tagName, t.userId, t.active, t.createdAt, t.updatedAt
		FROM tags t
		JOIN transaction_tags tt ON t.tagId = tt.tagId
		WHERE tt.transactionId = ?
	`, transactionId); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, schema.ListResponse[schema.Tag]{
		Count:   len(tags),
		Records: tags,
	})
}

func ResolveOrCreateTag(dbTx *sqlx.Tx, userId, name string) (string, error) {
	if name == "" {
		return "", nil
	}

	var id string
	err := dbTx.Get(&id,
		`SELECT tagId FROM tags WHERE userId = ? AND tagName = ? AND active = 1`,
		userId, name,
	)
	if err == nil {
		return id, nil
	}
	if err != sql.ErrNoRows {
		return "", fmt.Errorf("lookup tag: %w", err)
	}

	id = uuid.New().String()
	_, err = dbTx.Exec(
		`INSERT INTO tags (tagId, userId, tagName) VALUES (?, ?, ?)`,
		id, userId, name,
	)
	if err != nil {
		return "", fmt.Errorf("create tag: %w", err)
	}
	return id, nil
}

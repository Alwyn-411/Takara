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

func CreateTag(ctx *gin.Context, dbInstance *sqlx.DB) (string, error) {
	TagId := uuid.New().String()

	data := schema.Tag{}
	err := ctx.ShouldBindJSON(&data)

	if err != nil {
		return "", err
	}

	data.TagId = TagId

	query := `INSERT INTO tags (userId, tagId, tagName) VALUES (:userId, :tagId, :tagName)`

	_, err = dbInstance.NamedExec(query, data)
	if err != nil {
		return "", err
	}

	return TagId, nil
}

func UpdateTagbyId(tagId string, ctx *gin.Context, dbInstance *sqlx.DB) (string, error) {
	data := schema.Tag{}
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		return "", err
	}

	namedArgs := gin.H{
		"tagId":     data.TagId,
		"tagName":   data.TagName,
		"updatedAt": time.Now(),
	}

	query := `UPDATE tags SET tagName = :tagName, updatedAt = :updatedAt WHERE tagId = :tagId`

	_, err = dbInstance.NamedExec(query, namedArgs)
	if err != nil {
		return "", err
	}

	return tagId, nil
}

func GetTagById(tagId string, ctx *gin.Context, dbInstance *sqlx.DB) (schema.Tag, error) {
	query := `SELECT tagId, tagName, active FROM tags WHERE tagId = ? AND active = ?`

	tag := schema.Tag{}
	err := dbInstance.Get(&tag, query, tagId, true)

	if err == sql.ErrNoRows {
		return tag, err
	}

	if err != nil {
		return tag, err
	}

	return tag, nil
}

func GetTagsByUserId(userId string, ctx *gin.Context, dbInstance *sqlx.DB) {
	query := `SELECT userId, tagId, tagName, active FROM tags WHERE userId = ? AND active = ? LIMIT 20 `

	tags := []schema.Tag{}
	err := dbInstance.Select(&tags, query, userId, true)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if len(tags) == 0 {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}

	response := schema.ListResponse[schema.Tag]{}

	response.Count = len(tags)
	response.Records = tags

	ctx.JSON(http.StatusOK, response)
}

func AddTagsToTransaction(transactionId string, ctx *gin.Context, dbInstance *sqlx.DB) {
	var req struct {
		TagIds []string `json:"tagIds" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, tagId := range req.TagIds {
		_, err := dbInstance.Exec(
			`INSERT OR IGNORE INTO transaction_tags (transactionId, tagId) VALUES (?, ?)`,
			transactionId, tagId,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add tag: " + tagId})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "tags added"})
}

func RemoveTagFromTransaction(transactionId string, tagId string, ctx *gin.Context, dbInstance *sqlx.DB) {
	_, err := dbInstance.Exec(
		`DELETE FROM transaction_tags WHERE transactionId = ? AND tagId = ?`,
		transactionId, tagId,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove tag"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "tag removed"})
}

func GetTagsForTransaction(transactionId string, ctx *gin.Context, dbInstance *sqlx.DB) {
	tags := []schema.Tag{}
	err := dbInstance.Select(&tags, `
		SELECT t.tagId, t.tagName, t.userId, t.active, t.createdAt, t.updatedAt
		FROM tags t
		JOIN transaction_tags tt ON t.tagId = tt.tagId
		WHERE tt.transactionId = ?
	`, transactionId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := schema.ListResponse[schema.Tag]{
		Count:   len(tags),
		Records: tags,
	}
	ctx.JSON(http.StatusOK, response)
}

package holdings

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"takara/services/middleware"

	"github.com/gin-gonic/gin"
)


func (handler *HoldingsHandler) DeleteHoldingById(ctx *gin.Context) {
	userId, ok := middleware.CurrentUserId(ctx)
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	holdingId := ctx.Param("id")

	query := `
		UPDATE holdings
		SET active = 0, updatedAt = strftime('%s','now')
		WHERE holdingId = :holdingId AND userId = :userId AND active = 1
	`

	result, err := handler.dbInstance.NamedExec(query, gin.H{
		"holdingId":    holdingId,
		"userId": userId,
	})
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if rows == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf(
				"holding with holdingId=%s does not exist",
				holdingId,
			),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (handler *HoldingsHandler) DeleteValuation(ctx *gin.Context) {
	userId, ok := middleware.CurrentUserId(ctx)
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	valuationId := ctx.Param("id")

	tx, err := handler.dbInstance.Beginx()
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer tx.Rollback() // no-op once Commit succeeds

	var holdingId string
	err = tx.Get(&holdingId, `
		SELECT holdingId FROM holding_valuations
		WHERE valuationId = ? AND userId = ? AND active = 1`,
		valuationId, userId)
	if errors.Is(err, sql.ErrNoRows) {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("valuation with valuationId=%s does not exist", valuationId),
		})
		return
	}
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var activeCount int
	if err = tx.Get(&activeCount, `
		SELECT COUNT(*) FROM holding_valuations
		WHERE holdingId = ? AND active = 1`, holdingId); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if activeCount <= 1 {
		ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"error": "cannot delete the only valuation for this holding; edit it or delete the holding instead",
		})
		return
	}

	if _, err = tx.Exec(`
		UPDATE holding_valuations
		SET active = 0, updatedAt = strftime('%s','now')
		WHERE valuationId = ? AND userId = ? AND active = 1`,
		valuationId, userId); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}
package holdings

import (
	"fmt"
	"net/http"
	"takara/services/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateHoldingRequest struct {
	Type        *string `json:"type,omitempty"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	ClosedAt    *int64  `json:"closedAt,omitempty"`
}

func (handler *HoldingsHandler) UpdateHoldingById(ctx *gin.Context) {
	userId, ok := middleware.CurrentUserId(ctx)
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	holdingId := ctx.Param("id")

	var req UpdateHoldingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := map[string]any{
		"holdingId":   holdingId,
		"userId":      userId,
		"type":        req.Type,
		"name":        req.Name,
		"description": req.Description,
		"closedAt":    req.ClosedAt,
	}

	query := `
		UPDATE holdings
		SET type        = COALESCE(:type, type),
		    name        = COALESCE(:name, name),
		    description = COALESCE(:description, description),
		    closedAt    = COALESCE(:closedAt, closedAt),
		    updatedAt   = strftime('%s','now')
		WHERE holdingId = :holdingId AND userId = :userId AND active = 1`

	result, err := handler.dbInstance.NamedExec(query, params)
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
			"error": fmt.Sprintf("holding with holdingId=%s does not exist", holdingId),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

type RecordValuationRequest struct {
	HoldingId  string  `json:"holdingId" binding:"required"`
	Value      string  `json:"value" binding:"required"`
	Quantity   *string `json:"quantity,omitempty"`
	UnitPrice  *string `json:"unitPrice,omitempty"`
	ObservedAt int64   `json:"observedAt,omitempty"`
	Note       *string `json:"note,omitempty"`
}

func (handler *HoldingsHandler) CreateValuation(ctx *gin.Context) {
	userId, ok := middleware.CurrentUserId(ctx)
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var req RecordValuationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	observedAt := req.ObservedAt
	if observedAt == 0 {
		observedAt = time.Now().Unix()
	}
	valuationId := uuid.New().String()

	params := map[string]any{
		"valuationId": valuationId,
		"userId":      userId,
		"holdingId":   req.HoldingId,
		"value":       req.Value,
		"quantity":    req.Quantity,
		"unitPrice":   req.UnitPrice,
		"observedAt":  observedAt,
		"note":        req.Note,
	}

	query := `
		INSERT INTO holding_valuations (valuationId, userId, holdingId, value, quantity, unitPrice, observedAt, note)
		SELECT :valuationId, :userId, :holdingId, :value, :quantity, :unitPrice, :observedAt, :note
		WHERE EXISTS (
			SELECT 1 FROM holdings WHERE holdingId = :holdingId AND userId = :userId AND active = 1
		)`

	result, err := handler.dbInstance.NamedExec(query, params)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rows, err := result.RowsAffected()
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("holding with holdingId=%s does not exist", req.HoldingId),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"holdingId": req.HoldingId, "valuationId": valuationId})
}
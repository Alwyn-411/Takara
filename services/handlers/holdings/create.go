package holdings

import (
	"net/http"
	"takara/services/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateHoldingRequest struct {
	Kind        string `json:"kind"        binding:"required"`
	Type        string `json:"type"        binding:"required,oneof=Asset Liability"`
	Name        string `json:"name"        binding:"required"`
	Description string `json:"description,omitempty"`
	Currency    string `json:"currency"    binding:"required"`
	OpenedAt    int64  `json:"openedAt,omitempty"`

	// Seed valuation, written in the same transaction as the holding.
	Value      string  `json:"value"      binding:"required"` // magnitude, in the holding's currency
	Quantity   *string `json:"quantity,omitempty"`            // optional: units held (shares, grams)
	UnitPrice  *string `json:"unitPrice,omitempty"`           // optional: price per unit at observedAt
	ObservedAt int64   `json:"observedAt,omitempty"`
	Note       *string `json:"note,omitempty"`
}

func (handler *HoldingsHandler) CreateHolding(ctx *gin.Context) {
	userId, ok := middleware.CurrentUserId(ctx)
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var data CreateHoldingRequest
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now().Unix()
	openedAt := data.OpenedAt
	if openedAt == 0 {
		openedAt = now
	}
	observedAt := data.ObservedAt
	if observedAt == 0 {
		observedAt = now
	}

	holdingId := uuid.New().String()
	valuationId := uuid.New().String()

	tx, err := handler.dbInstance.Beginx()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer tx.Rollback()

	if _, err = tx.Exec(`
		INSERT INTO holdings (holdingId, userId, kind, type, name, description, currency, openedAt)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		holdingId, userId, data.Kind, data.Type, data.Name, data.Description, data.Currency, openedAt,
	); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if _, err = tx.Exec(`
		INSERT INTO holding_valuations
			(valuationId, userId, holdingId, value, quantity, unitPrice, observedAt, note)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		valuationId, userId, holdingId, data.Value, data.Quantity, data.UnitPrice, observedAt, data.Note,
	); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err = tx.Commit(); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"holdingId": holdingId, "valuationId": valuationId})
}
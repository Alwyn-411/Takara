package holdings

import (
	"database/sql"
	"errors"
	"net/http"
	"takara/services/middleware"
	"takara/services/schema"

	"github.com/gin-gonic/gin"
)


type HoldingWithValue struct {
	schema.Holding
	CurrentValue *string `db:"currentValue" json:"currentValue,omitempty"`
	ValuedAt     *int64  `db:"valuedAt" json:"valuedAt,omitempty"`
}

func (handler *HoldingsHandler) GetHoldingById(ctx *gin.Context) {
	userId, ok := middleware.CurrentUserId(ctx)
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var holding schema.Holding
	holdingId := ctx.Param("id")

	query := `SELECT holdingId, kind, type, name, description, currency, openedAt, closedAt 
				FROM holdings  WHERE userId = ? AND holdingId = ? AND active = 1`

	err := handler.dbInstance.Get(&holding, query, userId, holdingId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	
	ctx.JSON(http.StatusOK, holding)
}


func (handler *HoldingsHandler) ListHoldings(ctx *gin.Context) {
	userId, ok := middleware.CurrentUserId(ctx)
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var holdings []HoldingWithValue
	query := `
		SELECT h.*, v.value AS currentValue, v.observedAt AS valuedAt
		FROM holdings h
		LEFT JOIN holding_valuations v
		  ON v.holdingId = h.holdingId AND v.active = 1
		 AND v.observedAt = (
		     SELECT MAX(v2.observedAt) FROM holding_valuations v2
		     WHERE v2.holdingId = h.holdingId AND v2.active = 1
		 )
		WHERE h.userId = ? AND h.active = 1
		ORDER BY h.kind, h.name`

	err := handler.dbInstance.Select(&holdings, query, userId)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if len(holdings) == 0 {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}

	var response schema.ListResponse[HoldingWithValue]
	response.Count = len(holdings)
	response.Records = holdings
	
	ctx.JSON(http.StatusOK, response)
}

type ValuationResponse struct {
	ValuationId string  `db:"valuationId" json:"valuationId"`
	HoldingId   string  `db:"holdingId" json:"holdingId"`
	Value       string  `db:"value" json:"value"`
	Quantity    *string `db:"quantity" json:"quantity,omitempty"`
	UnitPrice   *string `db:"unitPrice" json:"unitPrice,omitempty"`
	ObservedAt  int64   `db:"observedAt" json:"observedAt"`
	Note        *string `db:"note" json:"note,omitempty"`
	CreatedAt   int64   `db:"createdAt" json:"createdAt"`
}

func (handler *HoldingsHandler) ListValuations(ctx *gin.Context) {
	userId, ok := middleware.CurrentUserId(ctx)
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	holdingId := ctx.Param("id")

	query := `
		SELECT valuationId, holdingId, value, quantity, unitPrice, observedAt, note, createdAt
		FROM holding_valuations
		WHERE holdingId = ? AND userId = ? AND active = 1
		ORDER BY observedAt ASC`

	var valuations []ValuationResponse
	if err := handler.dbInstance.Select(&valuations, query, holdingId, userId); err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var response schema.ListResponse[ValuationResponse]

	response.Count = len(valuations)
	response.Records = valuations

	ctx.JSON(http.StatusOK, response)
}
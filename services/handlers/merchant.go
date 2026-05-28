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

type MerchantHandler struct {
	dbInstance *sqlx.DB
}

func NewMerchantHandler(db *sqlx.DB) *MerchantHandler {
	return &MerchantHandler{dbInstance: db}
}

func (handler *MerchantHandler) ListMerchants(ctx *gin.Context) {
	userId, ok := middleware.CurrentUserId(ctx)
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	term := ctx.Query("q")

	merchants := []schema.Merchant{}
	var err error

	if term != "" {
		err = handler.dbInstance.Select(&merchants,
			`SELECT userId, merchantId, merchantName, active, createdAt, updatedAt
			 FROM merchants WHERE userId = ? AND active = 1 AND merchantName LIKE ? LIMIT 20`,
			userId, "%"+term+"%")
	} else {
		err = handler.dbInstance.Select(&merchants,
			`SELECT userId, merchantId, merchantName, active, createdAt, updatedAt
			 FROM merchants WHERE userId = ? AND active = 1 LIMIT 20`, userId)
	}

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, schema.ListResponse[schema.Merchant]{
		Count:   len(merchants),
		Records: merchants,
	})
}

func ResolveOrCreateMerchant(dbTx *sqlx.Tx, userId, name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", nil
	}

	var id string
	err := dbTx.Get(&id,
		`SELECT merchantId FROM merchants WHERE userId = ? AND merchantName = ? AND active = 1`,
		userId, name,
	)
	if err == nil {
		return id, nil
	}
	if err != sql.ErrNoRows {
		return "", fmt.Errorf("lookup merchant: %w", err)
	}

	id = uuid.New().String()
	_, err = dbTx.Exec(
		`INSERT INTO merchants (merchantId, userId, merchantName) VALUES (?, ?, ?)`,
		id, userId, name,
	)
	if err != nil {
		return "", fmt.Errorf("create merchant: %w", err)
	}
	return id, nil
}

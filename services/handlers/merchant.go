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

type MerchantHandler struct {
	dbInstance *sqlx.DB
}

func NewMerchantHandler(db *sqlx.DB) *MerchantHandler {
	return &MerchantHandler{dbInstance: db}
}

type CreateMerchantRequest struct {
	UserId       string `json:"userId" binding:"required"`
	MerchantName string `json:"merchantName" binding:"required"`
}

type UpdateMerchantRequest struct {
	MerchantName *string `json:"merchantName"`
}

func (handler *MerchantHandler) CreateMerchant(ctx *gin.Context) {
	var req CreateMerchantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	merchantId := uuid.New().String()
	if _, err := handler.dbInstance.NamedExec(
		`INSERT INTO merchants (userId, merchantId, merchantName) VALUES (:userId, :merchantId, :merchantName)`,
		gin.H{
			"userId":       req.UserId,
			"merchantId":   merchantId,
			"merchantName": req.MerchantName,
		},
	); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": merchantId})
}

func (handler *MerchantHandler) UpdateMerchant(ctx *gin.Context) {
	merchantId := ctx.Param("merchantId")

	var req UpdateMerchantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	if req.MerchantName == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "nothing to update"})
		return
	}

	res, err := handler.dbInstance.NamedExec(
		`UPDATE merchants SET merchantName = :merchantName, updatedAt = :updatedAt
		 WHERE merchantId = :merchantId AND active = 1`,
		gin.H{
			"merchantId":   merchantId,
			"merchantName": *req.MerchantName,
			"updatedAt":    time.Now().Unix(),
		},
	)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "merchant not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (handler *MerchantHandler) DeleteMerchant(ctx *gin.Context) {
	merchantId := ctx.Param("merchantId")

	res, err := handler.dbInstance.NamedExec(
		`UPDATE merchants SET active = 0, updatedAt = :updatedAt WHERE merchantId = :merchantId AND active = 1`,
		gin.H{
			"merchantId": merchantId,
			"updatedAt":  time.Now().Unix(),
		},
	)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "merchant not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (handler *MerchantHandler) GetMerchantById(ctx *gin.Context) {
	merchantId := ctx.Param("merchantId")

	merchant := schema.Merchant{}
	err := handler.dbInstance.Get(&merchant,
		`SELECT merchantId, merchantName, userId, active, createdAt, updatedAt
		 FROM merchants WHERE merchantId = ? AND active = 1`, merchantId)
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "merchant not found"})
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, merchant)
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

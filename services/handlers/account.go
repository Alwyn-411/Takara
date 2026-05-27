package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"takara/services/middleware"
	"takara/services/schema"
	"time"

	"github.com/bojanz/currency"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AccountHandler struct {
	dbInstance *sqlx.DB
}

func NewAccountHandler(db *sqlx.DB) *AccountHandler {
	return &AccountHandler{
		dbInstance: db,
	}
}

type UpdateAccountRequest struct {
	Type          *string `json:"type"`
	Name          *string `json:"name"`
	Active        *int    `json:"active"`
	Balance       *string `json:"balance"`
	Currency      *string `json:"currency"`
	Interest      *string `json:"interest"`
	Description   *string `json:"description"`
	AccountNumber *string `json:"accountNumber"`
}

func (handler *AccountHandler) GetAccountById(ctx *gin.Context) {
	id := ctx.Param("accountId")

	query := `
		SELECT userId, accountId, type, name, accountNumber, description, currency, interest, balance, active
		FROM accounts
		WHERE accountId = ? AND active = ?
	`

	var account schema.Account
	err := handler.dbInstance.Get(&account, query, id, true)

	if err == sql.ErrNoRows {
		ctx.AbortWithError(http.StatusNoContent, sql.ErrNoRows)
		return
	}

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (handler *AccountHandler) CreateAccount(ctx *gin.Context) {
	accountId := uuid.New().String()

	var data schema.Account
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	data.AccountID = accountId
	newBalance, err := currency.NewAmount(data.Balance, data.Currency)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid amount or currency"})
		return
	}

	data.Balance = newBalance.Number()

	query := `INSERT 
	INTO accounts (
		userId, accountId, type, name, accountNumber, description, currency, interest, balance
	) VALUES (
		:userId, :accountId, :type, :name, :accountNumber, :description, :currency, :interest, :balance
	 )`

	_, err = handler.dbInstance.NamedExec(query, data)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": data.UserID})
}

func (handler *AccountHandler) UpdateAccountById(ctx *gin.Context) {
	accountId := ctx.Param("accountId")

	var req UpdateAccountRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if req.Balance != nil && req.Currency != nil {
		newBalance, err := currency.NewAmount(*req.Balance, *req.Currency)
		if err != nil {
			ctx.AbortWithError(400, err)
			return
		}

		b := newBalance.Number()
		req.Balance = &b
	}

	query := `
		UPDATE accounts
		SET type = COALESCE(:type, type),
			name = COALESCE(:name, name),
			active = COALESCE(:active, active),
			balance = COALESCE(:balance, balance),
			currency = COALESCE(:currency, currency),
			interest = COALESCE(:interest, interest),
			description = COALESCE(:description, description),
			accountNumber = COALESCE(:accountNumber, accountNumber),
			updatedAt = :updatedAt
		WHERE accountId = :accountId
	`

	params := gin.H{
		"accountId":     accountId,
		"updatedAt":     time.Now().Unix(),
		"type":          req.Type,
		"name":          req.Name,
		"active":        req.Active,
		"balance":       req.Balance,
		"currency":      req.Currency,
		"interest":      req.Interest,
		"description":   req.Description,
		"accountNumber": req.AccountNumber,
	}

	result, err := handler.dbInstance.NamedExec(query, params)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	rows, _ := result.RowsAffected()

	if rows == 0 {
		ctx.AbortWithError(http.StatusNoContent, fmt.Errorf("Account with accountId = %s does not Exist", accountId))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (handler *AccountHandler) DeleteAccountById(ctx *gin.Context) {
	accountId := ctx.Param("accountId")

	query := `
		UPDATE accounts SET active = :active, updatedAt = :updatedAt WHERE accountId = :accountId
	`
	_, err := handler.dbInstance.NamedExec(query, gin.H{"accountId": accountId, "active": 0, "updatedAt": time.Now().Unix()})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (handler *AccountHandler) ListAccounts(ctx *gin.Context) {
	clientId, ok := middleware.CurrentUserId(ctx)
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	query := `
		SELECT userId, accountId, type, name, accountNumber, description, currency, interest, balance, active
		FROM accounts
		WHERE userId = ? AND active = ? 
		LIMIT 10
	`

	accounts := []schema.Account{}
	err := handler.dbInstance.Select(&accounts, query, clientId, true)

	if len(accounts) == 0 {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var response schema.ListResponse[schema.Account]

	response.Count = len(accounts)
	response.Records = accounts

	ctx.JSON(http.StatusOK, response)
}

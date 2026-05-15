package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"takara/services/schema"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func GetAccountById(ctx *gin.Context, dbInstance *sqlx.DB) {
	id := ctx.Param("accountId")

	query := `
		SELECT user_id, account_id, type, name, account_number, description, currency, interest, balance, active
		FROM accounts
		WHERE account_id = ? AND active = ?
	`

	var account schema.Account
	err := dbInstance.Get(&account, query, id, true)

	if err == sql.ErrNoRows {
		ctx.AbortWithError(http.StatusNotFound, sql.ErrNoRows)
		return
	}

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func CreateAccount(ctx *gin.Context, dbInstance *sqlx.DB) {
	accountId := uuid.New().String()

	var data schema.Account
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	data.AccountID = accountId

	query := `INSERT 
	INTO accounts (
		user_id, account_id, type, name, account_number, description, currency, interest, balance, active
	) VALUES (
		:user_id, :account_id, :type, :name, :account_number, :description, :currency, :interest, :balance, :active
	 )`

	_, err = dbInstance.NamedExec(query, data)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to create Account"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": data.UserID})
}

func UpdateAccountById(ctx *gin.Context, dbInstance *sqlx.DB) {
	accountId := ctx.Param("accountId")

	var data map[string]any

	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	set := []string{}
	args := []any{}

	allowedFields := map[string]bool{
		"type":           true,
		"name":           true,
		"active":         true,
		"balance":        true,
		"currency":       true,
		"interest":       true,
		"description":    true,
		"account_number": true,
	}

	for key, value := range data {
		if !allowedFields[key] {
			continue
		}
		set = append(set, fmt.Sprintf("%s = ?", key))
		args = append(args, value)
	}

	set = append(set, "updatedTimeStamp = ?")
	args = append(args, time.Now())

	args = append(args, accountId)

	query := fmt.Sprintf(`UPDATE accounts SET %s WHERE account_id = ?`, strings.Join(set, ", "))

	result, err := dbInstance.Exec(query, args...)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if rows == 0 {
		ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("Account with account_id = %s does not Exist", accountId))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func DeleteAccountById(ctx *gin.Context, dbInstance *sqlx.DB) {
	accountId := ctx.Param("accountId")

	query := `
		UPDATE accounts SET active = :active WHERE account_id = :id
	`
	_, err := dbInstance.NamedExec(query, gin.H{"account_id": accountId, "active": false})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

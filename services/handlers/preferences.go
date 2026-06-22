package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"takara/services/middleware"
	"takara/services/schema"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type UserPrefHandler struct {
	dbInstance *sqlx.DB
}

func NewUserPrefHandler(db *sqlx.DB) *UserPrefHandler {
	return &UserPrefHandler{
		dbInstance: db,
	}
}

type UserPrefRequest struct {
	Currency string `json:"currency"`
	Theme    string `json:"theme"`
}

type CreateUserPrefRequest struct {
	UserId string `json:"userId"`

	UserPrefRequest
}

func (handler *UserPrefHandler) GetPref(ctx *gin.Context) {
	clientId, ok := middleware.CurrentUserId(ctx)
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	query := `
		SELECT userId, currency, theme
		FROM userPrefs
		WHERE userId = ? AND active = 1
	`

	var prefs schema.Preferences
	err := handler.dbInstance.Get(&prefs, query, clientId)

	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "UserPrefNotSet"})
		return
	}

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, prefs)
}

func (handler *UserPrefHandler) CreatePref(ctx *gin.Context) {
	var prefs CreateUserPrefRequest
	if err := ctx.ShouldBindJSON(&prefs); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	params := gin.H{
		"userId":   prefs.UserId,
		"currency": prefs.Currency,
		"theme":    prefs.Theme,
	}

	insertQuery := `
		INSERT INTO userPrefs (userId, currency, theme) VALUES (:userId, :currency, :theme)
	`

	_, err := handler.dbInstance.NamedExec(insertQuery, params)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (handler *UserPrefHandler) UpdateUserPref(ctx *gin.Context) {
	clientId, ok := middleware.CurrentUserId(ctx)
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var req UserPrefRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	params := gin.H{
		"userId":   clientId,
		"currency": req.Currency,
		"theme":    req.Theme,
	}

	query := `
		UPDATE userPrefs
		SET currency = COALESCE(:currency, currency),
			theme = COALESCE(:theme, theme)
		WHERE userId = :userId
	`

	_, err := handler.dbInstance.NamedExec(query, params)

	if err != nil {
		log.Println(err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

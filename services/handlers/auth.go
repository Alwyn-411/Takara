package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"takara/services/middleware"
	"takara/services/schema"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	dbInstance *sqlx.DB
	tokens     *middleware.TokenService
}

func NewAuthHandler(db *sqlx.DB, tokens *middleware.TokenService) *AuthHandler {
	return &AuthHandler{
		dbInstance: db,
		tokens:     tokens,
	}
}

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func checkPassword(hashed string, plaintext string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plaintext))
}

type LoginRequest struct {
	UserName string `json:"userName" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (handler *AuthHandler) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var user schema.User
	err := handler.dbInstance.Get(&user,
		`SELECT userId, userName, password FROM users WHERE userName = ?`, req.UserName)
	if errors.Is(err, sql.ErrNoRows) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err := checkPassword(user.Password, req.Password); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := handler.tokens.CreateToken(user.UserID)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
		"id":    user.UserID,
	})
}

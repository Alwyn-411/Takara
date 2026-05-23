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

type UserHandler struct {
	dbInstance *sqlx.DB
}

func NewUserHandler(db *sqlx.DB) *UserHandler {
	return &UserHandler{
		dbInstance: db,
	}
}

type UpdateUserRequest struct {
	UserName *string `json:"userName"`
	AltName  *string `json:"altName"`
	Email    *string `json:"email"`
	AltEmail *string `json:"altEmail"`
	Password *string `json:"password"`
}

func (handler *UserHandler) GetUserById(ctx *gin.Context) {
	id := ctx.Param("id")

	query := `
		SELECT 
			userId, active, userName,
			altName, email, altEmail,
			password, createdAt, updatedAt
		FROM users
		WHERE userId = ? AND active = ?
	`

	var user schema.User
	err := handler.dbInstance.Get(&user, query, id, true)

	if err == sql.ErrNoRows {
		ctx.AbortWithError(http.StatusNotFound, sql.ErrNoRows)
		return
	}

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (handler *UserHandler) CreateUser(ctx *gin.Context) {
	id := uuid.New().String()

	var data schema.User

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	data.UserID = id

	// Hash password
	hashedPassword, err := HashPassword(data.Password)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	data.Password = string(hashedPassword)

	insertQuery := `
		INSERT INTO users (
			userId, userName, altName,
			email, altEmail, password
		) VALUES (
			:userId, :userName, :altName,
			:email, :altEmail, :password
		)
	`

	_, err = handler.dbInstance.NamedExec(insertQuery, data)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"error": "username or email already exists",
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id": data.UserID,
	})
}

func (handler *UserHandler) UpdateUserById(ctx *gin.Context) {
	userId := ctx.Param("id")

	var req UpdateUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	var hashedPassword *string

	if req.Password != nil {
		hash, err := HashPassword(*req.Password)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		hashStr := string(hash)
		hashedPassword = &hashStr
	}

	query := `
		UPDATE users
		SET userName = COALESCE(:userName, userName),
			altName = COALESCE(:altName, altName),
			email = COALESCE(:email, email),
			altEmail = COALESCE(:altEmail, altEmail),
			password = COALESCE(:password, password),
			updatedAt = :updatedAt
		WHERE userId = :userId
	`

	params := gin.H{
		"userId":    userId,
		"updatedAt": time.Now().Unix(),

		"userName": req.UserName,
		"altName":  req.AltName,
		"email":    req.Email,
		"altEmail": req.AltEmail,
		"password": hashedPassword,
	}

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
			"error": fmt.Sprintf(
				"user with userId=%s does not exist",
				userId,
			),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (handler *UserHandler) DeleteUserById(ctx *gin.Context) {
	id := ctx.Param("id")

	query := `
		UPDATE users
		SET active = :active
		WHERE userId = :userId
	`

	result, err := handler.dbInstance.NamedExec(query, gin.H{
		"userId": id,
		"active": false,
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
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

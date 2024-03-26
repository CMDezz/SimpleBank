package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/util"
)

type createUserRequest struct {
	UserName string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	FullName          string    `json:"full_name"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

type refreshAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpriesAt time.Time `json:"access_token_expires_at"`
}

type refreshAccessTokenRequest struct {
	RefeshToken string `json:"refresh_token" binding:"required"`
}

func NewUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		Email:             user.Email,
		FullName:          user.FullName,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		fmt.Println("err 1")
		return
	}

	hased_password, err := util.HashPassword(req.Password)
	if err != nil {
		fmt.Println("err 2")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:      req.UserName,
		FullName:      req.FullName,
		Email:         req.Email,
		HasedPassword: hased_password,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				fmt.Println("err 3")
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := NewUserResponse(user)
	ctx.JSON(http.StatusOK, resp)

}

type getUserByUserNameRequest struct {
	UserName string `uri:"username" binding:"required,alphanum"`
}

func (server *Server) getUserByUserName(ctx *gin.Context) {
	var req getUserByUserNameRequest
	if err := ctx.ShouldBindQuery(req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.UserName)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
	return
}

type loginUserRequest struct {
	UserName string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required"`
}

type loginUserRespone struct {
	SessionId           uuid.UUID    `json:"session_id"`
	User                userResponse `json:"user"`
	AccessToken         string       `json:"access_token"`
	AccessTokenExpired  time.Time    `json:"access_token_expired"`
	RefreshToken        string       `json:"refresh_token"`
	RefreshTokenExpired time.Time    `json:"refresh_token_expired"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	//check request
	var req loginUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	fmt.Println(req.UserName)
	//get user
	user, err := server.store.GetUser(ctx, req.UserName)
	if err != nil {
		//norow || server errpr
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//check password
	err = util.CheckPassword(req.Password, user.HasedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	//create token
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//refresh token
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.store.CreateSessions(ctx, db.CreateSessionsParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	userResponse := NewUserResponse(user)
	ctx.JSON(http.StatusOK, loginUserRespone{
		SessionId:           session.ID,
		User:                userResponse,
		AccessToken:         accessToken,
		AccessTokenExpired:  accessPayload.ExpiredAt,
		RefreshToken:        refreshToken,
		RefreshTokenExpired: refreshPayload.ExpiredAt,
	})
	return
}

func (server *Server) RefreshAccessToken(ctx *gin.Context) {
	var req refreshAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	refreshPayload, err := server.tokenMaker.ValidToken(req.RefeshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.store.GetSessions(ctx, refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if session.IsBlocked {
		err = fmt.Errorf("token is blocked")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if session.Username != refreshPayload.Username {
		err = fmt.Errorf("token mismatch payload infomation")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if time.Now().After(refreshPayload.ExpiredAt) {
		err = fmt.Errorf("refresh token is expired")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(refreshPayload.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, &refreshAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpriesAt: accessTokenPayload.ExpiredAt,
	})
	return
}

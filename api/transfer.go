package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/techschool/simplebank/db/sqlc"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,oneof=USD VND EUR"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest

	//validate req json
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//check valid account
	if !server.validAccount(ctx, req.FromAccountID, req.Currency) {
		// ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	//check valid account
	if !server.validAccount(ctx, req.ToAccountID, req.Currency) {
		// ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//create account
	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)

}

func (server *Server) validAccount(ctx *gin.Context, accountId int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("Account [%d] currency mismatch [%s] vs [%s]", accountId, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}

	return true
}

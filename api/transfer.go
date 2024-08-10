package api

import (
	db "banking_application/db/sqlc"
	"banking_application/token"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountNumber int64  `json:"from_account_number" binding:"required,min=1"`
	ToAccountNumber   int64  `json:"to_account_number" binding:"required,min=1"`
	Amount            int64  `json:"amount" binding:"required,gt=0"`
	Currency          string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, valid := server.validAccount(ctx, req.FromAccountNumber, req.Currency)
	if !valid {
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.UserName != authPayload.Username {
		err := errors.New("from account does not belong to the authenitcated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, valid = server.validAccount(ctx, req.ToAccountNumber, req.Currency)
	if !valid {
		return
	}

	arg := db.TransferTxParams{
		FromAccountNumber: req.FromAccountNumber,
		ToAccountNumber:   req.ToAccountNumber,
		Amount:            req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountNumber int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.AccountNumber, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}

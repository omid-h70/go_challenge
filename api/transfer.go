package api

import (
	"github.com/gin-gonic/gin"
	db "go_challenge/db/sqlc"
	"net/http"
)

type createTransferRequest struct {
	FromAccountID int64 `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64 `json:"to_account_id" binding:"required,min=1"`
	Amount        int64 `json:"amount" binding:"required,gt=0"`
	//Currency      string `json:"currency" binding:"required,oneof=USD EUR CAD"`
	Currency string `json:"currency" binding:"required,currency"` // use custom validator to act instead of hardcoded currency
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var createTransferRequest createTransferRequest
	if err := ctx.ShouldBindJSON(&createTransferRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	/*
			arg := db.CreateTransferTxParams{
				FromAccountID: req.FromAccountID,
				ToAccountID: req.ToAccountID,
				Amount: req.Amount,
			}

		    if !isAccountValid(ctx, request.FromAccountID, req.Currency){
				ctx.JSON(http.StatusUnauthorized, errorResponse(err))
				return
			}

			authPayLoad := ctx.MustGet(authorizationPayLoadKey).(*token.Payload)
			if fromAccount.Owner != authPayLoad.Username {
				err := errors.New("from account doesn't belong to authenticated user")
			}

		   if !isAccountValid(ctx, request.ToAccountID, req.Currency){
				return
			}

			result, err := server.store.TransferTx(ctx, arg)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusOK, result)
	*/
}

func (server *Server) isAccountValid(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	/*
		account, err := server.store.GetAccount(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, account)

		if account.Currency != currency {
			err := fmt.Errorf("account [%d] currency mismatch %s vs %s", accountID, account.Currency, currency)
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return account, false
		}
	*/
	return db.Account{}, false
}

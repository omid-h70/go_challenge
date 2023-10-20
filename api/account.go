package api

import (
	"github.com/gin-gonic/gin"
	db "go_challenge/db/sqlc"
	"go_challenge/token"
	"net/http"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (server *Server) createAccount(ctx *gin.Context) {

	var createAccountReq createAccountRequest
	if err := ctx.ShouldBindJSON(&createAccountReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	/* v1
	arg := db.CreateAccountParams{
		Owner: createAccountReq.Owner,
		Currency: createAccountReq.Currency,
		Balance: 0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
	*/

	/* v2 -------------------------- */

	authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateAccountParams{
		Owner:    authPayLoad.Username,
		Currency: createAccountReq.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		/*
			errCode := db.Error(err)
			if errCode == db.ForeignKeyViolation || errCode == db.UniqueViolation {
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		*/
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
	/* v2 -------------------------- */
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required, min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var getAccountReq getAccountRequest
	if err := ctx.ShouldBindJSON(&getAccountReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	/*
		account, err := server.store.GetAccount(ctx, req.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
		if account.Owner != authPayLoad.userName{
			err := error.New("account doesn't belong to authenticated user")
		    ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, account)
	*/

}

type listAccountRequest struct {
	PageID   int32 `from:"page_id" binding:"required, min=1"`
	PageSize int32 `from:"page_size" binding:"required, min=5, max=10"`
}

func (server *Server) getAccountList(ctx *gin.Context) {
	var listAccountReq listAccountRequest

	/*
		arg := db.ListAccountsParams{
			Limit: listAccountReq.PageSize
			Offset: (listAccountReq.PageId - 1)*listAccountReq.PageSize
		}
	*/

	//Check For Query Parameters, as it is a GET request with additional params
	if err := ctx.ShouldBindQuery(&listAccountReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	/*
		// ------------- v2 After Authorization Added
		authPayLoad := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
		arg := db.ListAccountsParams{
			Owner:authPayload.Username,
			Limit: req.PageSize,
			Offset: (req.PageID -1) * req.PageSize,
		}
		// ------------- v2

		account, err := server.store.ListAccounts(ctx, req.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, account)
	*/

}

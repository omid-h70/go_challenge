package api

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewAccessToken(ctx *gin.Context) {
	var renewReq renewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&renewReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	payload, err := server.tokenMaker.VerifyToken(renewReq.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	session, err := server.store.GetSession(ctx, payload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, err)
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if session.IsBlocked {
		err := errors.Errorf("session is blocked")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if session.UserName != payload.Username {
		err := errors.Errorf("incorrect session user")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if session.RefreshToken != renewReq.RefreshToken {
		err := errors.Errorf("mismatched refresh token")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		err := errors.Errorf("session is expired")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	accessToken, accessPayLoad, err := server.tokenMaker.CreateToken(
		payload.Username,
		server.config.AccessTokenDuration,
	)

	if err != nil {
		err := errors.Errorf("session is expired")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayLoad.ExpiresAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}

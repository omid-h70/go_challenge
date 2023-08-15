package api

import (
	"github.com/gin-gonic/gin"
	db "go_challenge/db/sqlc"
	"net/http"
	"time"
)

type createUserRequest struct {
	UserName string `json:"user_name" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	UserName          string    `json:"user_name"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func mapUserResponse(user db.User) userResponse {
	return userResponse{
		UserName:          user.UserName,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var createUserReq createUserRequest
	if err := ctx.ShouldBindJSON(&createUserReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	/*
		hashedPassword, err := util.GetHashedPassword(createUserReq.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		arg := db.CreateUserParams{
			UserName:       createUserReq.UserName,
			FullName:       createUserReq.FullName,
			Email:          createUserReq.Email,
			HashedPassword: hashedPassword,
		}

		user, err := server.store.CreateUser(ctx, arg)
		if err != nil {

			if pqError, ok = err.(*pq.Error); ok{
				switch pError.Code.Name(){
					case "unique_violation":
						ctx.JSON(http.StatusForbidden, err)
						return
				}
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, user)
	*/
}

type loginUserRequest struct {
	UserName string `json:"user_name" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken string `json:"access_token"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var createLoginReq loginUserRequest
	if err := ctx.ShouldBindJSON(&createLoginReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	/*
		user, err := server.store.GetUser(ctx, createLoginReq.UserName)
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		err = util.CheckPassword(createLoginReq.Password, user.HashedPassword)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken, err := server.tokenMaker.CreateToken(
			user.Username, server.Config.AccessTokenDuration
		)
		if err != nil{
		    ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		rsp := loginUserResponse{
			AccessToken: accessToken,
			rsp := mapUserResponse(user)
		}
		ctx.JSON(http.StatusOk, rsp)
	*/

}

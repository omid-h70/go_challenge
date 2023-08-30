package gapi

import (
	"context"
	"database/sql"
	db "go_challenge/db/sqlc"
	"go_challenge/pb"
	"go_challenge/util"
	customValidator "go_challenge/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {

	violations := ValidateLoginUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	user, err := s.store.GetUser(ctx, req.GetUserName())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find user")
	}

	err = util.CheckPassword(req.GetPassword(), user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "incorrect password")
	}

	//Create Access Token
	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
		user.UserName, s.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token")
	}

	//Create Session Token
	refreshToken, refreshPayLoad, err := s.tokenMaker.CreateToken(
		user.UserName, s.config.RefreshTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "passwords doesn't match  %s", err)
	}

	session, err := s.store.CreateSession(ctx, db.CreateSessionParams{
		SessionUuid:  refreshPayLoad.ID,
		UserName:     user.UserName,
		RefreshToken: refreshToken,
		UserAgent:    s.extractMetaData(ctx).UserAgent,
		ClientIp:     s.extractMetaData(ctx).ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayLoad.ExpiresAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session")
	}

	rsp := &pb.LoginUserResponse{
		User:                  convertUser(user),
		SessionId:             session.SessionUuid.String(),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiresAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayLoad.ExpiresAt),
	}
	return rsp, nil
}

func ValidateLoginUserRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := customValidator.ValidateUserName(req.GetUserName()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := customValidator.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return
}

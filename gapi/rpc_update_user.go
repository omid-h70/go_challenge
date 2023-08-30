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
	"time"
)

func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {

	payload, err := s.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := ValidateUpdateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	if payload.Username != req.GetUserName() {
		return nil, status.Errorf(codes.PermissionDenied, "cant update other users info")
	}

	arg := db.UpdateUserParams{
		UserName: req.GetUserName(),
		FullName: sql.NullString{
			Valid:  req.FullName != nil,
			String: req.GetFullName(),
		},
		Email: sql.NullString{
			Valid:  req.Email != nil,
			String: req.GetEmail(),
		},
	}

	if req.Password != nil {
		hashedPassword, err := util.GetHashedPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password %s", err)
		}
		arg.HashedPassword = sql.NullString{
			Valid:  true,
			String: hashedPassword,
		}

		arg.PasswordChangedAt = sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		}
	}

	user, err := s.store.UpdateUser(ctx, arg)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to create user %s", err)
	}
	rsp := &pb.UpdateUserResponse{
		User: convertUser(user),
	}
	return rsp, nil
}

func ValidateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := customValidator.ValidateUserName(req.GetUserName()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if req.Password != nil {
		if err := customValidator.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, fieldViolation("password", err))
		}
	}

	if req.FullName != nil {
		if err := customValidator.ValidateFullName(req.GetFullName()); err != nil {
			violations = append(violations, fieldViolation("full_name", err))
		}
	}

	if req.Email != nil {
		if err := customValidator.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, fieldViolation("email", err))
		}
	}
	return
}

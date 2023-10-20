package gapi

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/lib/pq"
	db "go_challenge/db/sqlc"
	"go_challenge/pb"
	"go_challenge/util"
	customValidator "go_challenge/validator"
	"go_challenge/worker"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	violations := ValidateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hashedPassword, err := util.GetHashedPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password %s", err)
	}

	arg := db.CreateUserParams{
		UserName:       req.GetUserName(),
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
		HashedPassword: hashedPassword,
	}

	//-------------- We must put it in one single db transaction
	// what if one fails ? duplicated user or what ?
	// you can manually bring  redis image down to see what will happen, wtf man
	// the problem is happening in different form
	// what if commit takes long time ?
	// what if creating user was successful - task added and yet we're on high traffic server and commits took longer than expected ?
	// retry mechaniem in queue will always help you !
	// also consider some delay !

	txArgs := db.CreateUserTxParams{
		CreateUserParams: arg,
		AfterCreate: func(user db.User) error {
			tskPayLoad := &worker.PayLoadSendVerifyEmail{
				UserName: arg.UserName,
			}

			opt := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second), // make a delay about 10 seconds
				asynq.Queue("critical"),           // =>send it to critical queue
			}

			return s.taskDistributor.DistributeTaskSendVerifyEmail(ctx, tskPayLoad, opt...)
		},
	}

	txnResult, err := s.store.CreateUserTx(ctx, txArgs)
	if err != nil {

		if pqError, ok := err.(*pq.Error); ok {
			switch pqError.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "user already exist %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user %s", err)
	}
	//====================== V1
	//tskPayLoad := &worker.PayLoadSendVerifyEmail{
	//	UserName: arg.UserName,
	//}
	//
	//opt := []asynq.Option{
	//	asynq.MaxRetry(10),
	//	asynq.ProcessIn(10 * time.Second), // make a delay about 10 seconds
	//	asynq.Queue("critical"),           // =>send it to critical queue
	//}
	//
	//err = s.taskDistributor.DistributeTaskSendVerifyEmail(ctx, tskPayLoad, opt...)
	//if err != nil {
	//	return nil, status.Errorf(codes.Internal, "failed to create user %s", err)
	//}
	//-------------- We must put it in one single db transaction

	rsp := &pb.CreateUserResponse{
		User: convertUser(txnResult.User),
	}
	return rsp, nil
}

func ValidateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := customValidator.ValidateUserName(req.GetUserName()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := customValidator.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	if err := customValidator.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}
	if err := customValidator.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}
	return
}

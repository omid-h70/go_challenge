package worker

import (
	"context"
	"github.com/hibiken/asynq"
	zlog "github.com/rs/zerolog/log"
	db "go_challenge/db/sqlc"
	"go_challenge/mail"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
	mailer mail.EmailSender
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, mailer mail.EmailSender) *RedisTaskProcessor {
	return &RedisTaskProcessor{
		server: asynq.NewServer(redisOpt, asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10, //Set Priority for each queue
				QueueDefault:  5,  //Set Priority for each queue
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				zlog.Error().Err(err).
					Str("type", task.Type()).
					Bytes("payload", task.Payload()).
					Msg("process task failed")
			}),
			Logger: NewLogger(), //add Custom Logger
		}),
		store:  store,
		mailer: mailer,
	}
}

func (p *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerifyEmail, p.ProcessTaskSendVerifyEmail)
	return p.server.Start(mux)
}

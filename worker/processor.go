package worker

import (
	"context"
	"github.com/hibiken/asynq"
	db "go_challenge/db/sqlc"
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
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) TaskProcessor {
	return &RedisTaskProcessor{
		server: asynq.NewServer(redisOpt, asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10, //Set Priority for each queue
				QueueDefault:  5,  //Set Priority for each queue
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *Task, err error) {
				log.Error().Msg("process Task failed")
				.Str("type", task.Type())
				.Bytes("payload", task.PayLoad())
				.Msg("process task failed")
			}),
			Logger: NewLogger(),//add Custom Logger
		}),
		store: store,
	}
}

func (p *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerifyEmail, p.ProcessTaskSendVerifyEmail)
	return p.server.Start(mux)
}

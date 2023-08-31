package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayLoadSendVerifyEmail struct {
	UserName string `json:"user_name"`
}

func (r RedisTaskDistributor) DistributeTaskSendVerifyEmail(ctx context.Context, payload *PayLoadSendVerifyEmail, opts ...asynq.Option) error {

	jsonPayLoad, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("unable to marshal payload %w", err)
	}

	tsk := asynq.NewTask(TaskSendVerifyEmail, jsonPayLoad, opts...)
	tskInfo, err := r.client.EnqueueContext(ctx, tsk)
	if err != nil {
		return fmt.Errorf("failed to enqueue task %w", err)
	}

	log.Info().Str("type", tsk.Type()).Bytes("payload", tsk.Payload()).
		Str("queue", tskInfo.Queue).Int("max_retry", tskInfo.MaxRetry).Msg("enqueued task")

	return nil
}

func (p *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, tsk *asynq.Task) error {
	var payload PayLoadSendVerifyEmail
	if err := json.Unmarshal(tsk.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload %w", asynq.SkipRetry)
	}

	user, err := p.store.GetUser(ctx, payload.UserName)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user doesnt exist %w", asynq.SkipRetry)
		}
		return fmt.Errorf("failed to get user %w", err)
	}

	log.Info().Str("type", tsk.Type()).Bytes("payload", tsk.Payload()).
		Str("email", user.Email).Msg("enqueued task")

	return nil
}

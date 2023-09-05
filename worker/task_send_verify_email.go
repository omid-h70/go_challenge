package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	db "go_challenge/db/sqlc"
	"go_challenge/util"
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
		/* Retry it some Other Time !
		if err == sql.ErrNoRows {
			return fmt.Errorf("user doesnt exist %w", asynq.SkipRetry)
		}
		*/
		return fmt.Errorf("failed to get user %w", err)
	}

	verifyEmail, err := p.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		Username:   user.UserName,
		Email:      user.Email,
		SecretCode: util.RandomString(32),
	})

	if err != nil {
		return fmt.Errorf("failed to create verify email %w", err)
	}

	subject := "welcome to simple bank"
	to := []string{""}
	verifyUrl := fmt.Sprintf("http://simple-bank.org?id=%d&secred_code=%s", verifyEmail.ID, verifyEmail.SecretCode)
	content := fmt.Sprintf(`
	<h1> wola wola </h1> 
	<a href="%s"> </a>
	`, verifyUrl)
	err = p.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email %w", err)
	}

	log.Info().Str("type", tsk.Type()).Bytes("payload", tsk.Payload()).
		Str("email", user.Email).Msg("enqueued task")

	return nil
}

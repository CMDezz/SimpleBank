package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	db "github.com/techschool/simplebank/db/sqlc"
)

const (
	CriticalQueue = "critical"
	DefaultQueue  = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessSendTaskVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) TaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{

		Queues: map[string]int{
			CriticalQueue: 10,
			DefaultQueue:  5,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Error().Err(err).Str("type :", task.Type()).Bytes("payload :", task.Payload()).Msg("failed to process task")
		}),
		Logger: NewLogger(),
	})
	return &RedisTaskProcessor{
		server: server,
		store:  store,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessSendTaskVerifyEmail)

	err := processor.server.Start(mux)
	if err != nil {
		return fmt.Errorf("cannot start processor %s: %s", TaskSendVerifyEmail, err)
	}

	return nil
}

func (processor *RedisTaskProcessor) ProcessSendTaskVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user does not exist: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("failed to get user %w", err)
	}
	//send email here
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("email", user.Email).Msg("proccessed task")
	return nil
}

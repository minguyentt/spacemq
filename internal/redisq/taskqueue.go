package redisq

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	DefaultNameSpace = "spacemq"
)

var (
	//go:embed scripts/enqueue.lua
	enqScript  string
	enqueueCmd = redis.NewScript(enqScript)

	//go:embed scripts/dequeue.lua
	deqScript  string
	dequeueCmd = redis.NewScript(deqScript)
)

type task struct {
	ID          string `redis:"id"`
	QueueID     string `redis:"queue_id"`
	Payload     []byte `redis:"payload"`
	MaxAttempts int    `redis:"max_attempts"`
	Attempts    int    `redis:"attempts"`
	CreatedAt   int64  `redis:"created_at"`
	CompletedAt int64  `redis:"completed_at"`
	Timeout     int64  `redis:"timeout"`
	TTL         int64  `redis:"ttl"`
}

func TaskKey(q string, taskID string) string {
	return fmt.Sprintf("%s:%s:task-id:%s", DefaultNameSpace, q, taskID)
}

func TaskKeyPrefix(qID string) string {
	return fmt.Sprintf("%s:%s:task-id:", DefaultNameSpace, qID)
}

func PendingKey(qID string) string {
	return fmt.Sprintf("%s:%s:pending", DefaultNameSpace, qID)
}

func ActiveKey(qID string) string {
	return fmt.Sprintf("%s:%s:active", DefaultNameSpace, qID)
}

type taskRunner struct {
	client redis.UniversalClient
}

func NewTaskRunner(client redis.UniversalClient) *taskRunner {
	return &taskRunner{client: client}
}

// spacemq:queueID:task-id:<task-id>
func (t *taskRunner) Enqueue(ctx context.Context, task *task) error {
	keys := []string{
		TaskKey(task.QueueID, task.ID),
		PendingKey(task.QueueID),
	}

	argv := []any{
		task.Payload,
		task.ID,
		time.Now().UnixNano(),
	}

	err := enqueueCmd.Run(ctx, t.client, keys, argv...).Err()
	if err != nil {
		return fmt.Errorf("redis eval err: %v", err)
	}

	return nil
}

func (t *taskRunner) Dequeue(ctx context.Context, queueID string) (*task, error) {
	// create the keys and argv to satisfy the dequeue commander
	// run the script
	// handle errors with custom error types
	// return the updated task info
}

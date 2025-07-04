package redisdbq

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
	spacemq "github.com/minguyentt/spaceMQ"
	"github.com/redis/go-redis/v9"
)

var (
	//go:embed scripts/enqueue.lua
	enqScript  string
	enqueueCmd = redis.NewScript(enqScript)

	//go:embed scripts/dequeue.lua
	deqScript  string
	dequeueCmd = redis.NewScript(deqScript)

	//go:embed scripts/acknowledge.lua
	ackScript   string
	ackQueueCmd = redis.NewScript(ackScript)

	//go:embed scripts/retryqueue.lua
	retryScript   string
	retryQueueCmd = redis.NewScript(retryScript)

	//go:embed scripts/requeue.lua
	requeScript string
	requeueCmd  = redis.NewScript(requeScript)

	//go:embed scripts/deadqueue.lua
	dlqScript    string
	deadQueueCmd = redis.NewScript(dlqScript)
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type (
	Timer   interface{ Now() time.Time }
	rdbTime struct{}
)

func NewTimer() *rdbTime          { return &rdbTime{} }
func (r *rdbTime) Now() time.Time { return time.Now() }

type Runner interface {
	Enqueue(ctx context.Context) error
	Dequeue(queueName string, infDur time.Duration) (*Task, error)
	Acknowledge(ctx context.Context)
	Retry()
	MoveToDeadQueue()

	Ping() error
	Close() error
	ClientInfo() string
}

type taskRunner struct {
	*DBQ
	Timer
}

type DBQ struct {
	client redis.UniversalClient
}

func New(c redis.UniversalClient) *DBQ {
	return &DBQ{client: c}
}

func (r *taskRunner) Close() error {
	return r.client.Close()
}

func (r *taskRunner) Ping() error {
	return r.client.Ping(context.Background()).Err()
}

func (r *taskRunner) Info() string {
	return r.client.ClientInfo(context.Background()).String()
}

func (r *taskRunner) GetClient() redis.UniversalClient {
	return r.client
}

type taskMeta struct {
	ID          string `redis:"id"           json:"id,omitempty"`
	QueueName   string `redis:"queue_name"   json:"queue_name,omitempty"`
	Payload     []byte `redis:"payload"      json:"payload,omitempty"`
	MaxAttempts uint8  `redis:"max_attempts" json:"max_attempts,omitempty"`
	Attempts    uint8  `redis:"attempts"     json:"attempts,omitempty"`
	State       string `redis:"state"        json:"state,omitempty"`

	CreatedAt   int64 `redis:"created_at"   json:"created_at,omitempty"`
	StartedAt   int64 `redis:"started_at"   json:"started_at,omitempty"`
	CompletedAt int64 `redis:"completed_at" json:"completed_at,omitempty"`

	FailedAt int64  `redis:"failed_at" json:"failed_at,omitempty"`
	FailedBy string `redis:"failed_by" json:"failed_by,omitempty"`

	LastFailedAt int64  `redis:"last_failed_at" json:"last_failed_at,omitempty"`
	LastFailedBy string `redis:"last_failed_by" json:"last_failed_by,omitempty"`

	Timeout int64 `redis:"timeout" json:"timeout,omitempty"`
	TTL     int64 `redis:"ttl"     json:"ttl,omitempty"`
}

type Task struct {
	ID          string
	QueueName   string
	Payload     []byte
	MaxAttempts uint8
	Attempts    uint8
	State       string

	CreatedAt   int64
	StartedAt   int64
	CompletedAt int64

	FailedAt int64
	FailedBy string

	LastFailedAt int64
	LastFailedBy string

	Timeout int64
	TTL     int64
}

func encodeTask(t *Task) ([]byte, error) {
	meta := &taskMeta{
		ID:           t.ID,
		Payload:      t.Payload,
		QueueName:    t.QueueName,
		MaxAttempts:  t.MaxAttempts,
		Attempts:     t.Attempts,
		State:        t.State,
		CreatedAt:    t.CreatedAt,
		CompletedAt:  t.CompletedAt,
		FailedAt:     t.FailedAt,
		FailedBy:     t.FailedBy,
		LastFailedAt: t.LastFailedAt,
		LastFailedBy: t.LastFailedBy,
		Timeout:      t.Timeout,
		TTL:          t.TTL,
	}

	b, err := json.Marshal(meta)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func decode(data []byte) (*Task, error) {
	var meta taskMeta
	err := json.Unmarshal(data, &meta)
	if err != nil {
		return nil, err
	}

	t := &Task{
		ID:           meta.ID,
		Payload:      meta.Payload,
		QueueName:    meta.QueueName,
		MaxAttempts:  meta.MaxAttempts,
		Attempts:     meta.Attempts,
		State:        meta.State,
		CreatedAt:    meta.CreatedAt,
		CompletedAt:  meta.CompletedAt,
		FailedAt:     meta.FailedAt,
		FailedBy:     meta.FailedBy,
		LastFailedAt: meta.LastFailedAt,
		LastFailedBy: meta.LastFailedBy,
		Timeout:      meta.Timeout,
		TTL:          meta.TTL,
	}

	return t, nil
}

func NewTaskRunner(rdbq *DBQ) *taskRunner {
	return &taskRunner{
		rdbq,
		NewTimer(),
	}
}

func (r *taskRunner) Enqueue(ctx context.Context, task *Task) error {
	enc, err := encodeTask(task)
	if err != nil {
		return err
	}

	keys := []string{
		FullTaskKey(task.QueueName, task.ID),
		PendingKey(task.QueueName),
	}

	argv := []any{
		task.ID,
		r.Now().Unix(),
		enc,
	}

	err = enqueueCmd.Run(ctx, r.client, keys, argv...).Err()
	if err != nil {
		return fmt.Errorf("redis eval err: %v", err)
	}

	return nil
}

func (r *taskRunner) Dequeue(queueName string, infDur time.Duration) (*Task, error) {
	inflightDuration := r.Now().Add(infDur)

	keys := []string{
		PendingKey(queueName),
		ActiveKey(queueName),
		PausedKey(queueName),
		InFlightKey(queueName),
	}

	argv := []any{
		TaskKeyPrefix(queueName),
		r.Now().Unix(),
		inflightDuration,
	}

	res, err := dequeueCmd.Run(context.Background(), r.client, keys, argv...).Slice()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, fmt.Errorf("rdb returned empty array from dequeue")
	}

	data := make(map[string]interface{})
	for i := 0; i < len(res); i += 2 {
		data[res[i].(string)] = data[res[i+1].(string)]
	}

	task, err := decode(data["meta"].([]byte))
	if err != nil {
		return nil, err
	}

	task.State = data["state"].(string)
	task.StartedAt = data["started_at"].(int64)

	return task, nil
}

func (r *taskRunner) Acknowledge(ctx context.Context, job *spacemq.Job) error {
	keys := []string{
		ActiveKey(job.QueueName()),
		InFlightKey(job.QueueName()),
		FullTaskKey(job.QueueName(), job.TaskID()),
	}

	argv := []any{
		job.TaskID(),
		r.Now().Unix(),
		job.GetTTL(),
	}

	err := ackQueueCmd.Run(ctx, r.client, keys, argv...).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *taskRunner) Retry(
	ctx context.Context,
	task *Task,
	retryTimer time.Time,
	failMsg string,
) error {
	keys := []string{}

	argv := []any{}

	return retryQueueCmd.Run(ctx, r.client, keys, argv...).Err()
}

func (r *taskRunner) Requeue(ctx context.Context, job *spacemq.Job, batchLimit int) error {
	keys := []string{
		RetryKey(job.QueueName()),
		PendingKey(job.QueueName()),
		FullTaskKey(job.QueueName(), job.TaskID()),
	}

	argv := []any{
		r.Now().Unix(),
		int64(batchLimit),
	}

	return requeueCmd.Run(ctx, r.client, keys, argv...).Err()
}

func (r *taskRunner) MoveToDeadQueue(ctx context.Context, job *spacemq.Job) error {
	return nil
}

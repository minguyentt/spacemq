package redisdbq

import (
	"context"
	"crypto/rand"
	j "encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

var defaultTTL = 10 * time.Minute

type payloadTest map[string]any

func generateID() string {
	return ulid.MustNew(ulid.Now(), rand.Reader).String()
}

func marsh(p payloadTest) []byte {
	b, _ := j.Marshal(p)
	return b
}

func mockEncoderData(id string, data map[string]any) []struct {
	in  *Task
	out *Task
} {
	return []struct {
		in  *Task
		out *Task
	}{
		{
			in: &Task{
				ID:           id,
				QueueName:    "test_queue",
				Payload:      marsh(data),
				MaxAttempts:  5,
				Attempts:     0,
				State:        "pending",
				CreatedAt:    1751605477,
				CompletedAt:  1751580307,
				FailedAt:     1751580307,
				FailedBy:     "worker1",
				LastFailedAt: 1751580307,
				LastFailedBy: "worker2",
				Timeout:      2000,
				TTL:          10000,
			},
			out: &Task{
				ID:           id,
				QueueName:    "test_queue",
				Payload:      marsh(data),
				MaxAttempts:  5,
				Attempts:     0,
				State:        "pending",
				CreatedAt:    1751605477,
				CompletedAt:  1751580307,
				FailedAt:     1751580307,
				FailedBy:     "worker1",
				LastFailedAt: 1751580307,
				LastFailedBy: "worker2",
				Timeout:      2000,
				TTL:          10000,
			},
		},
	}
}

func newSimclient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:7777",
		DB:   7,
	})
}

func setupTaskRunner(c *redis.Client) *taskRunner {
	// client + rdb + in runner
	return NewTaskRunner(New(c))
}

func flush(tb testing.TB, c redis.UniversalClient) {
	tb.Helper()

	if err := c.FlushDB(context.Background()).Err(); err != nil {
		tb.Fatalf("rdb: flush error: %v", err)
	}
}

func TestTaskEncoder(t *testing.T) {
	id := generateID()

	pl := payloadTest{
		"user_id":     2,
		"foo":         "bar",
		"hello_world": "hey there",
	}

	tt := mockEncoderData(id, pl)

	for _, s := range tt {
		enc, err := encodeTask(s.in)
		if err != nil {
			t.Errorf("Encoding Task: returned error: %v", err)
			continue
		}

		dec, err := decode(enc)
		if err != nil {
			t.Errorf("Decoding Task: returned error: %v", err)
			continue
		}

		if diff := cmp.Diff(s.out, dec); diff != "" {
			t.Errorf("(diff) Decoded message == %+v, want %+v;(-want,+got)\n%s",
				dec, s.out, diff)
		}
	}
}

func TestEnqueueScript(t *testing.T) {
	client := client()
	id := generateID()

	tq := setupTaskRunner(client)
	flush(t, tq.db.client)
	fmt.Println("(TestEnqueueScript) FLUSHING REDIS DATABASE...")

	data := map[string]any{
		"to":      "some_token",
		"from":    "task_queue_1",
		"message": "enqueue test notification",
	}

	in := &Task{
		ID:          id,
		QueueName:   "task_queue",
		Payload:     marsh(data),
		MaxAttempts: 7,
		Attempts:    0,
		Timeout:     1751580307,
		TTL:         tq.sign.Now().Add(defaultTTL).Unix(),
	}

	out := &Task{
		ID:          id,
		QueueName:   "task_queue",
		Payload:     marsh(data),
		MaxAttempts: 7,
		Attempts:    0,
		Timeout:     1751580307,
		TTL:         tq.sign.Now().Add(defaultTTL).Unix(),
	}

	enc, err := encodeTask(in)
	assert.NoErrorf(t, err, "encoding in error", "output", enc)

	if err := tq.Enqueue(context.Background(), in); err != nil {
		t.Fatalf("Task Runner Enqueue script error: got %v, want nil", err)
	}

	// check pending queue list has taskid
	pqKey := PendingKey(in.QueueName)
	t.Run("should have one pending queue in list", func(t *testing.T) {
		t.Helper()
		n := tq.db.client.LLen(context.Background(), pqKey).Val()
		assert.Equal(t, int64(1), n, "(REDIS LIST) length of pending queue list should be 1")
	})
	t.Run("should retrieve the val as ID in list", func(t *testing.T) {
		t.Helper()
		val := tq.db.client.LRange(context.Background(), pqKey, 0, -1).Val()
		assert.Equal(t, in.ID, val[0], "(REDIS LIST) id does not exist in the pending queue list")
	})

	// retrieve the meta data and check diff
	taskKey := FullTaskKey(in.QueueName, in.ID)
	meta := tq.db.client.HGet(context.Background(), taskKey, "meta").Val()

	dec, _ := decode([]byte(meta))
	if diff := cmp.Diff(out, dec); diff != "" {
		t.Errorf("mismatch meta data (-want, +got):\n%s", diff)
	}

	stateVal := tq.db.client.HGet(context.Background(), taskKey, "state").Val()
	assert.Equal(t, "pending", stateVal, "(REDIS LIST) should set state to pending")
}

func TestDequeueScript(t *testing.T) {
	client := client()
	id := generateID()

	tq := setupTaskRunner(client)
	flush(t, tq.db.client)
	fmt.Println("(TestEnqueueScript) FLUSHING REDIS DATABASE...")
}
// TODO: implement tests for...
// acknowledge
// deadqueue
// dequeue
// requeue
// retryqueue

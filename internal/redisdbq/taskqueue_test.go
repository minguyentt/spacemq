package redisdbq

import (
	"crypto/rand"
	j "encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/oklog/ulid/v2"
)

type payloadTest map[string]interface{}

func genTestID() string {
	return ulid.MustNew(ulid.Now(), rand.Reader).String()
}

func marsh(p payloadTest) []byte {
	b, _ := j.Marshal(p)
	return b
}

func TestTaskEncoder(t *testing.T) {
	id := genTestID()

	pl := payloadTest{
		"user_id":     2,
		"foo":         "bar",
		"hello_world": "hey there",
	}

	tt := []struct {
		in  *Task
		out *Task
	}{
		{
			in: &Task{
				ID:           id,
				QueueName:    "test_queue",
				Payload:      marsh(pl),
				MaxAttempts:  5,
				Attempts:     0,
				State:        "pending",
				CreatedAt:    1751605477,
				StartedAt:    1751580307,
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
				Payload:      marsh(pl),
				MaxAttempts:  5,
				Attempts:     0,
				State:        "pending",
				CreatedAt:    1751605477,
				StartedAt:    1751580307,
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
			t.Errorf("Decoded message == %+v, want %+v;(-want,+got)\n%s",
				dec, s.out, diff)
		}
	}
}

func TestRedisEnqueueCommand(t *testing.T) {
	// create redis client using port 7777
	// construct the taskqueue
	// enqueue the test table
}

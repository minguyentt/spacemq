package spacemq

import (
	"fmt"
	"time"
)

type Job struct {
	queueID     string
	payload     []byte
	maxAttempts int
	timeout     time.Duration
	ttl         time.Duration
}

type JobOpts struct {
	MaxAttempts int
	Timeout     time.Duration
	TTL         time.Duration
}

type JobInfo struct {
	ID          string
	QueueID     string
	Payload     []byte
	State       JobState
	MaxAttempts int
	Attempts    int
	CreatedAt   time.Time
	CompletedAt time.Time
	// want durationin minutes from creation to job completion
}

type JobState int8

const (
	JobStatePending JobState = iota + 1
	JobStateActive
	JobStateCompleted
	JobStateFailed
	JobStateDead
)

func (s JobState) String() string {
	switch s {
	case JobStatePending:
		return "pending"
	case JobStateActive:
		return "active"
	case JobStateCompleted:
		return "completed"
	case JobStateFailed:
		return "failed"
	case JobStateDead:
		return "dead"
	}

	return ""
}

// NOTE:
// * Handling the job creation will be the phase 1 of building the
// 	necessary requirements to start the queue sequence
//
// * Also will contain JobInfo structure to retrieve cross-communication responses
// from the business to app layer
func NewJob(group, qname string, payload []byte, opts JobOpts) *Job {
	queueID := fmt.Sprintf("%s:%s", group, qname)

	return &Job{
		queueID:     queueID,
		payload:     payload,
		maxAttempts: opts.MaxAttempts,
		timeout:     opts.Timeout,
		ttl:         opts.TTL,
	}
}

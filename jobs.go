package spacemq

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
)

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

const (
	DefaultMaxAttempts      = 25
	DefaultTTL              = 24 * time.Hour * 3 // time to live
	DefaultInFlightDuration = 30 * time.Second   // computing inflight operations
	DefaultTimeout          = 30 * time.Second   // used to compute enqueue operations
)

type Job struct {
	id               string
	queueName        string
	payload          []byte
	maxAttempts      int
	inflightDuration time.Duration
	timeout          time.Duration
	ttl              time.Duration
}

type JobOpts func(*Job)

func WithQueueName(s string) JobOpts              { return func(j *Job) { j.queueName = s } }
func WithMaxAttempts(n int) JobOpts               { return func(j *Job) { j.maxAttempts = n } }
func WithTimeout(d time.Duration) JobOpts         { return func(j *Job) { j.timeout = d } }
func WithTimeToLive(d time.Duration) JobOpts      { return func(j *Job) { j.ttl = d } }
func WithInFlightTimeout(d time.Duration) JobOpts { return func(j *Job) { j.inflightDuration = d } }

func NewJob(payload []byte, opts ...JobOpts) (*Job, error) {
	job := Job{
		id:               ulid.MustNew(ulid.Now(), rand.Reader).String(),
		payload:          payload,
		maxAttempts:      DefaultMaxAttempts,
		inflightDuration: DefaultInFlightDuration,
		timeout:          DefaultTimeout,
		ttl:              DefaultTTL,
	}

	if len(job.queueName) == 0 {
		return nil, fmt.Errorf("queue name must be greater than 0")
	}

	for _, opt := range opts {
		opt(&job)
	}

	return &job, nil
}

func (j *Job) TaskID() string { return j.id }

func (j *Job) QueueName() string { return j.queueName }

func (j *Job) GetINFDur() time.Duration { return j.inflightDuration }

func (j *Job) GetTimeout() time.Duration { return j.timeout }

func (j *Job) GetTTL() time.Duration { return j.ttl }

func (j *Job) MaxAttempts() int { return j.maxAttempts }

type JobInfo struct {
	ID          string
	QueueName   string
	Payload     []byte
	State       JobState
	MaxAttempts int
	Attempts    int

	CreatedAt   time.Time
	StartedAt   time.Time
	CompletedAt time.Time

	FailedAt time.Time
	FailedBy string

	LastFailedAt time.Time
	LastFailedBy string

	Timeout time.Duration
	TTL     time.Duration
}

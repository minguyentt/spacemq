package spacemq

import "github.com/redis/go-redis/v9"

type TaskQueue struct {
	client redis.UniversalClient
	job *Job
}

func NewTaskQueue(client redis.UniversalClient, job *Job) *TaskQueue {
	return &TaskQueue{
		client: client,
		job: job,
	}
}


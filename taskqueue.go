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

// TODO: 07-02-25
//	need client class to construct the Universal Client options.
//  figure out what other configs is needed for this client
//
//	TaskQueue class will inherit the client also the components to build the taskqueue.
//	Jobs, Workers, Custom error type, task queue runner ...
//
//	TaskQueue:
//	+Enqueue(ctx, job ...) error
//	...
//
//	This will be the external implementation where it'll modify catered to specific requirements for
//	task queue runner.


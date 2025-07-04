package redisdbq

import "fmt"

const (
	DefaultNameSpace = "spacemq"
)
// spacemq:<queueName>:task-id:<task-id>
func FullTaskKey(q string, taskID string) string {
	return fmt.Sprintf("%s:%s:task-id:%s", DefaultNameSpace, q, taskID)
}

// spacemq:<queueName>:pending
func PendingKey(queueName string) string {
	return fmt.Sprintf("%s:%s:pending", DefaultNameSpace, queueName)
}

// spacemq:<queueName>:active
func ActiveKey(queueName string) string {
	return fmt.Sprintf("%s:%s:active", DefaultNameSpace, queueName)
}

// spacemq:<queueName>:paused
func PausedKey(queueName string) string {
	return fmt.Sprintf("%s:%s:paused", DefaultNameSpace, queueName)
}

func InFlightKey(queueName string) string {
	return fmt.Sprintf("%s:%s:in-flight", DefaultNameSpace, queueName)
}

func RetryKey(queueName string) string {
	return fmt.Sprintf("%s:%s:retry", DefaultNameSpace, queueName)
}

func DeadQueueKey(queueName string) string {
	return fmt.Sprintf("%s:%s:dead", DefaultNameSpace, queueName)
}

func FailedKey(queueName string) string {
	return fmt.Sprintf("%s:%s:failed", DefaultNameSpace, queueName)
}

// spacemq:<queueName>:task-id:
func TaskKeyPrefix(queueName string) string {
	return fmt.Sprintf("%s:%s:task-id:", DefaultNameSpace, queueName)
}

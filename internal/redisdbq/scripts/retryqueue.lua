local activeQueueKey = KEYS[1]
local inflightQueueKey = KEYS[2]
local retryQueueKey = KEYS[3]
local taskKey = KEYS[5]

local taskID = ARGV[1]
local timeNowUnix = ARGV[2]
local retryAfter_score = ARGV[3]
local failedReason = ARGV[5]

redis.call("LREM", activeQueueKey, 0, taskID)
redis.call("ZREM", inflightQueueKey, taskID)

-- add to retry queue
redis.call("ZADD", retryQueueKey, retryAfter_score, taskID)
redis.call("HSET", taskKey, "state", "retry", "last_failed_at", timeNowUnix, "last_failed_by", failedReason)

return redis.status_reply("RETRY_SCHEDULED")

local queuePendingKey = KEYS[1]
local queueActiveKey = KEYS[2]
local queuePausedKey = KEYS[3]
local inflightKey = KEYS[4]

local taskPrefixKey = ARGV[1]
local timeNow = ARGV[2]
local flightDuration = ARGV[3]

-- if pending key exists, then queue is paused, return nil
if redis.call("EXISTS", queuePausedKey) == 1 then
	return nil
end

-- LMOVE source destination <LEFT | RIGHT> <LEFT | RIGHT>
-- pops the last element<taskid> in the queue list and push it to queue active list
local taskID = redis.call("LMOVE", queuePendingKey, queueActiveKey, "RIGHT", "LEFT")
if not taskID then
	return nil
end

local fullTaskKey = taskPrefixKey .. taskID

-- set state to active and append the time of dequeue
-- NOTE: for latency observability
redis.call("HSET", fullTaskKey, "state", "active", "started_at", timeNow)
redis.call("HDEL", fullTaskKey, "pending_since")

-- add the taskid to the inflight queue in a sorted set (ZADD) with expiration timestamp in unix secs.
-- NOTE: have the worker be configurable setting the expiration
redis.call("ZADD", inflightKey, flightDuration)

return {
	redis.call("HMGET", fullTaskKey, "meta", "state", "started_at"),
}

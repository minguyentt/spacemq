local activeQueueKey = KEYS[1]
local inflightKey = KEYS[2]
local fullTaskKey = KEYS[3]

local taskID = ARGV[1]
local timeNowUnix = ARGV[2]
local taskTTLSecs = tonumber(ARGV[3])

-- remove taskid from active queue
if redis.call("LREM", activeQueueKey, 0, taskID) == 0 then
    return redis.error_reply("ERR TASKID NOT FOUND IN ACTIVE QUEUE")
end

-- remove taskid from inflight queue
if redis.call("ZREM", inflightKey, taskID) == 0 then
    return redis.error_reply("ERR TASKID NOT FOUND IN IN-FLIGHT QUEUE")
end

-- set task meta data state to complete
redis.call("HSET", fullTaskKey, "state", "completed", "completed_at", timeNowUnix)

if tonumber(taskTTLSecs) and tonumber(taskTTLSecs) > 0 then
    redis.call("EXPIRE", fullTaskKey, tonumber(taskTTLSecs))
end

return redis.status_reply("ACKNOWLEDGED")

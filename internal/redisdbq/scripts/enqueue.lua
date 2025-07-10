-- enqueue commander
local fullTaskKey = KEYS[1]
local queuePendingKey = KEYS[2]

local taskid = ARGV[1]
local timeNow = ARGV[2]
local metadata = ARGV[3]

-- check if taskid exists in the queue
if redis.call("EXISTS", fullTaskKey) == 1 then
	return 0
end

redis.call("HSET", fullTaskKey,
    "meta", metadata,
    "state", "pending",
    "pending_since", timeNow)

-- push the taskid to the queue pending list
redis.call("LPUSH", queuePendingKey, taskid)

return 1

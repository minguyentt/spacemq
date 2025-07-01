-- enqueue commander
local taskkey = KEYS[1]
local queuetypeKey = KEYS[2]

local taskid = ARGV[1]
local payload = ARGV[2]
local timeNow = ARGV[3]

-- check if taskid exists in the queue
if redis.call('EXISTS', taskkey) == 1 then
    return 0
end

redis.call(
    "HSET", taskid,
    "payload", payload,
    "state", "pending",
    "task_pending_since", timeNow
)

redis.call("LPUSH", queuetypeKey, taskid)

return 1

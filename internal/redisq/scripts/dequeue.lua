local queue_pending_key = KEYS[1]
local queue_active_key = KEYS[2]
local queue_paused_key = KEYS[3]

local task_prefix_key = ARGV[1]
-- if pending key exists, then queue is paused, return nil
if redis.call("EXISTS", queue_paused_key) == 1 then
	return nil
end

-- pop queue pending item and push the item to the end of the active list
-- LMOVE source destination <LEFT | RIGHT> <LEFT | RIGHT>
local task_id = redis.call("LMOVE", queue_pending_key, queue_active_key, "RIGHT", "LEFT")
local full_key = task_prefix_key .. task_id

redis.call("HSET", full_key, "state", "active")
redis.call("HDEL", full_key, "task_pending_since")

return redis.call("HGET", full_key, "payload")

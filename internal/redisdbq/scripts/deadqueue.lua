local activeKey = KEYS[1]
local inflightKey = KEYS[2]
local deadKey = KEYS[3]
local taskKey = KEYS[4]

local taskID = ARGV[1]
local timeNowUnix = ARGV[2]
local failedReason = ARGV[3]

if redis.call("LREM", activeKey, 0, taskID) == 0 then
	return redis.error_reply("NOT_FOUND_IN_ACTIVE")
end
if redis.call("ZREM", inflightKey, taskID) == 0 then
	return redis.error_reply("NOT_FOUND_IN_LEASE")
end

redis.call("LPUSH", deadKey, taskID)
redis.call("HSET", taskKey, "state", "dead", "failed_at", timeNowUnix, "failed_by", failedReason)

return redis.status_reply("MOVED_TO_DEAD")

local retryKey = KEYS[1]
local pendingKey = KEYS[2]
local taskPrefixKey = KEYS[3]

local currentTimeUnix = ARGV[1]
local batchLimit = tonumber(ARGV[2])

local moved = {}
-- NOTE:
local batch = redis.call("ZRANGEBYSCORE", retryKey, '-inf', currentTimeUnix, 'LIMIT', 0, batchLimit)

for _, id in ipairs(batch) do
    local taskKey = taskPrefixKey .. id

    if redis.call("ZREM", retryKey, id) == 1 then
        redis.call("LPUSH", pendingKey, id)
        redis.call("HSET", taskKey, "state", "pending")

        table.insert(moved, id)
    end
end

return batch


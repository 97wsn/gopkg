local key = KEYS[1]
local burst = tonumber(ARGV[1])
local limitSecond = tonumber(ARGV[2])
local useNum = tonumber(ARGV[3])
local currentNow = tonumber(ARGV[4])

local res = redis.call('HGETALL', key)
if #res == 0 then
    res = {"tokens", burst, "last_time", currentNow }
    redis.call('HMSET', key, unpack(res))
end

local reservation = {}
for i = 1, #res, 2 do
    reservation[res[i]] = tonumber(res[i + 1])
end

local elapsed = currentNow - reservation["last_time"]
local newTokens = reservation["tokens"] + (elapsed / 1000000000.0) * limitSecond
if newTokens > burst then
    newTokens = burst
end

local remainingTokens = newTokens - useNum

if remainingTokens >= 0 then
    redis.call('HMSET', key, "tokens", remainingTokens, "last_time", currentNow)
    return 1
else
    return 0
end
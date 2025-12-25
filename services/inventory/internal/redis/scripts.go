package redisclient

// LockSeatsLua is an atomic Lua script that:
// 1. Checks if any seat key already exists
// 2. If yes → return 0 (failure)
// 3. If no → set all keys with TTL → return 1 (success)
const LockSeatsLua = `
-- KEYS   = seat keys (seat:{eventId}:{seatId})
-- ARGV[1] = owner id (userId or bookingId)
-- ARGV[2] = ttl in seconds

for i, key in ipairs(KEYS) do
  if redis.call("EXISTS", key) == 1 then
    return 0
  end
end

for i, key in ipairs(KEYS) do
  redis.call("SET", key, ARGV[1], "EX", ARGV[2])
end

return 1
`

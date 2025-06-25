local val = redis.call("GET", KEYS[1])
if val == ARGV[1] then
    redis.call("EXPIRE", KEYS[1], ARGV[2])
    return true
else
    return false
end

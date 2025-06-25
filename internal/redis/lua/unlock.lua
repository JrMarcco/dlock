local val = redis.call("GET", KEYS[1])
if val == ARGV[1] then
    redis.call("DEL", KEYS[1])
    return true
else
    return false
end

local val = redis.call("GET", KEYS[1])
if val == ARGV[1] then
    return redis.call("EXPIRE", KEYS[1], ARGV[2])
else
    return 0
end

local val = redis.call("GET", KEYS[1])
if val == false or val == nil then
    -- key 不存在
    return redis.call("SET", KEYS[1], ARGV[1], "EX", ARGV[2])
elseif val == ARGV[1] then
    -- key 存在且值相等，刷新过期时间
    redis.call("EXPIRE", KEYS[1], ARGV[2])
    return "OK"
else
    -- 锁被其他线程持有
    return ""
end

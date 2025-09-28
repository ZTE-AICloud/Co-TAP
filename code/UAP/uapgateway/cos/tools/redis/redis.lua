local zredis = require "kong.tools.redis.zredis"
local utils = require "kong.tools.utils"
local resty_lock = require "resty.lock"

local ngx = ngx
local kong = kong
local tonumber = tonumber
local tostring = tostring
local ipairs = ipairs
local null = ngx.null
local fmt = string.format

local REDIS_CONN = {}
local MAX_WAIT_TIME_MILIS = 1500  -- 获取连接，1500ms超时
local TIME_BETWEEN_EVICTION_RUNS_MILIS = 600000 -- 10min检查一次存活
local MIN_EVICTABLE_IDLE_TIME_MILIS = 7200000 -- 2h过期
local EXPIRE_TIME = 30
local TIMEOUT = 5
local DEFAULT_LOCKER

local zredis_cmds = {
  ["auth"] = function(self, pwd)
    local ok, reply = zredis.command(self.context, "AUTH " .. pwd,2)
    if not ok and reply then
      local err_info = ffi.string(reply.str)
      zredis.free_reply(reply)
      return nil, err_info
    end

    zredis.free_reply(reply)
    return nil
  end,
  ["select"] = function(self, database)
    local ok, reply = zredis.command(self.context, "SELECT " .. database,2)
    if not ok and reply then
      local err_info = ffi.string(reply.str)
      zredis.free_reply(reply)
      return nil, err_info
    end

    zredis.free_reply(reply)
    return nil
  end,
  ["get"] = function(self, key)
    local ok, reply = zredis.command(self.context, "GET " .. key, 2)

    if ok and reply then
      local out
      if tonumber(reply.type) == 1 then
        out = ffi.string(reply.str)
      elseif tonumber(reply.type) == 3 then
        out = tonumber(reply.integer)
      elseif tonumber(reply.type) == 4 then
        out = null
      end

      zredis.free_reply(reply)
      return out
    else
      zredis.free_reply(reply)
      return null
    end
  end,
  ["set"] = function(self, key, value, ...)
    local args = {...}
    local len = 3
    local ok, reply
    if #args == 0 then
      ok, reply = zredis.command(self.context, "SET " .. key .. " " .. value, len)
    else
      ok, reply = zredis.command(self.context, "SET " .. key .. " " .. value .. " " .. table.concat(args, " "), len + #args)
    end

    if not ok and reply then
      local err_info = ffi.string(reply.str)
      zredis.free_reply(reply)
      return nil, err_info
    end

    zredis.free_reply(reply)
    return nil
  end,
  ["setnx"] = function(self, key, value)
    local ok, reply
    ok, reply = zredis.command(self.context, "SETNX " .. key .. " " .. value,3)

    if ok and reply then
      local out
      if tonumber(reply.type) == 1 then
        out = ffi.string(reply.str)
      elseif tonumber(reply.type) == 3 then
        out = tonumber(reply.integer)
      elseif tonumber(reply.type) == 4 then
        out = nil
      end

      zredis.free_reply(reply)
      return out
    else
      zredis.free_reply(reply)
      return nil
    end
  end,
  ["pipeline"] = function(self, cmd_table)
    local add_ok
    for _,cmd in pairs(cmd_table) do
      add_ok= zredis.addPipeline(self.context, cmd)
      if not add_ok then
        local err_info = "add pipeline failed, cmd: " .. cmd
        return nil, err_info
      end
    end

    for i = 1, #cmd_table do
      local reply_ok, reply = zredis.getPipelineReply(self.context)
      if not reply_ok then
        local err_info = "get pipeline reply failed, cmd: " .. cmd_table[i]
        zredis.free_reply(reply)
        return nil, err_info
      end
      zredis.free_reply(reply)
    end
    return true
  end,
  ["script"] = function(self, cmd, script)
    local ok, reply = zredis.execScript(self.context, "script " .. cmd .. " %s", script)
    if ok and reply then
      local out
      if tonumber(reply.type) == 1 then
        out = ffi.string(reply.str)
      elseif tonumber(reply.type) == 3 then
        out = tonumber(reply.integer)
      elseif tonumber(reply.type) == 4 then
        out = null
      end
      zredis.free_reply(reply)
      return out
    else
      zredis.free_reply(reply)
      return null
    end
  end,
  ["evalsha"] = function(self, sha1, keylen, ...)
    local args = {...}
    local cmd = "evalsha " .. sha1 .. " " .. keylen ..  " " .. table.concat(args, " ")
    local ok, reply = zredis.command(self.context, cmd)
    local reply_is_null = reply == ffi.NULL 
    if ok and reply and not reply_is_null then
      local out
      local ret_type = tonumber(reply.type)
      if ret_type == 1 then
        out = ffi.string(reply.str)
      elseif ret_type == 2 then
        out = {}
        local count = reply.elements or 0
        for i = 1, tonumber(count) do
          local e = reply.element[i-1]
          if tonumber(e.type) == 1 then
            utils.insert_tail(out, ffi.string(e.str))
          elseif tonumber(e.type) == 3 then
            utils.insert_tail(out, tonumber(e.integer))
          elseif tonumber(e.type) == 4 then
            kong.log.err("pipeline get reply is nil")
          end
        end
      elseif ret_type == 3 then
        out = tonumber(reply.integer)
      elseif ret_type == 4 then
        kong.log.err("evalsha get reply nil")
      end
      zredis.free_reply(reply)
      return out
    else
      local err_info = reply_is_null and "evalsha failed, reply return null" or ffi.string(reply.str)
      zredis.free_reply(reply)
      return nil,err_info
    end
  end,
  ["del"] = function(self, key)
    local ok, reply = zredis.command(self.context, "DEL " .. key, 2)
    if not ok and reply then
      local err_info = ffi.string(reply.str)
      zredis.free_reply(reply)
      return nil, err_info
    end

    zredis.free_reply(reply)
    return nil
  end,
  ["incrby"] = function(self, key, value)
    local ok, reply = zredis.command(self.context, "INCRBY " .. key .. " " .. value,3)
    if not ok and reply then
      local err_info = ffi.string(reply.str)
      zredis.free_reply(reply)
      return nil, err_info
    elseif reply then
      local res = reply.integer and tonumber(reply.integer) or 0
      zredis.free_reply(reply)
      return res
    end

    zredis.free_reply(reply)
  end,
  ["zscore"] = function(self, key, member)
    local ok, reply = zredis.command(self.context, "ZSCORE " .. key .. " " .. member, 3)

    if ok and reply then
      local out
      if tonumber(reply.type) == 1 then
        out = ffi.string(reply.str)
      elseif tonumber(reply.type) == 3 then
        out = tonumber(reply.integer)
      elseif tonumber(reply.type) == 4 then
        out = null
      end

      zredis.free_reply(reply)
      return out
    else
      zredis.free_reply(reply)
      return null
    end
  end,
  ["zrem"] = function(self, key, member)
    local ok, reply = zredis.command(self.context, "ZREM " .. key .. " " .. member, 3)
    if not ok and reply then
      local err_info = ffi.string(reply.str)
      zredis.free_reply(reply)
      return nil, err_info
    end

    zredis.free_reply(reply)
    return nil
  end,
  ["zincrby"] = function(self, key, value, member)
    local ok, reply = zredis.command(self.context, "ZINCRBY " .. key .. " " .. value .. " " .. member,4)
    if not ok and reply then
      local err_info = ffi.string(reply.str)
      zredis.free_reply(reply)
      return nil, err_info
    elseif reply then
      local res = reply.integer and tonumber(reply.integer) or 0
      zredis.free_reply(reply)
      return res
    end

    zredis.free_reply(reply)
  end,
  ["scan"] = function(self, offset, match, pattern, count, limit)
    local ok, reply = zredis.command(self.context, "SCAN " .. offset .. " " .. match .. " " .. pattern .. " " .. count .. " " .. limit,6)
    if not ok and reply then
      local err_info = ffi.string(reply.str)
      zredis.free_reply(reply)
      return 0, err_info
    end

    local res = reply.element
    local next_offset = tonumber(res[0].integer)
    local len = tonumber(string.sub(tostring(res[1].elements), 1, -4))
    local keys = res[1].element
    local results = {next_offset}
    local out = {}
    for i = 0, len - 1 do
      local key = ffi.string(keys[i].str, keys[i].len)
      table.insert(out, key)
    end
    table.insert(results, out)
    zredis.free_reply(reply)
    return results, nil
  end,
  ["ttl"] = function(self, key)
    local ok, reply = zredis.command(self.context, "TTL " .. key, 2)
    if not ok and reply then
      local err_info = ffi.string(reply.str)
      zredis.free_reply(reply)
      return nil, err_info
    end

    zredis.free_reply(reply)
    return reply and tonumber(reply.integer) or -1
  end,
  ["expire"] = function(self, key, ttl)
    local ok, reply = zredis.command(self.context, "EXPIRE  " .. key .. " " .. ttl,3)
    if not ok and reply then
      local err_info = ffi.string(reply.str)
      zredis.free_reply(reply)
      return nil, err_info
    end

    zredis.free_reply(reply)
    return nil
  end,
  ["exists"] = function(self, key)
    local ok, reply = zredis.command(self.context, "EXISTS  " .. key,2)
    if not ok and reply then
      local err_info = ffi.string(reply.str)
      zredis.free_reply(reply)
      return nil, err_info
    elseif reply then
      zredis.free_reply(reply)
      return reply.integer and tonumber(reply.integer)
    end

    zredis.free_reply(reply)
  end,
  ["init_pipeline"] = function(self, idle_timeout, pool_size)
    -- do nothing
    return true
  end,
  ["commit_pipeline"] = function(self, idle_timeout, pool_size)
    -- do nothing
    return true
  end,
  ["set_keepalive"] = function(self, idle_timeout, pool_size)
    -- zredis.free_connect(self)
    -- unlock
    shm_unlock(locker)
    return true
  end,
  ["ping"] = function(self)
    local ok, reply = zredis.command(self.context, "ping")
    if ok then
      zredis.free_reply(reply)

      return true
    end
    zredis.free_reply(reply)

    return false
  end,
  ["xadd"] = function(self, stream_name, ...)
    local args = {...}
    local msg = table.concat(args, " ")
    local ok, reply
    local cmd = "XADD " .. stream_name .. " " .. "maxlen 1000 * " .. msg
    ok, reply = zredis.command(self.context, cmd)

    if not ok and reply then
      local err_info = ffi.string(reply.str)
      zredis.free_reply(reply)
      return nil, err_info
    end

    zredis.free_reply(reply)
    return nil
  end,
  ["xread"] = function(self, stream_name, block_time, last_id)
    block_time = block_time or 1000
    last_id = last_id or "$"
    local cmd = "XREAD COUNT 1 BLOCK " .. block_time .. " STREAMS " .. stream_name .. " " .. last_id
    local ok, reply = zredis.command(self.context, cmd)
    local reply_is_null = reply == ffi.NULL 
    if ok and reply and not reply_is_null then
      local out
      local ret_type = tonumber(reply.type)
      if ret_type == 1 then
        out = ffi.string(reply.str)
      elseif ret_type == 2 then
        out = getArrayReply(reply)
      elseif ret_type == 3 then
        out = tonumber(reply.integer)
      end
      zredis.free_reply(reply)
      return out
    else
      local err_info = reply_is_null and "xread failed, reply return null" or ffi.string(reply.str)
      zredis.free_reply(reply)
      return nil,err_info
    end
  end,
}

local function now()
  ngx.update_time()
  return ngx.now() * 1000
end

local function get_lock_key(id, name)
  return fmt("redis:pool:%s:%s", id, name)
end

local function shm_lock(tenant, timeout, exptime)
  -- step1:acquire lock
  if not DEFAULT_LOCKER then
    local opts = {["timeout"] = timeout or 0,["exptime"] = exptime or 0.05}--this can be set using the conf file
    local rlock, err = resty_lock:new("kong_locks", opts)
    if not rlock then
      kong.log.warn("failed to create lock in cos bill control:", tostring(err))
      return nil, err
    end
    DEFAULT_LOCKER = rlock
  end

  local ctx = ngx.ctx
  local _tenant = tenant or ctx.tenant_in_used or ctx.tenant or "NONE"
  local lock_key = get_lock_key(ngx.worker.id(), _tenant)
  -- acquire lock
  local elapsed, err = DEFAULT_LOCKER:lock(lock_key)
  if not elapsed then
    if err == "timeout" then
      return nil, err
    end
    return nil, "failed to acquire worker lock: " .. err
  end
end

local function shm_unlock(lock, name)
  if not lock then
    return
  end

  local ctx = ngx.ctx
  local lock_key = get_lock_key(ngx.worker.id(), name)

  lock:unlock(lock_key)
end

local function is_valid(red)
  if not red or red.context == nil or red.context[0] == nil then
    return
  end

  local ok, err = red:ping()
  if not ok then
    return nil, err
  end
  return true
end

local function discard(red, name)
  if red and red.context then
    local _, err = shm_lock(name, TIMEOUT, EXPIRE_TIME)
    if err then
      kong.log.err("failed to lock: ", tostring(err))
      return
    end

    REDIS_CONN[name] = nil
    zredis.free_connect(red.context)
    red.context = nil
    red = nil
    shm_unlock(DEFAULT_LOCKER, name)
  end
end

local function get_zredis_connection(conf)
  local context, err
  local hosts = utils.split(conf.host, ",")

  for _, host in ipairs(hosts) do
    if conf and conf.user then
        context, err = zredis.connect_with_user(host, conf.port, conf.user, conf.pwd, conf.timeout)
    else
        context, err = zredis.connect(host, conf.port, conf.pwd, conf.timeout)
    end
  end

  if not context then   
    return nil, err
  end

  if conf.database ~= 0 then
    local ok, reply = zredis.command(context, "select " .. conf.database, 2)
    if not ok then
      local error_info = ffi.string(reply.str)
      zredis.free_reply(reply)

      return nil, error_info
    end

    zredis.free_reply(reply)
  end

  return setmetatable(
    {
        context = context,
        last_active_time = now()
    }, 
    {
        __index = function(self, cmd)
            return zredis_cmds[cmd]
        end
    })
end

local function get_connection_from_pool(conf)
  local red, err, _

  local start_time = now()
  local instance_name = conf and conf.name or "default"
  while (now() - start_time) < MAX_WAIT_TIME_MILIS do
    if not REDIS_CONN[conf.name] then
      -- get lock, first
      local _, lerr = shm_lock(instance_name, 0, 10)
      if lerr and lerr ~= "timeout" then
        return nil,  "can not get zredis connection lock: " .. tostring(lerr)
      end

      -- get lock success
      if not lerr then
        red, err = get_zredis_connection(conf)
        if err then
          shm_unlock(DEFAULT_LOCKER, instance_name)
          return nil, err
        end

        red.last_active_time = now()
        REDIS_CONN[instance_name] = red
        -- unlock
        shm_unlock(DEFAULT_LOCKER, instance_name)

        -- get redis context success
        return red
      end
    else
      red = REDIS_CONN[instance_name]
      if red and (now() - red.last_active_time) > TIME_BETWEEN_EVICTION_RUNS_MILIS and not is_valid(red) then
        discard(red, instance_name)
      elseif red then
        -- get redis context success
        red.last_active_time = now()
        return red
      end
    end

    sleep(0.1)
  end

  return nil, err or "can not get redis connection"
end

local function connect(conf)
    return get_connection_from_pool(conf)
end

return {
    connect = connect
}
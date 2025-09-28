local ffi = require 'ffi'
local utils = require "kong.tools.utils"

local redisffi = ffi.load("zredis")
local type = type

ffi.cdef [[
struct sockadr;
struct zredisReplica;
struct ev_Pool;
struct depth_Ctl;
struct serverList;
struct redisContext;

typedef int redisFD;
typedef void (redisPushFn)(void *, void *);

enum redisConnectionType {
    REDIS_CONN_TCP,
    REDIS_CONN_UNIX,
    REDIS_CONN_USERFD
};

typedef struct timeval{
    long tv_sec;
    long tv_usec;
} timeval;

typedef struct __pthread_internal_list
{
  struct __pthread_internal_list *__prev;
  struct __pthread_internal_list *__next;
} __pthread_list_t;

typedef union
{
  struct __pthread_mutex_s
  {
    int __lock;
    unsigned int __count;
    int __owner;
    unsigned int __nusers;
    int __kind;
    short __spins;
    short __elision;
    __pthread_list_t __list;
  } __data;
  char __size[40];
  long int __align;
} pthread_mutex_t;

typedef union
{
  struct
  {
    int __lock;
    unsigned int __nr_readers;
    unsigned int __readers_wakeup;
    unsigned int __writer_wakeup;
    unsigned int __nr_readers_queued;
    unsigned int __nr_writers_queued;
    int __writer;
    int __shared;
    signed char __rwelision;
    unsigned char __pad1[7];
    unsigned long int __pad2;
    /* FLAGS must stay at this position in the structure to maintain
       binary compatibility.  */
    unsigned int __flags;
  } __data;
  char __size[56];
  long int __align;
} pthread_rwlock_t;

typedef struct redisContextFuncs {
    void (*free_privctx)(void *);
    void (*async_read)(struct redisAsyncContext *);
    void (*async_write)(struct redisAsyncContext *);
    ssize_t (*read)(struct redisContext *, char *, size_t);
    ssize_t (*write)(struct redisContext *);
} redisContextFuncs;

typedef struct redisReadTask {
    int type;
    long long elements; /* number of elements in multibulk container */
    int idx; /* index in parent (array) object */
    void *obj; /* holds user-generated value for a read task */
    struct redisReadTask *parent; /* parent task */
    void *privdata; /* user-settable arbitrary field */
} redisReadTask;

typedef struct redisReplyObjectFunctions {
    void *(*createString)(const redisReadTask*, char*, size_t);
    void *(*createArray)(const redisReadTask*, size_t);
    void *(*createInteger)(const redisReadTask*, long long);
    void *(*createDouble)(const redisReadTask*, double, char*, size_t);
    void *(*createNil)(const redisReadTask*);
    void *(*createBool)(const redisReadTask*, int);
    void (*freeObject)(void*);
} redisReplyObjectFunctions;

typedef struct redisReader {
    int err; /* Error flags, 0 when there is no error */
    char errstr[128]; /* String representation of error when applicable */

    char *buf; /* Read buffer */
    size_t pos; /* Buffer cursor */
    size_t len; /* Buffer length */
    size_t maxbuf; /* Max length of unused buffer */
    long long maxelements; /* Max multi-bulk elements */

    redisReadTask **task;
    int tasks;

    int ridx; /* Index of current read task */
    void *reply; /* Temporary reply pointer */

    redisReplyObjectFunctions *fn;
    void *privdata;
} redisReader;

typedef struct redisContext {
    const redisContextFuncs *funcs;   /* Function table */

    int err; /* Error flags, 0 when there is no error */
    char errstr[128]; /* String representation of error when applicable */
    redisFD fd;
    int flags;
    char *obuf; /* Write buffer */
    redisReader *reader; /* Protocol reader */

    enum redisConnectionType connection_type;
    struct timeval *connect_timeout;
    struct timeval *command_timeout;

    struct {
        char *host;
        char *source_addr;
        int port;
    } tcp;

    struct {
        char *path;
    } unix_sock;

    /* For non-blocking connect */
    struct sockadr *saddr;
    size_t addrlen;

    /* Optional data and corresponding destructor users can use to provide
     * context to a given redisContext.  Not used by hiredis. */
    void *privdata;
    void (*free_privdata)(void *);

    /* Internal context pointer presently used by hiredis to manage
     * SSL connections. */
    void *privctx;

    /* An optional RESP3 PUSH handler */
    redisPushFn *push_cb;
    char salt[16];
} redisContext;

typedef struct redisReply {
    int type; /* REDIS_REPLY_* */
    long long integer; /* The integer when type is REDIS_REPLY_INTEGER */
    double dval; /* The double when type is REDIS_REPLY_DOUBLE */
    size_t len; /* Length of string */
    char *str; /* Used for REDIS_REPLY_ERROR, REDIS_REPLY_STRING
                  REDIS_REPLY_VERB, and REDIS_REPLY_DOUBLE (in additional to dval). */
    char vtype[4]; /* Used for REDIS_REPLY_VERB, contains the null
                      terminated 3 character content type, such as "txt". */
    size_t elements; /* number of elements, for REDIS_REPLY_ARRAY */
    struct redisReply **element; /* elements vector for REDIS_REPLY_ARRAY */
} redisReply;

typedef struct redis_Ha_Connection{
	struct redisContext *redis_HA_handle;
	struct serverList *svrlist;
	int flag;
	struct timeval timeout;
	struct ev_Pool *workerEvPool;
} redisHaConnection;

typedef struct zredishandle_t zredis_handle_t;
typedef struct redis_Async_Connection{
    struct zredisReplica** replicas;

    pthread_mutex_t  replicaslock;
    unsigned int replicasCount;
    char passwd[256];
    struct depth_Ctl* depth;
    struct ev_Pool* workerEvPool;
    struct ev_Pool* managerEvPool;
    struct zredisReplica **SlotToReplicate;
    struct redisContext* haCon;
    int clusterEnable;
    int routeAgain;
    zredis_handle_t *zredis_handle;
    int flag;
    struct serverList *svrlist;
} redisAsyncConnection;

struct zredishandle_t {
	int  type;//0:HA 1:single 2:cluster
	redisAsyncConnection *redis_async_handle;
	redisHaConnection *redis_Ha_Con;
	redisReply *reply;

	char ip[128];
	int port;
	char passwd[256];

	int error_code;
	struct timeval timeout;

	char details_info[512];
	int APItype;
};

void freeReplyObject(void *reply);
/*********zredis*********/
bool redisCreateHandle(const char *ip,int port,const char* passwd,zredis_handle_t **redis_handle);
bool zRedisvCommand(zredis_handle_t *redis_handle, const char *format, ...);
bool zRedisPipelineAdd(zredis_handle_t *redis_handle,const char *format, ...);
bool zRedisGetReply(zredis_handle_t*redis_handle);
void zredisFree(zredis_handle_t *redis_handle);
]]

local ok, new_tab = pcall(require, "table.new")
if not ok or type(new_tab) ~= "function" then
  new_tab = function(narr, nrec)
    return {}
  end
end

local _M = new_tab(0, 54)
_M._VERSION = '1.20.20.02'

local function trim(s)
  if not s then
    return ""
  end

  local idx = #s
  for i=#s, 1, -1 do
    local b = string.byte(s, i, i)
    -- ascii<=32或=127时为特殊字符
    if b > 32 and b ~= 127 then
      break
    end
    idx = idx - 1
  end
  return s:sub(1, idx)
end

local function remove_ipv6_bracket(ip)
    if not ip then
        return nil
    end
    -- for ipv6
    if utils.hostname_type(ip) == "ipv6" and sub(ip, 1, 1) == "[" then
        ip = sub(ip, 2, -2)
    end

    return ip
end

function _M.connect(ip, port, pwd, msec)
  if not ip or ip == "" or ip == ngx.null then
    return nil, "ip is invalid"
  end

  ip = remove_ipv6_bracket(ip)
  port = tonumber(port)
  if not port then
    return nil, "port is invalid"
  end

  local handle_t = ffi.new("zredis_handle_t *[1]")
  local ok = redisffi.redisCreateHandle(ip, port, pwd, handle_t)
  if not ok then
    _M.free_connect(handle_t)
    return nil, "connect to redis error"
  end

  return handle_t
end

function _M.connect_with_user(ip, port, user, pwd)
  pwd = trim(pwd)
  if user and #user > 0 then
    pwd = user .. " " .. pwd
  end

  return _M.connect(ip, port, pwd)
end

function _M.command(handle, command, len)
  local ok = redisffi.zRedisvCommand(handle[0], command)
  return ok, handle[0].reply
end

function _M.free_reply(reply)
  if reply ~= nil then
    redisffi.freeReplyObject(reply)
  end
end

function _M.addPipeline(handle, command)
    return redisffi.zRedisPipelineAdd(handle[0], command)
end

function _M.getPipelineReply(handle)
  if handle ~= nil and handle[0] ~= nil then
    local ok = redisffi.zRedisGetReply(handle[0])
    return ok, handle[0].reply
  end
end

function _M.execScript(handle, cmd, script)
  if handle ~= nil and handle[0] ~= nil then
    local ok = redisffi.zRedisvCommand(handle[0], cmd, script)
    return ok, handle[0].reply
  end
end

function _M.free_connect(handle)
  if handle ~= nil and handle[0] ~= nil then
    redisffi.zredisFree(handle[0])
    handle[0] = nil
    handle = nil
  end
end

return _M
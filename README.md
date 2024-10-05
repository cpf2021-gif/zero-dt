# 使用 go-zero + dtm 实现分布式事务
## 前提
### 1. 启动etcd
### 2. 启动dtm
- clone 项目
```sh
git clone git@github.com:dtm-labs/dtm.git
```
- 复制配置文件
```sh
cp conf.sample.yml conf.yml
```
- 修改配置文件
```yaml
#####################################################################
### dtm can be run without any config.
### all config in this file is optional. the default value is as specified in each line
### all configs can be specified from env. for example:
### MicroService.EndPoint => MICRO_SERVICE_END_POINT
#####################################################################

Store: # specify which engine to store trans status
  Driver: 'mysql'
  Host: 'localhost'
  User: 'root'
  Password: 'xxxxxxxxxxxx' # 修改为你的数据库密码
  Port: 3306
  Db: 'dtm'

#   Driver: 'boltdb' # default store engine

#   Driver: 'redis'
#   Host: 'localhost' # host1:port1,host2:port2 for cluster connection
#   User: ''
#   Password: ''
#   Port: 6379 # required but won't be used for cluster connection

#   Driver: 'postgres'
#   Host: 'localhost'
#   User: 'postgres'
#   Password: 'mysecretpassword'
#   Port: '5432'
#   Db: 'postgres'
#   Schema: 'public' # default value is 'public'

### following config is for only Driver postgres/mysql
#   MaxOpenConns: 500
#   MaxIdleConns: 500
#   ConnMaxLifeTime: 5 # default value is 5 (minutes)

### flollowing config is only for some Driver
#   DataExpire: 604800 # Trans data will expire in 7 days. only for redis/boltdb.
#   FinishedDataExpire: 86400 # finished Trans data will expire in 1 days. only for redis.
#   RedisPrefix: '{a}' # default value is '{a}'. Redis storage prefix. store data to only one slot in cluster

MicroService: # gRPC/HTTP based microservice config
  Driver: 'dtm-driver-gozero' # name of the driver to handle register/discover
  Target: 'etcd://localhost:2379/dtmservice' # register dtm server to this url
  EndPoint: 'localhost:36790'

### the unit of following configurations is second
# TransCronInterval: 3 # the interval to poll unfinished global transaction for every dtm process
# TimeoutToFail: 35 # timeout for XA, TCC to fail. saga's timeout default to infinite, which can be overwritten in saga options
# RetryInterval: 10 # the subtrans branch will be retried after this interval
# RequestTimeout: 3 # the timeout of HTTP/gRPC request in dtm

# LogLevel: 'info'              # default: info. can be debug|info|warn|error
# Log:
#    Outputs: 'stderr'           # default: stderr, split by ",", you can append files to Outputs if need. example:'stderr,/tmp/test.log'
#    RotationEnable: 0           # default: 0
#    RotationConfigJSON: '{}'    # example: '{"maxsize": 100, "maxage": 0, "maxbackups": 0, "localtime": false, "compress": false}'
#
# HttpPort: 36789
# GrpcPort: 36790
# JsonRpcPort: 36791

### advanced options
# UpdateBranchAsyncGoroutineNum: 1 # num of async goroutine to update branch status
# TimeZoneOffset: '' #default '' using system default. '+8': Asia/Shanghai; '0': GMT
# AdminBasePath: '' #default '' set admin access base path

# ConfigUpdateInterval: 10   # the interval to update configuration in memory such as topics map... (seconds)
# TimeZoneOffset: '' # default '' using system default. '+8': Asia/Shanghai; '0': GMT
# AlertRetryLimit: 3 # default 3; if a transaction branch has been retried 3 times, the AlertHook will be called
# AlertWebHook: '' # default ''; sample: 'http://localhost:8080/dtm-hook'. this hook will be called like this:
## curl -H "Content-Type: application/json" -d '{"gid":"xxxx","status":"submitted","retry_count":3}' http://localhost:8080/dtm-hook

```
- 初始化数据库
```sh
mysql -uroot -p < sqls/dtmcli.barrier.mysql.sql
mysql -uroot -p < sqls/dtmsvr.storage.mysql.sql
```
- 启动dtm
```sh
go run main.go -c conf.yml
```
### 3. 启动go-zero
- 初始化数据库
```sh
mysql -uroot -p < order/sql/order.sql
mysql -uroot -p < order/sql/stock.sql
```
- 启动服务
```sh
go run order/app/order/rpc/order.go
go run order/app/stock/rpc/stock.go

go run order/app/order/api/order.go
```
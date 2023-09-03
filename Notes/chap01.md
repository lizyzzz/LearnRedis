### 1. NoSQL 数据库概述
NoSQL(Not only SQL)，泛指非关系型的数据库。通过key-value存储。
* 不支持ACID
* 远超SQL的性能
* 支持持久化
* 适用场景：
    (1)对数据的高并发读写；  
    (2)海量数据的读写；  
    (3)对数据高扩展性的要求。  
* 不适用场景：
    (1)需要事务支持；  
    (2)基于sql的结构化查询存储，处理复杂的关系，需要关系查询。  
### 2. Redis 概述
(1) 一个开源的 key-value 存储系统；  
(2) 支持的数据类型有： string(字符串), list(链表), set(集合), zset(sorted set --有序集合)和 hash(哈希类型)；  
(3) 数据类型支持 push/pop, add/remove 及取交集并集和差集及更丰富的操作，而且这些操作都是 `原子性`的；  
(4) 支持不同方式的排序；  
(5) 数据都是缓存在内存中；  
(6) Redis 会`周期性`地把更新的`数据写入磁盘`或者把修改操作写入追加的记录文件；  
(7) 在此基础上实现`主从同步`。  
### 3. Redis 服务端和客户端
文件路径一般是：/usr/local/bin  
服务端：redis-server /etc/redis.conf  
客户端：redis-cli  
#### 3.1 Redis 相关知识介绍
(1) 端口号是6379；  
(2) 默认16个数据库，类似数组下标从0开始，`初始默认使用0号库`；  
(3) 使用命令 `select <dbid>`来切换数据库。如select 8；  
(4) 统一密码管理，所有库同样密码；  
(5) `dbsize` 查看当前数据库的 key 的数量；  
(6) flushdb 清空当前库；  
(7) flushall 杀掉全部库；  
(8) Redis 是单线程+IO多路复用技术。  
### 4. 五大常用数据类型
#### 4.1 键(key)
* `keys *` 查看当前库所有 key (匹配: keys*1)
* `exists key` 判断某个 key 是否存在
* `type key` 查看 key 的类型
* `del key` 删除指定的 key 
* `unlink key` 根据 value 选择非阻塞删除(仅将 keys 从 keyspace 元数据中删除, 真正的删除会在后续异步操作中。)
* `expire key 10` 为指定的 key 设置过期时间 10s
* `ttl key` 查看 key 的剩余生存时间, -1 表示永不过期, -2 表示已过期。
#### 4.2 字符串(String)
* String 类型是`二进制安全的`，意味着 Redis 的 String 可以包含任何数据。比如 jpg 图片或者序列化的对象。
* String 类型是 Redis 最基本的数据类型，一个字符串 value 最多可以是 `512M`。
* `set key value [EX s|PX ms] [NX|XX]` 添加一个key-value  
    EX: 过期的秒数  
    PX: 过期的毫秒数  
    NX: 当key不存在时添加到数据库中  
    XX: 当key存在时也添加到数据库中  
* `get key` 查询对应的键值
* `append key value` 将指定的 value 追加到原值的末尾
* `strlen key` 获取 key 对应的 value 的长度
* `setnx key value` 只有 key 不存在时，设置 key 的值
* `incr key` 将 key 中存储的数字值增加 1, 只能对数字值操作，如果为空，新增值为 1
* `decr key` 将 key 中存储的数字值减少 1
* `incrby/decrby key step` 将 key 中存储的数字值增加或减少自定义的步长
* 
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
* 以上操作具有 `原子性`  
* `mset key1 value1 key2 value2...` 同时设置一个或多个 key-value 对
* `meget key1 key2 key3...` 同时获取一个或多个value
* `msetnx key1 value1 key2 value2...` 同时设置一个或多个 key-value 对，当且仅当所有给定 key 都不存在。有一个失败则都失败。
* `getrange key begin end` 获取 key 对应的 value 的[begin, end]范围, 左右都是闭
* `setrange key offset value` 设置 key 对应的 value 的 offset 位置为 value
* `setex key EX value` 设置键值的同时设置过期时间，单位秒
* `getset key value` 以新换旧，设置了新值同时获得旧值
#### 4.3 列表(List)
* 单键多值(同一个 key 可以对应多个 value，用 list 存储，类似 unordered_map)
* Redis 列表就是一个双向链表，可以对头部或尾部进行插入或删除。通过索引下标的操作中间的节点性能会较差。
* `lpush/rpush key value value ...` 从左边/右边插入一个或多个值。
* `lpop/rpop key [count]` 从左边/右边取出 count 个值(默认是 1)`值在键在，值空键亡`。
* `rpoplpush key1 key2` 从 key1 列表右边取出一个值并插入到 key2 列表的左边。
* `lrange key begin end` 按照索引下标获得元素，从左到右, 左右都是闭, `0 -1` 显示所有。
* `lindex key index` 按照索引下标获得元素(下标从0开始，从左到右)
* `llen key` 获得列表的长度
* `linsert key [before|after] value newvalue` 在 value 前面/后面插入新值 newvalue (从左到右)
* `lrem key n value` 从左边开始删除 n 个 指定的 value, 返回值表示删除的个数
* `lset key index newvalue` 将列表 key 下标为 index 的值替换成 newvalue
* 底层数据结构: `quickList`, 首先在列表元素较少时使用一块连续的内存存储，这个结构是`ziplist`, 也即是`压缩列表`。当数据量较多时，将多个`ziplist`使用双向指针串起来使用，这样既满足了快速的插入删除性能，又不会出现太大的空间冗余(每个节点都需要一个前指针和后指针，`ziplist`内部空间连续不需要，从而节省了空间。)
#### 4.4 集合(Set)
* set 对外提供的功能与 list 类似，特殊之处在于 set 是可以**自动排重**的, 即元素不可重复, 并且提供了判断某个成员是否在一个 set 集合内的重要接口。
* set 底层是 **string 类型的无序集合**。它底层是一个以 value 作为 key 的 hash 表。所以增删查改都是 O(1) 复杂度。 
* `sadd key value1 value2...` 将一个或多个 value 加到集合 key 中，已经存在的 value 将被忽略(去重)。
* `smembers key` 取出该集合中的所有值
* `sismember key value` 判断集合 key 是否含有该 value 值，有返回 1，无返回 0。
* `scard key` 返回该集合的元素个数
* `srem key value value...` 删除集合中的一个或多个元素
* `spop key [count]` **随机从该集合中取出count个值(默认是 1)**
* `srandmember key n` 随即从该集合中获取 n 个值，但是不会从集合中删除
* `smove source dest value` 把集合中一个值从一个集合移动到另一个集合
* `sinter key1 key2` 返回两个集合的**交集**元素
* `sunion key1 key2` 返回两个集合的**并集**元素
* `sdiff key1 key2` 返回两个集合的**差集**元素(key1 - key2)
* set 的底层数据结构就是**哈希表**, 类似 cpp 的 unordered_map(只有 key, value 为定值)
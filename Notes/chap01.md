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
### 4. 五大常用数据类型(都可以作为key对应的value)
#### 4.0 键(key)操作
* `keys *` 查看当前库所有 key (匹配: keys*1)
* `exists key` 判断某个 key 是否存在
* `type key` 查看 key 的类型
* `del key` 删除指定的 key 
* `unlink key` 根据 value 选择非阻塞删除(仅将 keys 从 keyspace 元数据中删除, 真正的删除会在后续异步操作中。)
* `expire key 10` 为指定的 key 设置过期时间 10s
* `ttl key` 查看 key 的剩余生存时间, -1 表示永不过期, -2 表示已过期。
#### 4.1 字符串(String)
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
#### 4.2 列表(List)
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
#### 4.3 集合(Set)
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
#### 4.4 哈希(Hash)
* Redis hash 是一个键值对集合, 是一个 string 类型的 field 和 value 的映射表。因此可以以 field 作为字段名, value 作为字段值, 从而适合存储对象。
* `hset key field value...` 给 key 集合中的多个 field 键赋值 value
* `hget key field` 从 key 集合中取出 field 对应的 value
* `hmset key field1 value1 field2 value2...` 批量设置 hash 的值
* `hmget key field1 field2...` 批量获取哈希集合多个字段的值
* `hexists key field` 查看 key 集合中，给定 field 是否存在
* `hkeys key` 列出该哈希集合的所有 field
* `hvals key` 列出该哈希集合的所有 value
* `hincrby key field increment` 为哈希集合中的 field 对应的 value 加上增量 increment
* `hsetnx key field value` 将哈希集合中的 field 对应的值设置为 value, 当且仅当 field 不存在。
* hash 底层的数据结构是: 当数据量较少时, 采用 ziplist (压缩列表), 当数据量较大时采用 hashmap
#### 4.5 有序集合(Zset, sorted set)
* Redis 有序集合 zset 和普通集合 set 非常类似，是一个没有重复元素的字符串集合。不同之处是每个成员都关联了一个**优先级值**，这个**优先级值**用于对集合进行排序。**集合的成员是唯一的，但是优先级值可以是重复的**。
* `zadd key score member [score member]` 将一个或多个 member 元素及其 score 值加入到有序集合 key 中
* `zrange key start stop [WITHSCORES]` 返回有序集 key 中，下标在 [start, stop] (左闭右闭)之间的 member 元素, [WITHSCORES]表示将分数一起返回。
* `zrangebyscore key min max [withscores]` 返回有序集合 key 中, 所有 score 值位于[min, max] (左闭右闭)之间的 member, member 按 score 值递增排列
* `zrevrangebyscore key [max min] [withscores]` 同上, 但次序递减排列
* `zincrby key increment member` 为 member 的 score 增加增量 increment
* `zrem key member` 删除有序集合 key 的 member
* `zcount key min max` 统计有序集合 key 在 [min, max]分数区间内的 member 个数
* `zrank key member` 返回 member 在有序集合 key 中的排名, 从 0 开始
* zset 底层数据结构包括两部分:  
    (1) hash 表: 关联元素 member 和 分数 score, 保障元素 member 的唯一性，可以通过元素 member 找到相应的 score 值; 
    (2) 跳跃表: 跳跃表的目的在于给元素 member 排序, 根据 score 的范围获取元素列表.、
### 5. Redis6 配置文件详解
* 空间单位
```conf
# 1k => 1000 bytes
# 1kb => 1024 bytes
# 1m => 1000000 bytes
# 1mb => 1024*1024 bytes
# 1g => 1000000000 bytes
# 1gb => 1024*1024*1024 bytes
```
* 忽略大小写
* 网络相关
```conf
# `bind 127.0.0.1 -::1` 默认是本机访问；
# `port 6379` 端口号
# tcp-backlog: 连接队列: 总和 = 未完成三次握手队列 + 已完成三次握手的队列。高backlog值(默认511) 从而避免慢客户端连接问题
# timeout 0 (默认 0 for never)
# tcp-keepalive 300 (心跳机制)
```
* `daemonize yes` 后台启动进程
* `pidfile /var/run/redis_6379.pid` 保存进程号
* `loglevel notice` 日志级别
* `logfile ""` 日志路径, 默认为空
* `database 16` 16个库
* `requirepass foobared` 密码(默认没有)
* `maxclients 10000` 设置 redis 客户端的最大连接数
* `maxmemory <bytes>` 设置最大内存(内存数据库), 达到最大内存后根据内存规则移除(如LRU)
* `dbfilename dump.rdb` rdb持久化保存的文件
* `dir ./` rdb 文件保存的目录
* `rdbcompression yes` 压缩 rdb 文件
* `rdbchecksum yes` 保持数据完整性
* `appendonly no` 默认关闭 AOF 持久化
* `appendfilename "appendonly.aof"` AOF 持久化的文件名
### 6. Redis6 的发布和订阅
* 类似 ROS 的发布订阅模式
* Redis 客户端可以订阅任意数量的频道(channel)
* `subscribe channel1 channel2...` 客户端订阅多个 channel 消息
* `publish channel1 message` 客户端向 channel1 发布消息 message
### 7. Redis6 新数据类型
#### 7.1 位图(Bitmaps)
* `setbit key offset value` 设置 bitmaps 中偏移量为 offset(从0开始) 的值为value(0或1)
* **当第一次初始化 bitmaps 时，如果偏移量很大，那么整个初始化过程执行会很慢**
* `getbit key offset` 获取 bitmaps 中偏移量为 offset(从0开始) 的值
* `bitcount key start end [BYTE|BIT]` 统计在 [start, end] 中 1 的数量 [BYTE] 是计算字节
* `bitop [and|or|not|xor] destkey [key1 key2...]` 对 key1 和 key2 ... 做操作后存储到 destkey
#### 7.2 HyperLogLog
* 多用于解决基数统计问题：`HyberLogLog`的优点是：在输入元素的数量或者体积非常大时，计算基数所需的空间总是固定的、并且是很小的。`HyberLogLog`只需要花费12 KB内存就可以计算接近 2^64 个不同元素的基数。
* `pfadd key element1 element2...` 添加指定元素到 HyberLogLog 中
* `pfcount key1 key2...` 统计多个 HyberLogLog 总共的近似基数
* `pfmerge destkey srckey1 srckey2...` 将一个或多个 `HyberLogLog` 合并后的结果存放到目标 key 中
#### 7.3 Geospatial
&emsp;&emsp;对 GEO 类型的支持，就是地理信息的表达，可表示为元素的 2 维坐标，在地图上就是经纬度。Redis 基于该类型，提供了经纬度设置，查询，范围查询，距离查询，经纬度 Hash 等常见操作。
* `geoadd key longtitude latitude member [longtitude latitude member...]` 添加多个地理位置(精度,维度,名称)；有效经度[-180,180], 维度[-85.05112878,85.05112878]
* `geopos key member [member...]` 获取 key 对应的 GEO 类型值中的 member 的坐标值
* `geodist key member1 member2 [m|km|ft|mi]` 获取两个位置之间的直线距离。单位： m(米, 默认值); km(千米); mi(英里); ft(英尺)
* `georadius key longtitude latitude radius m|km|ft|mi` 以给定的经纬度为圆心，找出半径 radius m|km|ft|mi 内的元素。例如`georadius china:city 111 46 1000 km`
#### 8. C++/Go 连接 Redis
* C++ 库 `hiredis`
```cpp
// 注意要安装动态链接库
// sudo apt-get install libhiredis-dev
// 编译时 -lhiredis
#include <iostream>
#include <hiredis/hiredis.h>

int main() {
    // 创建 Redis 上下文
    redisContext* context = redisConnect("localhost", 6379);
    if (context == nullptr || context->err) {
        if (context) {
            std::cout << "连接 Redis 失败：" << context->errstr << std::endl;
            redisFree(context);
        } else {
            std::cout << "无法分配 Redis 上下文" << std::endl;
        }
        return 1;
    }

    // 测试连接
    redisReply* reply = static_cast<redisReply*>(redisCommand(context, "PING"));
    if (reply == nullptr) {
        std::cout << "PING 命令执行失败" << std::endl;
        redisFree(context);
        return 1;
    }
    std::cout << "连接 Redis 成功：" << reply->str << std::endl;
    freeReplyObject(reply);

    // 关闭连接
    redisFree(context);

    return 0;
}
```
* Go 语言包 `github.com/go-redis/redis`
```Go
package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func main() {
	pv := CreatePhoneVerify("localhost:6379", "lizy66")

	phone := "15013044875"
	pv.GenVerifyCode(phone)

	var input string
	fmt.Printf("input verify code:\n")
	_, err := fmt.Scan(&input)
	if err != nil {
		fmt.Println("input error:", err)
		return
	}

	success := pv.CheckVerifyCode(phone, input)
	if success {
		fmt.Println("success")
	} else {
		fmt.Println("fail")
	}
	pv.Close()
}

type PhoneVerify struct {
	client *redis.Client
	gen    *rand.Rand
}

func CreatePhoneVerify(addr, pw string) *PhoneVerify {
	result := &PhoneVerify{}
	result.client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pw,
		DB:       0,
	})
	result.gen = rand.New(rand.NewSource(time.Now().UnixNano()))
	return result
}

func (pv *PhoneVerify) Close() {
	err := pv.client.Close()
	if err != nil {
		fmt.Println("close redis failed: ", err)
		return
	}
	fmt.Println("close redis success")
}

func (pv *PhoneVerify) CheckVerifyCode(phone, code string) bool {
	// code key
	var codeKey string = "VerifyCode" + phone + ":code"
	redisCode, _ := pv.client.Get(codeKey).Result()
	if redisCode == code {
		return true
	} else {
		return false
	}
}

func (pv *PhoneVerify) GenVerifyCode(phone string) {
	// phone count key
	var countKey string = "VerifyCode" + phone + ":count"
	// code key
	var codeKey string = "VerifyCode" + phone + ":code"

	count, _ := pv.client.Get(countKey).Result()
	cnt, _ := strconv.Atoi(count)
	if count == "" {
		pv.client.Set(countKey, 1, time.Duration(time.Duration(24*60*60).Seconds()))
	} else if cnt <= 2 {
		pv.client.Incr(countKey)
	} else if cnt > 2 {
		fmt.Println("this day can't send again")
		return
	}

	// set code key
	vcode := pv.GetCode()
	pv.client.Set(codeKey, vcode, time.Duration(time.Duration(120).Seconds()))
}

func (pv *PhoneVerify) GetCode() string {
	var code string
	for i := 0; i < 6; i++ {
		ch := pv.gen.Intn(10)
		code += strconv.Itoa(ch)
	}
	return code
}
```
### 9. Redis 事务操作
&emsp;&emsp;Redis事务是一个单独的隔离操作：事务中的所有命令都会序列化、按顺序地执行。事务在执行的过程中，不会被其他客户端发送来的命令请求所打断。  
&emsp;&emsp;Redis事务的主要作用就是**串联多个命令**防止别的命令插队。
#### 9.1 Multi、Exec、discard
&emsp;&emsp;从输入`Multi`命令开始，输入的命令都会依次进入命令队列中，但不会执行，直到输入`Exec`后，Redis会将之前的命令队列中的命令依次执行。  
&emsp;&emsp;组队的过程中可以通过`discard`来放弃组队。  
![image-1](https://github.com/lizyzzz/LearnRedis/blob/main/images/1.png)  
```
multi
set key1 value1
set key2 value2
set key3 value3
exec
multi
set key4 value4
set key6 value6
set key5 value5
discard
```
#### 9.2 事务的错误处理
* 组队中某个命令出现了报告错误，执行时整个的所有队列都会被取消。  
![image-2](https://github.com/lizyzzz/LearnRedis/blob/main/images/2.png)  
* 如果执行阶段某个命令报出了错误，则只有报错的命令不会被执行，而其他的命令都会执行，不会回滚。  
![image-3](https://github.com/lizyzzz/LearnRedis/blob/main/images/3.png)
#### 9.3 事务冲突的问题
* **悲观锁(Pessimistic Lock)**, 顾名思义，就是很悲观，每次去拿数据的时候都认为别人会修改，所以每次在拿数据的时候都会上锁，这样别人想拿这个数据就会block直到它拿到锁。传统的关系型数据库里边就用到了很多这种锁机制，比如行锁，表锁等，读锁，写锁等，都是在做操作之前先上锁。
* **乐观锁(Optimistic Lock)**, 顾名思义，就是很乐观，每次去拿数据的时候都认为别人不会修改，所以不会上锁，但是在更新的时候会判断一下在此期间别人有没有更新这个数据，如果没有就可以写数据，可以使用版本号等机制。乐观锁适用于多读的应用类型，这样可以提高吞吐量。**Redis就是利用这种check-and-set机制实现事务的**。
* `WATCH key [key ...]` 在执行multi之前，先执行`watch key1 [key2]`,可以监视一个(或多个) key ，如果在事务执行时这个(或这些) key 被其他命令所改动，那么事务将被discard。
```
---------- client 1 --------
set balance 100     // 1
watch balance       // 2
multi               // 4
incrby balance 10   // 6
exec                // 7
// (integer) 110
---------- client 1 --------
watch balance        // 3
multi                // 5
incrby balance 100   // 8
exec                 // 9
// (nil)
```
* `UNWATCH` 取消 WATCH 命令对所有 key 的监视。如果在执行 WATCH 命令之后，EXEC 命令或DISCARD 命令先被执行了的话，那么就不需要再执行UNWATCH 了。
#### 9.4 Redis 事务三大特性
* 单独的隔离操作  
&emsp;&emsp;事务中的所有命令都会序列化、按顺序地执行。事务在执行的过程中，不会被其他客户端发送来的命令请求所打断。 
* 没有隔离级别的概念  
&emsp;&emsp;队列中的命令没有提交之前都不会实际被执行，因为事务提交前任何指令都不会被实际执行
* 不保证原子性  
&emsp;&emsp;事务中如果有一条命令执行失败，其后的命令仍然会被执行，没有回滚  
**案例请看秒杀go案例**
### 8. Redis 持久化
#### 8.1 持久化之 RDB(Redis DataBase)
* **简介**：在指定的时间间隔内将内存中的**数据集快照**写入磁盘， 也就是行话讲的**Snapshot快照**，它恢复时是将快照文件直接读到内存里。
* 备份的执行原理  
&emsp;&emsp;Redis会单独创建（fork）一个子进程来进行持久化，会先将数据写入到 **一个临时文件中**，待持久化过程都结束了，再用这个**临时文件替换上次持久化好的文件**。 整个过程中，主进程是不进行任何IO操作的，这就确保了极高的性能 如果需要进行大规模数据的恢复，且对于数据恢复的完整性不是非常敏感，那RDB方式要比AOF方式更加的高效。**RDB的缺点是最后一次持久化后的数据可能丢失**。  
![image-4](https://github.com/lizyzzz/LearnRedis/blob/main/images/4.png)  
* `save` ：save时只管保存，其它不管，全部阻塞。手动保存。不建议。
* `bgsave`：Redis会在后台异步进行快照操作，快照同时还可以响应客户端请求。可以通过 `lastsave` 命令获取最后一次成功执行快照的时间
* 持久化 rdb 的**优势**  
&emsp;&emsp;(1)适合大规模的数据恢复  
&emsp;&emsp;(2)对数据完整性和一致性要求不高更适合使用  
&emsp;&emsp;(3)节省磁盘空间  
&emsp;&emsp;(4)恢复速度快  
* 持久化 rdb 的**劣势**  
&emsp;&emsp;(1)Fork的时候，内存中的数据被克隆了一份，大致2倍的膨胀性需要考虑。  
&emsp;&emsp;(2)虽然Redis在fork时使用了写时拷贝技术,但是如果数据庞大时还是比较消耗性能。  
&emsp;&emsp;(3)在备份周期在一定间隔时间做一次备份，所以如果Redis意外down掉的话，就会丢失最后一次快照后的所有修改。  
#### 8.2 持久化之 AOF(Append Only File)
* **简介**：**以日志的形式来记录每个写操作（增量保存）**，将Redis执行过的所有写指令记录下来(**读操作不记录**)， **只许追加文件但不可以改写文件**，redis启动之初会读取该文件重新构建数据，换言之，redis 重启的话就根据日志文件的内容将写指令从前到后执行一次以完成数据的恢复工作。
* 工作流程  
&emsp;&emsp;(1)客户端的请求写命令会被append追加到AOF缓冲区内；  
&emsp;&emsp;(2)AOF缓冲区根据AOF持久化策略[always,everysec,no]将操作sync同步到磁盘的AOF文件中；  
&emsp;&emsp;(3)AOF文件大小超过重写策略或手动重写时，会对AOF文件rewrite重写，压缩AOF文件容量；  
&emsp;&emsp;(4)Redis服务重启时，会重新load加载AOF文件中的写操作达到数据恢复的目的（**当 AOF和 RDB 都启用时，恢复时首选AOF**）  
![image-5](https://github.com/lizyzzz/LearnRedis/blob/main/images/5.png)  
* 当AOF文件出现异常/损坏时，通过`/usr/local/bin/redis-check-aof--fix appendonly.aof` 进行恢复
* 同步频率：  
&emsp;&emsp;(1) appendfsync always: 始终同步，每次Redis的写入都会立刻记入日志；性能较差但数据完整性比较好.  
&emsp;&emsp;(2) appendfsync everysec: 每秒同步，每秒记入日志一次，如果宕机，本秒的数据可能丢失。  
&emsp;&emsp;(3) appendfsync no: redis不主动进行同步，把同步时机交给操作系统。  
* `Rewrite 压缩`  
&emsp;&emsp;AOF采用文件追加方式，文件会越来越大为避免出现此种情况，新增了重写机制, 当AOF文件的大小超过所设定的阈值时，Redis就会启动AOF文件的内容压缩， 只保留可以恢复数据的最小指令集.可以使用命令`bgrewriteaof`  
* `Rewrite 原理`  
&emsp;&emsp;AOF文件持续增长而过大时，会fork出一个新进程来将文件重写(也是先写临时文件最后再rename)，redis4.0版本后的重写，就是**把 rdb 的快照**，以**二级制的形式**附在新的aof头部，作为已有的历史数据，替换掉原来的流水账操作。  
* `no-appendfsync-on-rewrite`：如果 no-appendfsync-on-rewrite=yes ,不写入 aof 文件只写入缓存，用户请求不会阻塞，但是在这段时间如果宕机会丢失这段时间的缓存数据。（降低数据安全性，提高性能）；如果 no-appendfsync-on-rewrite=no,  还是会把数据往磁盘里刷，但是遇到重写操作，可能会发生阻塞。（数据安全，但是性能降低）
* 何时重写  
&emsp;&emsp;(1)Redis会记录上次重写时的AOF大小，默认配置是当AOF文件大小是上次rewrite后大小的一倍且文件大于64M时触发  
&emsp;&emsp;(2)重写虽然可以节约大量磁盘空间，减少恢复时间。但是每次重写还是有一定的负担的，因此设定Redis要满足一定条件才会进行重写。  
&emsp;&emsp;(3)`auto-aof-rewrite-percentage`：设置重写的基准比例，文件达到100%时开始重写（文件是原来重写后文件的2倍时触发）  
&emsp;&emsp;(4)`auto-aof-rewrite-min-size`：设置重写的基准值，最小文件64MB。达到这个值开始重写。  
* **重写流程**  
&emsp;&emsp;(1)bgrewriteaof触发重写，判断是否当前有bgsave或bgrewriteaof在运行，如果有，则等待该命令结束后再继续执行。  
&emsp;&emsp;(2)主进程fork出子进程执行重写操作，保证主进程不会阻塞。  
&emsp;&emsp;(3)子进程遍历redis内存中数据到临时文件，客户端的写请求同时写入aof_buf缓冲区和aof_rewrite_buf重写缓冲区保证原AOF文件完整以及新AOF文件生成期间的新的数据修改动作不会丢失。  
&emsp;&emsp;(4) 1).子进程写完新的AOF文件后，向主进程发信号，父进程更新统计信息。2).主进程把aof_rewrite_buf中的数据写入到新的AOF文件。  
&emsp;&emsp;(5)使用新的AOF文件覆盖旧的AOF文件，完成AOF重写。  
![image-6](https://github.com/lizyzzz/LearnRedis/blob/main/images/6.png)  
* 持久化 AOF 的**优势**  
&emsp;&emsp;(1)备份机制更稳健，丢失数据概率更低。  
&emsp;&emsp;(2)可读的日志文本，通过操作AOF稳健，可以处理误操作。  
* 持久化 AOF 的**劣势**   
&emsp;&emsp;(1)比起RDB占用更多的磁盘空间。  
&emsp;&emsp;(2)恢复备份速度要慢。  
&emsp;&emsp;(3)每次读写都同步的话，有一定的性能压力。  
&emsp;&emsp;(4)存在个别Bug，造成恢复不能。  
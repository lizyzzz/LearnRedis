package main

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func Cluster() {
	// 创建 Redis 集群客户端
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			"192.168.3.190:6379",
			"192.168.3.190:6380",
			"192.168.3.190:6381",
			"192.168.3.190:6389",
			"192.168.3.190:6390",
			"192.168.3.190:6391",
		},
		/* 添加所有节点 */
	})
	defer client.Close()

	// 使用 ctx.Background() 作为上下文
	ctx := context.Background()

	// 测试连接
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Println("ping: ", err)
		return
	}
	fmt.Println(pong) // 应该打印 "PONG"

	// 设置值
	err = client.Set(ctx, "my_key", "my_value", 0).Err()
	if err != nil {
		fmt.Println("set: ", err)
		return
	}

	// 获取值
	val, err := client.Get(ctx, "my_key").Result()
	if err != nil {
		fmt.Println("get: ", err)
		return
	}
	fmt.Println("my_key:", val) // 应该打印 "my_key: my_value"
}

func main() {
	Cluster()

	// pv := CreatePhoneVerify("localhost:6379", "lizy66")

	// phone := "15013044875"
	// pv.GenVerifyCode(phone)

	// var input string
	// fmt.Printf("input verify code:\n")
	// _, err := fmt.Scan(&input)
	// if err != nil {
	// 	fmt.Println("input error:", err)
	// 	return
	// }

	// success := pv.CheckVerifyCode(phone, input)
	// if success {
	// 	fmt.Println("success")
	// } else {
	// 	fmt.Println("fail")
	// }
	// pv.Close()
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
	redisCode, _ := pv.client.Get(context.Background(), codeKey).Result()
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

	count, _ := pv.client.Get(context.Background(), countKey).Result()
	cnt, _ := strconv.Atoi(count)
	if count == "" {
		pv.client.Set(context.Background(), countKey, 1, time.Duration(time.Duration(24*60*60).Seconds()))
	} else if cnt <= 2 {
		pv.client.Incr(context.Background(), countKey)
	} else if cnt > 2 {
		fmt.Println("this day can't send again")
		return
	}

	// set code key
	vcode := pv.GetCode()
	pv.client.Set(context.Background(), codeKey, vcode, time.Duration(time.Duration(120).Seconds()))
}

func (pv *PhoneVerify) GetCode() string {
	var code string
	for i := 0; i < 6; i++ {
		ch := pv.gen.Intn(10)
		code += strconv.Itoa(ch)
	}
	return code
}

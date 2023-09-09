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

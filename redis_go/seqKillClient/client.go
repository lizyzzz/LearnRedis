package main

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		index := i
		go func() {
			wg.Add(1)
			defer wg.Done()
			seqClient(index)
		}()
		fmt.Printf("id %d start\n", i)
	}

	wg.Wait()
}

func seqClient(userId int) {
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Printf("-%d- connect failed\n", userId)
		return
	} else {
		fmt.Printf("-%d- connect success\n", userId)
	}
	time.Sleep(2 * time.Second)

	var str string = strconv.Itoa(userId)

	data := []byte(str)
	data = append(data, ',')

	for {
		_, errw := conn.Write(data) // send
		if errw != nil {
			fmt.Println(errors.New(fmt.Sprintf("write error: %s", errw)))
			break
		}

		buff := make([]byte, 2)
		n := 2
		for n > 0 {
			k, err := conn.Read(buff[2-n:])
			if err != nil {
				fmt.Printf("userId -%d- Read failed\n", userId)
				return
			}
			n -= k
		}
		// fmt.Println(string(buff))

		result, _ := strconv.Atoi(string(buff[:len(buff)-1]))
		if result == 1 {
			fmt.Printf("userId -%d- seqKill success\n", userId)
			break
		} else if result == 0 {
			// fmt.Printf("userId -%d- seqKill failed\n", userId)
		} else if result == 3 {
			fmt.Printf("userId -%d- seqKill stop\n", userId)
			break
		} else if result == 2 {
			fmt.Printf("userId -%d- has seqKill\n", userId)
			break
		}
	}
	conn.Close()
}

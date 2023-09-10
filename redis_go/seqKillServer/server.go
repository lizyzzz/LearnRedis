package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-redis/redis"
)

func main() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	s := CreateSeqServer("127.0.0.1:8888", "127.0.0.1:6379", "lizy66")
	go func() {
		<-sigChan
		fmt.Println("catch ctrl c..")
		s.Close()
		os.Exit(0)
	}()

	s.SetCargoNum("cargo1", 5)
	for {
		err := s.Accept()
		if err != nil {
			break
		}
	}
	s.Close()
}

type SeqServer struct {
	// TODO: redis connect pool
	client     *redis.Client
	listener   net.Listener
	cargoName  string
	successKey string
}

func CreateSeqServer(serverAddr, redisAddr, redisPw string) *SeqServer {
	result := &SeqServer{}
	var err error
	result.listener, err = net.Listen("tcp", serverAddr)
	if err != nil {
		fmt.Println("listen failed")
		return nil
	} else {
		fmt.Println("serve listening....")
	}
	result.client = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPw,
		DB:       0,
	})
	fmt.Println("redis connected....")

	return result
}

func (s *SeqServer) handleConnection(conn net.Conn) {
	// fmt.Printf("%s connected\n", conn.RemoteAddr().String())

	buff := bytes.NewBuffer(make([]byte, 128))
	for {
		_, err := conn.Read(buff.Bytes())
		if err != nil {
			if err == io.EOF {
				fmt.Printf("peer disconnected: %s\n", conn.RemoteAddr().String())
			} else {
				fmt.Printf("peer read failed: %s\n", conn.RemoteAddr().String())
			}
			conn.Close()
			return
		}

		res, _ := buff.ReadBytes(',')
		userId, err := strconv.Atoi(string(res[:len(res)-1]))
		success := s.SeqKill(userId)

		buff.Reset()
		if success == 1 {
			s.client.SAdd(s.successKey, userId)
			fmt.Printf("seqKill success: %d\n", userId)
			buff.WriteString("1,")
		} else if success == 0 {
			fmt.Printf("seqKill failed: %d\n", userId)
			buff.WriteString("0,")
		} else if success == 3 {
			fmt.Printf("seqKill stop: %d\n", userId)
			buff.WriteString("3,")
		} else {
			fmt.Printf("userId %d has seqKill\n", userId)
			buff.WriteString("2,")
		}
		// send result
		n, err := conn.Write(buff.Bytes())
		if err != nil {
			fmt.Printf("peer write failed: %s\n", conn.RemoteAddr().String())
			conn.Close()
			return
		} else {
			fmt.Printf("write %d bytes to %d\n", n, userId)
		}
	}

}

func (s *SeqServer) SeqKill(userId int) int {
	cargoKey := s.cargoName + ":count"

	txf := func(tx *redis.Tx) error {
		// check has kill
		hasKill := s.client.SIsMember(s.successKey, userId).Val()
		if hasKill {
			return fmt.Errorf("userId has seqKill")
		}
		// check cargo count
		n, err := tx.Get(cargoKey).Int()
		if err != nil {
			if err == redis.Nil {
				err = errors.New("seqKill not start")
			}
			return err
		} else {
			if n <= 0 {
				return errors.New("seqKill has end")
			}
		}
		// add command
		_, err = tx.Pipelined(func(pipe redis.Pipeliner) error {
			pipe.Decr(cargoKey)
			return nil
		})
		return err
	}
	// watch and exec
	err := s.client.Watch(txf, cargoKey)
	if err == redis.TxFailedErr {
		return 0
	} else if err != nil {
		// fmt.Printf("other err: %d, err: %v\n", userId, err)
		if err.Error() == "userId has seqKill" {
			return 2
		} else {
			// seqKill has end
			return 3
		}
	} else {
		return 1
	}
}

func (s *SeqServer) SetCargoNum(cargo string, cargoNum int) {
	s.cargoName = cargo
	s.successKey = "successUser"
	// cargoKey
	cargoKey := s.cargoName + ":count"
	success, err := s.client.SetNX(cargoKey, cargoNum, 0).Result()
	if err != nil {
		fmt.Println("set key failed")
	} else {
		if success {
			fmt.Printf("set cargokey: %d\n", cargoNum)
		} else {
			cntstr, _ := s.client.Get(cargoKey).Result()
			cnt, _ := strconv.Atoi(cntstr)
			fmt.Printf("cargokey is existed: %d\n", cnt)
		}
	}
}

func (s *SeqServer) Accept() error {
	conn, err := s.listener.Accept()
	if err != nil {
		// handle error
		return err
	}
	go s.handleConnection(conn)
	return nil
}

func (s *SeqServer) Close() {
	err := s.listener.Close()
	if err != nil {
		fmt.Println("close server failed: ", err)
		return
	} else {
		fmt.Println("close server success")
	}
	err = s.client.Close()
	if err != nil {
		fmt.Println("close redis failed: ", err)
		return
	}
	fmt.Println("close redis success")
}

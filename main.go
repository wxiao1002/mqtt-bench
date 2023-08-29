package main

import (
	"context"
	"flag"
	"log"
	"strconv"
	"sync/atomic"
	"time"

	c "mqtt-bench/client"
	"mqtt-bench/csv"
)

func main() {
	var (
		broker               = flag.String("broker", "tcp://127.0.0.1:1883", "MQTT broker 地址")
		csvPath              = flag.String("csv", "device_secret.csv", "设备用户密码配置csv文件地址")
		clients              = flag.Int("clients", 20, "客户端数量")
		benchmarkTime        = flag.Int("benchmarkTime", 1, "mqtt 压测时间，分钟")
		messageIntervalInSec = flag.Int("message-interval", 1, "生成消息间隔")
	)
	var clientPrefix string = "mqtt-benchmark"
	var qos int = 1
	var wait int = 1000
	flag.Parse()
	if *csvPath == "" {
		log.Fatalf("Invalid arguments: csv  should be is file path, given: %v", *csvPath)
		return
	}

	clientCSV, err := csv.ReaderCSV(*csvPath)
	if err != nil {
		panic(err)
	}
	if *clients < 1 {
		log.Fatalf("Invalid arguments: number of clients should be > 1, given: %v", clients)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Duration(*benchmarkTime) * time.Minute)
		cancel()
	}()
	for i, r := range clientCSV {
		if i >= *clients {
			break
		}
		c := &c.Client{
			ID:              i + 1,
			ClientID:        clientPrefix + strconv.Itoa(i),
			BrokerURL:       *broker,
			BrokerUser:      r.Username,
			BrokerPass:      r.Password,
			MsgQoS:          byte(qos),
			WaitTimeout:     time.Duration(wait) * time.Millisecond,
			MessageInterval: *messageIntervalInSec,
		}
		go c.RunBench(ctx)
	}
	<-ctx.Done()
	log.Printf("总消息数据:%v,Succ:%v,Error:%v,Timeout:%v", atomic.LoadInt64(&c.MsgSeq), atomic.LoadInt64(&c.Succ), atomic.LoadInt64(&c.Failure), atomic.LoadInt64(&c.Timeout))
	log.Println("exit program")
}

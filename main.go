package main

import (
	"flag"
	"log"
	"strconv"
	"time"

	. "mqtt-bench/client"
	"mqtt-bench/csv"
)

func main() {
	var (
		broker               = flag.String("broker", "tcp://10.50.6.1:1883", "MQTT broker endpoint as scheme://host:port")
		csvPath              = flag.String("csv", "device_secret.csv", "device Csv file path")
		clients              = flag.Int("clients", 1000, "client number")
		benchmarkTime        = flag.Int("benchmarkTime", 2, "mqtt benchmark time, in minutes")
		messageIntervalInSec = flag.Int("message-interval", 1, "Time interval in seconds to publish message")
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

	timeout := make(chan bool, 1)
	exit := make(chan bool)
	go func() {
		time.Sleep(time.Duration(*benchmarkTime) * time.Minute)
		timeout <- true
		exit <- true
	}()
	for i, r := range clientCSV {
		if i >= *clients {
			break
		}
		c := &Client{
			ID:              i,
			ClientID:        clientPrefix + strconv.Itoa(i),
			BrokerURL:       *broker,
			BrokerUser:      r.Username,
			BrokerPass:      r.Password,
			MsgQoS:          byte(qos),
			WaitTimeout:     time.Duration(wait) * time.Millisecond,
			MessageInterval: *messageIntervalInSec,
		}
		go c.RunBench(timeout)
	}
	<-exit
	log.Panicln("exit ing")
}

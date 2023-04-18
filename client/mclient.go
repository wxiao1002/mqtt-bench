package client

import (
	"log"
	"time"

	"github.com/GaryBoone/GoStats/stats"
	. "github.com/wxiao1002/mqtt-bench/core"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Mclient struct {
	ClientID    string
	Broker      string
	Username    string
	Password    string
	MsgTopic    string
	MsgPayload  string
	MsgSize     int
	MsgCount    int
	MsgQoS      byte
	WaitTimeout time.Duration
	MsgInterval int
}

func (c *Mclient) RunBench(res chan *BenchResults) {
	newMsgs := make(chan *Mmsg)
	pubMsgs := make(chan *Mmsg)
	doneGen := make(chan bool)
	donePub := make(chan bool)
	benchResults := new(BenchResults)
	started := time.Now()
	go c.createMsg(newMsgs, doneGen)
	go c.pubMsg(newMsgs, pubMsgs, doneGen, donePub)

	benchResults.ClientId = c.ClientID
	times := []float64{}
	for {
		select {
		case m := <-pubMsgs:
			if m.Failure {
				log.Printf("Client %v error publishing message: %v: at %v\n", c.ClientID, m.Topic, m.Sent.Unix())
				benchResults.Failures++
			} else {
				benchResults.Successes++
				times = append(times, m.Arrive.Sub(m.Sent).Seconds()*1000) // in milliseconds
			}
		case <-donePub:

			duration := time.Since(started)
			benchResults.MsgTimeMin = stats.StatsMin(times)
			benchResults.MsgTimeMax = stats.StatsMax(times)
			benchResults.MsgTimeMean = stats.StatsMean(times)
			benchResults.RunTime = duration.Seconds()
			benchResults.MsgsPerSec = float64(benchResults.Successes) / duration.Seconds()
			if c.MsgCount > 1 {
				benchResults.MsgTimeStd = stats.StatsSampleStandardDeviation(times)
			}
			res <- benchResults
			return
		}
	}
}

func (c *Mclient) createMsg(ch chan *Mmsg, done chan bool) {
	var payload interface{}
	if c.MsgPayload != "" {
		payload = c.MsgPayload
	} else {
		payload = make([]byte, c.MsgSize)
	}

	for i := 0; i < c.MsgCount; i++ {
		ch <- &Mmsg{
			Topic:   c.MsgTopic,
			QoS:     c.MsgQoS,
			Payload: payload,
		}
		time.Sleep(time.Duration(c.MsgInterval) * time.Second)
	}
	done <- true
}

func (c *Mclient) pubMsg(in, out chan *Mmsg, doneGen, donePub chan bool) {
	publishMsg := func(client mqtt.Client) {
		ctr := 0
		for {
			select {
			case m := <-in:
				m.Sent = time.Now()
				token := client.Publish(m.Topic, m.QoS, false, m.Payload)
				res := token.WaitTimeout(time.Second * 30)
				if !res {
					log.Printf("Client %v Timeout sending message: %v\n", c.ClientID, token.Error())
					m.Failure = true
				} else if token.Error() != nil {
					log.Printf("Client %v Error sending message: %v\n", c.ClientID, token.Error())
					m.Failure = true
				} else {
					m.Arrive = time.Now()
					m.Failure = false
				}
				out <- m
				ctr++
			case <-doneGen:
				donePub <- true
				return
			}
		}
	}

	opts := mqtt.NewClientOptions().
		AddBroker(c.Broker).
		SetClientID(c.ClientID).
		SetCleanSession(true).
		SetAutoReconnect(true).
		SetOnConnectHandler(publishMsg).
		SetConnectionLostHandler(func(client mqtt.Client, reason error) {
			log.Printf("Client %v lost connection to the broker: %v. Will reconnect...\n", c.ClientID, reason.Error())
		})
	opts.SetUsername(c.Username)
	opts.SetPassword(c.Password)

	client := mqtt.NewClient(opts)
	token := client.Connect()
	token.Wait()
	if token.Error() != nil {
		log.Printf("CLIENT %v had error connecting to the broker: %v\n", c.ClientID, token.Error())
	}
}

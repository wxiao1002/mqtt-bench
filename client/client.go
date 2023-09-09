package client

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"sync/atomic"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Message
type Message struct {
	Topic   string
	QoS     byte
	Payload string
}

var MsgSeq int64
var Succ int64
var Failure int64
var Timeout int64

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// mqtt client
type Client struct {
	ID              int
	Topic           string
	ClientID        string
	BrokerURL       string
	BrokerUser      string
	BrokerPass      string
	MsgQoS          byte
	WaitTimeout     time.Duration
	MessageInterval int
}

// 运行压测
func (c *Client) RunBench(ctx context.Context) {
	message := make(chan *Message,1)
	go c.generateMessages(message, ctx)
	go c.pubishMessages(message, ctx)
}

// 生成 消息
func (c *Client) generateMessages(out chan<- *Message, quit context.Context) {
	for {
		select {
		case <-quit.Done():
			return
		default:
			var payload = c.generatePayload()
			out <- &Message{
				Topic:   c.Topic,
				QoS:     c.MsgQoS,
				Payload: payload,
			}
			time.Sleep(time.Duration(c.MessageInterval) * time.Second)
		}
	}

}

// 发送消息
func (c *Client) pubishMessages(in <-chan *Message, quit context.Context) {
	onConnected := func(client mqtt.Client) {
		log.Printf("CLIENT %v  connected to the broker,Will publish msg\n", c.ClientID)
		for {
			select {
			case m := <-in:
				token := client.Publish(m.Topic, m.QoS, false, m.Payload)
				res := token.WaitTimeout(c.WaitTimeout)
				if !res {
					atomic.AddInt64(&Timeout, 1)
				} else if token.Error() != nil {
					atomic.AddInt64(&Failure, 1)
				} else {
					atomic.AddInt64(&Succ, 1)
				}
			case <-quit.Done():
				return
			}
		}
	}

	opts := mqtt.NewClientOptions().
		AddBroker(c.BrokerURL).
		SetClientID(c.ClientID).
		SetCleanSession(true).
		SetAutoReconnect(true).
		SetOnConnectHandler(onConnected).
		SetConnectRetryInterval(3).
		SetKeepAlive(60).
		SetConnectTimeout(time.Duration(20) * time.Second).
		SetConnectionLostHandler(func(client mqtt.Client, reason error) {
			log.Printf("CLIENT %v lost connection to the broker: %v. Will reconnect\n", c.ClientID, reason.Error())
		})
	if c.BrokerUser != "" && c.BrokerPass != "" {
		opts.SetUsername(c.BrokerUser)
		opts.SetPassword(c.BrokerPass)
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Printf("CLIENT %v connecting to the broker,has error: %v\n", c.ClientID, token.Error())
	}
}

type payload struct {
	N1  int64   `json:"N1"`
	N2  int64   `json:"N2"`
	N3  int64   `json:"N3"`
	N4  int64   `json:"N4"`
	N5  int64   `json:"N5"`
	N6  int64   `json:"N6"`
	N7  int64   `json:"N7"`
	N8  int64   `json:"N8"`
	N9  int64   `json:"N9"`
	N10 int64   `json:"N10"`
	F1  float64 `json:"F1"`
	F2  float64 `json:"F2"`
	F3  float64 `json:"F3"`
	F4  float64 `json:"F4"`
	F5  float64 `json:"F5"`
	F6  float64 `json:"F6"`
	F7  float64 `json:"F7"`
	F8  float64 `json:"F8"`
	F9  float64 `json:"F9"`
	F10 float64 `json:"F10"`
	S1  string  `json:"S1"`
	S2  string  `json:"S2"`
	S3  string  `json:"S3"`
	S4  string  `json:"S4"`
	S5  string  `json:"S5"`
	S6  string  `json:"S6"`
	S7  string  `json:"S7"`
	S8  string  `json:"S8"`
	S9  string  `json:"S9"`
	S10 string  `json:"S10"`
}

func (c *Client) generatePayload() string {
	atomic.AddInt64(&MsgSeq, 1)
	p := &payload{
		N1:  rand.Int63(),
		N2:  rand.Int63(),
		N3:  rand.Int63(),
		N4:  rand.Int63(),
		N5:  rand.Int63(),
		N6:  rand.Int63(),
		N7:  rand.Int63(),
		N8:  rand.Int63(),
		N9:  rand.Int63(),
		N10: rand.Int63(),
		F1:  rand.Float64(),
		F2:  rand.Float64(),
		F3:  rand.Float64(),
		F4:  rand.Float64(),
		F5:  rand.Float64(),
		F6:  rand.Float64(),
		F7:  rand.Float64(),
		F8:  rand.Float64(),
		F9:  rand.Float64(),
		F10: rand.Float64(),
		S1:  RandStringBytes(10),
		S2:  RandStringBytes(11),
		S3:  RandStringBytes(20),
		S4:  RandStringBytes(30),
		S5:  RandStringBytes(50),
		S6:  RandStringBytes(55),
		S7:  RandStringBytes(70),
		S8:  RandStringBytes(80),
		S9:  RandStringBytes(89),
		S10: RandStringBytes(8),
	}
	b1, err := json.Marshal(p)
	if err == nil {
		return string(b1)
	}
	return ""
}

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

package client

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync/atomic"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Message describes a message
type Message struct {
	Topic   string
	QoS     byte
	Payload string
}

var MSG_SEQ_F float64
var MSG_SEQ int64

// mqtt client
type Client struct {
	ID              int
	ClientID        string
	BrokerURL       string
	BrokerUser      string
	BrokerPass      string
	MsgQoS          byte
	WaitTimeout     time.Duration
	MessageInterval int
}

func (c *Client) RunBench(quit <-chan bool) {
	newMsgs := make(chan *Message)
	// start generator msg
	go c.genMessages(newMsgs, quit)
	// start publisher msg
	go c.pubMessages(newMsgs, quit)
}

func (c *Client) genMessages(ch chan *Message, quit <-chan bool) {
	for {
		select {
		case <-quit:
			return
		default:
			var payload = c.generatePayload()
			ch <- &Message{
				Topic:   "api/" + c.BrokerUser + "/attributes",
				QoS:     c.MsgQoS,
				Payload: payload,
			}
			time.Sleep(time.Duration(c.MessageInterval) * time.Second)
		}
	}

}

func (c *Client) pubMessages(in chan *Message, quit <-chan bool) {
	onConnected := func(client mqtt.Client) {
		ctr, successNum, errorNum, timeoutNum := 0, 0, 0, 0
		for {
			select {
			case m := <-in:
				token := client.Publish(m.Topic, m.QoS, false, m.Payload)
				res := token.WaitTimeout(c.WaitTimeout)
				if !res {
					log.Printf("CLIENT %v Timeout sending message: %v\n", c.ID, token.Error())
					timeoutNum++
				} else if token.Error() != nil {
					log.Printf("CLIENT %v Error sending message: %v\n", c.ID, token.Error())
					errorNum++
				} else {
					log.Printf("CLIENT %v is  send messages,seq: %d\n", c.ID, atomic.LoadInt64(&MSG_SEQ))
					successNum++
				}
				if ctr > 0 && ctr%100 == 0 {
					log.Printf("CLIENT %v published %v messages:[succ=%v,err=%v,timeout=%v] and keeps publishing...\n", c.ID, ctr, successNum, errorNum, timeoutNum)
				}
				ctr++
			case <-quit:
				client.Disconnect(0)
				return
			}
		}
	}

	opts := mqtt.NewClientOptions().
		AddBroker(c.BrokerURL).
		SetClientID(fmt.Sprintf("%s-%v", c.ClientID, c.ID)).
		SetCleanSession(true).
		SetAutoReconnect(true).
		SetOnConnectHandler(onConnected).
		SetConnectionLostHandler(func(client mqtt.Client, reason error) {
			log.Printf("CLIENT %v lost connection to the broker: %v. Will reconnect...\n", c.ID, reason.Error())
		})
	if c.BrokerUser != "" && c.BrokerPass != "" {
		opts.SetUsername(c.BrokerUser)
		opts.SetPassword(c.BrokerPass)
	}

	client := mqtt.NewClient(opts)
	token := client.Connect()
	token.Wait()

	if token.Error() != nil {
		log.Printf("CLIENT %v connecting to the broker,has error: %v\n", c.ID, token.Error())
	}
}

type payload struct {
	Seq int64   `json:"Seq"`
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
	atomic.AddInt64(&MSG_SEQ, 1)
	p := &payload{
		Seq: MSG_SEQ,
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
		F1:  MSG_SEQ_F + 0.000001,
		F2:  MSG_SEQ_F + 0.000002,
		F3:  MSG_SEQ_F + 0.000003,
		F4:  MSG_SEQ_F + 0.000004,
		F5:  MSG_SEQ_F + 0.000005,
		F6:  MSG_SEQ_F + 0.000006,
		F7:  MSG_SEQ_F + 0.000007,
		F8:  MSG_SEQ_F + 0.000008,
		F9:  MSG_SEQ_F + 0.000009,
		F10: MSG_SEQ_F + 0.00001,
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
	log.Printf("CLIENT %v SEQ: %v\n", c.ID, &MSG_SEQ)
	return ""
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

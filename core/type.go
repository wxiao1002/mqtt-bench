package core

import "time"

// Mmsg mqtt message
type Mmsg struct {
	Topic   string
	QoS     byte
	Payload interface{}
	Sent    time.Time
	Arrive  time.Time
	Failure bool
}

type BenchResults struct {
	ClientId    string  
	Successes   int64   
	Failures    int64   
	RunTime     float64 
	MsgTimeMin  float64 
	MsgTimeMax  float64 
	MsgTimeMean float64 
	MsgTimeStd  float64 
	MsgsPerSec  float64 
}

type TotalResults struct {
	Ratio           float64 
	Successes       int64   
	Failures        int64   
	TotalRunTime    float64 
	AvgRunTime      float64 
	MsgTimeMin      float64 
	MsgTimeMax      float64 
	MsgTimeMeanAvg  float64 
	MsgTimeMeanStd  float64 
	TotalMsgsPerSec float64 
	AvgMsgsPerSec   float64 
}


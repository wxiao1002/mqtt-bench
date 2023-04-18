package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/GaryBoone/GoStats/stats"
	client "github.com/wxiao1002/mqtt-bench/client"
	core "github.com/wxiao1002/mqtt-bench/core"
	csv "github.com/wxiao1002/mqtt-bench/csv"
)

func main() {
	var (
		broker               = flag.String("broker", "tcp://localhost:1883", "mqtt broker ")
		topic                = flag.String("topic", "/test", "MQTT topic for outgoing messages")
		payload              = flag.String("payload", "", "MQTT message payload")
		csvPath              = flag.String("csv", "", "Csv file path")
		size                 = flag.Int("size", 100, "Size of the messages payload (bytes)")
		count                = flag.Int("count", 100, "Number of messages to send per client")
		messageIntervalInSec = flag.Int("message-interval", 1, "Time interval in seconds to publish message")
	)
	flag.Parse()
	if *count < 1 {
		log.Fatalf("Invalid arguments: messages count should be > 1, given: %v", *count)
	}
	resCh := make(chan *core.BenchResults)
	start := time.Now()
	if *csvPath == "" {
		log.Fatalf("Invalid arguments: csv  should be is file path, given: %v", *csvPath)
		return
	}

	clients, err := csv.Reader(*csvPath)
	if err != nil {
		panic(err)
	}
	for _, r := range clients {
		c := &client.Mclient{
			ClientID:    r.ClientId,
			Broker:      *broker,
			Username:    r.Username,
			Password:    r.Password,
			MsgTopic:    *topic,
			MsgPayload:  *payload,
			MsgSize:     *size,
			MsgCount:    *count,
			MsgQoS:      1,
			MsgInterval: *messageIntervalInSec,
		}
		go c.RunBench(resCh)
	}
	clientNum := len(clients)
	results := make([]*core.BenchResults, clientNum)
	for i := 0; i < clientNum; i++ {
		results[i] = <-resCh
	}
	totalTime := time.Since(start)
	totals := calculateTotalResults(results, totalTime, clientNum)
	printResults(results, totals)
}

func calculateTotalResults(results []*core.BenchResults, totalTime time.Duration, sampleSize int) *core.TotalResults {
	totals := new(core.TotalResults)
	totals.TotalRunTime = totalTime.Seconds()

	msgTimeMeans := make([]float64, len(results))
	msgsPerSecs := make([]float64, len(results))
	runTimes := make([]float64, len(results))
	bws := make([]float64, len(results))

	totals.MsgTimeMin = results[0].MsgTimeMin
	for i, res := range results {
		totals.Successes += res.Successes
		totals.Failures += res.Failures
		totals.TotalMsgsPerSec += res.MsgsPerSec

		if res.MsgTimeMin < totals.MsgTimeMin {
			totals.MsgTimeMin = res.MsgTimeMin
		}

		if res.MsgTimeMax > totals.MsgTimeMax {
			totals.MsgTimeMax = res.MsgTimeMax
		}

		msgTimeMeans[i] = res.MsgTimeMean
		msgsPerSecs[i] = res.MsgsPerSec
		runTimes[i] = res.RunTime
		bws[i] = res.MsgsPerSec
	}
	totals.Ratio = float64(totals.Successes) / float64(totals.Successes+totals.Failures)
	totals.AvgMsgsPerSec = stats.StatsMean(msgsPerSecs)
	totals.AvgRunTime = stats.StatsMean(runTimes)
	totals.MsgTimeMeanAvg = stats.StatsMean(msgTimeMeans)
	if sampleSize > 1 {
		totals.MsgTimeMeanStd = stats.StatsSampleStandardDeviation(msgTimeMeans)
	}
	return totals
}

func printResults(results []*core.BenchResults, totals *core.TotalResults) {
	for _, res := range results {
		fmt.Printf("======= CLIENT %s =======\n", res.ClientId)
		fmt.Printf("Ratio:               %.3f (%d/%d)\n", float64(res.Successes)/float64(res.Successes+res.Failures), res.Successes, res.Successes+res.Failures)
		fmt.Printf("Runtime (s):         %.3f\n", res.RunTime)
		fmt.Printf("Msg time min (ms):   %.3f\n", res.MsgTimeMin)
		fmt.Printf("Msg time max (ms):   %.3f\n", res.MsgTimeMax)
		fmt.Printf("Msg time mean (ms):  %.3f\n", res.MsgTimeMean)
		fmt.Printf("Msg time std (ms):   %.3f\n", res.MsgTimeStd)
		fmt.Printf("Bandwidth (msg/sec): %.3f\n\n", res.MsgsPerSec)
	}
	fmt.Printf("========= TOTAL (%d) =========\n", len(results))
	fmt.Printf("Total Ratio:                 %.3f (%d/%d)\n", totals.Ratio, totals.Successes, totals.Successes+totals.Failures)
	fmt.Printf("Total Runtime (sec):         %.3f\n", totals.TotalRunTime)
	fmt.Printf("Average Runtime (sec):       %.3f\n", totals.AvgRunTime)
	fmt.Printf("Msg time min (ms):           %.3f\n", totals.MsgTimeMin)
	fmt.Printf("Msg time max (ms):           %.3f\n", totals.MsgTimeMax)
	fmt.Printf("Msg time mean mean (ms):     %.3f\n", totals.MsgTimeMeanAvg)
	fmt.Printf("Msg time mean std (ms):      %.3f\n", totals.MsgTimeMeanStd)
	fmt.Printf("Average Bandwidth (msg/sec): %.3f\n", totals.AvgMsgsPerSec)
	fmt.Printf("Total Bandwidth (msg/sec):   %.3f\n", totals.TotalMsgsPerSec)
}

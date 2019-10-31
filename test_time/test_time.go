package main

import (
	"fmt"
	"log"
	"time"
)

// timeTrack displays the elapsed time
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s = %d ms", name, elapsed/1000000)
}

func main() {
	start := time.Now()
	time.Sleep(time.Millisecond * 100)
	timeTrack(start, "elapsed:")
	time.Sleep(time.Second * 1)
	location, _ := time.LoadLocation("GMT")
	asciiTime := time.Now().In(location).Format(time.RFC1123)
	//------------
	ftime, _ := time.Parse(time.RFC1123, asciiTime)
	start = time.Now()
	downloadInterval := time.Duration(time.Second * 30)
	waitDuration := ftime.Add(downloadInterval).Sub(start)
	//latency = MaxDuration(latency, 0)
	//waitDuration := MaxDuration(time.Duration(downloadInterval)*time.Second-latency, 0)
	fmt.Println(waitDuration)

}

package main

import (
	"fmt"
	"time"

	//"os"
	"github.com/beevik/ntp"
)

func main() {
	p := fmt.Println
	start := time.Now()
	ntpTime, _ := ntp.Time("pool.ntp.org") //0.beevik-ntp.pool.ntp.org")
	p(ntpTime, time.Since(start))
	p(time.Now(), time.Since(ntpTime))
	time.Sleep(time.Millisecond * 200)
	p(time.Now(), time.Since(ntpTime))
	response, _ := ntp.Query("0.beevik-ntp.pool.ntp.org")
	aNtpTime := time.Now().Add(response.ClockOffset)
	p(aNtpTime, time.Since(start), response.ClockOffset)
	p(time.Now(), time.Since(aNtpTime))
	p(response)
}

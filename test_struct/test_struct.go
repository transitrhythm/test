package main

import (
	"fmt"
)

type httpSchedule struct {
	hourTime string
	url      string
	dst      string
}

func main() {
	testSchedule := httpSchedule{"09:30:00", "http://victoria.mapstrat.com/current/trip_reference.txt", "../data/Victoria/trip_reference.txt"}
	fmt.Println(testSchedule)

}

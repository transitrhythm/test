package main

import (
	"github.com/patrickbr/gtfsparser/gtfs" //"github.com/geops/gtfsparser/gtfs"
	"sort"
	"time"
	"fmt"
)

// Block -
type Block struct {
	BlockID string
	StartAt gtfs.Time
	EndAt   gtfs.Time
	Trips   []*gtfs.Trip
}

// BlockDay -
type BlockDay struct {
	Blocks []*Block
}

// BlockCalendar -
type BlockCalendar struct {
	BlockDays []*BlockDay
}

// Trips -
type Trips []*gtfs.Trip

// Len -
func (s Trips) Len() int { return len(s) }

// Swap -
func (s Trips) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// ByDepartureTime -
type ByDepartureTime struct{ Trips }

// Less -
func (s ByDepartureTime) Less(i, j int) bool {
	return toSeconds(s.Trips[i].StopTimes[0].Departure_time) < toSeconds(s.Trips[j].StopTimes[0].Departure_time)
}

// Blocks -
type Blocks []*Block

// Len -
func (s Blocks) Len() int { return len(s) }

// Swap -
func (s Blocks) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// ByStartTime -
type ByStartTime struct{ Blocks }

// toSeconds -
func toSeconds(input gtfs.Time) int {
	return (3600 * int(input.Hour)) + (60 * int(input.Minute)) + int(input.Second)
}

// Less -
func (s ByStartTime) Less(i, j int) bool {
	return toSeconds(s.Blocks[i].StartAt) < toSeconds(s.Blocks[j].StartAt)
}

func sortBlockCalendar(blockCalendar BlockCalendar) {

	for _, day := range blockCalendar.BlockDays {
		if day != nil {
			if day.Blocks != nil {
				sort.Sort(ByStartTime{day.Blocks})
			}
			for _, block := range day.Blocks {
				if block.Trips != nil {
					sort.Sort(ByDepartureTime{block.Trips})
				}
			}
		}
	}
}

func daysThisMonth() (days int) {
	daysInMonth := [12]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	year, month, _ := time.Now().Date()
	days = daysInMonth[month-1]
	isLeapYear := time.Date(year, time.December, 31, 0, 0, 0, 0, time.Local).YearDay() == 366
	if month == time.February && isLeapYear {
		days = 29
	}
	return days
}

// Convert GTFS Date to Time @ 12:00:00 noon local time
func toTime(date gtfs.Date) time.Time {
	return time.Date(int(date.Year), time.Month(date.Month), int(date.Day), 12, 0, 0, 0, time.Local)
}

func dayOfWeek(date gtfs.Date) (day int) {
	return int(toTime(date).Weekday())
}

func weeksThisMonth() (weeks int) {
	year, month, _ := time.Now().Date()
	first := dayOfWeek(toDate(year, int(month), 1))
	last := dayOfWeek(toDate(year, int(month), daysThisMonth()))
	weeks = (daysThisMonth() + (first + last)) / 7
	return weeks
}

func toDate(year, month, day int) gtfs.Date {
	return gtfs.Date{int8(day), int8(month), int16(year)}
}

func main() {
	weeks := weeksThisMonth()
	fmt.Println(weeks)
	blockCalendar := BlockCalendar{}
	sortBlockCalendar(blockCalendar)
}

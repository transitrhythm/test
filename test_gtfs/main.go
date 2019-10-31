package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	proto "github.com/golang/protobuf/proto"
	"github.com/patrickbr/gtfsparser" // "github.com/geops/gtfsparser" // "github.com/geops/gtfsparser/gtfs" //

	// "transitrhythm.com/gtfs/realtime/server/transit_realtime"
	"github.com/thingful/transit_realtime"
)

const defaultDownloadInterval = time.Duration(time.Second * 1)
const filespec = "test_vehicle_position.pb"
const floatPrecision = 3

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage : %s <ZIPfile> <StopCode> \n", os.Args[0])
		os.Exit(0)
	}

	testSuite2D(floatPrecision)

	// Identify the input GTFS Zip file
	zipFile := os.Args[1]

	// Parse the GTFS Zip file
	feed := gtfsparser.NewFeed()
	feed.Parse(zipFile)
	fmt.Printf("Done, parsed: %d feeds, %d agencies, %d stops, %d routes, %d trips, %d shapes\n\n",
		len(feed.FeedInfos), len(feed.Agencies), len(feed.Stops), len(feed.Routes), len(feed.Trips), len(feed.Shapes))

	feeds := ParseFeeds(feed)
	for _, v := range feeds {
		if v.Active == true {
			agencies := ParseAgencies(feed)
			for k, v := range agencies {
				fmt.Printf("[%d] %s : <%s> - %s\n", k, v.ID, v.Name, v.Timezone)
			}
			routes, _ := ParseRoutes(feed)
			for k, v := range routes {
				fmt.Printf("[%d] %s - %s : %s - <%s> - <%s>\n", k, v.Type, v.ID, v.ShortName, v.LongName, v.Desc)
			}
			stops, _ := ParseStops(feed, floatPrecision)
			for _, v := range stops {
				fmt.Printf("[%s] %s : <%s> - <%s> - <%s> - <%s> - <%s> (@ %.3f,%.3f)\n", v.ID, v.Code, v.Desc, v.Name, v.LocationLabel, v.Parent, v.Zone, v.Location.x, v.Location.y)
			}
			shapes, _ := ParseShapes(feed, floatPrecision)
			// ParseTrips -
			trips, _ := ParseTrips(feed, shapes)
			for _, trip := range trips {
				speed, duration := AverageSpeed(trip.Length, trip.StartTime, trip.EndTime)
				fmt.Printf("[%s] %s : <%s> - %s %s %s =%.1f min, %.4f km (%.4f kph) <%s> - <%s> - <%s>\n", trip.ID, trip.BlockID, trip.RouteID, GtfsServiceInOperationStatus(*trip.Service), Timestamp(trip.StartTime), Timestamp(trip.EndTime), duration, trip.Length, speed, trip.Headsign, GtfsTripDirectionStatus(int(trip.DirectionID)), trip.ShapeID)
			}
			for _, shape := range feed.Shapes {
				for _, trip := range trips {
					if shape.Id == trip.ShapeID {
						kmlFileWrite("../../data/Vancouver/kml/kml_"+shape.Id+".kml", feed, trip.ID, shape.Id)
						break
					}
				}
			}

			stopTimes, _ := ParseStopTimes(feed, stops, shapes, floatPrecision)
			for _, v := range stopTimes {
				//fmt.Printf("Trip:[%s] %s : %s - %s - %s\n", v.tripID, v.Headsign, GtfsTripDirectionStatus(int(v.Direction_id)), GtfsServiceInOperationStatus(*v.Service))
				times := v.times
				for i := 0; i < len(times); i++ {
					fmt.Printf("StopTimes:[%s] %d : %s - %s\n", times[i].StopID, times[i].Sequence, Timestamp(times[i].ArrivalTime), Timestamp(times[i].DepartureTime))
				}
			}

			// ActiveTrips -
			activeTrips, _ := ActiveTrips(trips)
			fmt.Println(activeTrips)
		}
	}
	fmt.Println("Start")
	_, err := os.Stat(filespec)
	if err == nil {
		data, err := ioutil.ReadFile(filespec)
		if err != nil {
			fmt.Printf("Error reading vehicle_position.pb file: %s\n", err)
		}
		fmt.Println("File size:", len(data))
		go Process(data, len(data))
	}

	for {
		time.Sleep(defaultDownloadInterval)
	}
}

// Process -
func Process(data []byte, size int) {
	var (
		expectedEntityLength = 1
		expectedEntityID     = "1"
		// expectedTripID       = "t0"
	)

	feed := transit_realtime.FeedMessage{}

	err := proto.Unmarshal(data, &feed)
	if err != nil {
		fmt.Printf("Error unmarshaling data: %s\n", err)
	}

	if len(feed.Entity) != expectedEntityLength {
		fmt.Printf("Expected entity length: %d, got: %d\n", expectedEntityLength, len(feed.Entity))
	}

	entity := feed.Entity[0]
	if *entity.Id != expectedEntityID {
		fmt.Printf("Expected entity id: %v, got: %v\n", expectedEntityID, entity.Id)
	}
}

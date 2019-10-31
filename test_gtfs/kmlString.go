package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/gershwinlabs/gokml"
	"github.com/patrickbr/gtfsparser" // "github.com/geops/gtfsparser" //
)

func toKmlPoint(lat float64, lon float64, alt float64) (point gokml.Point) {
	point.Lat = lat
	point.Lon = lon
	point.Alt = alt
	return point
}

// KmlStopPoint -
func KmlStopPoint(feed *gtfsparser.Feed, stopID string) (point gokml.Point) {
	for _, v := range feed.Stops {
		if stopID == v.Id {
			point = toKmlPoint(float64(v.Lat), float64(v.Lon), 0.0)
			break
		}
	}
	return point
}

// KmlShapePoint -
func KmlShapePoint(feed *gtfsparser.Feed, shapeID string, sequence int) (point gokml.Point) {
	for _, v := range feed.Shapes {
		if shapeID == v.Id {
			for _, v1 := range v.Points {
				if sequence == v1.Sequence {
					point = toKmlPoint(float64(v1.Lat), float64(v1.Lon), 0.0)
				}
			}
		}
	}
	return point
}

// KmlString -
func KmlString(feed *gtfsparser.Feed, tripID, shapeID string) (result string) {
	kml := gokml.NewKML("GTFS Trip #" + tripID + " Route + Stops")
	folder := gokml.NewFolder("Trip #"+tripID+" Folder", "Contains a list of route waypoints")
	kml.AddFeature(folder)

	stopStyle := gokml.NewStyle("StopStyle", 240, 0, 255, 0)
	stopStyle.SetIconURL("http://maps.google.com/mapfiles/kml/shapes/bus.png")
	folder.AddFeature(stopStyle)

	for _, v := range feed.Trips {
		if tripID == v.Id {
			for _, v1 := range v.StopTimes {
				fmt.Print(";", v1.Sequence)
				point := KmlStopPoint(feed, v1.Stop.Id)
				pm := gokml.NewPlacemark(v1.Stop.Name, "#"+strconv.Itoa(v1.Sequence)+"-"+v1.Stop.Code+v1.Headsign+fmt.Sprintf("-%.4f", v1.Shape_dist_traveled)+"km", &point)
				pm.SetStyle("PlaceStyle")
				folder.AddFeature(pm)
			}
		}
	}

	routeStyle := gokml.NewStyle("RouteStyle", 240, 255, 0, 0)
	folder.AddFeature(routeStyle)
	routeShape := gokml.NewLineString()
	pm := gokml.NewPlacemark("Shape #"+shapeID, "", routeShape)

	for _, v := range feed.Trips {
		if tripID == v.Id {
			for _, v1 := range feed.Shapes {
				if v.Shape.Id == v1.Id {
					for i := 1; i <= len(v1.Points); i++ {
						point := KmlShapePoint(feed, v.Shape.Id, i)
						routeShape.AddPoint(&point)
					}
					break
				}
			}
			break
		}
	}
	folder.AddFeature(pm)
	return kml.Render()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func kmlFileRead(filespec string) (data []byte, err error) {
	_, err = os.Stat(filespec)
	if err == nil {
		data, err := ioutil.ReadFile(filespec)
		if err != nil {
			fmt.Printf("Error reading file: %s\n", err)
		}
		fmt.Println("File size:", len(data))
	}
	return data, err
}

func kmlFileWrite(filename string, feed *gtfsparser.Feed, tripID string, shapeID string) (written int, err error) {
	file, err := os.Create(filename)
	check(err)
	defer file.Close()

	kmlString := KmlString(feed, tripID, shapeID)
	written, err = file.WriteString(kmlString)
	check(err)
	file.Sync()
	return written, err
}

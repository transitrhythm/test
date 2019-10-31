package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gershwinlabs/gokml"
)

// Place -
type Place struct {
	name  string
	desc  string
	coord gokml.Point
}

var places = []Place{
	Place{"Manhattan", "The Big Apple", gokml.Point{40.67, -73.9, 0.0}},
	Place{"London", "The City", gokml.Point{51.51, 0.1275, 0.0}},
	Place{"Paris", "The City of Light", gokml.Point{48.85, 2.35, 0.0}},
	Place{"Tokyo", "東京", gokml.Point{35.69, 139.7, 0.0}},
}

var coloradoShape = []gokml.Point{
	gokml.Point{41.071904, -101.868843, 0.0},
	gokml.Point{36.926393, -101.868843, 0.0},
	gokml.Point{36.926393, -109.279635, 0.0},
	gokml.Point{41.071904, -109.279635, 0.0},
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage : %s <KMLfile> \n", os.Args[0])
		os.Exit(0)
	}

	file, err := os.Create(os.Args[1])
	check(err)
	defer file.Close()

	k := gokml.NewKML("Test KML")
	f := gokml.NewFolder("Test Folder", "This is a test folder")
	k.AddFeature(f)

	placeStyle := gokml.NewStyle("PlaceStyle", 240, 0, 255, 0)
	placeStyle.SetIconURL("http://maps.google.com/mapfiles/kml/paddle/wht-circle.png")
	f.AddFeature(placeStyle)

	flightStyle := gokml.NewStyle("FlightStyle", 240, 255, 0, 0)
	f.AddFeature(flightStyle)

	for _, v := range places {
		point := gokml.NewPoint(v.coord.Lat, v.coord.Lon, v.coord.Alt)
		pm := gokml.NewPlacemark(v.name, v.desc, point)
		pm.SetStyle("PlaceStyle")
		f.AddFeature(pm)
	}

	flightPath := gokml.NewLineString()
	pm := gokml.NewPlacemark("Flight Path", "", flightPath)
	for k := range coloradoShape {
		flightPath.AddPoint(&coloradoShape[k])
	}
	pm.SetStyle("FlightStyle")
	f.AddFeature(pm)

	states := gokml.NewStyle("StateStyle", 240, 0, 0, 255)
	f.AddFeature(states)

	// Create polygon
	colorado := gokml.NewPolygon()
	pm = gokml.NewPlacemark("Colorado", "The Centennial State", colorado)
	pm.SetStyle("StateStyle")
	pm.SetTime(time.Now().Add(-10*time.Hour), time.Now())
	f.AddFeature(pm)
	for _, v := range coloradoShape {
		colorado.AddPoint(gokml.NewPoint(v.Lat, v.Lon, v.Alt))
	}
	// Save file
	file.WriteString(k.Render())
	file.Sync()
}

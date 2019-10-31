// Copyright 2019 Bevan Thistlethwaite. All rights reserved.

/*
	Package flag implements command-line flag parsing.

	Usage

	Define flags using flag.String(), Bool(), Int(), etc.

	This declares an integer flag, -flagname, stored in the pointer ip, with type *int.
		import "flag"
		var ip = flag.Int("flagname", 1234, "help message for flagname")
*/

package main // kml-xmldom

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/subchen/go-xmldom"
	//"gitlab.com/stone.code/xmldom-go"
	"github.com/im7mortal/UTM"
)

var (
	routeWaypoints   = []Cartesian{}
	routeWaypointMap = make(map[string]map[string][]Cartesian)
)

// Cartesian -
type Cartesian struct {
	x, y float64
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func toLocation(lat float32, lon float32, precision int) (location Cartesian, timeZone int, timeLetter string, err error) {
	x, y, timeZone, timeLetter, err := UTM.FromLatLon(float64(lat), float64(lon), lat > 0)
	location.x = toFixed(x, precision)
	location.y = toFixed(y, precision)
	return location, timeZone, timeLetter, err
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	log.SetFlags(0)
	flag.Parse()
	if flag.NArg() < 1 {
		log.Fatalf("Usage: %s [-ns xmlns] file.xml ...", os.Args[0])
	}

	var kml Kml
	err := kmlFileParseIntersections(os.Args[1], &kml)
	check(err)
	intersections, err := findIntersectionList(&kml, "DUNBAR ST", 4)
	intersections, err = sortIntersectionList(intersections, SortByLat)
	for _, v := range intersections {
		fmt.Printf("%s @ %s { %.4f, %.4f }\n", v.crossStreets[0], v.crossStreets[1], v.location.x, v.location.y)
	}

	var routeKml RouteKml
	kmlFileParseRoutes(os.Args[2], &routeKml)

	//docs := make([]xmldom.Document, 0, flag.NArg())
	var docs []*xmldom.Document
	start := time.Now()
	begin := start
	for _, filename := range flag.Args() {
		document := xmldom.Must(xmldom.ParseFile(filename))
		docs = append(docs, document)
		root := document.Root
		node := root
		fmt.Printf("name = %v\n", node.Name)
		fmt.Printf("attributes.len = %v\n", len(node.Attributes))
		fmt.Printf("children.len = %v\n", len(node.Children))
		fmt.Printf("root = %v\n", node == node.Root())
		// find all children
		fmt.Printf("children = %v\n", len(node.Query("//*")))
		// find node matched tag name
		nodeList := node.Query("//Folder")
		for _, node := range nodeList {
			nodeID := node.GetAttributeValue("id")
			routeWaypointMap[nodeID] = make(map[string][]Cartesian)
			//*fmt.Printf("%v: id = %v\n", node.Name, nodeId)
			placemarks := node.GetChildren("Placemark")
			for _, placemark := range placemarks {
				name := placemark.GetChild("name")
				// desc := placemark.GetChild("description")
				line := placemark.GetChild("LineString")
				point := placemark.GetChild("Point")
				var coord *xmldom.Node
				if line != nil {
					coord = line.GetChild("coordinates")
				} else if point != nil {
					coord = point.GetChild("coordinates")
				}
				//fmt.Printf("\n%v: id = %v desc: %v\n", placemark.Name, name.Text, desc.Text)
				s := strings.Fields(coord.Text)
				//fmt.Println("Len:", len(s), s)
				fmt.Println(time.Since(start), " [", nodeID, "][", name.Text, "]", s) //, desc.Text)
				start = time.Now()
				for i := 0; i < len(s); i++ {
					ll := strings.Split(s[i], ",")
					if longitude, err := strconv.ParseFloat(ll[0], 64); err == nil {
						if latitude, err := strconv.ParseFloat(ll[1], 64); err == nil {
							location, _, _, _ := toLocation(float32(latitude), float32(longitude), 4)
							//fmt.Printf("EN:[%v] %.2f;%.2f ", i, location.x, location.y)
							//*routeWaypoints = append(routeWaypoints, location)
							routeWaypointMap[nodeID][name.Text] =
								append(routeWaypointMap[nodeID][name.Text], location)
						}
					}
				}
			}
		}
		fmt.Println("\nTotal elapsed time: ", time.Since(begin))
	}

	for {

	}
}

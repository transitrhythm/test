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
	//"github.com/im7mortal/UTM"
	"github.com/wroge/wgs84"
)

type EastingNorthing struct {
	Easting  float64
	Northing float64
}

var (
	routeWaypoints   = []EastingNorthing{}
	routeWaypointMap = make(map[string]map[string][]EastingNorthing)
)

const (
	DEFAULT_XML_HEADER = `<?xml version="1.0" encoding="UTF-8"?>`
)

func main() {
	log.SetFlags(0)
	flag.Parse()
	if flag.NArg() < 1 {
		log.Fatalf("Usage: %s [-ns xmlns] file.xml ...", os.Args[0])
	}

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
			//*fmt.Printf("%v: id = %v\n", node.Name, nodeId)
			placemarks := node.GetChildren("Placemark")
			for _, placemark := range placemarks {
				name := placemark.GetChild("name")
				desc := placemark.GetChild("description")
				coord := placemark.GetChild("LineString").GetChild("coordinates")
				//fmt.Printf("\n%v: id = %v desc: %v\n", placemark.Name, name.Text, desc.Text)
				s := strings.Fields(coord.Text)
				//fmt.Println("Len:", len(s), s)
				baseCoord := strings.Split(s[0], ",")
				baseLongitude, _ := strconv.ParseFloat(baseCoord[0], 64)
				//mercatorZone := math.Floor((baseLongitude / 6) + 31)
				mercatorZone := math.Round((baseLongitude + 183) / 6)
				conversion := wgs84.ToUTM(mercatorZone, true).Round(2)
				fmt.Println(time.Since(start), " [", nodeID, "][", name.Text, "]", desc.Text)
				start = time.Now()
				routeWaypointMap[nodeID] = make(map[string][]EastingNorthing)
				for i := 0; i < len(s); i++ {
					ll := strings.Split(s[i], ",")
					if longitude, err := strconv.ParseFloat(ll[0], 64); err == nil {
						if latitude, err := strconv.ParseFloat(ll[1], 64); err == nil {
							easting, northing, _ := conversion(longitude, latitude, 0)
							eastingNorthing := EastingNorthing{easting, northing}
							//fmt.Printf("EN:[%v] %.2f;%.2f ", i, eastingNorthing.Easting, eastingNorthing.Northing)
							//*routeWaypoints = append(routeWaypoints, eastingNorthing)
							routeWaypointMap[nodeID][name.Text] =
								append(routeWaypointMap[nodeID][name.Text], eastingNorthing)
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

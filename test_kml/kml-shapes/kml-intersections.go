// Copyright 2019 Bevan Thistlethwaite. All rights reserved.
package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Point -
type Point struct {
	Coordinates string `xml:"coordinates" json:"coordinates"`
}

// Placemark -
type Placemark struct {
	ID          string `xml:"id" json:"id"`
	Name        string `xml:"name" json:"name"`
	Description string `xml:"description" json:"description"`
	Point       *Point `xml:"Point" json:"Point"`
}

// Folder -
type Folder struct {
	ID         string      `xml:"id" json:"id"`
	Name       string      `xml:"name" json:"name"`
	Placemarks []Placemark `xml:"Placemark" json:"Placemarks"`
}

// Document -
type Document struct {
	Folder Folder `xml:"Folder" json:"Folder"`
}

// Kml -
type Kml struct {
	XMLName  xml.Name `xml:"kml" json:"-"`
	Document Document `xml:"Document" json:"Document"`
}

// func parseIntersections(data []byte) (jsonData []byte, err error) {
// 	var kml Kml
// 	xml.Unmarshal(data, &kml)
// 	jsonData, err = json.Marshal(data)
// 	fmt.Println(jsonData)
// 	return jsonData, err
// }

// StreetIntersection -
type StreetIntersection struct {
	crossStreets [2]string
	location     Cartesian
}

func kmlFileParseIntersections(filespec string, kml *Kml) (err error) {
	_, err = os.Stat(filespec)
	if err == nil {
		data, err := ioutil.ReadFile(filespec)
		// fmt.Println(data)
		if err == nil {
			xml.Unmarshal(data, &kml)
		}
		fmt.Println("File size:", len(data))
	}
	return err
}

// SortKey -
type SortKey int

// SortKey Elements -
const (
	SortByName SortKey = iota
	SortByLat
	SortByLon
)

// By -
type By func(s1, s2 *StreetIntersection) bool

// StreetIntersections -
type StreetIntersections []*StreetIntersection

// Len -
func (s StreetIntersections) Len() int { return len(s) }

// Swap -
func (s StreetIntersections) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// ByName -
type ByName struct{ StreetIntersections }

// Less -
func (s ByName) Less(i, j int) bool {
	return s.StreetIntersections[i].crossStreets[1] < s.StreetIntersections[j].crossStreets[1]
}

// ByLat -
type ByLat struct{ StreetIntersections }

// Less -
func (s ByLat) Less(i, j int) bool {
	return s.StreetIntersections[i].location.y < s.StreetIntersections[j].location.y
}

// ByLon -
type ByLon struct{ StreetIntersections }

// Less -
func (s ByLon) Less(i, j int) bool {
	return s.StreetIntersections[i].location.x < s.StreetIntersections[j].location.x
}

func sortIntersectionList(intersections []*StreetIntersection, sortKey SortKey) (sorted []*StreetIntersection, err error) {
	sorted = intersections
	switch sortKey {
	case SortByName:
		sort.Sort(ByName{sorted})
	case SortByLat:
		sort.Sort(ByLat{sorted})
	case SortByLon:
		sort.Sort(ByLon{sorted})
	}
	return sorted, err
}

func toFloat32(ascii string) (value float32) {
	value1, _ := strconv.ParseFloat(ascii, 32)
	return float32(value1)
}

func toLocationfromString(coordinates string, precision int) (location Cartesian) {
	r := csv.NewReader(strings.NewReader(coordinates))
	coord, _ := r.Read()
	location, _, _, _ = toLocation(toFloat32(coord[1]), toFloat32(coord[0]), precision)
	return location
}

func findIntersectionList(kml *Kml, street string, precision int) (intersections []*StreetIntersection, err error) {
	for _, v := range kml.Document.Folder.Placemarks {
		if strings.Contains(v.Name, street) {
			split := strings.Split(v.Name, " AND ")
			intersection := StreetIntersection{}
			if split[0] == street {
				intersection.crossStreets[0] = split[0]
				intersection.crossStreets[1] = split[1]
			} else {
				intersection.crossStreets[0] = split[1]
				intersection.crossStreets[1] = split[0]
			}
			intersection.location = toLocationfromString(v.Point.Coordinates, precision)
			intersections = append(intersections, &intersection)
		}
	}
	return intersections, err
}

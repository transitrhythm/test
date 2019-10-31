// Copyright 2019 Bevan Thistlethwaite. All rights reserved.
package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

// LineString -
type LineString struct {
	Coordinates string `xml:"coordinates" json:"coordinates"`
}

// RoutePlacemark -
type RoutePlacemark struct {
	ID          string      `xml:"id" json:"id"`
	Name        string      `xml:"name" json:"name"`
	Description string      `xml:"description" json:"description"`
	LineString  *LineString `xml:"LineString" json:"LineString"`
}

// RouteFolder -
type RouteFolder struct {
	ID         string           `xml:"id" json:"id"`
	Name       string           `xml:"name" json:"name"`
	Placemarks []RoutePlacemark `xml:"Placemark" json:"Placemarks"`
}

// RouteDocument -
type RouteDocument struct {
	Folder RouteFolder `xml:"Folder" json:"Folder"`
}

// RouteKml -
type RouteKml struct {
	XMLName  xml.Name      `xml:"kml" json:"-"`
	Document RouteDocument `xml:"Document" json:"Document"`
}

func kmlFileParseRoutes(filespec string, kml *RouteKml) (err error) {
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

package main

import (
	"fmt"
)

var (
	a = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b = "0123456789"
)

type EastingNorthing struct {
	easting  float64
	northing float64
}

var (
	aWaypoint = EastingNorthing{1.0, 2.0}
	bWaypoint = EastingNorthing{3.0, 4.0}
)

var (
	routeWaypoints [][]EastingNorthing
	aPlaces        = []EastingNorthing{aWaypoint, bWaypoint, {9.0, 8.0}}
	bPlaces        = []EastingNorthing{bWaypoint, aWaypoint}
	cPlaces        = []EastingNorthing{bWaypoint, aWaypoint, {7.0, 6.0}}
)

func main() {
	fmt.Println("Hello, playground")
	a += b
	fmt.Println(a + b + a)
	fmt.Println(b)
	fmt.Println(aWaypoint)
	fmt.Println(bWaypoint)
	aPlaces = append(aPlaces, bWaypoint)
	fmt.Println(aPlaces)
	routeWaypoints = append(routeWaypoints, aPlaces)
	routeWaypoints = append(routeWaypoints, bPlaces)
	routeWaypoints = append(routeWaypoints, cPlaces)
	fmt.Println(routeWaypoints)
}

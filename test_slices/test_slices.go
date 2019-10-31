package main

import (
	"fmt"
	//"os"
	//"strings"
)

func main() {
	cities := []string{
		"victoria",
		"nanaimo",
		"comox",
		"kamloops",
		"kelowna",
		"squamish"}

	filenames := []string{
		"google_transit.zip",
		"trip_reference.txt",
		"gtfrealtime_TripUpdates.bin",
		"gtfrealtime_ServiceAlerts.bin",
		"gtfrealtime_VehiclePositions.bin"}

	for i := 0; i < len(cities); i++ {
		city := cities[i]
		folder := ".mapstrat.com/current/"
		for file := 0; file < len(filenames); file++ {
			url := "https://" + city + folder + filenames[file]
			fmt.Println("City =", city, "Url =", url)
		}
	}
}

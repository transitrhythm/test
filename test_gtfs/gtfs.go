package main

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/im7mortal/UTM"
	"github.com/patrickbr/gtfsparser"      // "github.com/geops/gtfsparser" //
	"github.com/patrickbr/gtfsparser/gtfs" // "github.com/geops/gtfsparser/gtfs" //
	"github.com/paulmach/orb"
)

// Cartesian -
type Cartesian struct {
	x, y float64
}

// GtfsStop -
type GtfsStop struct {
	ID            string
	Code          string
	Name          string
	Desc          string
	Parent        string
	Zone          string
	Location      Cartesian
	LocationLabel string
}

// GtfsShapePoint -
type GtfsShapePoint struct {
	Location     Cartesian
	Sequence     int
	DistTraveled float64
}

// GtfsShape -
type GtfsShape struct {
	ID     string
	Points []GtfsShapePoint
}

// GtfsAgency -
type GtfsAgency struct {
	ID       string
	Name     string
	URL      string
	Timezone string
	Lang     string
	Phone    string
	FareURL  string
}

// GtfsRoute -
type GtfsRoute struct {
	ID        string
	AgencyID  string
	ShortName string
	LongName  string
	Desc      string
	Type      string
	URL       string
}

// ServiceException -
type ServiceException struct {
	Date Date
	Type int8
}

// GtfsService -
type GtfsService struct {
	ID         string
	Daymap     [7]bool
	StartDate  Date
	EndDate    Date
	Exceptions []*ServiceException
}

// Date -
type Date struct {
	Day   int8
	Month int8
	Year  int16
}

// GtfsStopTime -
type GtfsStopTime struct {
	Sequence          int
	StopID            string
	StopName          string
	Headsign          string
	ArrivalTime       gtfs.Time
	DepartureTime     gtfs.Time
	PickupType        int
	DropOffType       int
	ShapeDistTraveled float64
	Timepoint         bool
}

// GtfsFrequency -
type GtfsFrequency struct {
	StartTime   string
	EndTime     string
	HeadwaySecs int
	ExactTimes  bool
}

// GtfsTrip -
type GtfsTrip struct {
	ID                   string
	RouteID              string
	Service              *gtfs.Service
	Headsign             string
	ShortName            string
	DirectionID          int
	BlockID              string
	ShapeID              string
	WheelchairAccessible int
	BikesAllowed         int
	StopTimes            gtfs.StopTimes
	StartTime            gtfs.Time
	EndTime              gtfs.Time
	Length               float64
	Frequencies          []*GtfsFrequency
}

// GtfsFeedInfo -
type GtfsFeedInfo struct {
	PublisherName string
	PublisherURL  string
	Lang          string
	StartDate     gtfs.Date
	EndDate       gtfs.Date
	Phone         string
	Version       string
	Active        bool
}

// GtfsVehicle -
type GtfsVehicle struct {
	DistTraveled float64
}

func toGtfsDate(year int, month time.Month, day int) (date gtfs.Date) {
	date.Year = int16(year)
	date.Month = int8(month)
	date.Day = int8(day)
	return date
}

func toDayOfYear()

// ParseFeeds -
func ParseFeeds(feed *gtfsparser.Feed) (feedInfo []GtfsFeedInfo) {
	fmt.Println("ParseFeeds:")
	today := toGtfsDate(time.Now().Date())
	for k, v := range feed.FeedInfos {
		feedData := GtfsFeedInfo{}
		feedData.PublisherName = v.Publisher_name
		feedData.Version = v.Version
		feedData.StartDate = v.Start_date
		feedData.EndDate = v.End_date
		fmt.Printf("[%d] %s : <%s> - %s - %s\n", k, v.Version, v.Publisher_name, Datestamp(feedData.StartDate), Datestamp(feedData.EndDate))
		feedData.Active = inDateRange(feedData.StartDate, today, feedData.EndDate)
		feedInfo = append(feedInfo, feedData)
	}
	return feedInfo
}

// ParseAgencies -
func ParseAgencies(feed *gtfsparser.Feed) (agencies []GtfsAgency) {
	fmt.Println("ParseAgencies:")
	for _, v := range feed.Agencies {
		// fmt.Printf("[%s] %s : <%s> - %s\n", k, v.Id, v.Name, v.Timezone)
		agency := GtfsAgency{}
		agency.ID = v.Id
		agency.Name = v.Name
		agency.Timezone = v.Timezone.GetTzString()
		agencies = append(agencies, agency)
	}
	return agencies
}

// FindLabel -
func FindLabel(list []GtfsRouteType, types int) string {
	for _, n := range list {
		if types == n.Type {
			return n.Label
		}
	}
	return ""
}

// AverageSpeed -
func AverageSpeed(distance float64, startTime, endTime gtfs.Time) (speed float64, minutes float64) {
	seconds := toSeconds(endTime) - toSeconds(startTime)
	return distance / ((float64(seconds)) / 3600), float64(seconds) / 60
}

// By -
type By func(r1, r2 *GtfsRoute) bool

// Routes -
type Routes []*GtfsRoute

// Len -
func (s Routes) Len() int { return len(s) }

// Swap -
func (s Routes) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// ByID -
type ByID struct{ Routes }

// Less -
func (s ByID) Less(i, j int) bool { return s.Routes[i].ID < s.Routes[j].ID }

// ParseRoutes -
func ParseRoutes(feed *gtfsparser.Feed) (routes []*GtfsRoute, err error) {
	fmt.Println("ParseRoutes:")
	for _, v := range feed.Routes {
		//fmt.Printf("[%s] %s - %s : <%s>\n", v.Id, label, v.Short_name, v.Desc)
		route := GtfsRoute{}
		route.ID = v.Id
		route.AgencyID = v.Agency.Id
		route.Type = FindLabel(gtfsRouteTypes, int(v.Type))
		route.Desc = v.Desc
		route.ShortName = v.Short_name
		route.LongName = v.Long_name
		routes = append(routes, &route)
	}
	sort.Sort(ByID{routes})
	return routes, err
}

// toSeconds -
func toSeconds(input gtfs.Time) int {
	return (3600 * int(input.Hour)) + (60 * int(input.Minute)) + int(input.Second)
}

// TripActive -
func TripActive(time, start, end gtfs.Time) bool {
	return toSeconds(time) >= toSeconds(start) && toSeconds(time) <= toSeconds(end)
}

// TripComplete -
func TripComplete(distanceTravelled float64, tripDistance float64) bool {
	return true
}

// Now -
func Now() (currentTime gtfs.Time) {
	now := time.Now()
	currentTime.Hour = int8(now.Hour())
	currentTime.Minute = int8(now.Minute())
	currentTime.Second = int8(now.Second())
	return currentTime
}

// VehicleAssignment -
type VehicleAssignment struct {
	Date   gtfs.Date
	TripID string
}

// FeedVehicleRoster -
type FeedVehicleRoster struct {
	agency []AgencyVehicleRoster
}

// AgencyVehicleRoster -
type AgencyVehicleRoster struct {
	AgencyID  string
	VehicleID string
	Schedule  []VehicleAssignment
}

// FeedRoster -
var FeedRoster FeedVehicleRoster

// SetVehicleAssignment -
func SetVehicleAssignment(agencyID string, vehicleID string, date gtfs.Date, tripID string) {
	key := -1
	for k, v := range FeedRoster.agency {
		if agencyID == v.AgencyID && vehicleID == v.VehicleID {
			key = k
		}
	}
	// If the agency/vehicle key pair are not currently assigned, then assign a new entry
	if key == -1 {
		roster := AgencyVehicleRoster{}
		roster.AgencyID = agencyID
		roster.VehicleID = vehicleID
		assignment := VehicleAssignment{date, tripID}
		roster.Schedule = append(roster.Schedule, assignment)
		FeedRoster.agency = append(FeedRoster.agency, roster)
	} else {
		roster := FeedRoster.agency[key]
		assignment := VehicleAssignment{date, tripID}
		roster.Schedule = append(roster.Schedule, assignment)
	}
}

func getDistanceTravelled(vehicleID string, vehicleLocation Cartesian) (distance float64) {
	return distance
}

func getVehicleLocation(vehicleID string) (vehicleLocation Cartesian) {
	return vehicleLocation
}

func getTripLength(tripID string) (length float64) {
	return length
}

func getVehicleID(tripID string) (vehicleID string) {
	return vehicleID
}

// ActiveTrips -
func ActiveTrips(trips []*GtfsTrip) (activeTrips []*GtfsTrip, err error) {
	fmt.Println("ActiveTrips:")
	for i := range trips {
		stopTimes := trips[i].StopTimes
		tripID := trips[i].ID
		if GtfsServiceInOperationToday(*trips[i].Service) &&
			TripActive(Now(), stopTimes[0].Arrival_time, stopTimes[len(stopTimes)-1].Arrival_time) &&
			TripComplete(getDistanceTravelled(getVehicleID(tripID), getVehicleLocation(getVehicleID(tripID))), getTripLength(tripID)) == false {
			activeTrips = append(activeTrips, trips[i])
		}
	}
	return activeTrips, err
}

// Trips -
type Trips []*GtfsTrip

// Len -
func (s Trips) Len() int { return len(s) }

// Swap -
func (s Trips) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// ByBlockID -
type ByBlockID struct{ Trips }

// Less -
func (s ByBlockID) Less(i, j int) bool { return s.Trips[i].BlockID < s.Trips[j].BlockID }

// ParseTrips -
func ParseTrips(feed *gtfsparser.Feed, shapes []*GtfsShape) (trips []*GtfsTrip, err error) {
	fmt.Println("ParseTrips:")
	for _, v := range feed.Trips {
		trip := GtfsTrip{}
		trip.ID = v.Id
		if v.Route != nil {
			trip.RouteID = v.Route.Id
		}
		trip.BlockID = v.Block_id
		trip.Headsign = v.Headsign
		trip.DirectionID = int(v.Direction_id)
		if v.Shape != nil {
			trip.ShapeID = v.Shape.Id
			for i := 0; i < len(shapes); i++ {
				if trip.ShapeID == shapes[i].ID {
					points := shapes[i].Points
					trip.Length = points[len(points)-1].DistTraveled
				}
			}
		}
		trip.StopTimes = v.StopTimes
		trip.StartTime = v.StopTimes[0].Departure_time
		trip.EndTime = v.StopTimes[len(v.StopTimes)-1].Departure_time
		trip.Service = v.Service
		trips = append(trips, &trip)
	}
	sort.Sort(ByBlockID{trips})
	return trips, err
}

// Stops -
type Stops []*GtfsStop

// LocationType -
type LocationType struct {
	LocationType        int
	LocationLabel       string
	LocationDescription string
}

var locationTypes = []LocationType{
	{0, "Stop (or Platform)", "A location where passengers board or disembark from a transit vehicle. Is called a platform when defined within a parent_station"},
	{1, "Station", "A physical structure or area that contains one or more platforms"},
	{2, "Entrance/Exit", "A location where passengers can enter or exit a station from the street. If an entrance/exit belongs to multiple stations, it can be linked by pathways to both, but the data provider must pick one of them as parent"},
	{3, "Generic Node", "A location within a station, not matching any other location_type, which can be used to link together pathways define in pathways.txt"},
	{4, "Boarding Area", "A specific location on a platform, where passengers can board and/or alight vehicles"},
}

// Len -
func (s Stops) Len() int { return len(s) }

// Swap -
func (s Stops) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// ByCode -
type ByCode struct{ Stops }

func stringtoInt(input string) (result int) {
	result, err := strconv.Atoi(input)
	if err != nil {
		result = 0
	}
	return result
}

// Less -
func (s ByCode) Less(i, j int) bool { return stringtoInt(s.Stops[i].ID) < stringtoInt(s.Stops[j].ID) }

func stopLocationLabel(index int) (label string) {
	return locationTypes[index].LocationLabel
}

// ParseStops -
func ParseStops(feed *gtfsparser.Feed, precision int) (stops []*GtfsStop, err error) {
	fmt.Println("ParseStops:")
	for _, v := range feed.Stops {
		var location Cartesian
		stop := GtfsStop{}
		stop.ID = v.Id
		stop.Code = v.Code
		stop.Zone = v.Zone_id
		if v.Parent_station != nil {
			stop.Parent = v.Parent_station.Id + "-" + v.Parent_station.Name
		}
		stop.Name = v.Name
		stop.Desc = v.Desc
		location, _, _, err = toLocation(v.Lat, v.Lon, precision)
		stop.Location = location
		stop.LocationLabel = stopLocationLabel(int(v.Location_type))
		stops = append(stops, &stop)
	}
	sort.Sort(ByCode{stops})
	return stops, err
}

func createSlices(a Cartesian) (routes [][]Cartesian) {
	route := []Cartesian{}
	route = append(route, a)
	routes = append(routes, route)
	return routes
}

func createStopTimes(a GtfsStopTime) (stopTimes [][]GtfsStopTime, err error) {
	stopTimesByTrip := []GtfsStopTime{}
	stopTimesByTrip = append(stopTimesByTrip, a)
	stopTimes = append(stopTimes, stopTimesByTrip)
	return stopTimes, err
}

// BySequence -
type BySequence struct{ gtfs.ShapePoints }

// GtfsShapePoints -
type GtfsShapePoints []GtfsShapePoint

// ByGtfsSequence -
type ByGtfsSequence struct{ GtfsShapePoints }

// Len -
func (shapePoints GtfsShapePoints) Len() int {
	return len(shapePoints)
}

// Less -
func (s ByGtfsSequence) Less(i, j int) bool {
	return s.GtfsShapePoints[i].Sequence < s.GtfsShapePoints[j].Sequence
}

// Swap -
func (shapePoints GtfsShapePoints) Swap(i, j int) {
	shapePoints[i], shapePoints[j] = shapePoints[j], shapePoints[i]
}

// ParseShapes -
func ParseShapes(feed *gtfsparser.Feed, precision int) (shapes []*GtfsShape, err error) {
	fmt.Println("ParseShapes:")
	for k, v := range feed.Shapes {
		fmt.Printf("[%s]%s\n", v.Id, k)
		sort.Sort(BySequence{v.Points})
		shape := GtfsShape{}
		shape.ID = v.Id
		var lastLocation, location Cartesian
		var lastCoord Coord
		var distTraveled, diffDistance, lastDistTraveled float64
		for i := range v.Points {
			shapePoint := v.Points[i]
			gtfsShapePoint := GtfsShapePoint{}
			gtfsShapePoint.Sequence = shapePoint.Sequence
			location, _, _, _ = toLocation(shapePoint.Lat, shapePoint.Lon, precision)
			if i == 0 {
				lastLocation = location
				lastCoord = toCoord(shapePoint.Lat, shapePoint.Lon, precision*2)
				distTraveled = 0.0
				lastDistTraveled = 0.0
			}
			if shapePoint.HasDistanceTraveled() {
				gtfsShapePoint.DistTraveled = toFixed(float64(shapePoint.Dist_traveled), 4)
				dist := hypotenuse(location.x-lastLocation.x, location.y-lastLocation.y, precision)
				fmt.Println(toFixed(dist, precision), "=", toFixed(location.x-lastLocation.x, precision), ";", toFixed(location.y-lastLocation.y, precision))
				coord := toCoord(shapePoint.Lat, shapePoint.Lon, precision*2)
				dist2 := SphericalDistance(coord, lastCoord, precision)
				lastCoord = coord
				diffDistance = (gtfsShapePoint.DistTraveled - lastDistTraveled) * 1000
				fmt.Println("Length:", toFixed(diffDistance, precision), "[", dist, PercentDiff(diffDistance, dist, precision), "% ", "]", "[", dist2, PercentDiff(diffDistance, dist2, precision), "% ]")
				lastDistTraveled = gtfsShapePoint.DistTraveled
			} else {
				diffDistance = hypotenuse(location.x-lastLocation.x, location.y-lastLocation.y, precision)
				distTraveled += diffDistance
				gtfsShapePoint.DistTraveled = toKm(distTraveled)
			}
			gtfsShapePoint.Location = location
			fmt.Printf("[%d] %.4f - %.4f : (@ %.3f,%.3f)\n", i, diffDistance, gtfsShapePoint.DistTraveled*1000, toFixed(location.x, precision), toFixed(location.y, precision))
			lastLocation = location
			shape.Points = append(shape.Points, gtfsShapePoint)
		}
		sort.Sort(ByGtfsSequence{shape.Points})
		shapes = append(shapes, &shape)
	}
	return shapes, err
}

// LineSegment -
type LineSegment struct {
	p1, p2 Cartesian
}

func toLocation(lat float32, lon float32, precision int) (location Cartesian, timeZone int, timeLetter string, err error) {
	x, y, timeZone, timeLetter, err := UTM.FromLatLon(float64(lat), float64(lon), lat > 0)
	location.x = toFixed(x, precision)
	location.y = toFixed(y, precision)
	return location, timeZone, timeLetter, err
}

func toPoint(cartesian Cartesian) (point orb.Point) {
	point[0] = cartesian.x
	point[1] = cartesian.y
	return point
}

func toLineSegment2D(a Cartesian, b Cartesian) (line Line2D) {
	line.p1 = toPoint(a)
	line.p2 = toPoint(b)
	return line
}

// ShapeLine -
type ShapeLine struct {
	sequence int
	line     Line2D
}

// BoundingBox -
type BoundingBox struct {
	location [4]orb.Point
}

// InRange -
func InRange(a, b, c float64) bool {
	return (a-b <= 0.0) && (b-c <= 0.0)
}

// Determines if a point is within the bounds of a box, created from a centreline and width
func foundBoundingBox(stopPoint orb.Point, segment ShapeLine, width float64, precision int) bool {
	centreline := toPolarReference(segment.line.p1, segment.line.p2, precision)
	newPoint := rotate2DPoint(segment.line.p2, segment.line.p1, -centreline.polarPoint.bearing, precision)
	newStop := rotate2DPoint(stopPoint, segment.line.p1, -centreline.polarPoint.bearing, precision)
	var box BoundingBox
	box.location[0] = orb.Point{segment.line.p1.X() - width/2, segment.line.p1.Y()}
	box.location[1] = orb.Point{segment.line.p1.X() + width/2, segment.line.p1.Y()}
	box.location[2] = orb.Point{newPoint.X() + width/2, newPoint.Y()}
	box.location[3] = orb.Point{newPoint.X() - width/2, newPoint.Y()}
	return (InRange(box.location[0].X(), newStop.X(), box.location[1].X()) && InRange(box.location[0].Y(), newStop.Y(), box.location[2].Y()))
}

// MaxRoadWidth -
const MaxRoadWidth = 20.0

func calculateStopTripDistance(shapes []*GtfsShape, shapeID string, stopPoint Cartesian, precision int) (distance float64, err error) {
	for _, v := range shapes {
		if shapeID == v.ID {
			// minimumDistance := 100.
			segmentSelected := 0
			polyLine := []ShapeLine{}
			shapeLine := ShapeLine{}
			// Find shape segment associated with a particular Stop
			for k1, v1 := range v.Points {
				if k1+1 == len(v.Points) {
					break
				}
				shapeLine.sequence = v1.Sequence
				shapeLine.line = toLineSegment2D(v.Points[k1].Location, v.Points[k1+1].Location)
				if foundBoundingBox(toPoint(stopPoint), shapeLine, MaxRoadWidth, precision) {
					polyLine = append(polyLine, shapeLine)
					segmentSelected = shapeLine.sequence - 1
					break
				}
				// polarReference := offsetTo2DLine(toPoint(stopPoint), shapeLine.line, precision)
				// if minimumDistance > polarReference.polarPoint.distance {
				// 	minimumDistance = polarReference.polarPoint.distance
				// 	segmentSelected = shapeLine.sequence
				// }
				// polyLine = append(polyLine, shapeLine)
			}
			fmt.Print(segmentSelected)
			distance += LengthLine2D(toLine2D(toPoint(stopPoint), polyLine[segmentSelected].line.p1), precision)
		}
	}
	return distance, err
}

func stopLocation(stops []*GtfsStop, stopID string) (location Cartesian) {
	for _, v := range stops {
		if v.ID == stopID {
			location = v.Location
			break
		}
	}
	return location
}

// GtfsStopTimes -
type GtfsStopTimes struct {
	tripID string
	times  []GtfsStopTime
}

// GtfsServiceExceptionStatus -
func GtfsServiceExceptionStatus(exceptionType int) (status string) {
	if exceptionType == 1 {
		status = "added"
	} else if exceptionType == 2 {
		status = "dropped"
	}
	return status
}

// GtfsServiceInOperationStatus -
func GtfsServiceInOperationStatus(service gtfs.Service) (status string) {
	dayString := []byte{'S', 'M', 'T', 'W', 'T', 'F', 'S'}
	for i := 0; i < len(service.Daymap); i++ {
		if service.Daymap[i] {
			status += string(dayString[i])
		} else {
			status += "-"
		}
	}
	status += " - " + Datestamp(service.Start_date) + " - " + Datestamp(service.End_date)
	if len(service.Exceptions) > 0 {
		for k, v := range service.Exceptions {
			status += " " + Datestamp(k) + ":" + GtfsServiceExceptionStatus(int(v)) + ";"
		}
	}
	return status
}

// GtfsServiceInOperationToday -
func GtfsServiceInOperationToday(service gtfs.Service) (status bool) {
	now := time.Now()
	today := gtfs.GetGtfsDateFromTime(now)
	status = service.IsActiveOn(today)
	fmt.Println(today)
	return status
}

// GtfsTripDirectionStatus -
func GtfsTripDirectionStatus(direction int) (status string) {
	status = "inbound"
	if direction == 1 {
		status = "outbound"
	}
	return status
}

// ParseStopTimes -
func ParseStopTimes(feed *gtfsparser.Feed, stops []*GtfsStop, shapes []*GtfsShape, precision int) (stopTimes []*GtfsStopTimes, err error) {
	fmt.Println("ParseStopTimes:")
	for _, v := range feed.Trips {
		stopTimesByTrip := GtfsStopTimes{}
		stopTimesByTrip.tripID = v.Id
		stopTimesByTrip.times = []GtfsStopTime{}
		if v.StopTimes != nil {
			for _, v1 := range v.StopTimes {
				stopTime := GtfsStopTime{}
				stopTime.Sequence = v1.Sequence
				stopTime.Headsign = v1.Headsign
				stopTime.ArrivalTime = v1.Arrival_time
				stopTime.DepartureTime = v1.Departure_time
				stopTime.Timepoint = v1.Timepoint
				stopTime.StopID = v1.Stop.Id
				stopTime.StopName = v1.Stop.Name
				location := stopLocation(stops, v1.Stop.Id)
				if err == nil {
					if v.Id == "10846175" && stopTime.StopID == "448" {
						fmt.Printf("...")
					}
					stopTime.ShapeDistTraveled, err = calculateStopTripDistance(shapes, v.Shape.Id, location, precision)
					fmt.Printf("StopTimes:[%s] %s %d : %s - %s\n", v.Id, stopTime.StopID, stopTime.Sequence, Timestamp(stopTime.ArrivalTime), Timestamp(stopTime.DepartureTime))
					stopTimesByTrip.times = append(stopTimesByTrip.times, stopTime)
				}
			}
		}
		stopTimes = append(stopTimes, &stopTimesByTrip)
	}
	return stopTimes, err
}

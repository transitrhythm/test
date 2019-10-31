package main

import (
	"fmt"
	"math"

	"github.com/golang/geo/r3"
	"github.com/golang/geo/s2"
)

// . "github.com/paulmach/orb"
// . "github.com/golang/geo/r3"
// . "github.com/golang/geo/s2"

// Segment -
type Segment struct {
	begin s2.Point
	end   s2.Point
}

// SegmentOffset -
type SegmentOffset struct {
	point    s2.Point
	distance float64
}

func minimumDistanceToS2Segment(p, p1, p2 s2.Point) (distance float64) {
	angle := s2.DistanceFromSegment(p, p1, p2)
	distance = 0
	fmt.Println(angle)
	return distance
}

// rotatePoint - Rotate 2D point through angle around origin.
func rotateS2Point(point s2.Point, origin s2.Point, degrees Degrees, precision int) (result s2.Point) {
	length := hypotenuse((origin.X - point.X), (origin.Y - point.Y), precision)
	degrees -= toDegrees((origin.X - point.X), (origin.Y - point.Y), precision)
	result.X = toFixed(origin.X+length*math.Cos(toRadians(degrees)), precision)
	result.Y = toFixed(origin.Y+length*math.Sin(toRadians(degrees)), precision)
	result.Z = origin.Z
	return result
}

// offset_to_segment - Rotate course start point through course angle to horizontal
func offsetToS2Segment(point s2.Point, segment Segment, precision int) (offset SegmentOffset) {
	angle := toDegrees(segment.begin.X-segment.end.X, segment.begin.Y-segment.end.Y, precision)
	newPoint := rotateS2Point(segment.begin, point, angle, precision)
	newPoint.X = point.X
	offset.distance = toFixed(point.Y-newPoint.Y, precision)
	offset.point = rotateS2Point(newPoint, point, -angle, precision)
	return offset
}

func testRotateS2Points() {
	startPoint := s2.Point{Vector: r3.Vector{X: 0, Y: 1, Z: 0}}
	originPoint := s2.Point{Vector: r3.Vector{X: 0, Y: 2, Z: 0}}
	// Test #1
	firstPoint := s2.Point{Vector: r3.Vector{X: 0, Y: 3, Z: 0}}
	resultPoint := rotateS2Point(startPoint, originPoint, math.Pi, 3)
	fmt.Println("Test #1", (resultPoint == firstPoint))
	// Test #2
	secondPoint := s2.Point{Vector: r3.Vector{X: 1, Y: 2, Z: 0}}
	resultPoint = rotateS2Point(startPoint, originPoint, math.Pi/2, 3)
	fmt.Println("Test #2", (resultPoint == secondPoint))
	// Test #3
	thirdPoint := s2.Point{Vector: r3.Vector{X: 1, Y: 2, Z: 0}}
	resultPoint = rotateS2Point(startPoint, originPoint, -math.Pi/2, 3)
	fmt.Println("Test #3", (resultPoint == thirdPoint))
	// Test #4
	resultPoint = rotateS2Point(startPoint, originPoint, math.Pi*2, 3)
	fmt.Println("Test #4", (resultPoint == startPoint))
}

func testOffsetToS2Segments() {

}

func testSuiteS2() {
	testRotateS2Points()
	testOffsetToS2Segments()
}

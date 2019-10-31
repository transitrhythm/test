package main

import (
	"fmt"
	"math"

	"github.com/paulmach/orb"
)

// . "github.com/paulmach/orb"
// . "github.com/golang/geo/r3"
// . "github.com/golang/geo/s2"

// Line2D -
type Line2D struct {
	p1, p2 orb.Point
}

var (
	f = fmt.Println
)

// PolarReference -
type PolarReference struct {
	point      orb.Point
	polarPoint PolarPoint
}

func toPolarReference(p1, p2 orb.Point, precision int) (polarReference PolarReference) {
	polarReference.polarPoint.distance = hypotenuse(p2[0]-p1[0], p2[1]-p1[1], precision)
	polarReference.polarPoint.bearing = toDegrees(p2[0]-p1[0], p2[1]-p1[1], precision)
	if polarReference.polarPoint.bearing < 0 {
		polarReference.polarPoint.bearing += 360
	}
	polarReference.point = p1
	return polarReference
}

// LengthLine2D -
func LengthLine2D(line Line2D, precision int) float64 {
	return hypotenuse(line.p1[0]-line.p2[0], line.p1[1]-line.p2[1], precision)
}

func toLine2D(a orb.Point, b orb.Point) (line Line2D) {
	line.p1 = a
	line.p2 = b
	return line
}

func toPolarPoint(point orb.Point, precision int) (polarPoint PolarPoint) {
	polarPoint.distance = hypotenuse(point[0], point[1], precision)
	polarPoint.bearing = toDegrees(point[0], point[1], precision)
	return polarPoint
}

// Segment2D -
type Segment2D struct {
	begin, end orb.Point
}

// Segment2DOffset -
type Segment2DOffset struct {
	point    orb.Point
	distance float64
}

// func minimumDistanceTo2DSegment(p, p1, p2 orb.Point) (distance float64) {
// 	angle := orb.DistanceFromSegment(p, p1, p2)
// 	distance = 0
// 	f.Println(angle)
// 	return distance
// }

// rotatePoint - Rotate 2D point through angle around origin.
func rotate2DPoint(point orb.Point, origin orb.Point, degrees Degrees, precision int) (result orb.Point) {
	length := LengthLine2D(toLine2D(origin, point), precision)
	angle := toDegrees(point[0]-origin[0], point[1]-origin[1], precision)
	degrees = Degrees(toFixed(math.Mod(float64(degrees+angle), 360), precision))
	result[0] = toFixed(origin[0]+length*math.Sin(toRadians(degrees)), precision)
	result[1] = toFixed(origin[1]+length*math.Cos(toRadians(degrees)), precision)
	return result
}

// offsetTo2DLine - Rotate course start point through course angle to horizontal
func offsetTo2DLine(point orb.Point, line Line2D, precision int) (polarReference PolarReference) {
	// Calculate line angle to positive Y axis
	degrees := toDegrees(line.p2[0]-line.p1[0], line.p2[1]-line.p1[1], precision)
	// Rotate line start through line angle from specified point
	newPoint := rotate2DPoint(line.p1, point, -degrees, precision)
	newPoint[1] = point[1]
	// Rotate new point by the calculated line angle to obtain the point on the line orthogonal to the original line
	polarReference.point = rotate2DPoint(newPoint, point, degrees, precision)
	polarReference.polarPoint.distance = hypotenuse(point[0]-polarReference.point[0], point[1]-polarReference.point[1], precision)
	polarReference.polarPoint.bearing = toDegrees(point[0]-polarReference.point[0], point[1]-polarReference.point[1], precision)
	return polarReference
}

func testRotate2DPoints(precision int) {
	f("testRotate2DPoints:")
	startPoint := orb.Point{0, 1}
	originPoint := orb.Point{1, 1}
	// Test #1
	firstPoint := orb.Point{2, 1}
	resultPoint := rotate2DPoint(startPoint, originPoint, 180, precision)
	f("Test #1", (resultPoint == firstPoint))
	// Test #2
	secondPoint := orb.Point{1, 2}
	resultPoint = rotate2DPoint(startPoint, originPoint, 90, precision)
	f("Test #2", (resultPoint == secondPoint))
	// Test #3
	thirdPoint := orb.Point{1, 0}
	resultPoint = rotate2DPoint(startPoint, originPoint, -90, precision)
	f("Test #3", (resultPoint == thirdPoint))
	// Test #4
	resultPoint = rotate2DPoint(startPoint, originPoint, 360, precision)
	f("Test #4", (resultPoint == startPoint))
}

func testOffsetTo2DLines(precision int) {
	f("testOffsetTo2DLines:")
	// Test #1
	firstPoint := orb.Point{1, 1}
	firstReference := PolarReference{orb.Point{0, 1}, PolarPoint{1, 90}}
	resultReference := offsetTo2DLine(firstPoint, toLine2D(orb.Point{0, -2}, orb.Point{0, 2}), precision)
	f("Test #1", (resultReference == firstReference))
	// Test #2
	secondPoint := orb.Point{-1, -3}
	secondReference := PolarReference{orb.Point{0, -3}, PolarPoint{1, 270}}
	resultReference = offsetTo2DLine(secondPoint, toLine2D(orb.Point{0, -5}, orb.Point{0, 5}), precision)
	f("Test #2", (resultReference == secondReference))
	// Test #3
	thirdPoint := orb.Point{1, 1}
	thirdReference := PolarReference{orb.Point{0, 0}, PolarPoint{1.414, 45}}
	resultReference = offsetTo2DLine(thirdPoint, toLine2D(orb.Point{-1, 1}, orb.Point{1, -1}), precision)
	f("Test #3", (resultReference == thirdReference))
	// Test #4
	fourthPoint := orb.Point{-1, 1}
	fourthReference := PolarReference{orb.Point{0, 0}, PolarPoint{1.414, 315}}
	resultReference = offsetTo2DLine(fourthPoint, toLine2D(orb.Point{-2, -2}, orb.Point{2, 2}), precision)
	f("Test #4", (resultReference == fourthReference))
	// Test #5
	fifthPoint := orb.Point{0, 0}
	fifthReference := PolarReference{orb.Point{0, 0}, PolarPoint{0, 0}}
	resultReference = offsetTo2DLine(fifthPoint, toLine2D(orb.Point{-2, -2}, orb.Point{2, 2}), precision)
	f("Test #5", (resultReference == fifthReference))
}

func testToPolarPoint(precision int) {
	f("testToPolar:")
	// Test #1
	firstPoint := orb.Point{0, 3}
	firstPolar := PolarPoint{3, 0}
	resultPolar := toPolarPoint(firstPoint, precision)
	f("Test #1", (resultPolar == firstPolar))
	// Test #2
	secondPoint := orb.Point{2, 0}
	secondPolar := PolarPoint{2, 90}
	resultPolar = toPolarPoint(secondPoint, precision)
	f("Test #2", (resultPolar == secondPolar))
	// Test #3
	thirdPoint := orb.Point{0, -3}
	thirdPolar := PolarPoint{3, 180}
	resultPolar = toPolarPoint(thirdPoint, precision)
	f("Test #3", (resultPolar == thirdPolar))
	// Test #4
	fourthPoint := orb.Point{-4, 0}
	fourthPolar := PolarPoint{4, 270}
	resultPolar = toPolarPoint(fourthPoint, precision)
	f("Test #4", (resultPolar == fourthPolar))
}

func testToPolarReference(precision int) {
	f("testToPolarReference:")
	originPoint := orb.Point{0, 2}
	// Test #1
	firstPoint := orb.Point{0, 3}
	firstPolarReference := PolarReference{orb.Point{0, 2}, PolarPoint{1, 0}}
	resultPolarReference := toPolarReference(originPoint, firstPoint, precision)
	f("Test #1", (resultPolarReference == firstPolarReference))
	// Test #2
	secondPoint := orb.Point{0, 1}
	secondPolarReference := PolarReference{orb.Point{0, 2}, PolarPoint{1, 180}}
	resultPolarReference = toPolarReference(originPoint, secondPoint, precision)
	f("Test #2", (resultPolarReference == secondPolarReference))
	// Test #3
	thirdPoint := orb.Point{2, 2}
	thirdPolarPolarReference := PolarReference{orb.Point{0, 2}, PolarPoint{2, 90}}
	resultPolarReference = toPolarReference(originPoint, thirdPoint, precision)
	f("Test #3", (resultPolarReference == thirdPolarPolarReference))
	// Test #4
	fourthPoint := orb.Point{-1, 2}
	fourthPolarReference := PolarReference{orb.Point{0, 2}, PolarPoint{1, 270}}
	resultPolarReference = toPolarReference(originPoint, fourthPoint, precision)
	f("Test #4", (resultPolarReference == fourthPolarReference))
}

func testBoundingBoxes(precision int) {
	f("testBoundingBoxes:")
	stopPoint := orb.Point{5.0, 5.5}
	firstLine := Line2D{orb.Point{1, 1}, orb.Point{11, 11}}
	segment := ShapeLine{1, firstLine}
	width := 1.0
	f("Test #1", foundBoundingBox(stopPoint, segment, width, precision))
}

func testSuite2D(precision int) {
	testToPolarPoint(precision)
	testToPolarReference(precision)
	testRotate2DPoints(precision)
	testOffsetTo2DLines(precision)
	testBoundingBoxes(precision)
}

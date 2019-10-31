package main

import (
	"fmt"
)

	
type EastingNorthing struct {
	Easting  float64
	Northing float64
}

func main() {
	fmt.Println("Hello, playground")
    	var myMap = map[string]map[string][]EastingNorthing{
        "foo": {
            "bar": {{Easting: 1.0, Northing: 2.0}},
            "baz": {{3.0, 4.0}}, //or by order...
        },
        "bar": {
            "gaz": {{5.0, 6.0}},
            "faz": {{7.0, 8.0}},
        },
	}
	
	//var en = EastingNorthing{6.6,5.5}

	myMap["foo"]["bar"] = append(myMap["foo"]["bar"],EastingNorthing{6.6,5.5})
	myMap["foo"]["bbb"] = append(myMap["foo"]["bbb"],EastingNorthing{6.06,5.05})

    fmt.Println(myMap["foo"]["bar"][0].Easting)
    fmt.Println(myMap["bar"]["gaz"][0].Northing)	
	fmt.Println(myMap["foo"]["bar"][1].Easting)
	fmt.Println(myMap["foo"]["bbb"][0].Easting)
}

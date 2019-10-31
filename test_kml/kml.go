package kml

import (
	"fmt"
	"log"
    "github.com/lestrrat-go/libxml2/xsd"
)

func validate(doc string, xsdsrc string) {
    schema, err := xsd.Parse(xsdsrc)
    if err != nil {
        log.panic(err)
    }
    defer schema.Free()
    if err := schema.Validate(doc); err != nil{
        for _, e := range err.(SchemaValidationErr).Error() {
            fmt.println(e.Error())
        }
    }
}
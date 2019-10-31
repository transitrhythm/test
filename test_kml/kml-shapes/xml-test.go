// Copyright 2019 Bevan Thistlethwaite. All rights reserved.

package main

import (
	"fmt"

	"github.com/subchen/go-xmldom"
)

func xmlTest() {
	xml := `<testsuite tests="2" failures="0" time="0.009" name="github.com/subchen/go-xmldom">
    <testcase classname="go-xmldom" name="ExampleParseXML" time="0.004"></testcase>
    <testcase classname="go-xmldom" name="ExampleParse" time="0.005"></testcase>
</testsuite>`

	doc := xmldom.Must(xmldom.ParseXML(xml))
	root := doc.Root

	name := root.GetAttributeValue("name")
	time := root.GetAttributeValue("time")
	fmt.Printf("testsuite: name=%v, time=%v\n", name, time)

	for _, node := range root.GetChildren("testcase") {
		name := node.GetAttributeValue("name")
		time := node.GetAttributeValue("time")
		fmt.Printf("testcase: name=%v, time=%v\n", name, time)
	}
}

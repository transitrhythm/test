package main

import (
	"fmt"
	"log"
	"net/rpc"
	"time"
)

// Item gRPC message structure.
type Item struct {
	Title string
	Body  string
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s = %dms", name, elapsed.Nanoseconds()/1000)
}

func main() {
	var reply Item
	var db []Item

	defer timeTrack(time.Now(), "gRPC Test duration:")

	client, err := rpc.DialHTTP("tcp", "localhost:4040")
	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	a := Item{"First", "A first item"}
	b := Item{"Second", "A second item"}
	c := Item{"Third", "A third item"}

	client.Call("API.AddItem", a, &reply)
	client.Call("API.AddItem", b, &reply)
	client.Call("API.AddItem", c, &reply)
	client.Call("API.GetDB", "", &db)
	fmt.Println("Database: ", db)

	client.Call("API.EditItem", Item{"Second", "A new second item"}, &reply)
	client.Call("API.DeleteItem", c, &reply)
	client.Call("API.GetDB", "", &db)
	fmt.Println("Database: ", db)

	client.Call("API.GetByName", "First", &reply)
	fmt.Println("First item: ", reply)
}

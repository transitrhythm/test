 package main

 import (
 	"fmt"
 	"net/http"
 	"os"
 )

 func main() {

 	if len(os.Args) != 2 {
 		fmt.Printf("Usage : %s <URL> \n", os.Args[0])
 		os.Exit(0)
 	}

 	url := os.Args[1]

 	fmt.Println("Heading ", url)

 	client := &http.Client{}

 	resp, err := client.Head(url)

 	if err != nil {
 		fmt.Println(err)
 		os.Exit(-1)
 	}

 	fmt.Println(url, " status is : ", resp.Status)

 }
 
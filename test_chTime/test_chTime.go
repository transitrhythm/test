package main

import (
   "fmt"
   "os"
   "time"
)

func main() {
	fmt.Println("Hello world!")
	if len(os.Args) != 2 {
		fmt.Printf("Usage : %s <filespec> \n", os.Args[0])
		os.Exit(0)
	}
    filename := os.Args[1]

     // get last modified time
	 file, err := os.Stat(filename)
     if err != nil {
        fmt.Println(err)
     }

     modifiedtime := file.ModTime()
	 fmt.Println("Last modified time : ", modifiedtime)
	 
	// get current timestamp
	currenttime := time.Now().Local()
	fmt.Println("Current time : ", currenttime.Format("2006-01-02 15:04:05 +0800"))

	// change both atime and mtime to currenttime
	err = os.Chtimes(filename, currenttime, currenttime)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Changed the file time information")
}
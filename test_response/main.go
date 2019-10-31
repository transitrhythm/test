package main

import (
	"net/http"
	"net/url"
)
import "io/ioutil"
import "fmt"
import "strings"

func keepLines(s string, n int) string {
	result := strings.Join(strings.Split(s, "\n")[:n], "\n")
	return strings.Replace(result, "\r", "", -1)
}

const getURL = "http://g.cn/robots.txt"
const postURL = "http://duckduckgo.com"

// const getURL = "https://victoria.mapstrat.com/gtfrealtime_VehiclePositions.bin"

func main() {
	resp, err := http.Get(getURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("get:\n", keepLines(string(body), 3))

	resp, err = http.PostForm(postURL,
		url.Values{"q": {"github"}})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	fmt.Println("post:\n", keepLines(string(body), 3))
}

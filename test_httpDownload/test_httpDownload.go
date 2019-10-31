package main

import (
	"fmt"
	"io"

	//	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/beevik/ntp"
)

type httpSchedule struct {
	//	hourTime string
	url string
	dst string
	//next     scheduler.Job
}

func timeTrack(start time.Time, name string) time.Duration {
	elapsed := time.Since(start)
	log.Println(name, "=", elapsed)
	return elapsed
}

// MaxDuration returns the larger of x or y.
func MaxDuration(x, y time.Duration) time.Duration {
	if x < y {
		return y
	}
	return x
}

// MinDuration returns the smaller of x or y.
func MinDuration(x, y time.Duration) time.Duration {
	if x > y {
		return y
	}
	return x
}

func main() {
	cityNames := []string{
		"victoria",
		"nanaimo",
		"comox",
		"kamloops",
		"kelowna",
		"squamish"}

	fileNames := []string{
		"google_transit.zip",
		"trip_reference.txt",
		"gtfrealtime_TripUpdates.bin",
		"gtfrealtime_ServiceAlerts.bin",
		"gtfrealtime_VehiclePositions.bin"}

	folder := ".mapstrat.com/current/"
	downloadSchedule := []httpSchedule{}

	job := func(cityIndex int, fileIndex int) {
		scheduleIndex := cityIndex*len(fileNames) + fileIndex
		go DownloadFile(downloadSchedule[scheduleIndex].dst, downloadSchedule[scheduleIndex].url)
	}

	ntpTime, _ := ntp.Time("0.beevik-ntp.pool.ntp.org")
	fmt.Println(time.Since(ntpTime))

	// Initialize the download schedule table with Url and Dst filespec, and a blank schedule.
	for cityIndex := 0; cityIndex < len(cityNames); cityIndex++ {
		for fileIndex := 0; fileIndex < len(fileNames); fileIndex++ {
			cityName := cityNames[cityIndex]
			url := "https://" + cityName + folder + fileNames[fileIndex]
			dst := "../data/" + cityName + "/" + fileNames[fileIndex]
			entrySchedule := httpSchedule{url, dst}
			downloadSchedule = append(downloadSchedule, entrySchedule)
			job(cityIndex, fileIndex)
			//gocron.Every(30).Second().Do(job, cityIndex, fileIndex)
		}
	}
	/*
		file, err := os.Stat(dst)
		if err == nil {
			entrySchedule.hourTime = file.ModTime().Format("15:04:05")
		}
	*/
	//	<-gocron.Start()

	//scheduler.Every(30).Seconds().At(entrySchedule.hourTime).NotImmediately().Run(job)

	//fmt.Println(downloadSchedule)

	// Run every 2 seconds but not now.
	//scheduler.Every(10).Seconds().Run(job)
	for {
		time.Sleep(time.Second * 1)
	}
}

// TimeIn returns the time in UTC if the name is "" or "UTC".
// It returns the local time if the name is "Local".
// Otherwise, the name is taken to be a location name in
// the IANA Time Zone database, such as "Africa/Lagos".
func TimeIn(t time.Time, name string) (time.Time, error) {
	loc, err := time.LoadLocation(name)
	if err == nil {
		t = t.In(loc)
	}
	return t, err
}

// SaveFile copies the file in the body of the reponse into a local file repository
func SaveFile(resp *http.Response, filepath string, timestamp string) time.Duration {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		log.Panic("File open error: ", err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Panic("File write error: ", err)
	}
	defer resp.Body.Close()

	ftime, err := time.Parse(time.RFC1123, timestamp)
	latency := time.Since(ftime)

	// change both atime and mtime to lastModifiedTime
	err = os.Chtimes(filepath, ftime, ftime)
	if err != nil {
		log.Panic("File timestamp error: ", err)
	}
	//location, _ := time.LoadLocation("UTC")
	//fmt.Println("Now =", time.Now().In(location), "; Filepath =", filepath, "; Last modified =", timestamp, "; Latency:", latency)
	fmt.Println("Now =", time.Now(), "; Filepath =", filepath, "; Last modified =", ftime.Local(), "; Latency:", latency)
	return latency
}

func loopDuration(datestamp string, interval time.Duration) time.Duration {
	timeValue, _ := time.Parse(time.RFC1123, datestamp)
	elapsed := time.Since(timeValue)
	return MinDuration(interval-elapsed, interval)
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) {
	downloadInterval := time.Duration(time.Second * 30)
	for {
		// Check for pre-existing local data file
		var lastModified string
		waitDuration := downloadInterval
		file, err := os.Stat(filepath)
		if err == nil {
			modTime := file.ModTime()
			location, err := time.LoadLocation("GMT")
			if err != nil {
				log.Panic("Timezone location error: ", err)
			}
			// Convert local file timestamp into GMT string
			modifiedTime := modTime.In(location).Format(time.RFC1123)
			client := &http.Client{}
			request, err := http.NewRequest("GET", url, nil)
			if err != nil {
				log.Panic("File request error: ", err)
			}
			// Add date flag to check for file changes since the last download
			request.Header.Set("If-Modified-Since", modifiedTime)
			response, err := client.Do(request)
			if response.StatusCode == http.StatusOK {
				lastModified = response.Header.Get("Last-Modified")
				SaveFile(response, filepath, lastModified)
				waitDuration = loopDuration(lastModified, downloadInterval)
			} else {
				lastModified = modifiedTime
				waitDuration = downloadInterval
			}
		} else {
			// Get the data
			resp, err := http.Get(url)
			if err != nil {
				log.Panic("HTTP GET error: ", err)
			}
			lastModified = resp.Header.Get("Last-Modified")
			SaveFile(resp, filepath, lastModified)
			waitDuration = loopDuration(lastModified, downloadInterval)
		}
		time.Sleep(waitDuration)
	}
}

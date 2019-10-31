package main
import (
	"fmt"
	"flag"
	//"net/http"
	"os"
	"time"
    "github.com/r3labs/sse"
)

func main() {
	args := os.Args
	stopID := flag.String("StopID","123456", "the unique stop identifying number")
	for i := range args {
		fmt.Println(args[i])
	}
	client := sse.NewClient("http://server/events")

	client.Subscribe(*stopID, func(msg *sse.Event) {
		// Got some data!
		fmt.Println(time.Now(), msg.Data)
	})
}
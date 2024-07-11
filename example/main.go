package main

import (
	"fmt"

	"github.com/seanakdh/watchit"
)

func main() {
	//Read files to depth 1, so all files contained in './' but not further
	files, _ := watchit.GetDirFiles("./", 1)
	//Create a response channel, maybe consider making it bufferd if expecting frequent changes
	changedch := make(chan watchit.WatchItInfo, 100)
	//Start watching the files with 0.1 second interval
	watchit.WatchMultipleFilesAsync(files, changedch, 100)
	//Just print the responses
	for msg := range changedch {
		fmt.Println(msg)
	}
}

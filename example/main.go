package main

import (
	"fmt"

	"github.com/seanakdh/diditchange"
)

func main() {
	//Read files to depth 1, so all files contained in './' but not further
	files, _ := diditchange.GetDirFiles("./", 2)
	//Create a response channel, maybe consider making it bufferd if expecting frequent changes
	changedch := make(chan diditchange.FileChangedInfo)
	//Start watching the files
	diditchange.WatchMultipleFilesAsync(files, changedch, 100)
	//Just print the responses
	for msg := range changedch {
		fmt.Println(msg)
	}
}

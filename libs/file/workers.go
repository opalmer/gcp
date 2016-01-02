package files

// workers.go - contains the functions which perform various units of 'work'
// to copy the file such as reading data from a file,
// compression and encryption

import (
	"os"
)

// Channels - Defines what channels we use to communicate with workers on.
type Channels struct {
	paths chan string
	files chan File
}

var channels Channels

func openpaths(channel <-chan string) {
	for path := range channel {
		channels.files <- File{sourcepath: path}
	}
}

func processfiles(channel <-chan File) {
	for file := range channel {
		// Open the file(s) for the given input
		err := file.open()
		if err != nil {
			if os.IsNotExist(err) {
				log.Warningf("%s does not exist", file.sourcepath)
			}
			file.clean()
			log.Fatalf("Failed to open %s (err: %s)", file.sourcepath, err)
		}

		err = file.process()
		if err != nil {
			file.clean()
			log.Fatalf("Failed to process %s (err: %s)", file.sourcepath, err)
		}

		err = file.save()
		if err != nil {
			file.clean()
			log.Fatalf("Failed to save %s (err: %s)", file.outpath, err)
		}

		filesProcessing.Done()
	}
}

package files

// io.go - Responsible for performing input and output operations.

import (
	"../config"
	"io"
	"io/ioutil"
	"os"
)

const readSize = 1024

// OutputHandler - The main object for controlling output to a file.
type OutputHandler struct {
	source      *os.File
	tmp         *os.File
	destination string
}

// NewOutput - Produce a copy of Output and prepares to write data to disk.
func NewOutput(source string, destination string) OutputHandler {
	sourceFile, err := os.Open(source)
	if err != nil {
		log.Fatalf("Failed to open %s (err: %s)", source, err)
	}

	tempfile, err := ioutil.TempFile(os.TempDir(), "gcp")

	if err != nil {
		defer sourceFile.Close()
		log.Fatalf("Failed to create temp file for %s (err: %s)", source, err)
	}

	out := OutputHandler{
		source: sourceFile, tmp: tempfile, destination: destination}

	return out
}

// Process - Opens the source file, performs operations (compress/encrypt) and
// writes the output to the temporary file.
func (out *OutputHandler) Process() {
	for {
		// Read n bytes from the source file
		data := make([]byte, readSize)
		_, err := out.source.Read(data)
		if err == io.EOF {
			break

		} else if err != nil {
			log.Fatalf(
				"Failed to read from %s (err: %s)", out.source.Name(), err)
		}

		out.tmp.Write(data)

		if config.Compress {

		}
		if config.Encrypt {

		}
	}

	defer out.Cleanup()
}

// Cleanup - Removes the tempfile on disk
func (out *OutputHandler) Cleanup() {
	errClose := out.tmp.Close()
	if errClose != nil {
		log.Fatalf("Failed to close %s (err: %s)", out.tmp.Name(), errClose)
	}

	errRemove := os.Remove(out.tmp.Name())
	if errRemove != nil {
		log.Fatalf("Failed to remove %s (err: %s)", out.tmp.Name(), errRemove)
	}

}

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
func NewOutput(source string, destination string) (OutputHandler, error) {
	sourceFile, err := os.Open(source)
	if err != nil {
		log.Errorf("Failed to open %s (err: %s)", source, err)
		return OutputHandler{}, err
	}

	tempfile, err := ioutil.TempFile(os.TempDir(), "gcp")

	if err != nil {
		defer sourceFile.Close()
		log.Errorf("Failed to create temp file for %s (err: %s)", source, err)
		return OutputHandler{}, err
	}

	out := OutputHandler{
		source: sourceFile, tmp: tempfile, destination: destination}

	return out, nil
}

// Done - Called when we've finished processing the file.
func (out *OutputHandler) Done() {
	processing--
	wait.Done()

	errClose := out.tmp.Close()
	if errClose != nil {
		log.Fatalf("Failed to close %s (err: %s)", out.tmp.Name(), errClose)
	}

	errRemove := os.Remove(out.tmp.Name())
	if errRemove != nil {
		log.Fatalf("Failed to remove %s (err: %s)", out.tmp.Name(), errRemove)
	}

	log.Debug("%s", out.destination)
}

// Process - Opens the source file, performs operations (compress/encrypt) and
// writes the output to the temporary file.
func (out *OutputHandler) Process() {
	defer out.Done()

	log.Debug("%s -> %s", out.source.Name(), out.destination)
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

		if config.Compress {

		}
		if config.Encrypt {

		}

		out.tmp.Write(data)
	}

}

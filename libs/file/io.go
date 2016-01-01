package files

// io.go - Responsible for performing input and output operations.

import (
	"../config"
	"bytes"
	"code.google.com/p/lzma"
	"io"
	"io/ioutil"
	"os"
)

const maxReadSize = 5e+7 // 50MB

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

	processing++
	return out, nil
}

// SourceSize - Returns the size of the source in bytes
func (out *OutputHandler) SourceSize() int64 {
	stat, err := out.source.Stat()
	if err != nil {
		log.Fatalf("Failed to stat %s (err: %s)", out.source.Name(), err)
	}
	return stat.Size()
}

// DestinationSize - The size of the destination file
func (out *OutputHandler) DestinationSize() int64 {
	stat, err := out.tmp.Stat()
	if err != nil {
		log.Fatalf("Failed to stat %s (err: %s)", out.tmp.Name(), err)
	}
	return stat.Size()
}

// ReadSize - Returns how many bytes we should read at one time.
func (out *OutputHandler) ReadSize() int64 {
	sourceSize := out.SourceSize()

	if sourceSize < maxReadSize {
		return sourceSize
	}

	autoReadSize := sourceSize / 50

	if autoReadSize > maxReadSize {
		return maxReadSize
	}
	return autoReadSize
}

// Done - Called when we've finished processing the file.
func (out *OutputHandler) Done() {
	processing--
	wait.Done()

	sourceSize := out.SourceSize()
	sizeDifference := sourceSize - out.DestinationSize()
	compressionRatio := int(
		(float64(sizeDifference) / float64(sourceSize)) * 100)

	tmpErrClose := out.tmp.Close()
	if tmpErrClose != nil {
		log.Fatalf(
			"Failed to close %s (err: %s)", out.tmp.Name(), tmpErrClose)
	}

	tmpErrRemove := os.Remove(out.tmp.Name())
	if tmpErrRemove != nil {
		log.Fatalf(
			"Failed to remove %s (err: %s)", out.tmp.Name(), tmpErrRemove)
	}

	srcCloseErr := out.source.Close()
	if srcCloseErr != nil {
		log.Fatalf(
			"Failed to close %s (err: %s)", out.source.Name(), srcCloseErr)
	}

	log.Infof("%s (%d%%)", out.destination, compressionRatio)
}

// Process - Opens the source file, performs operations (compress/encrypt) and
// writes the output to the temporary file.
func (out *OutputHandler) Process() {
	defer out.Done()
	readSize := out.ReadSize()

	log.Debug("%s -> %s", out.source.Name(), out.destination)
	for {
		// Read n bytes from the source file
		data := make([]byte, readSize)
		_, err := out.source.Read(data)

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf(
				"Failed to read from %s (err: %s)", out.source.Name(), err)
		}

		if config.Compress {
			var compressedData bytes.Buffer
			lzmaWriter := lzma.NewWriterSize(
				&compressedData, lzma.BestCompression)
			lzmaWriter.Write(data)
			lzmaWriter.Close()
			data = compressedData.Bytes()
		}

		if config.Encrypt {

		}

		out.tmp.Write(data)
	}

}

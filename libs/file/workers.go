package files

// workers.go - contains the functions which perform various units of 'work'
// to copy the file such as reading data from a file,
// compression and encryption

import (
	"bytes"
	"code.google.com/p/lzma"
	"io"
)

// Channels - Defines what channels we use to communicate with workers on.
type Channels struct {
	paths chan string
	files chan File
}

var channels Channels

func openpaths(channel <-chan string) {
	for path := range channel {
		channels.files <- NewFile(path)
	}
}

func processfiles(channel <-chan File) {
	for file := range channel {
		log.Debugf(file.source.Name())

		encrypt := file.ShouldEncrypt()
		compress := file.ShouldEncrypt()

		// We have nothing else to do here in this case.
		if !compress && !encrypt {
			file.Done()
			continue
		}

		for {
			data := make([]byte, file.readsize)
			bytesRead, err := file.source.Read(data)
			bytesRead64 := int64(bytesRead)

			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(
					"Failed to read %s (err: %s)", file.source.Name(), err)
			}

			// It's possible we didn't read as many bytes
			// from the file as we allocated for `data`.  If this
			// is the case, resize data so it matches the number
			// of bytes read.  Otherwise we end up with empty bytes
			// in the file we're writing to disk.
			if file.readsize > bytesRead64 {
				data = append([]byte(nil), data[:bytesRead]...)
			}

			if compress {
				var compressed bytes.Buffer
				lzmaWriter := lzma.NewWriterSizeLevel(
					&compressed, bytesRead64, lzma.BestCompression)
				_, err := lzmaWriter.Write(data)
				lzmaWriter.Close()

				if err != nil {
					log.Fatalf(
						"Compression failed for %s (err: %s)",
						file.source.Name(), err)
				}
				data = compressed.Bytes()
			}

			if encrypt {

			}

			file.tempout.Write(data)
		}
		file.Done()
	}
}

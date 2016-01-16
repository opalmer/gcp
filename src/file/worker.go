package files

// worker.go - contains the functions which perform various units of 'work'
// to copy the file such as reading data from a file,
// compression and encryption

import (
	"bytes"
	"config"
	"github.com/opalmer/lzma"
	"github.com/ryanuber/go-glob"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

var processPaths chan string
var processing = sync.WaitGroup{}

const maxReadSize = 5e+7 // 50MB

func operations(path string) (bool, bool) {
	compress := true
	encrypt := true

	for _, exclusionPattern := range config.ExcludeCompression {
		if glob.Glob(exclusionPattern, path) {
			compress = false
		}
	}
	for _, exclusionPattern := range config.ExcludeEncryption {
		if glob.Glob(exclusionPattern, path) {
			encrypt = false
		}
	}
	return compress, encrypt
}

func outputFilePath(path string, compress bool, encrypt bool) string {
	relative, err := filepath.Rel(config.Source, path)
	if err != nil {
		log.Fatal(err)
	}
	output := filepath.Join(config.Destination, relative)

	// Add file extensions if necessary
	if compress {
		output += ".lzma"
	}

	if encrypt {
		output += ".aes"
	}

	return output
}

func closeOrFatal(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Fatalf("Failed to close %s: %s", file.Name(), err)
	}
}

func mkdir(path string) {
	if Exists(path) {
		return
	}
	err := os.MkdirAll(path, 0700)
	if err != nil {
		log.Fatalf("Failed to create %s: %s", path, err)
	}
}

func rename(src string, dst string) {
	mkdir(filepath.Dir(dst))
	err := os.Rename(src, dst)
	if err != nil {
		log.Fatalf("Failed to rename %s -> %s: %s", src, dst, err)
	}
}

func compressBytes(data []byte, size int64) []byte {
	var compressed bytes.Buffer
	lzmaWriter := lzma.NewWriterSizeLevel(
		&compressed, size, lzma.BestCompression)
	_, err := lzmaWriter.Write(data)
	lzmaWriter.Close()

	if err != nil {
		log.Fatalf("Compression failed: %s", err)
	}

	return compressed.Bytes()
}

func handle(
	source *os.File, output *os.File, chunkSize int64,
	compress bool, encrypt bool) {
	for {
		buffer := make([]byte, chunkSize)
		bytesRead, err := source.Read(buffer)

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failed to read from file(err: %s)", err)
		}

		bytesRead64 := int64(bytesRead)

		// It's possible we didn't read as many bytes
		// from the file as we allocated for `data`.  If this
		// is the case, resize data so it matches the number
		// of bytes read.  Otherwise we end up with empty bytes
		// in the file we're writing to disk.
		if chunkSize > bytesRead64 {
			buffer = append([]byte(nil), buffer[:bytesRead]...)
		}

		if compress {
			buffer = compressBytes(buffer, bytesRead64)
		}

		if encrypt {

		}
		output.Write(buffer)
	}
}

func process(channel <-chan string) {
	for path := range channel {
		compress, encrypt := operations(path)
		output := outputFilePath(path, compress, encrypt)

		if Exists(output) {
			log.Debugf("Skip %s -> %s", path, output)

		} else if config.DryRun {
			log.Infof("[DRY-RUN] Copy %s -> %s", path, output)

		} else {
			log.Infof("Copy %s -> %s", path, output)

			// Open the temporary output file.
			tempfile, err := ioutil.TempFile(os.TempDir(), "gcp")
			if err != nil {
				log.Fatalf("Failed to create temp file: %s", err)
			}

			// Open the source file
			srcFile, err := os.Open(path)
			if err != nil {
				log.Fatalf("Failed to open %s: %s", path, err)
			}

			// Direct copy
			if !compress && !encrypt {
				_, err := io.Copy(tempfile, srcFile)

				if err != nil {
					log.Fatalf("Copy failed: %s", err)
				}

			} else {
				// Determine how large of a chunk we should be reading
				// from the file.
				stat, err := srcFile.Stat()
				if err != nil {
					log.Fatalf("Failed to stat %s: %s", path, err)
				}

				chunkSize := stat.Size()
				if chunkSize > maxReadSize {
					chunkSize = maxReadSize
				}

				handle(srcFile, tempfile, chunkSize, compress, encrypt)
			}

			tempfile.Close()
			srcFile.Close()
			rename(tempfile.Name(), output)

			// TODO This shouldn't be required?  Without this the above tends
			// to chew through memory :/  Probably has something to do with
			// the byte slices...
			runtime.GC()
		}

		processing.Done()
	}
}

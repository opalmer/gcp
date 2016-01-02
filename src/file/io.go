package files

// io.go - Responsible for performing input and output operations.

import (
	"bytes"
	"config"
	"github.com/opalmer/lzma"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const maxReadSize = 5e+7 // 50MB

// File - The main object used for storing and processing a single file.
type File struct {
	sourcepath     string
	source         *os.File
	tempout        *os.File
	readsize       int64
	outpath        string
	shouldCompress bool
	shouldEncrypt  bool
}

func outpath(file *File) string {
	// Setup the output path
	outpath := filepath.Join(config.Destination, file.source.Name())
	if file.shouldCompress {
		outpath += ".lzma"
	}
	if file.shouldEncrypt {
		outpath += ".aes"
	}
	return outpath
}

func readsize(file *File) int64 {
	stat, err := file.source.Stat()
	if err != nil {
		log.Fatal("Failed to start %s (err: %s)", file.source.Name(), err)
	}
	size := stat.Size()
	if size > maxReadSize {
		size = maxReadSize
	}
	return size
}

func shouldCompress(file *File) bool {
	if !config.Compress {
		return false
	}
	name := strings.ToLower(file.source.Name())
	if strings.HasSuffix(name, ".iso") {
		return false
	}
	return true
}

func shouldEncrypt(file *File) bool {
	if !config.Encrypt {
		return false
	}
	name := strings.ToLower(file.source.Name())
	if strings.HasSuffix(name, ".iso") {
		return false
	}
	return true
}

func compress(file *File, data []byte, bytesRead int64) ([]byte, error) {
	var compressed bytes.Buffer
	lzmaWriter := lzma.NewWriterSizeLevel(
		&compressed, bytesRead, lzma.BestCompression)
	_, err := lzmaWriter.Write(data)
	lzmaWriter.Close()

	if err != nil {
		log.Warningf(
			"Compression failed for %s (err: %s)",
			file.source.Name(), err)
		return nil, err
	}

	return compressed.Bytes(), nil
}

// TODO
func encrypt(file *File, data []byte) ([]byte, error) {
	return data, nil
}

// Open - Opens the input and output files where applicable, also sets up the
// output path.
func (file *File) open() error {
	// Open the source file
	source, err := os.Open(file.sourcepath)

	if err != nil {
		return err
	}

	// Open the temporary output file.
	tempout, err := ioutil.TempFile(os.TempDir(), "gcp")

	if err != nil {
		source.Close()
		return err
	}

	// Establish the attributes we'll need for working
	// with the file.
	//  NOTE: Order matters here.
	file.source = source
	file.tempout = tempout
	file.readsize = readsize(file)
	file.shouldCompress = shouldCompress(file)
	file.shouldEncrypt = shouldEncrypt(file)
	file.outpath = outpath(file)

	return nil
}

// Performs the main IO operations responsible for
// processing the file.  The results end up in the
// temporary output path.
func (file *File) process() error {
	log.Debugf("%s -> %s", file.source.Name(), file.tempout.Name())
	defer file.source.Close()

	// Files which are neither compressed or encrypted will
	// just be coped over to their temporary output.
	if !file.shouldCompress && !file.shouldEncrypt {
		io.Copy(file.tempout, file.source)
		return nil
	}

	// Iterate over the whole file and compress and/or encrypt
	for {
		data := make([]byte, file.readsize)
		bytesRead, err := file.source.Read(data)
		bytesRead64 := int64(bytesRead)

		if err == io.EOF {
			break
		} else if err != nil {
			log.Warningf(
				"Failed to read %s (err: %s)", file.source.Name(), err)
			return err
		}

		// It's possible we didn't read as many bytes
		// from the file as we allocated for `data`.  If this
		// is the case, resize data so it matches the number
		// of bytes read.  Otherwise we end up with empty bytes
		// in the file we're writing to disk.
		if file.readsize > bytesRead64 {
			data = append([]byte(nil), data[:bytesRead]...)
		}

		if file.shouldCompress {
			data, err = compress(file, data, bytesRead64)
			if err != nil {
				return err
			}
		}

		if file.shouldEncrypt {
			data, err = encrypt(file, data)
			if err != nil {
				return err
			}
		}

		file.tempout.Write(data)
	}

	return nil
}

// Responsible for saving the file to the final location.
func (file *File) save() error {

	log.Infof("%s -> %s", file.source.Name(), file.outpath)

	err := file.tempout.Sync()
	if err != nil {
		log.Warning("Failed to sync temp output")
		return err
	}

	err = file.tempout.Close()
	if err != nil {
		log.Warning("Failed to close temp output")
		return err
	}

	directory := filepath.Dir(file.outpath)
	err = os.MkdirAll(directory, 0700)
	if err != nil {
		log.Warningf("Failed to create %s", directory)
		return err
	}

	err = os.Rename(file.tempout.Name(), file.outpath)
	if err != nil {
		log.Warning("Failed to rename file")
		return err
	}

	return nil

}

// Performs some final cleanup in the event of an error.  This is mainly
// aimed at closing the file handles and removing the temp. output file.  We
// ignore errors in this block of code because we expect processfiles() to
// call log.Fatal* soon after this function.
func (file *File) clean() {
	defer filesProcessing.Done()

	if file.source != nil {
		file.source.Close()
	}

	if file.tempout != nil {
		file.tempout.Close()
		os.Remove(file.tempout.Name())
	}
}

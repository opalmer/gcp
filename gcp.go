package main

import (
	//	"./libs/file"
	"flag"
	"github.com/op/go-logging"
	"os"
)

var log *logging.Logger

func initializeLogging(levelInput string) {
	// Convert the input log level name to
	// a logging.Level instance.
	level, err := logging.LogLevel(levelInput)
	if err != nil {
		log.Fatal(err)
	}

	// Global format configuration
	formatter := logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} %{level:7s} %{color:reset} %{message}`,
	)
	logging.SetFormatter(formatter)

	// Setup stderr to handle ERROR and above
	stderr := logging.NewLogBackend(os.Stderr, "", 0)
	stderrLeveled := logging.AddModuleLevel(stderr)
	stderrLeveled.SetLevel(logging.ERROR, "gcp")

	// If the log level has been set to something larger
	// than we'd capture in stdout then just make stderr
	// the one backend and return.  This prevents us from
	// possibly duplicating logs.
	if level <= logging.WARNING {
		logging.SetBackend(stderrLeveled)
		return
	}

	stdout := logging.NewLogBackend(os.Stdout, "", 0)
	stdoutLeveled := logging.AddModuleLevel(stdout)
	logging.SetBackend(stdoutLeveled, stderrLeveled)
}

func copyPaths(
	encryption bool, compression bool, source string, destination string) {
	log.Debug("Copy %s -> %s", source, destination)
}

// TODO config handling
// TODO   add exclusions (fnmatch)
// TODO   add inclusions to ignore from exclusions (fnmatch)


func main() {
	disableCompression := flag.Bool(
		"disable-compression", false, "If provided files will not be zipped up")
	disableEncryption := flag.Bool(
		"disable-encryption", false, "If provided files will not be encrypted")
	logLevelInput := flag.String(
		"log", "debug", "The logging level")

	flag.Parse()

	initializeLogging(*logLevelInput)
	log = logging.MustGetLogger("gcp")

	if *disableCompression {
		log.Info("Compression has been disabled")
	}

	if *disableEncryption {
		log.Warning("Encryption has been disabled")
	}

	source := "/tmp/foo"
	destination := "/tmp/bar"

	copyPaths(*disableEncryption, *disableEncryption, source, destination)
}

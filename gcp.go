package main

import (
	"config"
	"file"
	"flag"
	"github.com/op/go-logging"
	"io/ioutil"
	"os"
	"runtime"
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
		`%{color}%{time:15:04:05.000} %{shortfunc} %{level} %{color:reset} %{message}`,
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

func main() {
	logLevelInput := flag.String(
		"log", "debug", "The logging level")
	configPath := flag.String(
		"config", "", "A direct path to a configuration file.")
	skipRelativeCheck := flag.Bool(
		"ignore-relative-check", false,
		`If provided, don't halt if the source and destination paths appear
		to relative to each other.`)
	encryptionKey := flag.String(
		"key", "",
		"A string to use as the encryption key or the path to a file")
	concurrency := flag.Int(
		"concurrency", runtime.NumCPU(),
		"A sudo-limit which tries to limit concurrency some.")
	dryRun := flag.Bool(
		"dry-run", false,
		"When provided, only print the operations to perform")

	flag.Parse()
	args := flag.Args()

	initializeLogging(*logLevelInput)
	log = logging.MustGetLogger("gcp")

	// Make sure we're not missing any input arguments
	if len(args) != 2 {
		log.Error(
			"Expected two input arguments: %s <source> <destination>",
			os.Args[0])
		flag.Usage()
		os.Exit(1)
	}

	// Get the config module setup
	config.Load(*configPath)
	config.Source = files.AbsolutePath(flag.Arg(0))
	config.Destination = files.AbsolutePath(flag.Arg(1))
	config.DryRun = *dryRun
	config.Concurrency = *concurrency
	config.CryptoKey = *encryptionKey
	config.Exclude = config.GetSlice("exclude")
	config.Include = config.GetSlice("include")
	config.ExcludeCompression = config.GetSlice("exclude_compression")
	config.ExcludeEncryption = config.GetSlice("exclude_encryption")

	if config.Concurrency < 1 {
		log.Fatalf("-concurrency must be at least one")
	}

	if !files.Exists(config.Source) {
		log.Fatalf("Source '%s' does not exist", config.Source)
	}

	if files.IsRelative(config.Source, config.Destination) {
		if !*skipRelativeCheck {
			log.Fatal(
				"Source and destination appear to be relative to one another")
		} else {
			log.Warning(
				"Source and destination appear to be relative to one another")
		}
	}

	if len(*encryptionKey) < 1 {
		log.Fatal("You must provide an encryption key to -key")
	}

	// Reading the encryption key
	data, err := ioutil.ReadFile(*encryptionKey)
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}

	if err == nil {
		log.Info("Reading encryption key from file '%s'", *encryptionKey)
		config.CryptoKey = string(data[:])
	} else {
		log.Info("Crypto key does not appear to be a file")
		config.CryptoKey = *encryptionKey
	}

	files.Copy()
}

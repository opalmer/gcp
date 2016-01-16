package config

// load.go - Module responsible for loading the configuration object

import (
	"github.com/go-ini/ini"
	"github.com/op/go-logging"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

var cfg ini.File

var log = logging.MustGetLogger("gcp")

func loadFile(path string) {
	if len(path) < 1 {
		return
	}

	log.Debugf("Loading %s", path)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read %s (err: %s)", path, err)
	}

	if cfg.Append(data) != nil {
		log.Fatalf("Failed to append bytes from %s", path)
	}
}

// Load - Loads the configuration
func Load(path string) {
	cfg = *ini.Empty()
	cfg.Append([]byte(`
		[gcp]
		include =
		exclude =
		exclude_compression = *.iso,*.png,*.jpg,*.jpeg
		exclude_encryption = *.iso,
	`))

	// Try to load from the environment
	value, found := os.LookupEnv("GCP_CONFIG")
	if found {
		loadFile(value)
	}

	// Load from $HOME/.gcp.ini
	user, err := user.Current()
	if err == nil {
		path := filepath.Join(user.HomeDir, ".gcp.ini")
		loadFile(path)
	}

	loadFile(path)
}

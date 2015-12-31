package config

// load.go - Module responsible for loading the configuration object

import (
	"github.com/go-ini/ini"
	"github.com/op/go-logging"
	"os"
	"os/user"
	"path/filepath"
)

var config = Default()
var log = logging.MustGetLogger("gcp")

// Default - returns the default configuration
func Default() *ini.File {
	cfg := ini.Empty()
	cfg.Append("default", `
		[gcp]
		crypto_key =
		include =
		exclude =
	`)
	return cfg
}

// LoadConfigFile - Update the configuration with data from `path`
func LoadConfigFile(path string) {
	if len(path) < 1 {
		return
	}
	err := config.Append(path, path)

	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	} else if err != nil {
		log.Warning("Failed to load config from file '%s'", path)
	} else {
		log.Debug("Loaded config from file '%s'", path)
	}
}

// Load - Loads the configuration
func Load(path string, encryptionKey string) {
	// Try to load from the environment
	value, found := os.LookupEnv("GCP_CONFIG")
	if found {
		LoadConfigFile(value)
	}

	// Load from $HOME/.gcp.ini
	user, err := user.Current()
	if err != nil {
		path := filepath.Join(user.HomeDir, ".gcp.ini")
		LoadConfigFile(path)
	}

	LoadConfigFile(path)

	// A specific encryption key was provided, use that instead.
	if len(encryptionKey) > 0 {
		section := config.Section("gcp")
		key := section.Key("crypto_key")
		key.SetValue(encryptionKey)
	}

}
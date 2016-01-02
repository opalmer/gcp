package config

// load.go - Module responsible for loading the configuration object

import (
	"github.com/go-ini/ini"
	"github.com/op/go-logging"
	"os"
	"os/user"
	"path/filepath"
)

var cfg = Default()
var log = logging.MustGetLogger("gcp")

// Default - returns the default configuration
func Default() *ini.File {
	cfg := ini.Empty()
	cfg.Append([]byte(`
		[gcp]
		encrypt = true
		compress = true
		dry_run = false
		crypto_key =
		include =
		exclude = .DS_Store,.git,.svn,.hg,.egg*,__pycache__,.idea,*.pyc
	`))
	return cfg
}

// LoadConfigFile - Update the configuration with data from `path`
func LoadConfigFile(path string) {
	if len(path) < 1 {
		return
	}

	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Warningf("File %s does not exist", path)
	} else if err != nil {
		log.Warningf("Can't load %s (err: %s)", path, err)
	} else if stat.IsDir() {
		log.Warningf("%s is a directory", path)
	} else {
		err := cfg.Append(path, path)
		if err != nil && !os.IsNotExist(err) {
			log.Fatal(err)
		} else if err != nil {
			log.Warning(
				"Failed to load config from file '%s' (err: %s)", path, err)
		} else {
			log.Debug("Loaded config from file '%s'", path)
		}
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
	if err == nil {
		path := filepath.Join(user.HomeDir, ".gcp.ini")
		LoadConfigFile(path)
	}

	LoadConfigFile(path)

	// A specific encryption key was provided, use that instead.
	if len(encryptionKey) > 0 {
		section := cfg.Section("gcp")
		key := section.Key("crypto_key")
		key.SetValue(encryptionKey)
	}

	CryptoKey = GetKey("crypto_key").Value()
}

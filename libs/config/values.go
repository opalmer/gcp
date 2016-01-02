package config

import (
	"github.com/go-ini/ini"
	"strings"
)

// GetKey - Retrieves a specific key from the configuration in the
// 'gcp' section and returns the value.
func GetKey(key string) *ini.Key {
	section, err := cfg.GetSection("gcp")
	if err != nil {
		log.Fatalf("Failed to retrieve config section gcp (err: %s)", err)
	}

	value, err := section.GetKey(key)
	if err != nil {
		log.Fatalf("Failed to retrieve key gcp.%s (err: %s)", key, err)
	}

	return value
}

// GetBool - Returns a boolean for the given ``key``
func GetBool(key string) bool {
	result, err := GetKey(key).Bool()
	if err != nil {
		log.Fatalf("Failed to retrieve gcp.encrypt as bool (err: %s)", err)
	}
	return result
}

// GetSlice - Converts a config key from a comma separated entry into a slice.
func GetSlice(key string) []string {
	source := GetKey(key).Value()
	output := []string{}

	for _, value := range strings.Split(source, ",") {
		trimmed := strings.TrimSpace(value)
		if len(trimmed) > 0 {
			output = append(output, trimmed)
		}
	}

	return output
}

// Encrypt - True if encryption is enabled in the config
var Encrypt = GetBool("encrypt")

// Compress - True if compression is enabled in the config
var Compress = GetBool("compress")

// Include -- All inclusion patterns
var Include = GetSlice("include")

// Exclude - All exclusion patterns
var Exclude = GetSlice("exclude")

// Destination - The root path where files will be copied to
var Destination string

// Source - The location to copy files from
var Source string

// CryptoKey - The key to use for encrypting things
var CryptoKey string

// DryRun - Disables certain opertations if True.
var DryRun = GetBool("dry_run")

// Concurrency - Limits concurrency in a few places
var Concurrency int

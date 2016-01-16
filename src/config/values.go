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
		log.Fatal(err)
	}

	value, err := section.GetKey(key)
	if err != nil {
		log.Fatal(err)
	}

	return value
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

// Include -- Path inclusion patterns
var Include []string

// Exclude - Path exclusion patterns
var Exclude []string

// ExcludeCompression - Patterns to dictate which files should not be compressed
var ExcludeCompression []string

// ExcludeEncryption - Patterns to dictate which files should not be encrypted
var ExcludeEncryption []string

// Destination - The root path where files will be copied to
var Destination string

// Source - The location to copy files from
var Source string

// CryptoKey - The key to use for encrypting things
var CryptoKey string

// Concurrency - Limits concurrency in a few places
var Concurrency int

// DryRun - Used for testing
var DryRun bool

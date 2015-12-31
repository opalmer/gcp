package config

import (
	"strings"
)

// Get - Retrieves a specific key from the configuration in the 'gcp' section.
func Get(key string) string {
	section, err := cfg.GetSection("gcp")
	if err != nil {
		log.Fatalf("Failed to retrieve config section gcp (err: %s)", err)
	}

	value, err := section.GetKey(key)
	if err != nil {
		log.Fatalf("Failed to retrieve key gcp.%s (err: %s)", key, err)
	}
	return value.Value()

}

// GetSlice - Converts a config key from a comma separated entry into a slice.
func GetSlice(key string) []string {
	source := Get(key)
	output := []string{}

	for _, value := range strings.Split(source, ",") {
		output = append(output, strings.TrimSpace(value))
	}

	return output
}

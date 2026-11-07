package config

import (
	"fmt"
	"strings"
)

const defaultPDS = "https://bsky.social"

// Config holds the only settings Biosky needs.
type Config struct {
	Handle      string
	AppPassword string
	PDS         string
}

// Load reads configuration from the provided environment lookup.
func Load(getenv func(string) string) (Config, error) {
	handle := strings.TrimSpace(getenv("BIOSKY_HANDLE"))
	password := strings.TrimSpace(getenv("BIOSKY_APP_PASSWORD"))
	pds := strings.TrimSpace(getenv("BIOSKY_PDS"))
	if pds == "" {
		pds = defaultPDS
	}
	var missing []string
	if handle == "" {
		missing = append(missing, "BIOSKY_HANDLE")
	}
	if password == "" {
		missing = append(missing, "BIOSKY_APP_PASSWORD")
	}
	if len(missing) > 0 {
		return Config{}, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}
	return Config{Handle: handle, AppPassword: password, PDS: pds}, nil
}

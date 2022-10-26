package testconfig

import (
	"admincheckapi/api/config"
	"embed"
	"io"
	"os"
	"testing"
)

const defaultConfig = "dev-config.yaml"

//go:embed *
var fs embed.FS

// SetTestSetup sets Setup to test configuration
func Set(t *testing.T) {
	path := os.Getenv("CONFIG")
	if path == "" {
		path = defaultConfig
	}
	r, err := fs.Open(path)
	if err != nil {
		t.Fatalf("Unable to open config file: %s", err)
	}
	input, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("Unable to read config file: %s", err)
	}
	s, err := config.NewSetupValueSet(input)
	if err != nil {
		t.Fatalf("Error loading valid setup: %s", err)
	}
	if s == nil {
		t.Fatalf("Empty Setup")
	}
	config.Setup = s
}

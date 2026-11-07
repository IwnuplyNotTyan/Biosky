package config

import (
	"strings"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	env := map[string]string{
		"BIOSKY_HANDLE":       " me.bsky.social ",
		"BIOSKY_APP_PASSWORD": " app-pass ",
	}
	got, err := Load(func(k string) string { return env[k] })
	if err != nil {
		t.Fatal(err)
	}
	if got.Handle != "me.bsky.social" || got.AppPassword != "app-pass" {
		t.Fatalf("unexpected config: %+v", got)
	}
	if got.PDS != "https://bsky.social" {
		t.Fatalf("PDS = %q", got.PDS)
	}
}

func TestLoadCustomPDS(t *testing.T) {
	env := map[string]string{
		"BIOSKY_HANDLE":       "me.bsky.social",
		"BIOSKY_APP_PASSWORD": "app-pass",
		"BIOSKY_PDS":          "https://pds.example",
	}
	got, err := Load(func(k string) string { return env[k] })
	if err != nil {
		t.Fatal(err)
	}
	if got.PDS != "https://pds.example" {
		t.Fatalf("PDS = %q", got.PDS)
	}
}

func TestLoadMissingRequired(t *testing.T) {
	_, err := Load(func(string) string { return "" })
	if err == nil {
		t.Fatal("expected an error")
	}
	for _, name := range []string{"BIOSKY_HANDLE", "BIOSKY_APP_PASSWORD"} {
		if !strings.Contains(err.Error(), name) {
			t.Errorf("error %q does not name %s", err, name)
		}
	}
}

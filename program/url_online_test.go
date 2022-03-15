package program_test

import (
	"context"
	"testing"

	"go.xargs.dev/bindl/download"
	"go.xargs.dev/bindl/program"
	"sigs.k8s.io/yaml"
)

// WARNING: tests this file requires connection to GitHub to succeed

func TestIntegrationDownloadArchy(t *testing.T) {
	if testing.Short() {
		t.Skipf("skipping integration test in short mode")
	}
	p := program.URLProgram{}
	err := yaml.Unmarshal([]byte(rawArchyLockManifest), &p)
	if err != nil {
		t.Fatalf("failed when expecting pass: %v", err)
	}
	_, err = p.DownloadArchive(context.Background(), &download.HTTP{}, "linux", "amd64")
	if err != nil {
		t.Fatal(err)
	}
}

package program_test

import (
	"context"
	"testing"

	"go.xargs.dev/bindl/program"
	"sigs.k8s.io/yaml"
)

// WARNING: tests this file requires connection to GitHub to succeed

func TestIntegrationConvertProgramConfigURLProviderToURL(t *testing.T) {
	if testing.Short() {
		t.Skipf("skipping integration test in short mode")
	}
	c := &program.Config{}
	if err := yaml.Unmarshal([]byte(rawArchyProgramManifest), c); err != nil {
		t.Fatal(err)
	}

	u, err := c.URLProgram(context.Background(), map[string][]string{
		"linux":  []string{"amd64"},
		"darwin": []string{"arm64"},
	})
	if err != nil {
		t.Fatal(err)
	}

	assert(t,
		"ffe042e8281d9b51d021e4551b1051edbd8d0bb64235f9ac3feb4495338c0b58",
		u.Checksums["archy_0.1.1_Darwin_arm64.tar.gz"].Archive)
}

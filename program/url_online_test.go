package program_test

import (
	"context"
	"testing"

	"go.xargs.dev/bindl/download"
	"go.xargs.dev/bindl/program"
	"sigs.k8s.io/yaml"
)

// WARNING: tests this file requires connection to GitHub to succeed

const rawArchyProgramManifest = `
name: archy
version: 0.1.1
url: 'https://github.com/xargs-dev/archy/releases/download/v{{ .Version }}/archy_{{ .Version }}_{{ .OS }}_{{ .Arch }}.tar.gz'
overlay:
  OS:
    linux: Linux
    darwin: Darwin
  Arch:
    amd64: x86_64
checksums:
  archy_0.1.1_Linux_x86_64.tar.gz:
    _archive: b999ac46efeb15ea1e304c732ef42a7a313a773c61deea2192d78025794939c2
`

func TestIntegrationDownloadArchy(t *testing.T) {
	if testing.Short() {
		t.Skipf("skipping integration test in short mode")
	}
	p := program.URLProgram{}
	err := yaml.Unmarshal([]byte(rawArchyProgramManifest), &p)
	if err != nil {
		t.Fatalf("failed when expecting pass: %v", err)
	}
	_, err = p.DownloadArchive(context.Background(), &download.HTTP{}, "linux", "amd64")
	if err != nil {
		t.Fatal(err)
	}
}

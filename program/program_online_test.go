package program_test

import (
	"testing"

	"go.xargs.dev/bindl/program"
	"sigs.k8s.io/yaml"
)

// WARNING: tests this file requires connection to GitHub to succeed

func TestIntegrationConvertProgramConfigURLProviderToURL(t *testing.T) {
	raw := `
name: archy
version: 0.1.1
provider: url
path: https://github.com/xargs-dev/archy/releases/download/v{{ .Version }}/archy_{{ .Version }}_{{ .OS }}_{{ .Arch }}.tar.gz
overlay:
  OS:
    linux: Linux
    darwin: Darwin
  Arch:
    amd64: x86_64
checksums:
  _src: https://github.com/xargs-dev/archy/releases/download/v{{ .Version }}/checksums.txt
  archy_0.1.1_Linux_x86_64.tar.gz: notrealchecksum
`
	c := &program.Config{}
	if err := yaml.Unmarshal([]byte(raw), c); err != nil {
		t.Fatal(err)
	}

	p, err := c.Program(map[string][]string{
		"linux":  []string{"amd64"},
		"darwin": []string{"arm64"},
	})
	assert(t, nil, err)

	u, ok := p.(*program.URLProgram)
	if !ok {
		t.Fatal("'p' implements Program, but not of type URLProgram")
	}

	assert(t, "notrealchecksum", u.Checksums["archy_0.1.1_Linux_x86_64.tar.gz"].Archive)
	assert(t,
		"ffe042e8281d9b51d021e4551b1051edbd8d0bb64235f9ac3feb4495338c0b58",
		u.Checksums["archy_0.1.1_Darwin_arm64.tar.gz"].Archive)
}

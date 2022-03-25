// Copyright 2022 Bindl Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package program_test

import (
	_ "embed"
	"os"
	"testing"
)

func assert[C comparable](t *testing.T, expected, got C) {
	if expected != got {
		t.Fatalf("mismatch!\nt-expected:\n\t%v\nt-got:\n\t%v", expected, got)
	}
}

//go:embed testdata/myprogram/myprogram
var myProgramBinary []byte

//go:embed testdata/myprogram.tar.gz
var myProgramTarGz []byte

func testDirFile(t *testing.T, filename string, data []byte) *os.File {
	f, err := os.CreateTemp(t.TempDir(), filename)
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.Write(data)
	if err != nil {
		t.Fatal(err)
	}
	// Reset to beginning of file, ready for reading
	f.Seek(0, 0)
	return f
}

const rawProgramManifest = `
name: myprogram
version: 0.1.0-rc.2
url: 'http://myurl.com/foo/{{ .Version }}/{{ .Name }}-{{ .OS }}-{{ .Arch }}.tar.gz'
overlay:
  OS:
    linux: Linux
  Arch:
    amd64: x86_64
checksums:
  myprogram-Linux-x86_64.tar.gz:
    archive: 61577c9d9010c0c7190428fe3c15f406209be3bd409c3b87fb767febd3a784b9
    binary: d5b12eda84454df3bf1a4729dc3cf39c124232f62bf2f33f4defb5432b60f08e
`

const rawArchyProgramManifest = `
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
`

const rawArchyLockManifest = `
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
    archive: b999ac46efeb15ea1e304c732ef42a7a313a773c61deea2192d78025794939c2
`

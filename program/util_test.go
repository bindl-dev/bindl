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
    _archive: 61577c9d9010c0c7190428fe3c15f406209be3bd409c3b87fb767febd3a784b9
    myprogram: d5b12eda84454df3bf1a4729dc3cf39c124232f62bf2f33f4defb5432b60f08e
`
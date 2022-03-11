package program_test

import (
	"bytes"
	"context"
	"testing"

	"sigs.k8s.io/yaml"

	"go.xargs.dev/bindl/download"
	"go.xargs.dev/bindl/program"
)

func TestProgramChecksumsYAMLUnmarshalJSON(t *testing.T) {
	p := program.URLProgram{}
	err := yaml.Unmarshal([]byte(rawProgramManifest), &p)
	if err != nil {
		t.Fatalf("failed when expecting pass: %v", err)
	}

	assert(t, "myprogram", p.Name())
	assert(t, "0.1.0-rc.2", p.Version)

	cs := p.Checksums["myprogram-Linux-x86_64.tar.gz"]
	assert(t, "61577c9d9010c0c7190428fe3c15f406209be3bd409c3b87fb767febd3a784b9", cs.Archive)
	assert(t, "d5b12eda84454df3bf1a4729dc3cf39c124232f62bf2f33f4defb5432b60f08e", cs.Binaries["myprogram"])
}

func TestProgramURL(t *testing.T) {
	p := &program.URLProgram{}
	err := yaml.Unmarshal([]byte(rawProgramManifest), p)
	if err != nil {
		t.Fatalf("failed when expecting pass: %v", err)
	}
	url, err := p.URL("linux", "amd64")
	if err != nil {
		t.Fatalf("failed to generate URL: %v", err)
	}
	assert(t, "http://myurl.com/foo/0.1.0-rc.2/myprogram-Linux-x86_64.tar.gz", url)
}

func TestDownloadChecksum(t *testing.T) {
	p := &program.URLProgram{}
	err := yaml.Unmarshal([]byte(rawProgramManifest), p)
	if err != nil {
		t.Fatalf("failed when expecting pass: %v", err)
	}

	f := testDirFile(t, "myprogram-Linux-x86_64.tar.gz", myProgramTarGz)
	defer f.Close()

	d := download.NewLocalFile(f)

	archive, err := p.DownloadArchive(context.Background(), d, "linux", "amd64")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(archive.Data, myProgramTarGz) {
		t.Fatalf("archive data mismatch:\nt-expected:\n\t%x\nt-got:\n\t%x", myProgramTarGz, archive.Data)
	}
	data, err := archive.Extract("myprogram")
	if err != nil {
		t.Fatalf("extracting from archive: %v", err)
	}
	if !bytes.Equal(data, myProgramBinary) {
		t.Fatalf("archive data mismatch:\nt-expected:\n\t%s\nt-got:\n\t%s", myProgramBinary, data)
	}
}

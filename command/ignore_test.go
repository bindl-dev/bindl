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

package command_test

import (
	"os"
	"testing"

	"go.xargs.dev/bindl/command"
	"go.xargs.dev/bindl/config"
)

func TestUpdateIgnoreFile(t *testing.T) {
	testCases := []struct {
		name         string
		binPathDir   string
		fileContents string
		want         string
	}{
		{
			name:         "Empty ignore file",
			binPathDir:   "bin",
			fileContents: "",
			want: `# Development and tool binaries
bin/*
`,
		},
		{
			name:         "NO ADD: bin/* variation",
			binPathDir:   "bin/",
			fileContents: "bin/*",
			want:         "bin/*",
		},
		{
			name:         "NO ADD: bin/ variation, binPathDir with trailing slash",
			binPathDir:   "bin/",
			fileContents: "bin/",
			want:         "bin/",
		},
		{
			name:         "NO ADD: bin/ variation, binPathDir no trailing slash",
			binPathDir:   "bin",
			fileContents: "bin/",
			want:         "bin/",
		},
		{
			name:         "NO ADD: bin variation",
			binPathDir:   "bin/",
			fileContents: "bin",
			want:         "bin",
		},
		{
			name:       "Ignore entry commented out",
			binPathDir: "binny",
			fileContents: `# binny/*
secret`,
			want: `# binny/*
secret

# Development and tool binaries
binny/*
`,
		},
		{
			name:       "End with newline",
			binPathDir: "binny",
			fileContents: `secret1
secret
`,
			want: `secret1
secret

# Development and tool binaries
binny/*
`,
		},
	}

	dir := t.TempDir()
	for _, tc := range testCases {
		conf := &config.Runtime{
			BinPathDir: tc.binPathDir,
		}
		f, err := os.CreateTemp(dir, "ignore*")
		if err != nil {
			t.Fatalf("[%s] failed to create temp file: %v", tc.name, err)
		}
		ignoreFileName := f.Name()
		if _, err := f.WriteString(tc.fileContents); err != nil {
			t.Fatalf("[%s] failed to write tc ignore file: %v", tc.name, err)
		}
		if err = f.Close(); err != nil {
			t.Fatalf("[%s] failed to close file: %v", tc.name, err)
		}
		err = command.UpdateIgnoreFile(conf, ignoreFileName)
		if err != nil {
			t.Fatalf("[%s] failed when expecting pass: %v", tc.name, err)
		}
		got, err := os.ReadFile(ignoreFileName)
		if err != nil {
			t.Fatalf("[%s] failed to read tc ignore file: %v", tc.name, err)
		}
		if string(got) != tc.want {
			t.Fatalf("[%s] got\n%s\nwant\n%s", tc.name, string(got), tc.want)
		}
	}
}

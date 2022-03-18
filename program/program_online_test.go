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

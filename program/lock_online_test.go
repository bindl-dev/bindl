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

	"go.xargs.dev/bindl/download"
	"go.xargs.dev/bindl/program"
	"sigs.k8s.io/yaml"
)

// WARNING: tests this file requires connection to GitHub to succeed

func TestIntegrationDownloadArchy(t *testing.T) {
	if testing.Short() {
		t.Skipf("skipping integration test in short mode")
	}
	p := program.Lock{}
	err := yaml.Unmarshal([]byte(rawArchyLockManifest), &p)
	if err != nil {
		t.Fatalf("failed when expecting pass: %v", err)
	}
	_, err = p.DownloadArchive(context.Background(), &download.HTTP{}, "linux", "amd64")
	if err != nil {
		t.Fatal(err)
	}
}

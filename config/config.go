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

package config

import "github.com/bindl-dev/bindl/program"

// Config is a configuration which is used to declare a project's dependencies.
// By default, this is the content of bindl.yaml
type Config struct {
	// Platforms is a matrix of OS and Arch for the binaries which
	// the project would like to save checksums on.
	// i.e. map["linux"][]{"amd64", "arm64"}
	Platforms map[string][]string `json:"platforms"`

	// Programs is a list of program specification
	Programs []*program.Config `json:"programs"`
}

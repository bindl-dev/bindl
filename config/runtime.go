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

// Runtime is a configuration which is primarily used by command/cli.
// All variables which users can change at runtime with global effect
// can be found here.
type Runtime struct {
	Path         string `envconfig:"CONFIG"`
	LockfilePath string `envconfig:"LOCK"`
	BinDir       string `envconfig:"BIN"`
	ProgDir      string `envconfig:"PROG"`

	OS   string `envconfig:"OS"`
	Arch string `envconfig:"ARCH"`

	UseCache bool `envconfig:"USE_CACHE"`
	Hardlink bool // Only use in containers

	Debug  bool `envconfig:"DEBUG"`
	Silent bool `envconfig:"SILENT"`
}

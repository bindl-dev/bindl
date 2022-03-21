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

import (
	"os"
	"time"

	"go.xargs.dev/bindl/internal"
	"go.xargs.dev/bindl/program"
	"sigs.k8s.io/yaml"
)

type Lock struct {
	Updated  time.Time             `json:"updated"`
	Programs []*program.URLProgram `json:"programs"`
}

func ParseLock(path string) (*Lock, error) {
	l := &Lock{}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(raw, l); err != nil {
		return nil, err
	}
	internal.Log().Debug().Int("programs", len(l.Programs)).Msg("parsing lockfile")
	if len(l.Programs) == 0 {
		internal.Log().Warn().Msg("no programs found in lockfile")
	}
	return l, nil
}

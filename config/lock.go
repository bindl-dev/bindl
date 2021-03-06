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
	"sort"
	"strings"

	"github.com/bindl-dev/bindl/internal"
	"github.com/bindl-dev/bindl/program"
	"sigs.k8s.io/yaml"
)

// Lock is a configuration which was generated from Config.
// By default, this is the content of .bindl-lock.yaml
type Lock struct {
	Programs ProgramLocks `json:"programs"`
}

func NewLock(p []*program.Lock) *Lock {
	l := &Lock{Programs: ProgramLocks(p)}
	sort.Stable(l.Programs)
	return l
}

// ParseLock reads a file from path and returns *Lock
func ParseLock(path string) (*Lock, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseLockBytes(raw)
}

// ParseLock reads a file from parameter and returns *Lock
func ParseLockBytes(b []byte) (*Lock, error) {
	l := &Lock{}
	if err := yaml.Unmarshal(b, l); err != nil {
		return nil, err
	}
	internal.Log().Debug().Int("programs", len(l.Programs)).Msg("parsing lockfile")
	if len(l.Programs) == 0 {
		internal.Log().Warn().Msg("no programs found in lockfile")
	}
	return l, nil
}

// ProgramLocks implements sort.Interface
// TODO: use generics?
type ProgramLocks []*program.Lock

func (p ProgramLocks) Len() int {
	return len(p)
}

func (p ProgramLocks) Less(i, j int) bool {
	return strings.Compare(p[i].Name, p[j].Name) < 0
}

func (p ProgramLocks) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

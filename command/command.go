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

package command

import (
	"errors"
	"fmt"

	"go.xargs.dev/bindl/config"
	"go.xargs.dev/bindl/internal"
	"go.xargs.dev/bindl/program"
)

var FailExecError = errors.New("failed to execute command, please troubleshoot logs")

func FilterPrograms(conf *config.Runtime, names []string) ([]*program.URLProgram, error) {
	l, err := config.ParseLock(conf.LockfilePath)
	if err != nil {
		return nil, fmt.Errorf("parsing lockfile: %w", err)
	}

	if len(names) == 0 {
		return l.Programs, nil
	}

	notFound := []string{}
	programs := []*program.URLProgram{}

	for _, name := range names {
		found := false
		for _, p := range l.Programs {
			if p.PName == name {
				programs = append(programs, p)
				found = true
				break
			}
		}
		if !found {
			notFound = append(notFound, name)
		}
	}

	if len(notFound) > 0 {
		internal.Log().Error().Strs("programs", notFound).Msg("undefined programs in lockfile")
		return nil, FailExecError
	}

	return programs, nil
}

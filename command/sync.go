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
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/bindl-dev/bindl/config"
	"github.com/bindl-dev/bindl/internal"
	"github.com/bindl-dev/bindl/program"
	"sigs.k8s.io/yaml"
)

// Sync reads the configuration file (conf.Path) and generates the lockfile (conf.LockfilePath).
// By default, it overwrites the lockfile blindly.
// If `writeToStdout` is true, then it writes to STDOUT and lockfile will not be touched.
func Sync(ctx context.Context, conf *config.Runtime, writeToStdout bool) error {
	c := &config.Config{}
	raw, err := os.ReadFile(conf.Path)
	if err != nil {
		return fmt.Errorf("reading config: %w", err)
	}
	if err := yaml.Unmarshal(raw, c); err != nil {
		return fmt.Errorf("parsing yaml: %w", err)
	}

	parsed := make(chan *program.Lock, 4)
	hasError := false

	var wg sync.WaitGroup

	for _, programConfig := range c.Programs {
		wg.Add(1)
		go func(prog *program.Config) {
			defer wg.Done()

			internal.Log().Info().Str("program", prog.Name).Msg("building program spec")
			p, err := prog.Lock(ctx, c.Platforms, conf.UseCache)
			if err != nil {
				internal.Log().Err(err).Str("program", prog.Name).Msg("parsing configuration")
				hasError = true
				return
			}
			parsed <- p
		}(programConfig)
	}

	go func() {
		wg.Wait()
		close(parsed)
	}()

	programs := []*program.Lock{}
	for p := range parsed {
		internal.Log().Info().Str("program", p.Name).Msg("built program spec")
		programs = append(programs, p)
	}

	if hasError {
		return fmt.Errorf("unsuccessful configuration parsing")
	}

	l := &config.Lock{
		Updated:  time.Now().UTC(),
		Programs: programs,
	}

	data, err := yaml.Marshal(l)
	if err != nil {
		return fmt.Errorf("marshaling yaml: %w", err)
	}
	if writeToStdout {
		_, err = os.Stdout.Write(data)
	} else {
		err = os.WriteFile(conf.LockfilePath, data, 0644)
	}
	return err
}

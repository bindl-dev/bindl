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
	"errors"
	"fmt"
	"sync"

	"go.xargs.dev/bindl/config"
	"go.xargs.dev/bindl/internal"
	"go.xargs.dev/bindl/program"
)

// ErrFailExec is used as generic failure for command line interface as
// preserving the real error can be difficult with concurrent setup.
var ErrFailExec = errors.New("failed to execute command, please troubleshoot logs")

// ProgramCommandFunc is a shorthand to command execution function signature,
// allowing the command to be run concurrently for each program.
type ProgramCommandFunc func(context.Context, *config.Runtime, *program.Lock) error

// IterateLockfilePrograms is an iterator which spawns a goroutine for each
// selected programs. Any subcommand can leverage this by honoring ProgramCommandFunc.
func IterateLockfilePrograms(ctx context.Context, conf *config.Runtime, names []string, fn ProgramCommandFunc) error {
	progs := make(chan *program.Lock, 1)
	errs := make(chan error, 1)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		errs <- filterPrograms(ctx, conf, names, progs)
		wg.Done()
	}()

	for p := range progs {
		wg.Add(1)
		go func(p *program.Lock) {
			errs <- fn(ctx, conf, p)
			wg.Done()
		}(p)
	}

	go func() {
		wg.Wait()
		close(errs)
	}()

	hasError := false
	for err := range errs {
		if err != nil {
			if !errors.Is(err, ErrFailExec) {
				internal.ErrorMsg(err)
			}
			hasError = true
		}
	}

	if hasError {
		return ErrFailExec
	}
	return nil
}

func filterPrograms(ctx context.Context, conf *config.Runtime, names []string, progs chan<- *program.Lock) error {
	defer close(progs)

	l, err := config.ParseLock(conf.LockfilePath)
	if err != nil {
		return fmt.Errorf("parsing lockfile: %w", err)
	}

	// In the event that no specific names were provided, then *all* programs in lockfile
	// will be included in the filter.
	if len(names) == 0 {
		for _, p := range l.Programs {
			if err := ctx.Err(); err != nil {
				return err
			}
			progs <- p
		}
		return nil
	}

	notFound := []string{}

	for _, name := range names {
		if err := ctx.Err(); err != nil {
			return err
		}
		found := false
		for _, p := range l.Programs {
			if p.Name == name {
				progs <- p
				found = true
				break
			}
		}
		if !found {
			internal.ErrorMsg(fmt.Errorf("program not found: %v", name))
			notFound = append(notFound, name)
		}
	}

	// This can probably be done with boolean, but leaving it here for now to
	// assist debugging as needed.
	if len(notFound) > 0 {
		return ErrFailExec
	}

	return nil
}

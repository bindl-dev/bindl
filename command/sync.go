package command

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"go.xargs.dev/bindl/config"
	"go.xargs.dev/bindl/internal"
	"go.xargs.dev/bindl/program"
	"sigs.k8s.io/yaml"
)

func Sync(ctx context.Context, conf *config.Runtime, writeToStdout bool) error {
	c := &config.Config{}
	raw, err := os.ReadFile(conf.Path)
	if err != nil {
		return fmt.Errorf("reading config: %w", err)
	}
	if err := yaml.Unmarshal(raw, c); err != nil {
		return fmt.Errorf("parsing yaml: %w", err)
	}

	parsed := make(chan *program.URLProgram, 4)
	hasError := false

	var wg sync.WaitGroup

	for _, programConfig := range c.Programs {
		wg.Add(1)
		go func(prog *program.Config) {
			defer wg.Done()

			internal.Log().Info().Str("program", prog.PName).Msg("building program spec")
			p, err := prog.URLProgram(ctx, c.Platforms)
			if err != nil {
				internal.Log().Err(err).Str("program", prog.PName).Msg("parsing configuration")
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

	programs := []*program.URLProgram{}
	for p := range parsed {
		internal.Log().Info().Str("program", p.PName).Msg("built program spec")
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
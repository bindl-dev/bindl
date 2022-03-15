package command

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"
	"sync"

	"go.xargs.dev/bindl/config"
	"go.xargs.dev/bindl/download"
	"go.xargs.dev/bindl/internal"
	"go.xargs.dev/bindl/program"
)

func GetAll(ctx context.Context, conf *config.Runtime) error {
	l, err := config.ParseLock(conf.LockfilePath)
	if err != nil {
		return err
	}

	errors := []error{}
	var wg sync.WaitGroup

	for _, prog := range l.Programs {
		if ctx.Err() != nil {
			break
		}
		internal.Log().Debug().Str("program", prog.Name()).Msg("found program spec")
		wg.Add(1)
		go func(prog *program.URLProgram) {
			err := downloadProgram(ctx, prog, conf.OutputDir)
			if err != nil {
				errors = append(errors, fmt.Errorf("downloading %s: %w", prog.Name(), err))
			}
			wg.Done()
		}(prog)
	}

	wg.Wait()

	for _, err := range errors {
		internal.ErrorMsg(err)
	}
	if len(errors) > 0 {
		return FailExecError
	}

	return nil
}

func Get(ctx context.Context, conf *config.Runtime, names ...string) error {
	internal.Log().Debug().Strs("programs", names).Msg("attempting to download programs")
	l, err := config.ParseLock(conf.LockfilePath)
	if err != nil {
		return err
	}

	errors := []error{}
	missing := []string{}
	var wg sync.WaitGroup

	for _, name := range names {
		if ctx.Err() != nil {
			break
		}
		found := false
		for _, prog := range l.Programs {
			if ctx.Err() != nil {
				break
			}
			if name == prog.Name() {
				found = true
				internal.Log().Debug().Str("program", name).Msg("found program spec")
				wg.Add(1)
				go func(name string, prog *program.URLProgram) {
					err := downloadProgram(ctx, prog, conf.OutputDir)
					if err != nil {
						errors = append(errors, fmt.Errorf("downloading %s: %w", name, err))
					}
					wg.Done()
				}(name, prog)
				break
			}
		}
		if !found {
			missing = append(missing, name)
		}
	}

	wg.Wait()
	if len(missing) > 0 {
		internal.Log().Error().Strs("programs", missing).Msg("missing program spec")
	}
	for _, err := range errors {
		internal.ErrorMsg(err)
	}
	if len(missing) > 0 || len(errors) > 0 {
		return FailExecError
	}
	return nil
}

func downloadProgram(ctx context.Context, p *program.URLProgram, outDir string) error {
	a, err := p.DownloadArchive(ctx, &download.HTTP{}, runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return err
	}
	internal.Log().Debug().Str("program", p.Name()).Msg("extracting archive")
	bin, err := a.Extract(p.Name())
	if err != nil {
		return err
	}
	internal.Log().Debug().Str("program", p.Name()).Msg("found binary")

	loc := path.Join(outDir, p.Name())
	err = os.WriteFile(loc, bin, 0755)
	if err != nil {
		return err
	}
	internal.Log().Info().Str("output", loc).Str("program", p.Name()).Msg("downloaded")
	return nil
}

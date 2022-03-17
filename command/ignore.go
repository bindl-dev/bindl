package command

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.xargs.dev/bindl/config"
	"go.xargs.dev/bindl/internal"
)

func getIgnoreEntries(conf *config.Runtime) map[string]bool {
	outdir := conf.OutputDir
	return map[string]bool{
		filepath.Join(outdir, "*"): false,
	}
}

func UpdateIgnoreFile(conf *config.Runtime, path string) error {
	internal.Log().Debug().Str("ignore", path).Msg("attempting to update ignore file")

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("opening file '%s': %w", path, err)
	}
	defer f.Close()

	fc := map[string]bool{}
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		l := strings.Trim(scanner.Text(), "\n")
		fc[l] = true
		internal.Log().Debug().Str("line", l).Msg("parsing ignore file contents")
	}

	desired := getIgnoreEntries(conf)
	for entry, _ := range desired {
		if fc[entry] {
			desired[entry] = true
		}
	}

	for entry, presence := range desired {
		if presence {
			continue
		}

		_, err := f.WriteString(entry + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

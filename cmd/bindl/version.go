package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = ""
	date    = ""
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show build version",
	Run:   printVersion,
}

func printVersion(*cobra.Command, []string) {
	fmt.Printf("version: %s\ncommit: %s\ndate: %s\n", version, commit, date)
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Printf("unable to retrieve build info")
		os.Exit(1)
	}
	if len(bi.Deps) > 0 {
		fmt.Printf("dependencies:\n")
	}
	for _, d := range bi.Deps {
		fmt.Printf("  %s @ %s\n", d.Path, d.Version)
		fmt.Printf("    %s\n", d.Sum)
	}
}

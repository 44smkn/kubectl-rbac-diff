package main

import (
	"fmt"
	"os"

	"github.com/44smkn/kubectl-role-diff/pkg/build"
	"github.com/44smkn/kubectl-role-diff/pkg/cmd"
	"github.com/44smkn/kubectl-role-diff/pkg/cmdutil"
	"github.com/44smkn/kubectl-role-diff/pkg/log"
)

const (
	ExitCodeOK int = 0

	// Errors start at 10
	ExitCodeError = 10 + iota
)

func main() {
	os.Exit(run(os.Args))
}

func run(args []string) int {
	buildDate := build.Date
	buildVersion := build.Version

	logLevel := "info"
	if os.Getenv("DEBUG") == "true" {
		logLevel = "debug"
	}
	logger, err := log.New(logLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate logger: %s", err.Error())
	}
	cmdFactory := cmdutil.NewFactory(buildVersion, logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize process: %s", err.Error())
	}

	rootCmd := cmd.NewCmdDiff(cmdFactory, buildVersion, buildDate)
	if cmd, err := rootCmd.ExecuteC(); err != nil {
		// TODO: error handling
		fmt.Fprintln(os.Stderr, cmd.UsageString())
	}
	return ExitCodeOK
}

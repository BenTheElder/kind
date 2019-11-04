/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package kind implements the root kind cobra command, and the cli Main()
package kind

import (
	"context"
	"io"
	"os"

	"github.com/spf13/cobra"

	"sigs.k8s.io/kind/pkg/cmd"
	"sigs.k8s.io/kind/pkg/cmd/kind/build"
	"sigs.k8s.io/kind/pkg/cmd/kind/completion"
	"sigs.k8s.io/kind/pkg/cmd/kind/create"
	"sigs.k8s.io/kind/pkg/cmd/kind/delete"
	"sigs.k8s.io/kind/pkg/cmd/kind/export"
	"sigs.k8s.io/kind/pkg/cmd/kind/get"
	"sigs.k8s.io/kind/pkg/cmd/kind/load"
	"sigs.k8s.io/kind/pkg/cmd/kind/version"
	"sigs.k8s.io/kind/pkg/errors"
	"sigs.k8s.io/kind/pkg/exec"
	"sigs.k8s.io/kind/pkg/globals"
	"sigs.k8s.io/kind/pkg/log"
)

// Flags for the kind command
type Flags struct {
	LogLevel  string
	Verbosity int32
	Quiet     bool
}

// NewCommand returns a new cobra.Command implementing the root command for kind
func NewCommand(ctx context.Context, streams cmd.IOStreams) *cobra.Command {
	flags := &Flags{}
	command := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "kind",
		Short: "kind is a tool for managing local Kubernetes clusters",
		Long:  "kind creates and manages local Kubernetes clusters using Docker container 'nodes'",
		PersistentPreRunE: func(command *cobra.Command, args []string) error {
			err := runE(flags, command, os.Stderr)
			if err != nil {
				logError(err)
			}
			return err
		},
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       version.Version(),
	}
	command.PersistentFlags().StringVar(
		&flags.LogLevel,
		"loglevel",
		"",
		"DEPRECATED: see -v instead",
	)
	command.PersistentFlags().Int32VarP(
		&flags.Verbosity,
		"verbosity",
		"v",
		0,
		"info log verbosity",
	)
	command.PersistentFlags().BoolVarP(
		&flags.Quiet,
		"quiet",
		"q",
		false,
		"silence all stderr output",
	)
	// add all top level subcommands
	command.AddCommand(build.NewCommand())
	command.AddCommand(completion.NewCommand())
	command.AddCommand(create.NewCommand())
	command.AddCommand(delete.NewCommand())
	command.AddCommand(export.NewCommand())
	command.AddCommand(get.NewCommand())
	command.AddCommand(version.NewCommand())
	command.AddCommand(load.NewCommand())
	return command
}

func runE(flags *Flags, command *cobra.Command, errOut io.Writer) error {
	// handle limited migration for --loglevel
	setLogLevel := command.Flag("loglevel").Changed
	setVerbosity := command.Flag("verbosity").Changed
	if setLogLevel && !setVerbosity {
		switch flags.LogLevel {
		case "debug":
			flags.Verbosity = 3
		case "trace":
			flags.Verbosity = 2147483647
		}
	}
	// normal logger setup
	if flags.Quiet {
		globals.SetLogger(log.NoopLogger{})
	} else {
		globals.UseCLILogger(errOut, log.Level(flags.Verbosity))
	}
	// warn about deprecated flag if used
	if setLogLevel {
		globals.GetLogger().Warn("WARNING: --loglevel is deprecated, please switch to -v and -q!")
	}
	return nil
}

// logError logs the error and the root stacktrace if there is one
func logError(err error) {
	globals.GetLogger().Errorf("ERROR: %v", err)
	// If debugging is enabled (non-zero verbosity), display more info
	if globals.GetLogger().V(1).Enabled() {
		// Display Output if the error was running a command ...
		if err := exec.RunErrorForError(err); err != nil {
			globals.GetLogger().Errorf("\nOutput:\n%s", err.Output)
		}
		// Then display stack trace if any (there should be one...)
		if trace := errors.StackTrace(err); trace != nil {
			globals.GetLogger().Errorf("\nStack Trace: %+v", trace)
		}
	}
}

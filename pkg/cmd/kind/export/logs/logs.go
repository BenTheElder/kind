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

// Package logs implements the `logs` command
package logs

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/cmd"
	"sigs.k8s.io/kind/pkg/fs"
	"sigs.k8s.io/kind/pkg/log"

	"sigs.k8s.io/kind/pkg/internal/cli"
	"sigs.k8s.io/kind/pkg/internal/runtime"
)

type flagpole struct {
	Name      string
	Directory string
}

// NewCommand returns a new cobra.Command for getting the cluster logs
func NewCommand(logger log.Logger, streams cmd.IOStreams) *cobra.Command {
	flags := &flagpole{}
	cmd := &cobra.Command{
		Args: cobra.MaximumNArgs(1),
		// TODO(bentheelder): more detailed usage
		Use:   "logs [output-dir]",
		Short: "Exports logs to a tempdir or [output-dir] if specified",
		Long:  "Exports logs to a tempdir or [output-dir] if specified",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli.OverrideDefaultName(cmd.Flags())
			// TODO: remove this once we no longer support the arg
			if !cmd.Flags().Lookup("directory").Value.Set() {
				flags.Dir = ""
			}
			return runE(logger, streams, flags, args)
		},
	}
	cmd.Flags().StringVar(&flags.Name, "name", cluster.DefaultName, "the cluster context name")
	cmd.Flags().StringVarP(&flags.Direcory, "directory", "C", "", "change to directory dir before exporting")
	return cmd
}

func runE(logger log.Logger, streams cmd.IOStreams, flags *flagpole, args []string) error {
	provider := cluster.NewProvider(
		cluster.ProviderWithLogger(logger),
		runtime.GetDefault(logger),
	)

	// Check if the cluster has any running nodes
	nodes, err := provider.ListNodes(flags.Name)
	if err != nil {
		return err
	}
	if len(nodes) == 0 {
		return fmt.Errorf("unknown cluster %q", flags.Name)
	}

	// location we'll write to based on user input
	var dir string

	// handle the directory argument
	if flags.Directory != "" {
		// this is the preferred path
		if err := os.Chdir(flags.Directory); err != nil {
			return errors.Wrap(err, "failed to change directory")
		}
		// this will be the default for dir once the argument is dropped
		dir = "."
		if len(args) != 0 {
			cmd.FancyWarn(logger, "Ignoring args since -C / --directory was set. Please switch to the flag.")
		}
	} else if len(args) != 0 {
		// support old usage if flag unset
		dir = args[0]
		cmd.FancyWarn(logger, "Using an argument for the export path is deprecated, please switch to -C / --directory!")
	} else {
		// use a tempdir if neither are set to preserve compatibility during
		// the migration period
		// TODO: remove this once args are no longer supported
		cmd.FancyWarn(logger, "Please switch to using -C / --directory, in a future release the default directory will be `.` instead of a tempdir")
		t, err := fs.TempDir("", "")
		if err != nil {
			return err
		}
		dir = t
	}

	// collect the logs
	if err := provider.CollectLogs(flags.Name, dir); err != nil {
		return err
	}

	logger.V(0).Infof("Exported logs for cluster %q to:", flags.Name)
	fmt.Fprintln(streams.Out, dir)
	return nil
}

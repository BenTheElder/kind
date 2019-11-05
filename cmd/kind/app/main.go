/*
Copyright 2019 The Kubernetes Authors.

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

package app

import (
	"os"

	"sigs.k8s.io/kind/pkg/cmd"
	"sigs.k8s.io/kind/pkg/exec"
	"sigs.k8s.io/kind/pkg/errors"
	"sigs.k8s.io/kind/pkg/cmd/kind"
	"sigs.k8s.io/kind/pkg/log"
)

// Main is the kind main(), it will invoke Run(), if an error is returned
// it will then call os.Exit
func Main() {
	if err := Run(cmd.NewLogger(), cmd.StandardIOStreams()); err != nil {
		os.Exit(1)
	}
}

// Run invokes the kind root command, returning the error.
// See: sigs.k8s.io/kind/pkg/cmd/kind
func Run(logger log.Logger, streams cmd.IOStreams) error {
	err := kind.NewCommand(logger, streams).Execute()
	if err != nil {
		logError(logger, err)
	}
	return err
}

// logError logs the error and the root stacktrace if there is one
func logError(logger log.Logger, err error) {
	logger.Errorf("ERROR: %v", err)
	// If debugging is enabled (non-zero verbosity), display more info
	if logger.V(1).Enabled() {
		// Display Output if the error was running a command ...
		if err := exec.RunErrorForError(err); err != nil {
			logger.Errorf("\nOutput:\n%s", err.Output)
		}
		// Then display stack trace if any (there should be one...)
		if trace := errors.StackTrace(err); trace != nil {
			logger.Errorf("\nStack Trace: %+v", trace)
		}
	}
}

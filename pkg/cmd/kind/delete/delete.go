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

// Package delete implements the `delete` command
package delete

import (
	"github.com/spf13/cobra"

	"sigs.k8s.io/kind/pkg/log"
	"sigs.k8s.io/kind/pkg/cmd"
	deletecluster "sigs.k8s.io/kind/pkg/cmd/kind/delete/cluster"
)

// NewCommand returns a new cobra.Command for cluster creation
func NewCommand(logger log.Logger, streams cmd.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Args: cobra.NoArgs,
		// TODO(bentheelder): more detailed usage
		Use:   "delete",
		Short: "Deletes one of [cluster]",
		Long:  "Deletes one of [cluster]",
	}
	cmd.AddCommand(deletecluster.NewCommand(logger, streams))
	return cmd
}

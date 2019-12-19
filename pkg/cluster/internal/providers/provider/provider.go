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

package provider

import (
	"sigs.k8s.io/kind/pkg/cluster/nodes"

	"sigs.k8s.io/kind/pkg/internal/apis/config"
	"sigs.k8s.io/kind/pkg/internal/cli"
)

// Provider represents a provider of cluster / node infrastructure
// This is an alpha-grade internal API
type Provider interface {
	// Provision should create and start the nodes, just short of
	// actually starting up Kubernetes, based on the given cluster config
	Provision(status *cli.Status, cluster string, cfg *config.Cluster) error
	// Deprovision should delete nodes and cleanup resources from the cluster
	Deprovision(cluster string) error
	// ListClusters discovers the clusters that currently have resources
	// under this providers
	ListClusters() ([]string, error)
	// ListNodes returns the nodes under this provider for the given
	// cluster name, they may or may not be running correctly
	ListNodes(cluster string) ([]nodes.Node, error)
	// GetAPIServerEndpoint returns the host endpoint for the cluster's API server
	GetAPIServerEndpoint(cluster string) (string, error)
}

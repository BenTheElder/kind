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

package cmd

import (
	"io"
	"os"
)

// IOStreams provides the standard names for iostreams.
// This is useful for embedding and for unit testing.
// Inconsistent and different names make it hard to read and review code
type IOStreams struct {
	// In think, os.Stdin
	In io.Reader
	// Out think, os.Stdout
	Out io.Writer
	// NOTE: we do NOT have ErrOut
	// This is because we exclusively use sigs.k8s.io/kind/pkg/log for that
}

// StandardIOStreams returns an IOStreams from os.Stdin, os.Stdout
func StandardIOStreams() IOStreams {
	return IOStreams{
		In:  os.Stdin,
		Out: os.Stdout,
	}
}

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

package log

// Writer represents command line output
//
// You can inject your own default Writer using SetDefault() with an
// implementation of this interface.
//
// See also klog.Use() in our klog package under this one
type Logger interface {
	// Printf should be used to write user facing messages
	Printf(format string, args ...interface{})
	// Warnf should be used to write user facing warnings
	Warnf(format string, args ...interface{})
	// V(N > 0) returns a logger at a particular verbosity level,
	// it should be used to write debug messages with detail increasing at
	// each higher level like: V(1).Infof("My Debug Message")
	V(Level) LevelLogger
}

// Level is a leveled logging level, see https://github.com/kubernetes/klog
type Level int32

// LevelLogger is an interface like a subset of klog.Verbose
// see: https://github.com/kubernetes/klog
type LevelLogger interface {
	// Enabled should return true if this verbosity level is enabled
	Enabled() bool
	// Infof should be used to write debug messages
	Infof(format string, args ...interface{})
}

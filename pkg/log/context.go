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

import (
	"context"
)

// contextKey is a private type used for context values
type contextKey struct{}

// loggerKey is the key used to identify the logger in a context.Context
var loggerKey contextKey

// noop is a logger returned when there is not a logger set
var noop Logger = NoopLogger{}

// NewContext returns a new context from parent with the Logger
// This can be retreived with FromContext
func NewContext(parent context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext returns the Logger set in NewContext or else a NoopLogger
func FromContext(ctx context.Context) Logger {
	l, ok := ctx.Value(loggerKey).(Logger)
	if !ok {
		l = noop
	}
	return l
}

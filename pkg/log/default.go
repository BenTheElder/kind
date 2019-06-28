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

import (
	"bytes"
	"fmt"
	"os"
	"io"
	"sync"
)

var std Logger = &DefaultLogger{
	writer: os.Stderr,
}

// SetDefault sets the default logger
func SetDefault(logger Logger) {
	std = logger
}

// Printf prints a user facing message
// This uses the default logger (See SetDefault)
func Printf(format string, args ...interface{}) {
	std.Printf(format, args...)
}

// Warnf prints a user facing warning message
// This uses the default logger (See SetDefault)
func Warnf(format string, args ...interface{}) {
	std.Warnf(format, args...)
}

// V returns a LevelLogger for the given verbosity level
// This uses the default logger (See SetDefault)
func V(level Level) LevelLogger {
	return std.V(level)
}

/* Default Implementation */

// DefaultLogger is the default Logger implementation
type DefaultLogger struct {
	mu           sync.Mutex
	writer       io.Writer
	enabledLevel Level
}

var _ Logger = &DefaultLogger{}

func (dl *DefaultLogger) printfln(format string, args ...interface{}) {
	dl.mu.Lock()
	defer dl.mu.Unlock()
	fprintfln(dl.writer, format, args...)
}

func (dl *DefaultLogger) enabled(level Level) bool {
	dl.mu.Lock()
	defer dl.mu.Unlock()
	return level >= dl.enabledLevel
}

func (dl *DefaultLogger) Printf(format string, args ...interface{}) {
	dl.printfln(format, args...)
}

func (dl *DefaultLogger) Warnf(format string, args ...interface{}) {
	dl.printfln("WARNING: "+format, args...)
}

func (dl *DefaultLogger) V(level Level) LevelLogger {
	dl.mu.Lock()
	defer dl.mu.Unlock()
	return &defaultLevelLogger{
		logger:  dl,
		level: level,
	}
}

type defaultLevelLogger struct {
	logger  *DefaultLogger
	level Level
}

var _ LevelLogger = &defaultLevelLogger{}

func (dll *defaultLevelLogger) Enabled() bool {
	return dll.logger.enabled(dll.level)
}

func (dll *defaultLevelLogger) Infof(format string, args ...interface{}) {
	if dll.Enabled() {
		dll.logger.printfln(fmt.Sprintf("INFO[%d]: ", dll.level)+format, args...)
	}
}

func fprintfln(writer io.Writer, format string, args ...interface{}) {
	var buff bytes.Buffer
	fmt.Fprintf(&buff, format, args...)
	if buff.Bytes()[buff.Len()-1] != '\n' {
		buff.WriteByte('\n')
	}
	writer.Write(
		buff.Bytes(),
	)
}

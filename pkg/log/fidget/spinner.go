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

// Package fidget implements CLI functionality for bored users waiting for results
package fidget

import (
	"fmt"
	"io"
	"time"
)

// custom CLI loading spinner for kind
var spinnerFrames = []string{
	"⠈⠁",
	"⠈⠑",
	"⠈⠱",
	"⠈⡱",
	"⢀⡱",
	"⢄⡱",
	"⢄⡱",
	"⢆⡱",
	"⢎⡱",
	"⢎⡰",
	"⢎⡠",
	"⢎⡀",
	"⢎⠁",
	"⠎⠁",
	"⠊⠁",
}

// Spinner is a simple and efficient CLI loading spinner used by kind
// It is simplistic and assumes that the line length will not change.
// It is best used indirectly via log.Status (see parent package)
type Spinner struct {
	stop    chan struct{}
	stopped chan struct{}
	ticker  *time.Ticker
	writer  io.Writer
}

// NewSpinner initializes and returns a new Spinner that will write to
func NewSpinner(w io.Writer) *Spinner {
	return &Spinner{
		stop:    make(chan struct{}),
		stopped: make(chan struct{}),
		ticker:  time.NewTicker(time.Millisecond * 80),
		writer:  w,
	}
}

// Start starts the spinner running
func (s *Spinner) Start() {
	go func() {
		fmt.Fprintf(s.writer, "  ")
		for {
			for _, frame := range spinnerFrames {
				select {
				case <-s.stop:
					fmt.Fprintf(s.writer, "\b\b  \b\b")
					s.stopped <- struct{}{}
					return
				case <-s.ticker.C:
					fmt.Fprintf(s.writer, "\b\b%s", frame)
				}
			}
		}
	}()
}

// Stop signals the spinner to stop
func (s *Spinner) Stop() {
	s.stop <- struct{}{}
	<-s.stopped
}

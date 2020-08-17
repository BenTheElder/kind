/*
Copyright 2020 The Kubernetes Authors.

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

package docker

import (
	"encoding/json"
	"strings"
	"time"

	"sigs.k8s.io/kind/pkg/errors"
	"sigs.k8s.io/kind/pkg/exec"
)

func removeVolume(name string) error {
	return exec.Command("docker", "volume", "rm", name).Run()
}

func createVolume(name string, labels []string) (string, error) {
	b, err := exec.Output(exec.Command("docker", "volume", "create", "--label", strings.Join(labels, ","), name))
	return string(b), err
}

type volumeInfo struct {
	CreatedAt time.Time         `json:"CreatedAt"`
	Labels    map[string]string `json:"Labels"`
	// fields we do not currently need are elided for now ...
}

func inspectVolume(name string) (*volumeInfo, error) {
	// TODO wrap errors
	b, err := exec.Output(exec.Command("docker", "volume", "inspect", name))
	if err != nil {
		return nil, err
	}
	var volumes []volumeInfo
	if err := json.Unmarshal(b, &volumes); err != nil {
		return nil, err
	}
	if len(volumes) == 0 {
		return nil, nil
	}
	if len(volumes) > 1 {
		return nil, errors.Errorf("received %d volumes for inspect volume %q, expected one", len(volumes), name)
	}
	return &volumes[0], nil
}

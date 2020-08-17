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
	"time"

	"github.com/google/uuid"
)

// volumeLock is similar to a file lock but based on a docker volume
//
// This allows us to centralize the lock state to the dockerd for protecting
// resources from concurrent mutation even if the clients are on different hosts
//
// This is safe across processes / hosts as long as there is not drastic clock skew.
//
// Additionally, the lock respects a timeout on
//
// Design:
// https://docs.google.com/document/d/1Q7Njyco2mAz66lS44pVV7ixT22RAkqBrmVMetG1zuT4/edit#
// (This doc is shared with members of kubernetes-dev@googlegroups.com and
// kubernetes-sig-testing@googlegroups.com. Join one for access.
// individual requests will be ignored.)
type volumeLock struct {
	name string
	uuid string
	// runtime state
	locked bool
}

func newVolumeLock(name string) *volumeLock {
	textID, _ := uuid.Must(uuid.NewRandom()).MarshalText()
	return &volumeLock{
		name: name,
		uuid: string(textID),
	}
}

func (l *volumeLock) Unlock() error {
	if !l.locked {
		return nil
	}
	if err := l.release(); err != nil {
		return err
	}
	l.locked = false
	return nil
}

func (l *volumeLock) release() error {
	return removeVolume(l.name)
}

func (l *volumeLock) TryLock() (bool, error) {
	// first opportunistically release the lock if expired
	if expired, _ := l.checkLockExpired(); expired {
		_ = l.release()
	}
	// try to acquire it
	if err := l.createVolume(); err != nil {
		return false, err
	}
	// check if we won the race acquiring
	l.locked, _ = l.checkHaveLock()
	return l.locked, nil
}

func (l *volumeLock) createVolume() error {
	_, err := createVolume(l.name, []string{lockUUIDLabel + "=" + l.uuid})
	return err
}

var lockTimeout = time.Second

func (l *volumeLock) checkLockExpired() (bool, error) {
	info, err := inspectVolume(l.name)
	if err != nil {
		return false, err
	}
	if info == nil {
		return false, nil
	}
	return time.Now().After(info.CreatedAt.Add(lockTimeout)), nil
}

func (l *volumeLock) checkHaveLock() (bool, error) {
	info, err := inspectVolume(l.name)
	if err != nil {
		return false, err
	}
	if info == nil {
		return false, nil
	}
	return info.Labels[lockUUIDLabel] == l.uuid, nil
}

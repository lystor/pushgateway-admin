// Copyright 2025 Mykola Ulianytskyi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pushgateway

import (
	"fmt"
	"strings"
)

type ErrInvalidGroupingKey struct {
	err string
}

func (e ErrInvalidGroupingKey) Error() string {
	return fmt.Sprintf("invalid grouping key (%s)", e.err)
}

type GroupingKey struct {
	*Labels
}

func (f *GroupingKey) Match(l *Labels) bool {
	if f == nil || len(f.m) == 0 {
		return true
	}

	for _, name := range f.Names() {
		x := f.Get(name)
		y, ok := l.GetOk(name)
		if !ok || x != y {
			return false
		}
	}

	return true
}

func (f *GroupingKey) URL() string {
	url := ""
	for _, name := range f.Names() {
		url += "/" + name + "/" + f.Get(name)
	}
	return url
}

func (f *GroupingKey) Validate() error {
	if len(f.m) == 0 {
		return ErrInvalidGroupingKey{err: "no labels"}
	}

	for name, value := range f.m {
		if name == "" {
			return ErrInvalidGroupingKey{err: "empty label name"}
		}
		if value == "" {
			return ErrInvalidGroupingKey{err: "empty label value"}
		}
	}

	if _, ok := f.m["job"]; !ok {
		return ErrInvalidGroupingKey{err: "job is absent"}
	}

	return nil
}

func NewGroupingKeyFromMap(labels *Labels) (*GroupingKey, error) {
	f := &GroupingKey{
		Labels: NewLabels(),
	}

	for name, value := range labels.AsMap() {
		f.Set(name, value)
	}

	if err := f.Validate(); err != nil {
		return nil, err
	}

	return f, nil
}

func NewGroupingKeyFromSlice(labels []string) (*GroupingKey, error) {
	f := &GroupingKey{
		Labels: NewLabels(),
	}

	for _, l := range labels {
		name, value, ok := strings.Cut(l, "=")
		if !ok {
			return nil, ErrInvalidGroupingKey{err: "invalid format"}
		}
		f.Set(name, value)
	}

	if err := f.Validate(); err != nil {
		return nil, err
	}

	return f, nil
}

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
	"encoding/json"
	"maps"
	"sort"
	"strings"
)

type Labels struct {
	m map[string]string
}

func (l *Labels) AsMap() map[string]string {
	return maps.Clone(l.m)
}

func (l *Labels) Get(name string) string {
	return l.m[name]
}

func (l *Labels) GetOk(name string) (string, bool) {
	v, ok := l.m[name]
	return v, ok
}

func (l *Labels) Names() []string {
	a := make([]string, 1, len(l.m))
	a[0] = "job"
	for name := range l.m {
		if name != "job" {
			a = append(a, name)
		}
	}
	sort.Strings(a[1:])
	return a
}

func (l *Labels) Set(name, value string) {
	l.m[name] = value
}

func (l Labels) String() string {
	s := strings.Builder{}
	s.WriteString("{")
	for i, name := range l.Names() {
		if i != 0 {
			s.WriteString(",")
		}
		s.WriteString(name + "=" + l.Get(name))
	}
	s.WriteString("}")
	return s.String()
}

func (l *Labels) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &l.m)
}

func NewLabels() *Labels {
	return &Labels{
		m: make(map[string]string),
	}
}

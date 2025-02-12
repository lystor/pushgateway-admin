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
	"time"
)

// var ignoredMetricNames = map[string]struct{}{
// 	"labels":                    {},
// 	"last_push_successful":      {},
// 	"push_failure_time_seconds": {},
// 	"push_time_seconds":         {},
// }

type MetricGroup struct {
	Labels Labels
	// MetricFamilies         map[string]MetricFamily
	PushFailureTimeSeconds MetricFamily
	PushTimeSeconds        MetricFamily
}

func (g *MetricGroup) UnmarshalJSON(data []byte) error {
	m := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if err := json.Unmarshal(m["labels"], &g.Labels); err != nil {
		return err
	}
	if err := json.Unmarshal(m["push_failure_time_seconds"], &g.PushFailureTimeSeconds); err != nil {
		return err
	}
	if err := json.Unmarshal(m["push_time_seconds"], &g.PushTimeSeconds); err != nil {
		return err
	}

	// g.MetricFamilies = make(map[string]MetricFamily, len(m))
	// for k, v := range m {
	// 	_, ok := ignoredMetricNames[k]
	// 	if ok {
	// 		continue
	// 	}

	// 	mf := MetricFamily{}
	// 	if err := json.Unmarshal(v, &mf); err != nil {
	// 		return err
	// 	}
	// 	g.MetricFamilies[k] = mf
	// }

	return nil
}

type MetricFamily struct {
	// Help      string    `json:"help"`
	Metrics   []Metric  `json:"metrics"`
	Timestamp time.Time `json:"time_stamp"`
	// Type      string    `json:"type"`
}

type Metric struct {
	Labels Labels `json:"labels"`
	// Value  string       `json:"value"`
}

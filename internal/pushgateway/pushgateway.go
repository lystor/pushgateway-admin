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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type PushgatewayResponse struct {
	Body       *bytes.Buffer
	Duration   time.Duration
	StatusCode int
}

func (r *PushgatewayResponse) PrintBody(pretty bool) error {
	if !pretty {
		log.Print(r.Body)
		return nil
	}

	var data any
	if err := json.Unmarshal(r.Body.Bytes(), &data); err != nil {
		return err
	}

	buf, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	log.Print(string(buf))
	return nil
}

func (r *PushgatewayResponse) PrintStatusCode() {
	log.Printf("Response: status_code=%d duration=%s\n", r.StatusCode, r.Duration)
}

type PushgatewayError PushgatewayResponse

func (e PushgatewayError) Error() string {
	return fmt.Sprintf("Pushgateway status_code=%d body='%v'", e.StatusCode, strings.TrimSpace(e.Body.String()))
}

type Pushgateway struct {
	Timeout uint
	Url     string
}

func (gw *Pushgateway) Request(method, url string, body io.Reader, verbose bool) (*PushgatewayResponse, error) {
	if verbose {
		log.Printf("Request:  %s %s\n", method, url)
	}

	client := &http.Client{
		Timeout: time.Duration(gw.Timeout) * time.Second,
	}

	req, err := http.NewRequest(method, gw.Url + url, body)
	if err != nil {
		return nil, err
	}

	start := time.Now()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	duration := time.Since(start)

	respBody := &bytes.Buffer{}
	if _, err = io.Copy(respBody, resp.Body); err != nil {
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		return nil, PushgatewayError{Body: respBody, StatusCode: resp.StatusCode}
	}

	return &PushgatewayResponse{Body: respBody, Duration: duration, StatusCode: resp.StatusCode}, nil
}

func (gw *Pushgateway) Delete(key *GroupingKey) (*PushgatewayResponse, error) {
	url := "/metrics" + key.URL()
	return gw.Request("DELETE", url, nil, true)
}

func (gw *Pushgateway) Metrics() (*PushgatewayResponse, error) {
	return gw.Request("GET", "/api/v1/metrics", nil, false)
}

func (gw *Pushgateway) MetricGroups() ([]MetricGroup, error) {
	resp, err := gw.Request("GET", "/api/v1/metrics", nil, false)
	if err != nil {
		return nil, err
	}

	var data struct {
		Data []MetricGroup `json:"data"`
	}

	if err = json.Unmarshal(resp.Body.Bytes(), &data); err != nil {
		return nil, err
	}

	return data.Data, err
}

func (gw *Pushgateway) Post(key *GroupingKey, metrics io.Reader) (*PushgatewayResponse, error) {
	url := "/metrics" + key.URL()
	return gw.Request("POST", url, metrics, true)
}

func (gw *Pushgateway) Put(key *GroupingKey, metrics io.Reader) (*PushgatewayResponse, error) {
	url := "/metrics" + key.URL()
	return gw.Request("PUT", url, metrics, true)
}

func (gw *Pushgateway) Status() (*PushgatewayResponse, error) {
	return gw.Request("GET", "/api/v1/status", nil, false)
}

func (gw *Pushgateway) Wipe() (*PushgatewayResponse, error) {
	return gw.Request("PUT", "/api/v1/admin/wipe", nil, true)
}

func New(url string, timeout uint) *Pushgateway {
	url = strings.TrimRight(url, "/")

	return &Pushgateway{
		Timeout: timeout,
		Url:     url,
	}
}

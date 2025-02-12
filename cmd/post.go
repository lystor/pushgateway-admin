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

package cmd

import (
	"io"
	"os"
	"strings"

	"github.com/lystor/pushgateway-admin/internal/pushgateway"
	"github.com/spf13/cobra"
)

func newPostCmd(rootOpts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "post",
		Short: "Send metrics to group via POST method",
		Long: "POST works exactly like the PUT " +
			"but only metrics with the same name as the newly pushed metrics are replaced.",
	}

	opts := struct {
		group       []string
		metric      string
		metric_file string
	}{}

	cmd.Flags().StringSliceVarP(&opts.group, "group", "g", []string{}, "Grouping key (label=value,...)")
	cmd.Flags().StringVarP(&opts.metric, "metric", "m", "", "Metric")
	cmd.Flags().StringVarP(&opts.metric_file, "metric-file", "f", "", "Read metrics from file")
	_ = cmd.MarkFlagRequired("group")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		key, err := pushgateway.NewGroupingKeyFromSlice(opts.group)
		if err != nil {
			return err
		}

		var mr io.Reader
		switch {
		case opts.metric != "":
			mr = strings.NewReader(opts.metric + "\n")
		case opts.metric_file != "":
			f, err := os.Open(opts.metric_file)
			if err != nil {
				return err
			}
			defer f.Close()
			mr = f
		default:
			mr = os.Stdin
		}

		gw := newPushgateway(rootOpts)

		resp, err := gw.Post(key, mr)
		if err != nil {
			return err
		}

		resp.PrintStatusCode()
		return nil
	}

	return cmd
}

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
	"github.com/spf13/cobra"
)

func newMetricsCmd(rootOpts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metrics",
		Short: "Returns all metrics",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		gw := newPushgateway(rootOpts)

		resp, err := gw.Metrics()
		if err != nil {
			return err
		}

		if err := resp.PrintBody(rootOpts.pretty); err != nil {
			return err
		}

		return nil
	}

	return cmd
}

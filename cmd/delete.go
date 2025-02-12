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
	"github.com/lystor/pushgateway-admin/internal/pushgateway"
	"github.com/spf13/cobra"
)

func newDeleteCmd(rootOpts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete all metrics in group",
	}

	opts := struct {
		group []string
	}{}

	cmd.Flags().StringSliceVarP(&opts.group, "group", "g", []string{}, "Grouping key (label=value,...)")
	_ = cmd.MarkFlagRequired("group")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		key, err := pushgateway.NewGroupingKeyFromSlice(opts.group)
		if err != nil {
			return err
		}

		gw := newPushgateway(rootOpts)

		resp, err := gw.Delete(key)
		if err != nil {
			return err
		}

		resp.PrintStatusCode()
		return nil
	}

	return cmd
}

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
	"log"
	"time"

	"github.com/lystor/pushgateway-admin/internal/pushgateway"
	"github.com/spf13/cobra"
)

type cleanOptions struct {
	dryrun bool
	group  []string
	sleep  uint
	stale  uint
}

func newCleanCmd(rootOpts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Clean stale metrics",
	}

	opts := cleanOptions{}
	cmd.Flags().BoolVarP(&opts.dryrun, "dry-run", "n", false, "Enable dry-run mode")
	cmd.Flags().StringSliceVarP(&opts.group, "group", "g", []string{}, "Grouping key (label=value,...)")
	cmd.Flags().UintVarP(&opts.sleep, "sleep", "s", 0, "Sleep before exit (sec)")
	cmd.Flags().UintVarP(&opts.stale, "stale", "l", 300, "Stale metrics threshold (sec)")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		var err error
		var key *pushgateway.GroupingKey
		if len(opts.group) > 0 {
			key, err = pushgateway.NewGroupingKeyFromSlice(opts.group)
			if err != nil {
				return err
			}
		}

		gw := newPushgateway(rootOpts)

		groups, err := gw.MetricGroups()
		if err != nil {
			return err
		}

		for _, group := range groups {
			if err := cleanMetricGroup(&opts, key, gw, &group); err != nil {
				return err
			}
		}

		if opts.sleep > 0 {
			log.Printf("Sleeping %d sec\n", opts.sleep)
			time.Sleep(time.Duration(opts.sleep) * time.Second)
		}

		return nil
	}

	return cmd
}

func cleanMetricGroup(opts *cleanOptions, key *pushgateway.GroupingKey, gw *pushgateway.Pushgateway, group *pushgateway.MetricGroup) error {
	const (
		cleanActionDelete       = "delete"
		cleanActionDeleteDryRun = "delete(dry-run)"
		cleanActionSkipAlive    = "skip(alive)"
		cleanActionSkipFilter   = "skip(filter)"
	)

	action := cleanActionSkipAlive

	lastseen := int(time.Now().UTC().Sub(group.PushTimeSeconds.Timestamp).Seconds())
	if lastseen > int(opts.stale) {
		action = cleanActionDelete
	}

	if opts.dryrun && action == cleanActionDelete {
		action = cleanActionDeleteDryRun
	}

	if !key.Match(&group.Labels) {
		action = cleanActionSkipFilter
	}

	log.Printf("MetricGroup: %s action=%s last_seen=%dsec ts=%s\n", group.Labels, action, lastseen, group.PushTimeSeconds.Timestamp.Format(time.DateTime))

	if action == cleanActionDelete {
		key, err := pushgateway.NewGroupingKeyFromMap(&group.Labels)
		if err != nil {
			return err
		}

		resp, err := gw.Delete(key)
		if err != nil {
			return err
		}

		resp.PrintStatusCode()
	}

	return nil
}

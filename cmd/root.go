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
	"os"

	"github.com/spf13/cobra"
)

type rootOptions struct {
	pretty  bool
	timeout uint
	url     string
}

func Execute() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	cmd := newRootCmd()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Short: "Prometheus Pushgateway admin tool",
		Use:   "pushgateway-admin",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		SilenceUsage: true,
	}

	opts := &rootOptions{}
	cmd.PersistentFlags().BoolVarP(&opts.pretty, "pretty", "p", false, "Enable pretty print")
	cmd.PersistentFlags().UintVarP(&opts.timeout, "timeout", "t", 5, "Pushgateway timeout")
	cmd.PersistentFlags().StringVarP(&opts.url, "url", "u", "", "Pushgateway URL")
	_ = cmd.MarkPersistentFlagRequired("url")

	cmd.AddCommand(newCleanCmd(opts))
	cmd.AddCommand(newDeleteCmd(opts))
	cmd.AddCommand(newMetricsCmd(opts))
	cmd.AddCommand(newPostCmd(opts))
	cmd.AddCommand(newPutCmd(opts))
	cmd.AddCommand(newStatusCmd(opts))
	cmd.AddCommand(newWipeCmd(opts))
	return cmd
}

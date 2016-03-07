// Copyright (c) 2016 Pulcy.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/spf13/cobra"

	"github.com/pulcy/pulsar/docker"
)

var (
	pushDockerRegistry string
	pushCmd            = &cobra.Command{
		Use:   "push",
		Short: "Push an image to the default registry",
		Long:  "Push an image to the default registry",
		Run:   runPush,
	}
)

func init() {
	pushCmd.Flags().StringVarP(&pushDockerRegistry, "registry", "r", defaultDockerRegistry(), "Specify docker registry")
	mainCmd.AddCommand(pushCmd)
}

func runPush(cmd *cobra.Command, args []string) {
	switch len(args) {
	case 0:
		CommandError(cmd, "Too few arguments\n")
	case 1:
		err := docker.Push(log, args[0], pushDockerRegistry)
		if err != nil {
			Quitf("%s\n", err)
		}
	default:
		CommandError(cmd, "Too many arguments\n")
	}
}
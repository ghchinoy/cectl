// Copyright © 2017 G. Hussain Chinoy <ghchinoy@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version string

// SetVersion sets the version for use in the cmd package
func SetVersion(v string) {
	version = v
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version of cectl",
	Long:  `List the version cectl`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("cectl version %s\n", version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

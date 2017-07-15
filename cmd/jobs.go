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
	"github.com/spf13/cobra"
)

// jobsCmd represents the jobs command
var jobsCmd = &cobra.Command{
	Use:   "jobs",
	Short: "Manage jobs on the platform",
	Long:  `Manage jobs on the platform`,
}

func init() {
	RootCmd.AddCommand(jobsCmd)

	jobsCmd.PersistentFlags().StringVar(&profile, "profile", "default", "profile name")
	jobsCmd.PersistentFlags().BoolVarP(&outputJSON, "json", "j", false, "output as json")
	jobsCmd.PersistentFlags().BoolVarP(&showCurl, "curl", "c", false, "show curl command")

}

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

// hubsCmd represents the hubs command
var hubsCmd = &cobra.Command{
	Use:   "hubs",
	Short: "Hub management",
	Long:  `List details about existing hubs`,
}

func init() {
	RootCmd.AddCommand(hubsCmd)

	hubsCmd.PersistentFlags().StringVar(&profile, "profile", "default", "profile name")
	hubsCmd.PersistentFlags().BoolVarP(&outputJSON, "json", "j", false, "output as json")
	hubsCmd.PersistentFlags().BoolVarP(&showCurl, "curl", "c", false, "show curl command")

}

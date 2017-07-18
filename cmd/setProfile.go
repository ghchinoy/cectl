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
	"sort"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// setProfileCmd represents the setProfile command
var setProfileCmd = &cobra.Command{
	Use:   "set <profile>",
	Short: "sets a profile to be the default profile",
	Long:  `Sets given profile name as the default profile`,
	Run: func(cmd *cobra.Command, args []string) {
		// check if specified profile exists
		// if not, end
		// copy profile data to default, update label
		// indicate change has occurred

		// check for args, if arg, then list the details for that particular profile
		if len(args) > 0 {
			profile = args[0]
		}
		fmt.Printf("%7s: %s\n", "profile", profile)
		if viper.IsSet(profile) {
			p := viper.GetStringMap(profile)
			for k, v := range p {
				if k == "base" {
					fmt.Printf("%7s: %s\n", k, v)
				}
				if k == "label" {
					fmt.Printf("%7s: %s\n", k, v)
				}
			}
		} else {
			fmt.Printf("No %s profile exists in config file %s.", profile, cfgFile)
		}
		fmt.Println()
		settings := viper.AllSettings()
		var profiles []string
		for k := range settings {
			profiles = append(profiles, k)
		}
		sort.Strings(profiles)
		posn := sort.SearchStrings(profiles, "profile")
		profiles = append(profiles[:posn], profiles[posn+1:]...)
	},
}

func init() {
	profilesCmd.AddCommand(setProfileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setProfileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setProfileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

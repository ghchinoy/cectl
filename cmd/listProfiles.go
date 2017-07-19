// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listProfilesCmd represents the listProfiles command
var listProfilesCmd = &cobra.Command{
	Use:   "list",
	Short: "lists available profiles",
	Long:  `Lists available profiles`,
	Run: func(cmd *cobra.Command, args []string) {
		// check for args, if arg, then list the details for that particular profile
		if len(args) > 0 {
			profile = args[0]
		}
		fmt.Printf("%7s: %s", "profile", profile)
		if viper.IsSet(profile) {
			p := viper.GetStringMap(profile)

			if val, ok := p["label"]; ok {
				fmt.Printf(" (%s)\n", val)
			} else {
				fmt.Println()
			}

			for k, v := range p {
				if k == "base" {
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
		fmt.Println("Valid profiles:", strings.Join(profiles, ", "))
	},
}

func init() {
	profilesCmd.AddCommand(listProfilesCmd)

}

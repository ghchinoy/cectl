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
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// profilesCmd represents the profile command
var profilesCmd = &cobra.Command{
	Use:   "profiles",
	Short: "Manage profiles",
	Long:  `Add, remove, list profiles to manage Cloud Elements access`,
	Run: func(cmd *cobra.Command, args []string) {
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
	RootCmd.AddCommand(profilesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// profileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// profileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	profilesCmd.PersistentFlags().StringVar(&cfgFile, "config", "", cfgHelp)
	profilesCmd.PersistentFlags().StringVar(&profile, "profile", "default", "profile name")
	// Set bash-completion
	validConfigFilenames := []string{"toml", ""}
	profilesCmd.PersistentFlags().SetAnnotation("config", cobra.BashCompFilenameExt, validConfigFilenames)

}

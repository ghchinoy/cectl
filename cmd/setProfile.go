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
	"os"
	"sort"
	"strings"

	toml "github.com/pelletier/go-toml"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// setProfileCmd represents the setProfile command
var setProfileCmd = &cobra.Command{
	Use:   "set <profile>",
	Short: "sets a profile to be the default profile",
	Long:  `Sets given profile name as the default profile`,
	Run: func(cmd *cobra.Command, args []string) {

		// check for args, if arg, then list the details for that particular profile
		if len(args) > 0 {
			profile = args[0]
		}
		fmt.Printf("%7s: %s\n", "profile", profile)
		// check if specified profile exists
		if viper.IsSet(profile) {
			p := viper.GetStringMap(profile)

			// copy profile data to default, update label
			// indicate change has occurred

			viper.Set("default.base", p["base"])
			viper.Set("default.org", p["org"])
			viper.Set("default.user", p["user"])
			viper.Set("default.label", profile)

			// writing back has a PR to make this more formal: https://github.com/spf13/viper/pull/287
			err := writeConfigFile(true)
			if err != nil {
				fmt.Println("Unable to write config file", err.Error())
				fmt.Printf("Config file %s unchanged.\n", viper.ConfigFileUsed())
			}
			fmt.Printf("Default profile set to %s\n", profile)

			for k, v := range p {
				if k == "base" {
					fmt.Printf("%7s: %s\n", k, v)
				}
				if k == "label" {
					fmt.Printf("%7s: %s\n", k, v)
				}
			}
		} else { // if not, end
			fmt.Printf("No %s profile exists in config file %s.\n", profile, cfgFile)
			fmt.Printf("Cannot set %s as default profile.\n", profile)
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
		}
	},
}

func writeConfigFile(force bool) error {
	filename := viper.ConfigFileUsed()
	//fmt.Println(configFile)
	//fmt.Printf("%+v\n", viper.AllSettings())
	t, err := toml.TreeFromMap(viper.AllSettings())
	if err != nil {
		return err
	}
	//fmt.Printf("%s", t)
	s := t.String()

	var flags int
	if force == true {
		flags = os.O_CREATE | os.O_TRUNC | os.O_WRONLY
	} else {
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			flags = os.O_WRONLY
		} else {
			return fmt.Errorf("file: %s exists - cannot overwrite, use force option", filename)
		}
	}

	var AppFs = afero.NewOsFs()
	f, err := AppFs.OpenFile(filename, flags, os.FileMode(0644))
	if err != nil {
		return err
	}

	_, err = f.WriteString(s)

	return nil
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

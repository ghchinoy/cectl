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
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// addProfileCmd represents the addProfile command
var addProfileCmd = &cobra.Command{
	Use:   "add <profile>",
	Short: "add a new profile",
	Long:  `Adds a new profile to the available profiles. Provide a name to get started.`,
	Run: func(cmd *cobra.Command, args []string) {
		// check for args, if arg, then list the details for that particular profile
		if len(args) > 0 {
			profile = args[0]
		} else {
			fmt.Println("please provide a profile name to add, profile add <name>")
			os.Exit(1)
		}
		fmt.Printf("%7s: %s\n", "profile", profile)
		if viper.IsSet(profile) {
			fmt.Printf("Profile %s exists.\n", profile)
			os.Exit(1)
		}
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("base URI: ")
		base, _ := reader.ReadString('\n')
		base = strings.Replace(base, "\n", "", -1)
		fmt.Print("user token: ")
		user, _ := reader.ReadString('\n')
		user = strings.Replace(user, "\n", "", -1)
		fmt.Print("org token: ")
		org, _ := reader.ReadString('\n')
		org = strings.Replace(org, "\n", "", -1)

		viper.Set(profile+".base", base)
		viper.Set(profile+".org", org)
		viper.Set(profile+".user", user)

		err := writeConfigFile(true)
		if err != nil {
			fmt.Println("Unable to write config file", err.Error())
			fmt.Printf("Config file %s unchanged.\n", viper.ConfigFileUsed())
		}
		fmt.Printf("Added profile %s\n", profile)
	},
}

func init() {
	profilesCmd.AddCommand(addProfileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addProfileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addProfileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

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

	"github.com/ghchinoy/ce-go/ce"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	profileSource, profileTarget string
)

// moleculesCmd is the top level command for managing integration assets
var moleculesCmd = &cobra.Command{
	Use:    "molecules",
	Short:  "Manage integration molecules from the Platform",
	Hidden: true,
	Long:   `Manage the integration assets of the Platform`,
}

// exportCmd is the command to export assets
var exportCmd = &cobra.Command{
	Use:   "export [formulas|resources|all (default)]",
	Short: "exports assets from the platform",
	Long:  "Exports a set of assets",
	Run: func(cmd *cobra.Command, args []string) {

		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		scope := "all"
		if len(args) > 0 {
			// args[0] should be either "formulas" | "resources"
			scope = args[0]
		}

		err = ce.ExportAllFormulasToDir(profilemap["base"], profilemap["auth"], "./formulas")
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		err = ce.ExportAllResourcesToDir(profilemap["base"], profilemap["auth"], "./resources")
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Println("Scope", scope)
	},
}

// cloneCmd is the command to clone assets between accounts
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "clone assets from one profile to another",
	Long:  "Clone exports assets from one account profile (--from) and imports them into another profile (--to)",
	Run: func(cmd *cobra.Command, args []string) {

		/*
			// check for profile
			profilemap, err := getAuth(profile)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		*/

		if viper.IsSet(profileTarget) && viper.IsSet(profileSource) {
			fmt.Printf("Exporting from profile '%s' into profile '%s'\n", profileSource, profileTarget)
		} else {
			if !viper.IsSet(profileSource) {
				fmt.Println("Cannot find profile named:", profileSource)
			}
			if !viper.IsSet(profileTarget) {
				fmt.Println("Cannot find profile named:", profileTarget)
			}
		}

	},
}

func init() {
	RootCmd.AddCommand(moleculesCmd)

	moleculesCmd.AddCommand(exportCmd)
	moleculesCmd.AddCommand(cloneCmd)
	cloneCmd.PersistentFlags().StringVar(&profileSource, "from", "default", "source profile name")
	cloneCmd.PersistentFlags().StringVar(&profileTarget, "to", "", "target profile name")
	cloneCmd.MarkPersistentFlagRequired("to")
}

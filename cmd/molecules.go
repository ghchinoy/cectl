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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

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

		scope := []string{"formulas", "resources"}
		if len(args) > 0 {
			// args[0] should be either "formulas" | "resources"
			if args[0] == "formulas" {
				scope = []string{"formulas"}
			}
			if args[0] == "resources" {
				scope = []string{"resources"}
			}
		}

		for _, v := range scope {
			if v == "formulas" {
				err = ExportAllFormulasToDir(profilemap["base"], profilemap["auth"], "./formulas")
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			}
			if v == "resources" {
				err = ExportAllResourcesToDir(profilemap["base"], profilemap["auth"], "./resources")
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			}
		}

	},
}

// ExportAllFormulasToDir creates a directory given and exports all Formula JSON files
func ExportAllFormulasToDir(base, auth string, dirname string) error {
	formulaListByes, _, _, err := ce.FormulasList(base, auth)
	if err != nil {
		return err
	}
	var formulas []Formula
	err = json.Unmarshal(formulaListByes, &formulas)
	if err != nil {
		return err
	}

	// create formulas dir
	err = os.MkdirAll(dirname, os.ModePerm)
	if err != nil {
		return err
	}
	for _, f := range formulas {
		name := fmt.Sprintf("%s.formula.json", strings.Replace(f.Name, " ", "", -1))
		formulaBytes, err := json.Marshal(f)
		if err != nil {
			break
		}
		fmt.Printf("Exporting '%s' to %s/%s\n", f.Name, dirname, name)
		err = ioutil.WriteFile(fmt.Sprintf("%s/%s", dirname, name), formulaBytes, 0644)
	}

	return nil
}

// ExportAllResourcesToDir writes out all the resources to the speceified irectory
func ExportAllResourcesToDir(base, auth string, dirname string) error {
	resourcesListBytes, _, _, err := ce.ResourcesList(base, auth)
	if err != nil {
		return err
	}
	var resources []CommonResource
	err = json.Unmarshal(resourcesListBytes, &resources)
	if err != nil {
		return err
	}
	err = os.MkdirAll(dirname, os.ModePerm)
	if err != nil {
		return err
	}
	for _, r := range resources {

		resourceBytes, _, _, err := ce.GetResourceDefinition(base, auth, r.Name)
		if err != nil {
			log.Println(err.Error())
			break
		}
		name := fmt.Sprintf("%s.cro.json", r.Name)
		fmt.Printf("Exporting %s to %s/%s\n", r.Name, dirname, name)
		err = ioutil.WriteFile(fmt.Sprintf("%s/%s", dirname, name), resourceBytes, 0644)
	}

	return nil
}

// cloneCmd is the command to clone assets between accounts
var cloneCmd = &cobra.Command{
	Use:    "clone",
	Short:  "clone assets from one profile to another",
	Long:   "Clone exports assets from one account profile (--from) and imports them into another profile (--to)",
	Hidden: true,
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

	moleculesCmd.PersistentFlags().StringVar(&profile, "profile", "default", "profile name")

	moleculesCmd.AddCommand(exportCmd)
	moleculesCmd.AddCommand(cloneCmd)
	cloneCmd.PersistentFlags().StringVar(&profileSource, "from", "default", "source profile name")
	cloneCmd.PersistentFlags().StringVar(&profileTarget, "to", "", "target profile name")
	cloneCmd.MarkPersistentFlagRequired("to")
}

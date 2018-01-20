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
	"log"
	"os"

	"github.com/ghchinoy/cectl/ce"
	"github.com/spf13/cobra"
)

// formulasCmd represents the formulas command
var formulasCmd = &cobra.Command{
	Use:   "formulas",
	Short: "Manage formulas on the platform",
	Long:  `Currently allows listing formulas and execution details`,
}

// listFormulasCmd represents the listFormulas command
var listFormulasCmd = &cobra.Command{
	Use:   "list",
	Short: "lists available formulas on the plaform",
	Long:  `returns a list of formulas on the platform`,
	Run: func(cmd *cobra.Command, args []string) {

		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		bodybytes, statuscode, curlcmd, err := ce.FormulasList(profilemap["base"], profilemap["auth"])
		if err != nil {
			fmt.Println("Unable to list formulas", err.Error())
			os.Exit(1)
		}

		// handle global options, curl
		if showCurl {
			log.Println(curlcmd)
		}
		// handle non 200
		if statuscode != 200 {
			log.Printf("HTTP Error: %v\n", statuscode)
			// handle this nicely, show error description
		}
		// handle global options, json
		if outputJSON {
			fmt.Printf("%s\n", bodybytes)
			return
		}

		err = ce.OutputFormulasList(bodybytes, profilemap["base"], profilemap["auth"])
		if err != nil {
			fmt.Println("Unable to render formula table", err.Error())
		}
	},
}

func init() {
	RootCmd.AddCommand(formulasCmd)
	formulasCmd.AddCommand(listFormulasCmd)

	formulasCmd.PersistentFlags().StringVar(&profile, "profile", "default", "profile name")
	formulasCmd.PersistentFlags().BoolVarP(&outputJSON, "json", "j", false, "output as json")
	formulasCmd.PersistentFlags().BoolVarP(&showCurl, "curl", "c", false, "show curl command")
}

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
	"log"
	"os"
	"strconv"

	"github.com/ghchinoy/cectl/ce"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// formulaActivateCmd represents the formulaActivate command
var formulaActivateCmd = &cobra.Command{
	Use:   "activate",
	Short: "Activate a Formula template",
	Long:  `Sets a Formula template to active state`,
	Run: func(cmd *cobra.Command, args []string) {
		if !viper.IsSet(profile + ".base") {
			fmt.Println("Can't find info for profile", profile)
			os.Exit(1)
		}

		if len(args) < 1 {
			fmt.Println("must supply an ID of a Formula")
			os.Exit(1)
		}

		base := viper.Get(profile + ".base")
		user := viper.Get(profile + ".user")
		org := viper.Get(profile + ".org")
		auth := fmt.Sprintf("User %s, Organization %s", user, org)

		// Get the Formula
		formulaResponseBytes, statuscode, curlcmd, err := ce.FormulaDetailsAsBytes(args[0], fmt.Sprintf("%s", base), auth)
		if statuscode != 200 {
			fmt.Printf("Unable to retrieve formula %s, %s\n", args[0], err.Error())
			os.Exit(1)
		}
		var formula ce.Formula
		err = json.Unmarshal(formulaResponseBytes, &formula)
		if err != nil {
			fmt.Println("Unable to understand formula response", err.Error())
			os.Exit(1)
		}

		// Change the Formula to Active
		formula.Active = true

		// PATCH to set the Formula back
		patchBytes, statuscode, err := ce.FormulaUpdate(args[0], base.(string), auth, formula)
		err = json.Unmarshal(patchBytes, &formula)
		if err != nil {
			fmt.Printf("Unable to retrieve formula, %s\n", err.Error())
			os.Exit(1)
		}

		if showCurl {
			log.Println(curlcmd)
		}

		if outputJSON {
			fmt.Printf("%s\n", patchBytes)
			return
		}

		if statuscode != 200 {
			fmt.Println(statuscode)
			var ficr ce.FormulaInstanceCreationResponse
			err = json.Unmarshal(patchBytes, &ficr)
			if err != nil {
				fmt.Println("Cannot process response, tried error message")
				os.Exit(1)
			}
			fmt.Println(ficr.Message)
			os.Exit(1)
		}

		var f ce.Formula
		err = json.Unmarshal(patchBytes, &f)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		formulaResponseBytes, statuscode, curlcmd, err = ce.FormulaDetailsAsBytes(strconv.Itoa(f.ID), fmt.Sprintf("%s", base), auth)
		if statuscode != 200 {
			fmt.Printf("Unable to retrieve updated formula %s, %s\n", args[0], err.Error())
			os.Exit(1)
		}
		err = json.Unmarshal(formulaResponseBytes, &f)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		var instancecount string
		instances, err := ce.GetInstancesOfFormula(f.ID, base.(string), auth)
		if err != nil {
			// unable to retrieve instances of formula!
			instancecount = "N/A"
		}
		instancecount = strconv.Itoa(len(instances))

		api := "N/A"
		if f.Triggers[0].Type == "manual" {
			api = f.API
		}

		data := [][]string{}
		data = append(data, []string{
			strconv.Itoa(f.ID),
			f.Name,
			strconv.FormatBool(f.Active),
			strconv.Itoa(len(f.Steps)),
			f.Triggers[0].Type,
			instancecount,
			api,
		})

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "active", "steps", "trigger", "instances", "api"})
		table.SetBorder(false)
		table.AppendBulk(data)
		table.Render()
	},
}

func init() {
	formulasCmd.AddCommand(formulaActivateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// formulaActivateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// formulaActivateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

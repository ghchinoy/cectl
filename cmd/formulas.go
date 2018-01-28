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
	"strconv"

	"github.com/ghchinoy/ce-go/ce"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// formulasCmd represents the formulas command
var formulasCmd = &cobra.Command{
	Use:   "formulas",
	Short: "Manage formulas on the platform",
	Long:  `Currently allows listing formulas and execution details`,
}

var importFormulaPath string

// importFormulaCmd represents the importFormula command
var importFormulaCmd = &cobra.Command{
	Use:   "import <filepath>",
	Short: "Imports a Formula to the platform",
	Long:  `Providing a Formula JSON, this command adds a Formula template to the platform`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("must supply a path to a Formula JSON")
			os.Exit(1)
		}

		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// read in file
		filebytes, err := ioutil.ReadFile(args[0])
		if err != nil {
			fmt.Println("unable to read file", args[0], err.Error())
			os.Exit(1)
		}

		// Check if can decode into formula struct
		var f ce.Formula
		err = json.Unmarshal(filebytes, &f)
		if err != nil {
			fmt.Println(args[0], "doesn't seem like a Formula")
			os.Exit(1)
		}

		bodybytes, status, curlcmd, err := ce.ImportFormula(
			profilemap["base"],
			profilemap["auth"],
			f,
		)

		if showCurl {
			log.Println(curlcmd)
		}
		if status != 200 {
			fmt.Println(status)
			fmt.Printf("%s\n", bodybytes)
			os.Exit(1)
		}

		if outputJSON {
			fmt.Printf("%s\n", bodybytes)
			return
		}

		if status == 200 {
			fmt.Println("Formula template added to Platform.")
			var f ce.Formula
			err = json.Unmarshal(bodybytes, &f)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			err = ce.FormulaDetailsTableOutput(f)
			if err != nil {
				fmt.Println("Unable to render Formula details")
				os.Exit(1)
			}
		}

	},
}

// formulaDeactivateCmd represents the formulaDeactivate command
var formulaDeactivateCmd = &cobra.Command{
	Use:   "deactivate",
	Short: "Deactivate a Formula template",
	Long:  `Sets a Formula template to an inactive state`,
	Run: func(cmd *cobra.Command, args []string) {
		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if len(args) < 1 {
			fmt.Println("must supply an ID of a Formula")
			os.Exit(1)
		}

		// Get the Formula
		formulaResponseBytes, statuscode, curlcmd, err := ce.FormulaDetailsAsBytes(args[0], profilemap["base"], profilemap["auth"])
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
		formula.Active = false

		// PATCH to set the Formula back
		patchBytes, statuscode, err := ce.FormulaUpdate(args[0], profilemap["base"], profilemap["auth"], formula)
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

		formulaResponseBytes, statuscode, curlcmd, err = ce.FormulaDetailsAsBytes(strconv.Itoa(f.ID), profilemap["base"], profilemap["auth"])
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
		instances, err := ce.GetInstancesOfFormula(f.ID, profilemap["base"], profilemap["auth"])
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

// formulaActivateCmd represents the formulaActivate command
var formulaActivateCmd = &cobra.Command{
	Use:   "activate",
	Short: "Activate a Formula template",
	Long:  `Sets a Formula template to active state`,
	Run: func(cmd *cobra.Command, args []string) {
		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if len(args) < 1 {
			fmt.Println("must supply an ID of a Formula")
			os.Exit(1)
		}

		base := profilemap["base"]
		auth := profilemap["auth"]

		// Get the Formula
		formulaResponseBytes, statuscode, curlcmd, err := ce.FormulaDetailsAsBytes(args[0], base, auth)
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
		patchBytes, statuscode, err := ce.FormulaUpdate(args[0], base, auth, formula)
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

		formulaResponseBytes, statuscode, curlcmd, err = ce.FormulaDetailsAsBytes(strconv.Itoa(f.ID), base, auth)
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
		instances, err := ce.GetInstancesOfFormula(f.ID, base, auth)
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

// deleteFormulaCmd represents the deleteFormula command
var deleteFormulaCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a specific Formula template",
	Long:  `Delete a Formula template on the platform`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("must supply an ID of a Formula")
			os.Exit(1)
		}

		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		bodybytes, status, curlcmd, err := ce.DeleteFormula(profilemap["base"], profilemap["auth"], args[0])
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		if showCurl {
			log.Println(curlcmd)
		}

		if outputJSON {
			fmt.Printf("%s\n", bodybytes)
			return
		}

		if status != 200 {
			fmt.Println(status)
			var ficr ce.FormulaInstanceCreationResponse
			err = json.Unmarshal(bodybytes, &ficr)
			if err != nil {
				fmt.Println("Cannot process response, tried error message")
				os.Exit(1)
			}
			fmt.Println(ficr.Message)
			os.Exit(1)
		}

		if resp.StatusCode == 200 {
			fmt.Printf("Formula %s deleted.\n", args[0])
			fmt.Printf("%s\n", bodybytes)
		}
	},
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
	formulasCmd.AddCommand(deleteFormulaCmd)
	formulasCmd.AddCommand(formulaActivateCmd)
	formulasCmd.AddCommand(formulaDeactivateCmd)
	formulasCmd.AddCommand(importFormulaCmd)

	formulasCmd.PersistentFlags().StringVar(&profile, "profile", "default", "profile name")
	formulasCmd.PersistentFlags().BoolVarP(&outputJSON, "json", "j", false, "output as json")
	formulasCmd.PersistentFlags().BoolVarP(&showCurl, "curl", "c", false, "show curl command")
}

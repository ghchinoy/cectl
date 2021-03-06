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

	"github.com/ghchinoy/ce-go/ce"
	"github.com/spf13/cobra"
)

// formula-instancesCmd represents the formula-instances command
var formulaInstancesCmd = &cobra.Command{
	Use:   "formula-instances",
	Short: "Manage Formula Instances",
	Long:  `Manage configured Instances of a Formula template`,
}

var formulaInstanceConfiguration string

// createInstanceCmd is the command for creating a Formula Instance
var createInstanceCmd = &cobra.Command{
	Use:   "create <id> [name]",
	Short: "creates an instance of a Formula, given a Formula ID",
	Long: `Given the ID of Formula template, create an Instance of the Formula
Optionally, provide the Formula configuration definition via a flag.
A name for the Formula Instance will be required when using a flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please supply an ID of a Formula template\ncectl formula-instance create <ID> [name]")
			cmd.Help()
			os.Exit(1)
		}

		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// if no instance config json given, check for name
		var config ce.FormulaInstanceConfig
		// formulaInstanceConfiguration will be set if the flag --configuration has a value (which should be a JSON filename)
		if formulaInstanceConfiguration == "" {
			if len(args) < 2 {
				fmt.Println("Please provide a name for the Instance if not submitting a Formula Instance configuration definition\ncectl formula-instance create <ID> [name]")
				os.Exit(1)
			}
			config = ce.FormulaInstanceConfig{Name: args[1], Active: true}
		} else {
			var raw map[string]interface{}
			_ = json.Unmarshal([]byte(formulaInstanceConfiguration), &raw)
			config = ce.FormulaInstanceConfig{Name: args[1], Active: true, Configuration: raw}
		}

		bodybytes, status, curlcmd, err := ce.CreateFormulaInstance(profilemap["base"], profilemap["auth"], args[0], config)

		if showCurl {
			log.Println(curlcmd)
		}
		if outputJSON {
			fmt.Printf("%s\n", bodybytes)
			return
		}
		// handle non 200
		if status != 200 {
			log.Printf("HTTP Error: %v\n", status)
			// handle this nicely, show error description
			os.Exit(1)
		}
		var response map[string]interface{}
		err = json.Unmarshal(bodybytes, &response)
		if err != nil {
			log.Println("Couldn't unmarshal response", err.Error())
			os.Exit(1)
		}
		fmt.Printf("Formula Instance ID %v created.\n", response["id"])
	},
}

var (
	triggerBody       string
	triggerTextOutput bool
)

// triggerCmd represents the trigger command
var triggerCmd = &cobra.Command{
	Use:   "trigger <id>",
	Short: "invoke a Formula Instance",
	Long: `Invokes a Formula Instance by ID, resulting in a Formula Instance Execution.
The trigger body is optional (defaults to: {}).
This will only invoke a manually triggerable Formula.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("must supply an ID of a Formula Instance")
			os.Exit(1)
		}

		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		bodybytes, status, curlcmd, err := ce.TriggerFormulaInstance(profilemap["base"], profilemap["auth"], args[0], triggerBody)

		if outputJSON {
			fmt.Printf("%s\n", bodybytes)
			return
		}

		if showCurl {
			log.Println(curlcmd)
		}

		if status != 200 {
			var ex ce.FormulaInstanceCreationResponse
			err = json.Unmarshal(bodybytes, &ex)
			fmt.Printf("%s\nID: %v (%s)\n", ex.Message, args[0], ex.RequestID)
			fmt.Println(status)
			return
		}

		var ex []ce.FormulaInstanceCreationResponse
		err = json.Unmarshal(bodybytes, &ex)
		if err != nil { // that's not an array of responses
			log.Println(err.Error())
			return
		}

		if triggerTextOutput {
			fmt.Print(ex[0].ID)
		} else {
			fmt.Printf("Execution ID: %v\n", ex[0].ID)
		}
	},
}

var deleteFormulaInstanceCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "deletes a Formula Instance",
	Long:  `Deletes a Formula Instance by ID`,
	Run: func(cmd *cobra.Command, args []string) {
		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// check for Instance ID & Operation name
		if len(args) < 1 {
			fmt.Println("Please provide an Instance ID ")
			return
		}
		if _, err := strconv.ParseInt(args[0], 10, 64); err != nil {
			fmt.Println("Please provide an Instance ID that is an integer")
			return
		}
		bodybytes, statuscode, curlcmd, err := ce.DeleteFormulaInstance(profilemap["base"], profilemap["auth"], args[0])
		if err != nil {
			if statuscode == -1 {
				fmt.Println("Unable to reach CE API. Please check your configuration / profile.")
			}
			fmt.Println(err)
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
		fmt.Printf("Deleted Formula Instance %s", args[0])
	},
}

func init() {
	RootCmd.AddCommand(formulaInstancesCmd)

	formulaInstancesCmd.PersistentFlags().StringVar(&profile, "profile", "default", "profile name")
	formulaInstancesCmd.PersistentFlags().BoolVarP(&outputJSON, "json", "j", false, "output as json")
	formulaInstancesCmd.PersistentFlags().BoolVarP(&showCurl, "curl", "c", false, "show curl command")

	formulaInstancesCmd.AddCommand(triggerCmd)
	triggerCmd.Flags().StringVarP(&triggerBody, "data", "d", "{}", "data for trigger body")
	triggerCmd.Flags().BoolVarP(&triggerTextOutput, "text", "t", false, "output trigger id as text")

	formulaInstancesCmd.AddCommand(deleteFormulaInstanceCmd)

	formulaInstancesCmd.AddCommand(createInstanceCmd)

	createInstanceCmd.Flags().StringVarP(&formulaInstanceConfiguration, "configuration", "", "", "instance configuration definition")
	// deprecated
	createInstanceCmd.Flags().StringVarP(&formulaInstanceConfiguration, "instance", "i", "", "instance configuration definition")

}

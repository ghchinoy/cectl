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
	"net/http"
	"os"
	"strconv"

	"github.com/ghchinoy/ce-go/ce"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	outputLimit                   int
	formulaExecutionQueryEventID  int
	formulaExecutionQueryObjectID int
)

// formulaInstanceExecutionsCmd represents the formula-instance-executions command
var formulaInstanceExecutionsCmd = &cobra.Command{
	Use:   "executions",
	Short: "Manage Formula Instance Executions",
	Long:  `Manage the Executions of a Formula Instance`,
}

// listFormulaInstanceExecutionsCmd represents the listFormulaInstanceExecutions command
var listFormulaInstanceExecutionsCmd = &cobra.Command{
	Use:   "list <id>",
	Short: "list executions for instance id",
	Long:  `Lists the Formula Instance Executions given an ID of a Formula Instance`,
	Run: func(cmd *cobra.Command, args []string) {

		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if len(args) < 1 {
			cmd.Usage()
			os.Exit(1)
		}

		// add eventId and/or objectId query params as needed
		// if formulaExecutionQueryObjectID | formulaExecutionQueryEventID > 0
		var query []string

		if formulaExecutionQueryEventID > 0 {
			query = append(query, fmt.Sprintf("eventId=%v", formulaExecutionQueryEventID))
		}
		if formulaExecutionQueryObjectID > 0 {
			query = append(query, fmt.Sprintf("objectId=%v", formulaExecutionQueryObjectID))
		}
		if len(query) > 0 {

		}

		bodybytes, status, curlcmd, err := ce.GetFormulaInstanceExecutions(profilemap["base"], profilemap["auth"], args[0])

		if showCurl {
			log.Println(curlcmd)
		}

		if outputJSON {
			fmt.Printf("%s\n", bodybytes)
			return
		}

		if status != 200 {
			fmt.Print(status)
			if status == 404 {
				fmt.Println("Unable to contact CE API")
				return
			}
			fmt.Println()
		}

		data := [][]string{}

		var executions []ce.FormulaInstanceExecution
		err = json.Unmarshal(bodybytes, &executions)
		if err != nil { // can't make an array of executions
			log.Println(err.Error())
		}
		i := 0

		for _, v := range executions {
			diff := v.UpdatedDate.Sub(v.CreateDate)
			var difftext string
			if diff < 0 {
				difftext = "pending"
			} else {
				difftext = fmt.Sprintf("%v s", diff.Seconds())
			}
			data = append(data, []string{
				strconv.Itoa(v.ID),
				strconv.Itoa(v.FormulaInstanceID),
				v.Status,
				v.CreateDate.String(),
				v.UpdatedDate.String(),
				difftext,
			})
			i++
			if outputLimit > 0 {
				if i > outputLimit {
					break
				}
			}
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Instance", "Status", "Created", "Updated", "Duration"})
		table.SetBorder(false)
		table.AppendBulk(data)
		table.Render()

	},
}

// cancelExecutionCmd represents the cancelExecution command
var cancelExecutionCmd = &cobra.Command{
	Use:   "cancel <id>",
	Short: "Cancel a Formula Instance Execution by ID",
	Long:  `Given an Execution ID, cancel the Formula Instance Execution`,
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

		bodybytes, status, curlcmd, err := ce.CancelFormulaExecution(
			profilemap["base"],
			profilemap["auth"],
			args[0],
		)

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
		fmt.Println(status)

	},
}

// retryFormulaInstanceExecutionCmd represents the retryFormulaInstanceExecution command
var retryFormulaInstanceExecutionCmd = &cobra.Command{
	Use:   "retry",
	Short: "retry an execution",
	Long:  `Retry a Formula Instance Execution that has previously run`,
	Run: func(cmd *cobra.Command, args []string) {
		if !viper.IsSet(profile + ".base") {
			fmt.Println("Can't find info for profile", profile)
			os.Exit(1)
		}

		if len(args) < 1 {
			fmt.Println("must supply an ID of a Formula Instance Execution")
			os.Exit(1)
		}

		base := viper.Get(profile + ".base")
		user := viper.Get(profile + ".user")
		org := viper.Get(profile + ".org")

		url := fmt.Sprintf("%s%s",
			base,
			fmt.Sprintf(ce.FormulaRetryExecutionURI, args[0]),
		)
		auth := fmt.Sprintf("User %s, Organization %s", user, org)

		client := &http.Client{}
		req, err := http.NewRequest("PUT", url, nil)
		if err != nil {
			fmt.Println("Can't construct request", err.Error())
			os.Exit(1)
		}
		req.Header.Add("Authorization", auth)
		req.Header.Add("Accept", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Cannot process response", err.Error())
			os.Exit(1)
		}
		bodybytes, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		if outputJSON {
			fmt.Printf("%s\n", bodybytes)
			return
		}

		if resp.StatusCode != 200 {
			var ex ce.FormulaInstanceCreationResponse
			err = json.Unmarshal(bodybytes, &ex)
			fmt.Printf("%s\nID: %v (%s)\n", ex.Message, args[0], ex.RequestID)
			fmt.Println(resp.Status)
			return
		}

	},
}

// detailExecutionIDCmd results in the specific execution details
var detailExecutionIDCmd = &cobra.Command{
	Use:   "details <id>",
	Short: "Show details of a Formula Instance Execution by ID",
	Long:  `Given an Execution ID, show details of each of the steps in a the Formula Instance Execution`,
	Run: func(cmd *cobra.Command, args []string) {

		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// guard
		if len(args) < 1 {
			fmt.Println("must supply an ID of a Formula Instance Execution")
			os.Exit(1)
		}

		// guard
		_, err = strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Please provide a number as a Formula Execution ID")
			os.Exit(1)
		}

		bodybytes, statuscode, curlcmd, err := ce.GetFormulaInstanceExecutionID(args[0], profilemap["base"], profilemap["auth"])
		if err != nil {
			fmt.Println("Unable to get formula execution id", err.Error())
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

	},
}

func init() {
	RootCmd.AddCommand(formulaInstanceExecutionsCmd)
	formulaInstanceExecutionsCmd.PersistentFlags().StringVar(&profile, "profile", "default", "profile name")
	formulaInstanceExecutionsCmd.PersistentFlags().BoolVarP(&outputJSON, "json", "j", false, "output as json")

	formulaInstanceExecutionsCmd.AddCommand(listFormulaInstanceExecutionsCmd)
	listFormulaInstanceExecutionsCmd.Flags().IntVarP(&outputLimit, "top", "t", 0, "output limit from latest")
	listFormulaInstanceExecutionsCmd.Flags().IntVarP(&formulaExecutionQueryEventID, "event", "e", 0, "event ID to search for")
	listFormulaInstanceExecutionsCmd.Flags().IntVarP(&formulaExecutionQueryObjectID, "object", "o", 0, "object ID to search for")

	formulaInstanceExecutionsCmd.AddCommand(cancelExecutionCmd)

	formulaInstanceExecutionsCmd.AddCommand(retryFormulaInstanceExecutionCmd)

	formulaInstanceExecutionsCmd.AddCommand(detailExecutionIDCmd)
}

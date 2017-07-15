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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/ghchinoy/cectl/ce"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// cancelExecutionCmd represents the cancelExecution command
var cancelExecutionCmd = &cobra.Command{
	Use:   "cancel <id>",
	Short: "Cancel a Formula Instance Execution by ID",
	Long:  `Given an Execution ID, cancel the Formula Instance Execution`,
	Run: func(cmd *cobra.Command, args []string) {

		// construct a fixed json body for sending cancelled status
		cancelmessage := struct {
			Status string `json:"status"`
		}{"cancelled"}
		cancelbytes, err := json.Marshal(cancelmessage)
		if err != nil {
			fmt.Println("Can't even")
			os.Exit(1)
		}

		if len(args) < 1 {
			fmt.Println("must supply an ID of a Formula")
			os.Exit(1)
		}

		if !viper.IsSet(profile + ".base") {
			fmt.Println("Can't find info for profile", profile)
			os.Exit(1)
		}

		base := viper.Get(profile + ".base")
		user := viper.Get(profile + ".user")
		org := viper.Get(profile + ".org")

		url := fmt.Sprintf("%s%s",
			base,
			fmt.Sprintf(ce.FormulaCancelExecutionURIFormat, args[0]),
		)
		auth := fmt.Sprintf("User %s, Organization %s", user, org)

		client := &http.Client{}
		req, err := http.NewRequest("PATCH", url, bytes.NewReader(cancelbytes))
		if err != nil {
			fmt.Println("Can't construct request", err.Error())
			os.Exit(1)
		}
		req.Header.Add("Authorization", auth)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
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
			fmt.Println(resp.Status)
			var ficr ce.FormulaInstanceCreationResponse
			err = json.Unmarshal(bodybytes, &ficr)
			if err != nil {
				fmt.Println("Cannot process response, tried error message")
				os.Exit(1)
			}
			fmt.Println(ficr.Message)
			os.Exit(1)
		}
		fmt.Println(resp.Status)

	},
}

func init() {
	formulaInstanceExecutionsCmd.AddCommand(cancelExecutionCmd)

}

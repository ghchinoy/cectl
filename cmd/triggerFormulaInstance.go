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
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"encoding/json"

	"github.com/ghchinoy/cectl/ce"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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

		if !viper.IsSet(profile + ".base") {
			fmt.Println("Can't find info for profile", profile)
			os.Exit(1)
		}

		base := viper.Get(profile + ".base")
		user := viper.Get(profile + ".user")
		org := viper.Get(profile + ".org")

		url := fmt.Sprintf("%s%s",
			base,
			fmt.Sprintf(ce.FormulaExecutionsURIFormat, args[0]),
		)
		auth := fmt.Sprintf("User %s, Organization %s", user, org)

		client := &http.Client{}
		req, err := http.NewRequest("POST", url, bytes.NewReader([]byte(triggerBody)))
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
			var ex ce.FormulaInstanceCreationResponse
			err = json.Unmarshal(bodybytes, &ex)
			fmt.Printf("%s\nID: %v (%s)\n", ex.Message, args[0], ex.RequestID)
			fmt.Println(resp.Status)
			return
		}

		var ex []ce.FormulaInstanceCreationResponse
		err = json.Unmarshal(bodybytes, &ex)

		if triggerTextOutput {
			fmt.Print(ex[0].ID)
		} else {
			fmt.Printf("Execution ID: %v\n", ex[0].ID)
		}

	},
}

func init() {
	formulaInstancesCmd.AddCommand(triggerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// triggerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// triggerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	triggerCmd.Flags().StringVarP(&triggerBody, "data", "d", "{}", "data for trigger body")
	triggerCmd.Flags().BoolVarP(&triggerTextOutput, "text", "t", false, "output trigger id as text")

}

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

var formulaInstance string

// createInstanceCmd represents the createInstance command
var createInstanceCmd = &cobra.Command{
	Use:   "create <id> [name]",
	Short: "creates an instance of a Formula, given a Formula ID",
	Long:  `Given a Formula ID, create an Instance of the Formula`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please supply an ID of a Formula template\ncectl formula-instance create <ID> [name]")
			os.Exit(1)
		}

		if !viper.IsSet(profile + ".base") {
			fmt.Println("Can't find info for profile", profile)
			os.Exit(1)
		}
		base := viper.Get(profile + ".base")
		user := viper.Get(profile + ".user")
		org := viper.Get(profile + ".org")

		createformat := "/formulas/%s/instances"
		url := fmt.Sprintf("%s%s", base, fmt.Sprintf(createformat, args[0]))

		auth := fmt.Sprintf("User %s, Organization %s", user, org)

		// if no instance config json given, check for name
		var fi ce.FormulaInstanceConfig
		if formulaInstance == "" {
			if len(args) < 2 {
				fmt.Println("Please provide a name if not submitting a Formula Instance configuration\ncectl formula-instance create <ID> [name]")
				os.Exit(1)
			}
			fi = ce.FormulaInstanceConfig{Name: args[1], Active: true}
		}

		fibytes, err := json.Marshal(fi)
		fmt.Println(url)
		fmt.Printf("%s\n", fibytes)
		if err != nil {
			fmt.Println("Unable to convert to Formula Instance configuration json", err.Error())
			os.Exit(1)
		}

		client := &http.Client{}
		req, err := http.NewRequest("POST", url, bytes.NewReader(fibytes))
		if err != nil {
			fmt.Println("Unable to create request", err.Error())
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
		fmt.Println(resp.Status)
	},
}

func init() {
	formulaInstancesCmd.AddCommand(createInstanceCmd)

	createInstanceCmd.Flags().StringVarP(&formulaInstance, "instance", "i", "", "instance configuration")

}

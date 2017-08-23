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
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ghchinoy/cectl/ce"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listFormualInstancesCmd represents the listFormualInstances command
var listFormulaInstancesCmd = &cobra.Command{
	Use:   "instances <id>",
	Short: "List Instances associated with a specific Formula",
	Long:  `Retrieve a list of all instances associated with a particular formula template.`,
	Run: func(cmd *cobra.Command, args []string) {

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
			fmt.Sprintf(ce.FormulaInstancesURIFormat, args[0]),
		)
		auth := fmt.Sprintf("User %s, Organization %s", user, org)

		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
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

		if resp.StatusCode == 200 {
			data := [][]string{}

			var instances []ce.FormulaInstance
			err = json.Unmarshal(bodybytes, &instances)
			for _, v := range instances {

				var configs []string
				if c, ok := v.Configuration.(map[string]interface{}); ok {
					for k, v := range c {
						configs = append(configs, fmt.Sprintf("%s:%s", k, v))
					}
				}

				data = append(data, []string{
					strconv.Itoa(v.ID),
					v.Name,
					strconv.FormatBool(v.Active),
					fmt.Sprintf("%v %s", v.Formula.ID, v.Formula.Name),
					strings.Join(configs, ", "),
					v.CreatedDate.String(),
				})
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Instance", "active", "Formula", "Configuration", "Created"})
			table.SetBorder(false)
			table.AppendBulk(data)
			table.Render()
		}

	},
}

func init() {
	formulasCmd.AddCommand(listFormulaInstancesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listFormualInstancesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listFormualInstancesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

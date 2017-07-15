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

	"github.com/ghchinoy/cectl/ce"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var outputLimit int

// listFormulaInstanceExecutionsCmd represents the listFormulaInstanceExecutions command
var listFormulaInstanceExecutionsCmd = &cobra.Command{
	Use:   "list <id>",
	Short: "list executions for instance id",
	Long:  `Lists the Formula Instance Executions given an ID of a Formula Instance`,
	Run: func(cmd *cobra.Command, args []string) {
		if !viper.IsSet(profile + ".base") {
			fmt.Println("Can't find info for profile", profile)
			os.Exit(1)
		}

		if len(args) < 1 {
			fmt.Println("must supply an ID of a Formula Instance")
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
		req, err := http.NewRequest("GET", url, nil)
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

		data := [][]string{}

		var executions []ce.FormulaInstanceExecution
		err = json.Unmarshal(bodybytes, &executions)
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

func init() {
	formulaInstanceExecutionsCmd.AddCommand(listFormulaInstanceExecutionsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listFormulaInstanceExecutionsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listFormulaInstanceExecutionsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	listFormulaInstanceExecutionsCmd.Flags().IntVarP(&outputLimit, "num", "n", 0, "output limit")

}

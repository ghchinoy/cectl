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
	"strconv"

	"github.com/ghchinoy/ce-go/ce"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listHubsCmd represents the listHubs command
var listHubsCmd = &cobra.Command{
	Use:   "list",
	Short: "list available hubs",
	Long:  `List hubs on the Platform`,
	Run: func(cmd *cobra.Command, args []string) {
		if !viper.IsSet(profile + ".base") {
			fmt.Println("Can't find info for profile", profile)
			os.Exit(1)
		}

		base := viper.Get(profile + ".base").(string)
		user := viper.Get(profile + ".user").(string)
		org := viper.Get(profile + ".org").(string)

		hubs, curlcmd, err := ce.ListHubs(base, user, org, outputJSON)
		if err != nil {
			fmt.Println("Unable to read hubs", err.Error())
			os.Exit(1)
		}

		if showCurl {
			log.Println(curlcmd)
		}

		if !outputJSON {
			data := [][]string{}

			for _, v := range hubs {
				data = append(data, []string{
					strconv.Itoa(v.ID),
					v.Name,
					v.Key,
					strconv.FormatBool(v.Active),
					v.Description,
				})
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Name", "Key", "Active", "Description"})
			table.SetBorder(false)
			table.AppendBulk(data)
			table.Render()
		}
	},
}

func init() {
	hubsCmd.AddCommand(listHubsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listHubsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listHubsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

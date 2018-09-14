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

	"github.com/ghchinoy/ce-go/ce"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// jobsCmd represents the jobs command
var jobsCmd = &cobra.Command{
	Use:   "jobs",
	Short: "Manage jobs on the platform",
	Long:  `Manage jobs on the platform`,
}

// listJobsCmd represents the listJobs command
var listJobsCmd = &cobra.Command{
	Use:   "list",
	Short: "list jobs on platform",
	Long:  `List jobs on the platform`,
	Run: func(cmd *cobra.Command, args []string) {
		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		bodybytes, status, curlcmd, err := ce.ListJobs(profilemap["base"], profilemap["auth"])

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
				fmt.Printf("Unable to contact CE API, %s\n", profilemap["base"])
				return
			}
			fmt.Println()
		}

		data := [][]string{}

		var jobs []ce.Job
		err = json.Unmarshal(bodybytes, &jobs)
		if err != nil {
			fmt.Printf("Response not a list of Jobs, %s", err.Error())
			return
		}
		for _, v := range jobs {
			data = append(data, []string{
				v.ID,
				v.Name,
				v.Description,
			})
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Description"})
		table.SetBorder(false)
		table.AppendBulk(data)
		table.Render()
	},
}

func init() {
	RootCmd.AddCommand(jobsCmd)

	jobsCmd.PersistentFlags().StringVar(&profile, "profile", "default", "profile name")
	jobsCmd.PersistentFlags().BoolVarP(&outputJSON, "json", "j", false, "output as json")
	jobsCmd.PersistentFlags().BoolVarP(&showCurl, "curl", "c", false, "show curl command")

	jobsCmd.AddCommand(listJobsCmd)
}

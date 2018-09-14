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

var jobsDeleteAll bool

var deleteJobCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a job on the platform",
	Long:  "Delete a job on the platform given the job ID",
	Run: func(cmd *cobra.Command, args []string) {
		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if jobsDeleteAll == true {
			bodybytes, status, _, err := ce.ListJobs(profilemap["base"], profilemap["auth"])
			if err != nil {
				fmt.Println("Unable to list jobs", err.Error())
				os.Exit(1)
			}
			if status != 200 {
				fmt.Println("non-200 response", status)
				os.Exit(1)
			}
			var jobs []ce.Job
			err = json.Unmarshal(bodybytes, &jobs)
			if err != nil {
				fmt.Printf("Response not a list of Jobs, %s", err.Error())
				return
			}
			if len(jobs) == 0 {
				fmt.Print("No jobs to delete\n")
				os.Exit(0)
			}
			results := make(chan DeleteJobCheck)
			var unabletodelete []string
			for _, v := range jobs {
				go deleteJob(profilemap["base"], profilemap["auth"], v.ID, results)
				//fmt.Println(v.ID)
			}
			var num int
			for i := range results {
				if i.StatusCode != 200 {
					fmt.Printf("%v: unable to delete %s\n", i.StatusCode, i.JobID)
					unabletodelete = append(unabletodelete, i.JobID)
				}
				num++
				if len(jobs) == num {
					close(results)
				}
			}

			fmt.Printf("%v/%v 200\n", len(jobs)-len(unabletodelete), len(jobs))
			os.Exit(0)
		}

		if len(args) == 0 {
			fmt.Println("Job ID must be provided")
			cmd.Help()
			os.Exit(1)
		}

		bodybytes, status, curlcmd, err := ce.DeleteJob(profilemap["base"], profilemap["auth"], args[0])

		if showCurl {
			log.Println(curlcmd)
		}

		if outputJSON {
			fmt.Printf("%s\n", bodybytes)
			return
		}

		if err != nil {
			fmt.Println("error", err)
			os.Exit(1)
		}

		if status != 200 {
			fmt.Print(status)
			if status == 404 {
				fmt.Printf("Unable to contact CE API, %s\n", profilemap["base"])
				return
			}
			fmt.Println()
		}

		fmt.Printf("%s deleted\n", args[0])
	},
}

// DeleteJobCheck is a struct to hold results of a multidelete
type DeleteJobCheck struct {
	JobID      string
	StatusCode int
}

func deleteJob(base, auth, jobID string, checks chan DeleteJobCheck) (DeleteJobCheck, error) {
	var d DeleteJobCheck
	log.Println(jobID)
	_, status, _, err := ce.DeleteJob(base, auth, jobID)
	if err != nil {
		checks <- d
		return d, err
	}
	d = DeleteJobCheck{jobID, status}
	checks <- d
	return d, nil
}

var jobJSONFile string

var createJobCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a job on the platform",
	Long:  "Create a job on the platform given a file path to a JSON job description, using --file",
	Run: func(cmd *cobra.Command, args []string) {
		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// check for file
		if jobJSONFile == "" {
			fmt.Printf(cmd.UsageString())
			return
		}
		var filebytes []byte
		if jobJSONFile != "" {
			// read in file
			filebytes, err = ioutil.ReadFile(jobJSONFile)
			if err != nil {
				fmt.Println("unable to read file", jobJSONFile, err.Error())
				os.Exit(1)
			}
		}

		bodybytes, status, curlcmd, err := ce.CreateJob(profilemap["base"], profilemap["auth"], filebytes)

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

		fmt.Printf("%s\n", bodybytes)

	},
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

		if err != nil {
			fmt.Println("error %s", err)
			os.Exit(1)
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

	createJobCmd.PersistentFlags().StringVar(&jobJSONFile, "file", "", "job configuration json file")
	createJobCmd.MarkFlagRequired("file")
	jobsCmd.AddCommand(createJobCmd)

	deleteJobCmd.PersistentFlags().BoolVar(&jobsDeleteAll, "all", false, "delete all jobs")
	jobsCmd.AddCommand(deleteJobCmd)
}

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
var maxConcurrentDeletes int

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

		// if --all, then do that
		if jobsDeleteAll == true {
			err := deleteAllJobs(profilemap["base"], profilemap["auth"])
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
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

func deleteAllJobs(base, auth string) error {
	// Get all jobs
	bodybytes, status, _, err := ce.ListJobs(base, auth)
	if err != nil {
		fmt.Println("Unable to list jobs", err.Error())
		return err
	}
	if status != 200 {
		fmt.Println("non-200 response", status)
		return err
	}
	var jobs []ce.Job
	err = json.Unmarshal(bodybytes, &jobs)
	if err != nil {
		fmt.Printf("Response not a list of Jobs, %s", err.Error())
		return err
	}
	if len(jobs) == 0 {
		fmt.Print("No jobs to delete\n")
		return nil
	}

	if maxConcurrentDeletes < 1 { // guard against 0
		maxConcurrentDeletes = 1
	}
	q := make(chan int)     // queue of tasks
	done := make(chan bool) // result of task

	for i := 0; i < maxConcurrentDeletes; i++ {
		go deleteJobWorker(q, i, done)
	}

	for j := 0; j < len(jobs); j++ {
		go func(j int) {
			err := deleteJob(base, auth, jobs[j].ID)
			if err != nil {
				log.Println(err)
			}
			q <- j
		}(j)
	}

	for c := 0; c < len(jobs); c++ {
		<-done
	}

	return nil
	/*
		// delete all jobs
		results := make(chan DeleteJobCheck)
		var unabletodelete []string
		for _, v := range jobs {
			go deleteJob(base, auth, v.ID, results)
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

		// print results
		fmt.Printf("%v/%v 200\n", len(jobs)-len(unabletodelete), len(jobs))
		return nil
	*/
}

func deleteJobWorker(queue chan int, worknumber int, done chan bool) {
	for true {
		select {
		case k := <-queue:
			// doing work k worknumber
			//log.Printf("work! %v for worknumber %v", k, worknumber)
			//log.Println(k, worknumber)
			_ = k
			done <- true
		}
	}
}

func deleteJob(base, auth, jobID string) error {
	_, status, _, err := ce.DeleteJob(base, auth, jobID)
	if err != nil {
		log.Println(jobID, "failed")
		return err
	}
	if status != 200 {
		log.Println(jobID, status)
		return fmt.Errorf("Non-200 error %v", status)
	}
	log.Println(jobID, "deleted")
	return nil
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
	deleteJobCmd.PersistentFlags().IntVar(&maxConcurrentDeletes, "max", 4, "max concurrent delete calls")
	jobsCmd.AddCommand(deleteJobCmd)
}

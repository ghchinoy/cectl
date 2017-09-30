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

	"github.com/ghchinoy/cectl/ce"
	"github.com/moul/http2curl"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// formulasCmd represents the formulas command
var formulasCmd = &cobra.Command{
	Use:   "formulas",
	Short: "Manage formulas on the platform",
	Long:  `Currently allows listing formulas and execution details`,
}

// listFormulasCmd represents the listFormulas command
var listFormulasCmd = &cobra.Command{
	Use:   "list",
	Short: "lists available formulas on the plaform",
	Long:  `returns a list of formulas on the platform`,
	Run: func(cmd *cobra.Command, args []string) {

		if !viper.IsSet(profile + ".base") {
			fmt.Println("Can't find info for profile", profile)
			os.Exit(1)
		}

		base := viper.Get(profile + ".base")
		user := viper.Get(profile + ".user")
		org := viper.Get(profile + ".org")

		url := fmt.Sprintf("%s%s", base, ce.FormulasURI)
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

		if showCurl {
			curlcmd, _ := http2curl.GetCurlCommand(req)
			log.Println(curlcmd)
		}

		if outputJSON {
			fmt.Printf("%s\n", bodybytes)
			return
		}

		if resp.StatusCode != 200 {
			fmt.Print(resp.Status)
			if resp.StatusCode == 404 {
				fmt.Printf("Unable to contact CE API, %s\n", url)
				return
			}
			fmt.Println()
		}

		data := [][]string{}

		var formulas []ce.Formula
		err = json.Unmarshal(bodybytes, &formulas)
		for _, v := range formulas {
			if len(v.Triggers) < 1 {
				fmt.Printf("Formula %v is malformed, no trigger present\n", v.ID)
				break
			}

			var instancecount string
			instances, err := ce.GetInstancesOfFormula(v.ID, base.(string), auth)
			if err != nil {
				// unable to retrieve instances of formula!
				instancecount = "N/A"
			}
			instancecount = strconv.Itoa(len(instances))

			for _, t := range v.Triggers {

				api := "N/A"
				if v.Triggers[0].Type == "manual" {
					api = v.API
				}

				data = append(data, []string{
					strconv.Itoa(v.ID),
					v.Name,
					strconv.FormatBool(v.Active),
					strconv.Itoa(len(v.Steps)),
					instancecount,
					t.Type,
					strconv.Itoa(t.ID),
					fmt.Sprintf("%s", t.OnSuccess),
					api,
				},
				)
			}
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "active", "steps", "instances", "trigger", "id", "success", "api"})
		table.SetBorder(false)
		table.SetAutoMergeCells(true)
		table.AppendBulk(data)
		table.Render()
	},
}

func init() {
	RootCmd.AddCommand(formulasCmd)
	formulasCmd.AddCommand(listFormulasCmd)

	formulasCmd.PersistentFlags().StringVar(&profile, "profile", "default", "profile name")
	formulasCmd.PersistentFlags().BoolVarP(&outputJSON, "json", "j", false, "output as json")
	formulasCmd.PersistentFlags().BoolVarP(&showCurl, "curl", "c", false, "show curl command")
}

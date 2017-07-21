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
	"strings"

	"github.com/ghchinoy/cectl/ce"
	"github.com/moul/http2curl"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var withInstances bool

// resourcesCmd represents the resources command
var resourcesCmd = &cobra.Command{
	Use:   "resources",
	Short: "Manage common resources",
	Long:  `List, add, remove and inspect common resource objects.`,
}

// listResourcesCmd represents the listResources command
var listResourcesCmd = &cobra.Command{
	Use:   "list",
	Short: "lists common object resources",
	Long:  `lists common object resources`,
	Run: func(cmd *cobra.Command, args []string) {
		if !viper.IsSet(profile + ".base") {
			fmt.Println("Can't find info for profile", profile)
			os.Exit(1)
		}

		base := viper.Get(profile + ".base")
		user := viper.Get(profile + ".user")
		org := viper.Get(profile + ".org")

		url := fmt.Sprintf("%s%s", base, "/common-resources")
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

		var commonResources []ce.CommonResource
		err = json.Unmarshal(bodybytes, &commonResources)
		if err != nil {
			fmt.Printf("Response not a list of Common Resources, %s", err.Error())
			return
		}

		for _, v := range commonResources {

			var fieldList string
			if len(v.Fields) > 0 {
				var fields []string
				for _, f := range v.Fields {
					fields = append(fields, f.Path)
				}
				fieldList = strings.Join(fields[:], ", ")
				fieldList = " [" + fieldList + "]"
			}

			var instanceList string
			if len(v.ElementInstanceIDs) > 0 {
				var ids []string
				for _, i := range v.ElementInstanceIDs {
					ids = append(ids, strconv.Itoa(i))
				}
				instanceList = strings.Join(ids[:], ", ")
				instanceList = " [" + instanceList + "]"
			}

			data = append(data, []string{
				v.Name,
				strconv.Itoa(len(v.ElementInstanceIDs)) + instanceList,
				strconv.Itoa(len(v.Fields)),
				fieldList,
			})
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Mapped Instances", "#", "Fields"})
		table.SetBorder(false)
		table.SetColWidth(40)
		table.AppendBulk(data)
		table.Render()

	},
}

func init() {
	RootCmd.AddCommand(resourcesCmd)

	resourcesCmd.PersistentFlags().StringVar(&profile, "profile", "default", "profile name")
	resourcesCmd.PersistentFlags().BoolVarP(&outputJSON, "json", "j", false, "output as json")
	resourcesCmd.PersistentFlags().BoolVarP(&showCurl, "curl", "c", false, "show curl command")

	resourcesCmd.AddCommand(listResourcesCmd)
	listResourcesCmd.Flags().BoolVarP(&withInstances, "with-instances", "i", false, "show mapped instances")
}

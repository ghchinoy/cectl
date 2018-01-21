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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/ghchinoy/ce-go/ce"
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
		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		bodybytes, statuscode, curlcmd, err := ce.ResourcesList(profilemap["base"], profilemap["auth"])
		if err != nil {
			if statuscode == -1 {
				fmt.Println("Unable to reach CE API. Please check your configuration / profile.")
			}
			fmt.Println(err)
			os.Exit(1)
		}
		// handle global options, curl
		if showCurl {
			log.Println(curlcmd)
		}
		// handle non 200
		if statuscode != 200 {
			log.Printf("HTTP Error: %v\n", statuscode)
			// handle this nicely, show error description
		}
		// handle global options, json
		if outputJSON {
			fmt.Printf("%s\n", bodybytes)
			return
		}
		err = ce.OutputResourcesList(bodybytes)
		if err != nil {
			fmt.Println("Unable to render resources", err.Error())
		}
		/*

			if !viper.IsSet(profile + ".base") {
				fmt.Println("Can't find info for profile", profile)
				os.Exit(1)
			}

			base := viper.Get(profile + ".base")
			user := viper.Get(profile + ".user")
			org := viper.Get(profile + ".org")

			url := fmt.Sprintf("%s%s", base, ce.CommonResourcesURI)
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
		*/

	},
}

var addResourceCmd = &cobra.Command{
	Use:   "add <name> <filepath.json>",
	Short: "add a common resource",
	Long:  "Add a Common Resource to the platform, given a name and a json definition of that Common Resource",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("must supply a Common Resource name and a filepath to a json definition")
			os.Exit(1)
		}

		// read in file
		filebytes, err := ioutil.ReadFile(args[1])
		if err != nil {
			fmt.Println("unable to read file", args[1], err.Error())
			os.Exit(1)
		}
		// Check if can decode into formula struct
		var cro ce.CommonResource
		err = json.Unmarshal(filebytes, &cro)
		if err != nil {
			fmt.Println(args[1], "doesn't seem like a Common Resource Object")
			os.Exit(1)
		}

		base := viper.Get(profile + ".base")
		user := viper.Get(profile + ".user")
		org := viper.Get(profile + ".org")

		url := fmt.Sprintf("%s%s",
			base,
			fmt.Sprintf(ce.CommonResourceDefinitionsFormatURI, args[0]),
		)
		auth := fmt.Sprintf("User %s, Organization %s", user, org)

		client := &http.Client{}
		req, err := http.NewRequest("POST", url, bytes.NewReader(filebytes))
		if err != nil {
			fmt.Println("Can't construct request", err.Error())
			os.Exit(1)
		}
		req.Header.Add("Authorization", auth)
		req.Header.Add("Content-Type", "application/json")
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
			if resp.StatusCode == 409 {
				var errorMessage struct {
					RequestID string `json:"requestId"`
					Message   string `json:"message"`
				}
				err := json.Unmarshal(bodybytes, &errorMessage)
				if err != nil {
					fmt.Printf("%s\n", bodybytes)
					return
				}
				fmt.Printf("\n%s\nrequest: %s\n", errorMessage.Message, errorMessage.RequestID)
				return
			}
			fmt.Println()
		}

		err = json.Unmarshal(bodybytes, &cro)
		if err != nil {
			fmt.Println("Unable to convert 200 reponse into a Common Resource Object")
			os.Exit(1)
		}

		fmt.Printf("Added common object resource: %s\n", args[0])
		fmt.Printf("Field levels: %s\n", cro.Level)
		data := [][]string{}
		for _, v := range cro.Fields {
			data = append(data, []string{
				v.Path,
				v.Type,
			})
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Path", "Type"})
		table.SetBorder(false)
		table.SetColWidth(40)
		table.AppendBulk(data)
		table.Render()

	},
}

var defineResourceCmd = &cobra.Command{
	Use:   "definition",
	Short: "detailed definition of a resource",
	Long:  "Display the definition a named common object resource",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			fmt.Println("must supply a name of a Common Resource")
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
			fmt.Sprintf(ce.CommonResourceDefinitionsFormatURI, args[0]),
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

		var cro ce.CommonResource
		err = json.Unmarshal(bodybytes, &cro)
		if err != nil {
			fmt.Println("Unable to convert 200 reponse into a Common Resource Object")
			os.Exit(1)
		}

		fmt.Printf("Common Resource Object: %s\n", args[0])
		fmt.Printf("Field levels: %s\n", cro.Level)
		data := [][]string{}
		for _, v := range cro.Fields {
			data = append(data, []string{
				v.Path,
				v.Type,
			})
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Path", "Type"})
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
	resourcesCmd.AddCommand(defineResourceCmd)
	resourcesCmd.AddCommand(addResourceCmd)
}

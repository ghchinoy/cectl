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

	"github.com/ghchinoy/cectl/ce"
	"github.com/spf13/cobra"
)

// infoCmd is a command to show information about the CE account in question
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Information about this account",
	Long:  `Provide info on this account`,
	Run: func(cmd *cobra.Command, args []string) {
		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		all := false
		if len(args) > 0 {
			if args[0] == "all" {
				all = true
			}
		}

		var allCurlCommands []string

		// List formulas
		bodybytes, statuscode, curlcmd, err := ce.FormulasList(profilemap["base"], profilemap["auth"])

		// handle global options, curl
		allCurlCommands = append(allCurlCommands, curlcmd)
		// handle non 200
		if statuscode != 200 {
			log.Printf("HTTP Error: %v\n", statuscode)
			// handle this nicely, show error description
		}

		formulas, err := ce.CombinedFormulaAndInstances(bodybytes, profilemap["base"], profilemap["auth"])
		if err != nil {
			fmt.Println("Unable to obtain instances for formulas", err.Error())
		}
		fmt.Printf("Formulas: %v\n", len(formulas))
		for _, v := range formulas {
			fmt.Printf("%6v %2v %s\n", v.ID, len(v.Instances), v.Name)
		}
		if all {
			err = ce.OutputFormulasList(bodybytes, profilemap["base"], profilemap["auth"])
			if err != nil {
				fmt.Println("Unable to render formula table", err.Error())
			}
		}

		// List Custom Elements
		// Get elements
		bodybytes, statuscode, curlcmd, err = ce.GetAllElements(profilemap["base"], profilemap["auth"])
		if err != nil {
			if statuscode == -1 {
				fmt.Println("Unable to reach CE API. Please check your configuration / profile.")
			}
			fmt.Println(err)
			os.Exit(1)
		}
		allCurlCommands = append(allCurlCommands, curlcmd)
		customElementsOnly, err := ce.FilterCustomElements(bodybytes)
		if err != nil {
			fmt.Println("Error filtering custom elements", err.Error())
		}
		var customElements ce.Elements
		err = json.Unmarshal(customElementsOnly, &customElements)
		if err != nil {
			fmt.Println("Not an elements list", err.Error())
		}
		fmt.Println()
		fmt.Printf("Custom Elements: %v\n", len(customElements))
		if len(customElements) > 0 {
			ce.OutputElementsTable(customElementsOnly, "", "")
		}

		// List Instances
		bodybytes, statuscode, curlcmd, err = ce.GetAllInstances(profilemap["base"], profilemap["auth"])
		if err != nil {
			if statuscode == -1 {
				fmt.Println("Unable to reach CE API. Please check your configuration / profile.")
			}
			fmt.Println(err)
			os.Exit(1)
		}
		allCurlCommands = append(allCurlCommands, curlcmd)
		if statuscode != 200 {
			log.Printf("HTTP Error: %v\n", statuscode)
			// handle this nicely, show error description
		}
		fmt.Println()
		var instances []ce.ElementInstance
		err = json.Unmarshal(bodybytes, &instances)
		if err != nil {
			fmt.Println("Unable to read Element instances")
		}
		fmt.Printf("Element Instances: %v\n", len(instances))
		for _, v := range instances {
			fmt.Printf("%7v %14s %s\n", v.ID, v.Element.Key, v.Name)
		}

		// List Common Resource Objects

		bodybytes, statuscode, curlcmd, err = ce.ResourcesList(profilemap["base"], profilemap["auth"])
		if err != nil {
			if statuscode == -1 {
				fmt.Println("Unable to reach CE API. Please check your configuration / profile.")
			}
			fmt.Println(err)
			os.Exit(1)
		}
		allCurlCommands = append(allCurlCommands, curlcmd)
		fmt.Println()
		fmt.Println("Common Resource Objects")
		err = ce.OutputResourcesList(bodybytes)
		if err != nil {
			fmt.Println("Unable to render resources", err.Error())
		}

		// List Users
		bodybytes, status, curlcmd, err := ce.GetAllUsers(profilemap["base"], profilemap["auth"])
		if err != nil {
			fmt.Println("Unable to even", status)
			return
		}
		allCurlCommands = append(allCurlCommands, curlcmd)
		if status != 200 {
			fmt.Print(status)
			if status == 404 {
				fmt.Println("Unable to contact CE API")
				return
			}
			fmt.Println()
		}
		fmt.Println()
		fmt.Println("Users")
		err = ce.FormatUserList(bodybytes)
		if err != nil {
			fmt.Println("Unable to format")
		}

		// handle global options, curl
		if showCurl {
			fmt.Println()
			for _, v := range allCurlCommands {
				log.Println(v)
			}
		}

	},
}

func init() {
	RootCmd.AddCommand(infoCmd)

	infoCmd.PersistentFlags().StringVar(&profile, "profile", "default", "profile name")
	//infoCmd.PersistentFlags().BoolVarP(&outputJSON, "json", "j", false, "output as json")
	infoCmd.PersistentFlags().BoolVarP(&showCurl, "curl", "c", false, "show curl command")

	//elementsCmd.AddCommand(anotherInfoCmd)

}

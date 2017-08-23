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
	"strconv"

	"github.com/ghchinoy/cectl/ce"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var orderBy, filterBy string

// elementsCmd represents the elements command
var elementsCmd = &cobra.Command{
	Use:   "elements",
	Short: "Manage Elements on the Platform",
	Long:  `Manage Elements on the Platform`,
}

// listElementsCmd represents the /elements API
var listElementsCmd = &cobra.Command{
	Use:   "list",
	Short: "List Elements on the Platform",
	Long: `Retrieve all available Elements;
Optionally, add in a keyfilter to filter out Elements by key.`,
	Run: func(cmd *cobra.Command, args []string) {
		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// Get elements
		bodybytes, statuscode, curlcmd, err := ce.GetAllElements(profilemap["base"], profilemap["auth"])
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
		// optional element key filter
		if args[0] != "" {
			filteredElementsBytes, err := ce.FilterElementFromList(args[0], bodybytes)
			if err != nil {
				log.Printf("Unable to filter by '%s'- %s\n", args[0], err.Error())
			}
			bodybytes = filteredElementsBytes
		}
		// handle global options, json
		if outputJSON {
			fmt.Printf("%s\n", bodybytes)
			return
		}
		// output
		ce.OutputElementsTable(bodybytes, orderBy, filterBy)
	},
}

// elementDocsCmd represents the /elements/{id}/docs API
var elementDocsCmd = &cobra.Command{
	Use:   "docs <id|key>",
	Short: "Output the OAI Specification documentation for the Element",
	Long:  `Outputs the JSON format of the OAI Specification for the indicated Element`,
	Run: func(cmd *cobra.Command, args []string) {
		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// check for Element ID
		if len(args) < 1 {
			fmt.Println("Please provide an Element ID or Element Key")
			return
		}

		elementid, err := ce.ElementKeyToID(args[0], profilemap)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// Get element OAI
		bodybytes, statuscode, curlcmd, err := ce.GetElementOAI(profilemap["base"], profilemap["auth"], strconv.Itoa(elementid))
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
		fmt.Printf("%s", bodybytes)
	},
}

var elementInstancesCmd = &cobra.Command{
	Use:   "instances",
	Short: "List the Instances of an Element",
	Long:  `List the Instances associated with an Element`,
	Run: func(cmd *cobra.Command, args []string) {
		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// check for Element ID
		if len(args) < 1 {
			// list all instances
			listInstancesCmd.Run(cmd, args)
			return
		}

		// Get element instances
		bodybytes, statuscode, curlcmd, err := ce.GetElementInstances(profilemap["base"], profilemap["auth"], args[0])
		if err != nil {
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
			var errmsg struct {
				RequestID string `json:"requestId"`
				Message   string `json:"message"`
			}
			if statuscode == 404 {
				err := json.Unmarshal(bodybytes, &errmsg)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Printf("%s\nRequest ID: %s\n", errmsg.Message, errmsg.RequestID)
				return
			}
		}
		// handle global options, json
		if outputJSON {
			fmt.Printf("%s\n", bodybytes)
			return
		}
		// output
		err = ce.OutputElementInstancesTable(bodybytes)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

	},
}

var elementModelValidation = &cobra.Command{
	Use:   "validate-models",
	Short: "Validate the models of this Element",
	Long: `Validates the models associated for this Element,
primarily for use for internal housekeeping. Requires a model ID.`,
	Run: func(cmd *cobra.Command, args []string) {
		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// check for Element ID
		if len(args) < 1 {
			fmt.Println("Please provide an Element ID or Element Key")
			return
		}

		elementid, err := ce.ElementKeyToID(args[0], profilemap)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// Get element model validation
		bodybytes, statuscode, curlcmd, err := ce.GetElementModelValidation(profilemap["base"], profilemap["auth"], strconv.Itoa(elementid))
		if err != nil {
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

		fmt.Printf("%s", bodybytes)
	},
}

// elementMetadataCmd provides the metadata for the Element
var elementMetadataCmd = &cobra.Command{
	Use:   "metadata",
	Short: "Display Metadata of an Element",
	Long:  `Display Metadata of an Element`,
	Run: func(cmd *cobra.Command, args []string) {
		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// check for Element ID
		if len(args) < 1 {
			fmt.Println("Please provide an Element ID or Element Key")
			return
		}

		elementid, err := ce.ElementKeyToID(args[0], profilemap)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// Get element metadata
		bodybytes, statuscode, curlcmd, err := ce.GetElementMetadata(profilemap["base"], profilemap["auth"], strconv.Itoa(elementid))
		if err != nil {
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
		var metadata interface{}
		err = json.Unmarshal(bodybytes, &metadata)
		if err != nil {
			fmt.Println("Can't unmarshal")
			os.Exit(1)
		}
		formattedbytes, err := json.MarshalIndent(metadata, "", "    ")
		if err != nil {
			fmt.Println("Can't format json")
			os.Exit(1)
		}
		fmt.Printf("%s", formattedbytes)
	},
}

func getAuth(profile string) (map[string]string, error) {

	profilemap := make(map[string]string)

	if !viper.IsSet(profile + ".base") {
		return profilemap, fmt.Errorf("can't find profile")
	}

	profilemap["base"] = viper.Get(profile + ".base").(string)
	profilemap["user"] = viper.Get(profile + ".user").(string)
	profilemap["org"] = viper.Get(profile + ".org").(string)
	profilemap["auth"] = fmt.Sprintf("User %s, Organization %s", profilemap["user"], profilemap["org"])

	return profilemap, nil
}

func init() {
	RootCmd.AddCommand(elementsCmd)

	elementsCmd.PersistentFlags().StringVar(&profile, "profile", "default", "profile name")
	elementsCmd.PersistentFlags().BoolVarP(&outputJSON, "json", "j", false, "output as json")
	elementsCmd.PersistentFlags().BoolVarP(&showCurl, "curl", "c", false, "show curl command")

	elementsCmd.AddCommand(listElementsCmd)
	elementsCmd.AddCommand(elementMetadataCmd)
	elementsCmd.AddCommand(elementDocsCmd)
	elementsCmd.AddCommand(elementInstancesCmd)
	//elementsCmd.AddCommand(elementModelValidation)

	// order-by flag: Order element list by
	// --order-by key|name|id|hub
	listElementsCmd.Flags().StringVarP(&orderBy, "order", "", "", "order element (hub, name)")
	// filter-by flag: Show only elements where filter is true
	// --filter-by active|deleted|private|beta|cloneable|extendable
	//listElementsCmd.Flags().StringVarP(&filterBy, "filter", "", "", "elements where filter is true")

}

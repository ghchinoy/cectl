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

	"github.com/ghchinoy/cectl/ce"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
	Long:  `Retrieve all available Elements`,
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
			fmt.Println(err)
			os.Exit(1)
		}
		// handle global options, curl
		if showCurl {
			log.Println(curlcmd)
		}
		// handle non 200
		if statuscode != 200 {
			fmt.Printf("HTTP Error: %v\n", statuscode)
			// handle this nicely, show
		}
		// handle global options, json
		if outputJSON {
			fmt.Printf("%s\n", bodybytes)
			return
		}
		// output
		ce.OutputElementsTable(bodybytes)
	},
}

// elementDocsCmd represents the /elements/{id}/docs API
var elementDocsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Output the OAI Specification documentation for the Element",
	Long:  `Outputs the JSON format of the OAI Specification for the indicated Element`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var elementInstancesCmd = &cobra.Command{
	Use:   "instances",
	Short: "List the Instances of an Element",
	Long:  `List the Instances associated with an Element`,
	Run: func(cmd *cobra.Command, args []string) {

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
	elementsCmd.AddCommand(listElementsCmd)
	elementsCmd.AddCommand(elementKeysCmd)
	elementsCmd.AddCommand(elementDocsCmd)
	elementsCmd.AddCommand(elementInstancesCmd)

	elementsCmd.PersistentFlags().StringVar(&profile, "profile", "default", "profile name")
	elementsCmd.PersistentFlags().BoolVarP(&outputJSON, "json", "j", false, "output as json")
	elementsCmd.PersistentFlags().BoolVarP(&showCurl, "curl", "c", false, "show curl command")

}

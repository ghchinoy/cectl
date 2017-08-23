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

	"github.com/ghchinoy/cectl/ce"
	"github.com/spf13/cobra"
)

// instancesCmd represents the root command for managing Element Instances
var instancesCmd = &cobra.Command{
	Use:   "instances",
	Short: "Manage Instances of Elements on the Platform",
	Long:  `Manage Element Instances on the Platform`,
}

var listInstancesCmd = &cobra.Command{
	Use:   "list",
	Short: "List Instances on Platform",
	Long:  "List Element Instances on Platform",
	Run: func(cmd *cobra.Command, args []string) {
		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// Get elements
		bodybytes, statuscode, curlcmd, err := ce.GetAllInstances(profilemap["base"], profilemap["auth"])
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
		// output
		ce.OutputElementInstancesTable(bodybytes)
	},
}

var listInstanceTransformationsCmd = &cobra.Command{
	Use:   "transformations",
	Short: "Show the transformations mapped to an Instance",
	Long:  "Show the Transformations associated with this Element Instance",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Unimplemented: Transformations per Instance")
	},
}

var instanceDocsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Output the OAI Specification documentation for the Element Instance",
	Long:  `Outputs the JSON format of the OAI Specification for the indicated Element Instance`,
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
		bodybytes, statuscode, curlcmd, err := ce.GetInstanceOAI(profilemap["base"], profilemap["auth"], strconv.Itoa(elementid))
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

var instanceDetailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Details about an Instance",
	Long:  `Provides details about an Instance, given an Instance ID`,
	Run: func(cmd *cobra.Command, args []string) {
		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// check for Instance ID & Operation name
		if len(args) < 1 {
			fmt.Println("Please provide an Instance ID ")
			return
		}
		if _, err := strconv.ParseInt(args[0], 10, 64); err != nil {
			fmt.Println("Please provide an Instance ID that is an integer")
			return
		}

		// Get schema definition for operation
		bodybytes, statuscode, curlcmd, err := ce.GetInstanceInfo(profilemap["base"], profilemap["auth"], args[0])
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
		// output
		//ce.OutputElementInstancesTable(bodybytes)
		ce.OutputInstanceDetails(bodybytes)
	},
}

var instanceOperationDefinitionCmd = &cobra.Command{
	Use:   "operation",
	Short: "Show operation schema definition",
	Long: `Shows the schema definition for an operation.
Provide an instance ID and an operation name to retrieve the
associated schema.`,
	Run: func(cmd *cobra.Command, args []string) {
		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// check for Instance ID & Operation name
		if len(args) < 2 {
			fmt.Println("Please provide an Instance ID and Operation name")
			return
		}
		if _, err := strconv.ParseInt(args[0], 10, 64); err != nil {
			fmt.Println("Please provide an Instance ID that is an integer")
			return
		}

		// Get schema definition for operation
		bodybytes, statuscode, curlcmd, err := ce.GetInstanceOperationDefinition(profilemap["base"], profilemap["auth"], args[0], args[1])
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
		// output
		ce.OutputElementInstancesTable(bodybytes)
	},
}

func init() {
	RootCmd.AddCommand(instancesCmd)
	instancesCmd.AddCommand(listInstancesCmd)
	instancesCmd.AddCommand(listInstanceTransformationsCmd)
	instancesCmd.AddCommand(instanceDocsCmd)
	instancesCmd.AddCommand(instanceDetailsCmd)

	instancesCmd.PersistentFlags().StringVar(&profile, "profile", "default", "profile name")
	instancesCmd.PersistentFlags().BoolVarP(&outputJSON, "json", "j", false, "output as json")
	instancesCmd.PersistentFlags().BoolVarP(&showCurl, "curl", "c", false, "show curl command")
}

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
	"strconv"

	"github.com/ghchinoy/ce-go/ce"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// instancesCmd represents the root command for managing Element Instances
var instancesCmd = &cobra.Command{
	Use:   "instances",
	Short: "Manage Instances of Elements on the Platform",
	Long:  `Manage Element Instances on the Platform`,
}

// testInstanceCmd tests all instances by hitting the /ping endpoint
var testInstancesCmd = &cobra.Command{
	Use:   "test",
	Short: "Test Element Instances",
	Long:  "Tests Element Instances by hitting /ping endpoint and reporting results",
	Run: func(cmd *cobra.Command, args []string) {
		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// Get instances
		bodybytes, statuscode, _, err := ce.GetAllInstances(profilemap["base"], profilemap["auth"])
		if err != nil {
			if statuscode == -1 {
				fmt.Println("Unable to reach CE API. Please check your configuration / profile.")
			}
			fmt.Println(err)
			os.Exit(1)
		}
		var instances []ce.ElementInstance
		err = json.Unmarshal(bodybytes, &instances)
		if err != nil { // can't umarshal into Instance objects
			fmt.Println(err)
			os.Exit(1)
		}

		// start a thread to check each Instance, aggregating to "results" channel
		results := make(chan PingCheck)

		for _, i := range instances {
			pingurl := fmt.Sprintf("%s/hubs/%s/ping", profilemap["base"], i.Element.Hub)
			ceauthtoken := fmt.Sprintf("%s, Element %s", profilemap["auth"], i.Token)
			// start a goroutine to /ping Instance
			go pingInstance(i.Element.Name, i.Name, i.ID, pingurl, ceauthtoken, results)
		}

		// as results come in, print out if necessary
		var num, bad int // keep track of how many have come in
		fmt.Printf("Checking %v instances\n", len(instances))
		for i := range results {
			if i.StatusCode != 200 {
				fmt.Printf("%5v %s (%s) %s\n", i.InstanceID, i.ElementName, i.InstanceName, i.Status)
				bad++
			}
			num++
			// if all expected results are in, close out the channel
			if len(instances) == num {
				close(results)
			}
		}
		fmt.Printf("%v/%v 200", len(instances)-bad, len(instances))
	},
}

// PingCheck is a struct to hold results of a /ping check on an Element Instance
type PingCheck struct {
	ElementName  string
	InstanceName string
	StatusCode   int
	Status       string
	Message      []byte
	InstanceID   int
}

// pingInstance makes an HTTP call to the Instances /ping endpoint
func pingInstance(elementName, instanceName string, instanceID int, url, auth string, checks chan PingCheck) (PingCheck, error) {

	var c PingCheck

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Can't construct request", err.Error())
		os.Exit(1)
	}
	req.Header.Add("Authorization", auth)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		// unable to reach CE API
		checks <- c
		return c, err
	}
	bodybytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	c = PingCheck{ElementName: elementName, InstanceName: instanceName, StatusCode: resp.StatusCode, Status: resp.Status, Message: bodybytes, InstanceID: instanceID}
	checks <- c
	return c, nil

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
		// Get instances
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
	Use:   "transformations <id>",
	Short: "Show the transformations mapped to an Instance",
	Long:  "Show the Transformations associated with this Element Instance",
	Run: func(cmd *cobra.Command, args []string) {
		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if len(args) == 0 {
			fmt.Println("ID must be provided")
			cmd.Help()
			os.Exit(1)
		}
		if _, err := strconv.Atoi(args[0]); err != nil {
			fmt.Println("ID must be an integer")
			cmd.Help()
			os.Exit(1)
		}
		bodybytes, statuscode, curlcmd, err := ce.GetInstanceTransformations(profilemap["base"], profilemap["auth"], args[0])
		// handle global options, curl
		if showCurl {
			log.Println(curlcmd)
		}
		// handle non 200
		if statuscode != 200 {
			log.Printf("No Transformations for %s\n", args[0])
			os.Exit(0)
			// handle this nicely, show error description
		}
		// handle global options, json
		if outputJSON {
			fmt.Printf("%s\n", bodybytes)
			os.Exit(1)
		}

		txs := make(map[string]ce.Transformation)
		err = json.Unmarshal(bodybytes, &txs)
		if err != nil {
			os.Exit(1)
		}

		data := [][]string{}
		for k, v := range txs {
			data = append(data, []string{
				k,
				v.VendorName,
				fmt.Sprintf("%v", len(v.Fields)),
				fmt.Sprintf("%v", len(v.Configuration)),
				fmt.Sprintf("%v", v.IsLegacy),
				v.StartDate,
			})
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"resource", "vendor", "# Fields", "# Configs", "Legacy", "Start Date"})
		table.SetBorder(false)
		table.AppendBulk(data)
		table.Render()

	},
}

var instanceDocsCmd = &cobra.Command{
	Use:   "docs [ID]",
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
			fmt.Println("Please provide an Element Instance ID")
			return
		}

		instanceid, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Element ID must be an integer")
			return
		}

		// Get element OAI
		bodybytes, statuscode, curlcmd, err := ce.GetInstanceOAI(profilemap["base"], profilemap["auth"], strconv.Itoa(instanceid))
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

var deleteElementInstanceCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an Element Instance",
	Long:  "Delete an Element Instance",
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
		bodybytes, statuscode, curlcmd, err := ce.DeleteElementInstance(profilemap["base"], profilemap["auth"], args[0])
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
		fmt.Printf("Deleted Element Instance %s", args[0])
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

var instanceDefinitionsCmd = &cobra.Command{
	Use:   "definitions",
	Short: "Show all definitions for this Instance",
	Long: `Retrieve all of the object definitions within a specific instance. 
If no object definitions exist, then this will result in an error response.`,
	Run: func(cmd *cobra.Command, args []string) {
		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// check for Instance ID & Operation name
		if len(args) < 1 {
			fmt.Println("Please provide an Instance ID")
			return
		}
		if _, err := strconv.ParseInt(args[0], 10, 64); err != nil {
			fmt.Println("Please provide an Instance ID that is an integer")
			return
		}

		// Get schema definition for operation
		bodybytes, statuscode, curlcmd, err := ce.GetInstanceObjectDefinitions(profilemap["base"], profilemap["auth"], args[0])
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
		//if outputJSON {
		var pretty bytes.Buffer
		err = json.Indent(&pretty, bodybytes, "", "  ")
		if err != nil {
			fmt.Printf("%s\n", bodybytes)
			return
		}
		fmt.Printf("%s\n", pretty.Bytes())

		//	return
		//}
	},
}

var instanceOperationDefinitionCmd = &cobra.Command{
	Use:   "operation [ID] [operationName]",
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
		/*
			if outputJSON {
				fmt.Printf("%s\n", bodybytes)
				return
			}

			// output
			ce.OutputElementInstancesTable(bodybytes)
		*/
		var pretty bytes.Buffer
		err = json.Indent(&pretty, bodybytes, "", "  ")
		if err != nil {
			fmt.Printf("%s\n", bodybytes)
			return
		}
		fmt.Printf("%s\n", pretty.Bytes())
	},
}

func init() {
	RootCmd.AddCommand(instancesCmd)
	instancesCmd.AddCommand(listInstancesCmd)
	instancesCmd.AddCommand(listInstanceTransformationsCmd)
	instancesCmd.AddCommand(instanceDocsCmd)
	instancesCmd.AddCommand(instanceDetailsCmd)
	instancesCmd.AddCommand(instanceOperationDefinitionCmd)
	instancesCmd.AddCommand(instanceDefinitionsCmd)
	instancesCmd.AddCommand(testInstancesCmd)
	instancesCmd.AddCommand(deleteElementInstanceCmd)

	instancesCmd.PersistentFlags().StringVar(&profile, "profile", "default", "profile name")
	instancesCmd.PersistentFlags().BoolVarP(&outputJSON, "json", "j", false, "output as json")
	instancesCmd.PersistentFlags().BoolVarP(&showCurl, "curl", "c", false, "show curl command")
}

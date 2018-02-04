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
	"sort"
	"strconv"

	"github.com/ghchinoy/ce-go/ce"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// elementsCmd represents the elements command
var transformationsCmd = &cobra.Command{
	Use:   "transformations",
	Short: "Manage Transformations on the Platform",
	Long:  `Manage Transformations on the Platform`,
}

// associateTransformationCmd adds a Transformation to an Element, given a Transformation JSON file
// This isn't ready - a Transformation requires a vendorName otherwise an added Transformation
// may not map to an Element's
var associateTransformationCmd = &cobra.Command{
	Use:    "associate <element_key | element_id> <transformation.json> [name]",
	Short:  "Associate a Transformation with an Element",
	Long:   "Associate a Transformation with an Element given a Transformation JSON file path",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if len(args) < 2 {
			fmt.Println("Please provide both an Element key|id and a path to a Transformation JSON file")
			os.Exit(1)
		}
		// validate Element ID
		elementid, err := ce.ElementKeyToID(args[0], profilemap)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		// validate Transformation json file
		var transformation ce.Transformation
		txbytes, err := ioutil.ReadFile(args[1])
		if err != nil {
			fmt.Println("Supplied file cannot be read", err.Error())
			os.Exit(1)
		}
		err = json.Unmarshal(txbytes, &transformation)
		if err != nil {
			fmt.Println("Supplied file does not contain a Transformation", err.Error())
			os.Exit(1)
		}
		// Provide a name for the object if supplied
		if len(args) == 3 {
			transformation.ObjectName = args[2]
		}

		bodybytes, status, curlcmd, err := ce.AssociateTransformationWithElement(
			profilemap["base"], profilemap["auth"],
			strconv.Itoa(elementid),
			transformation)
		if err != nil {
			fmt.Println("Unable to import Transformation", err.Error())
			os.Exit(1)
		}
		// handle global options, curl
		if showCurl {
			log.Println(curlcmd)
		}
		if status != 200 {
			fmt.Println("Non-200 status: ", status)
			var message interface{}
			json.Unmarshal(bodybytes, &message)
			fmt.Printf("%s\n", message)
			os.Exit(1)
		}
		fmt.Printf("%s\n", bodybytes)

	},
}

var withElementAssociations bool

// listTransformationsCmd is the command to list Transformations
// the flag --with-elements will also list the Elements the Transformation has associations with
var listTransformationsCmd = &cobra.Command{
	Use:   "list",
	Short: "List Transformations",
	Long:  "List Transformations on the Platform",
	Run: func(cmd *cobra.Command, args []string) {
		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		bodybytes, statuscode, curlcmd, err := ce.GetTransformations(profilemap["base"], profilemap["auth"])
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
		if outputJSON {
			// todo uplift to output package, output/FormattedJSON
			var transformations interface{}
			err = json.Unmarshal(bodybytes, &transformations)
			if err != nil {
				fmt.Println("Can't unmarshal")
				os.Exit(1)
			}
			formattedbytes, err := json.MarshalIndent(transformations, "", "    ")
			if err != nil {
				fmt.Println("Can't format json")
				os.Exit(1)
			}
			fmt.Printf("%s", formattedbytes)
			os.Exit(0)
		}

		txs := make(map[string]ce.Transformation)
		err = json.Unmarshal(bodybytes, &txs)
		if err != nil {
			fmt.Println("Unable to parse Transformations", err.Error())
			os.Exit(1)
		}
		// sort by key
		var keys []string
		for k := range txs {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		elementAssociations := make(map[string][]string)
		if withElementAssociations {
			for _, k := range keys {
				bodybytes, status, _, err := ce.GetTransformationAssocation(profilemap["base"], profilemap["auth"], k)
				if err != nil {
					break
				}
				if status != 200 {
					break
				}
				var associations []ce.AccountElement
				err = json.Unmarshal(bodybytes, &associations)
				if err != nil {
					break
				}
				var elements []string
				for _, e := range associations {
					elements = append(elements, e.Element.Key)
				}
				elementAssociations[k] = elements
			}
		}

		data := [][]string{}
		for _, k := range keys {
			v := txs[k]
			if withElementAssociations {
				data = append(data, []string{
					k,
					v.Level,
					fmt.Sprintf("%v", len(v.Fields)),
					fmt.Sprintf("%v %s", len(elementAssociations[k]), elementAssociations[k]),
				})
			} else {
				data = append(data, []string{
					k,
					v.Level,
					fmt.Sprintf("%v", len(v.Fields)),
				})
			}
		}
		table := tablewriter.NewWriter(os.Stdout)
		//table.SetHeader([]string{"Resource", "Vendor", "Level", "# Fields", "# Configs", "Legacy", "Start Date"})
		if withElementAssociations {
			table.SetHeader([]string{"Resource", "Level", "# Fields", "Elements"})
		} else {
			table.SetHeader([]string{"Resource", "Level", "# Fields"})
		}
		table.SetBorder(false)
		table.AppendBulk(data)
		table.Render()
	},
}

func init() {
	RootCmd.AddCommand(transformationsCmd)

	transformationsCmd.PersistentFlags().StringVar(&profile, "profile", "default", "profile name")
	transformationsCmd.PersistentFlags().BoolVarP(&outputJSON, "json", "j", false, "output as json")
	transformationsCmd.PersistentFlags().BoolVarP(&showCurl, "curl", "c", false, "show curl command")
	//transformationsCmd.PersistentFlags().BoolVarP(&outputCSV, "csv", "", false, "output as CSV")
	transformationsCmd.AddCommand(listTransformationsCmd)
	listTransformationsCmd.PersistentFlags().BoolVarP(&withElementAssociations, "with-elements", "", false, "show Element associations")
	transformationsCmd.AddCommand(associateTransformationCmd)
}

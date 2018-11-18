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
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ghchinoy/ce-go/ce"
	"github.com/ghchinoy/cectl/output"
	"github.com/gorilla/mux"
	"github.com/olekukonko/tablewriter"
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

// transformationsForElementCmd is the command to list Transformations associated with an Element
var transformationsForElementCmd = &cobra.Command{
	Use:   "transformations <id|key>",
	Short: "Show Transformations for the Element",
	Long:  "Given the Element key or ID, show associated Transformations",
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
		bodybytes, statuscode, curlcmd, err := ce.GetTransformationsPerElement(profilemap["base"], profilemap["auth"], strconv.Itoa(elementid))
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
			log.Printf("%s", bodybytes)
			os.Exit(1)
		}
		if outputJSON {
			fmt.Printf("%s\n", bodybytes)
			return
		}
		txs := make(map[string]ce.Transformation)
		err = json.Unmarshal(bodybytes, &txs)
		if err != nil {
			fmt.Println("Can't parse JSON response", err.Error())
			os.Exit(1)
		}
		data := [][]string{}
		for k, v := range txs {
			var script bool // determine if js exists for the Transformation
			if len(v.Script.Body) > 0 {
				script = true
			}
			data = append(data, []string{
				k,
				v.VendorName,
				v.Level,
				fmt.Sprintf("%v", len(v.Fields)),
				fmt.Sprintf("%v", len(v.Configuration)),
				fmt.Sprintf("%v", script),
				fmt.Sprintf("%v", v.IsLegacy),
				v.StartDate,
			})
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Resource", "Vendor", "Level", "# Fields", "# Configs", "Script", "Legacy", "Start Date"})
		table.SetBorder(false)
		table.AppendBulk(data)
		table.Render()
	},
}

// importElementCmd is a command to import an Element json
var importElementCmd = &cobra.Command{
	Use:   "import <path_to_element_json>",
	Short: "Import an Element json to the Platform",
	Long:  "Given an Element json, import into the Platform",
	Run: func(cmd *cobra.Command, args []string) {
		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if len(args) < 1 {
			fmt.Println("Please provide an element json")
			cmd.Help()
			os.Exit(1)
		}
		// read in file
		filebytes, err := ioutil.ReadFile(args[0])
		if err != nil {
			fmt.Println("unable to read file", args[0], err.Error())
			os.Exit(1)
		}
		// validate it's an Element
		var e ce.Element
		err = json.Unmarshal(filebytes, &e)
		if err != nil {
			fmt.Println("JSON doesn't appear to be an Element")
			os.Exit(1)
		}

		bodybytes, statuscode, curlcmd, err := ce.ImportElement(profilemap["base"], profilemap["auth"], e)
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
			log.Printf("%s", bodybytes)
			os.Exit(1)
		}
		fmt.Printf("Element %s imported", e.Name)
	},
}

var forROI bool

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
		if len(args) > 0 && args[0] != "" {
			filteredElementsBytes, err := ce.FilterElementFromList(args[0], bodybytes)
			if err != nil {
				log.Printf("Unable to filter by '%s'- %s\n", args[0], err.Error())
			}
			bodybytes = filteredElementsBytes
		}
		// handle global options, json
		if outputCSV {
			ce.OutputElementsTableAsCSV(bodybytes, orderBy, filterBy)
			fmt.Println()
			return
		}
		if outputJSON {
			fmt.Printf("%s\n", bodybytes)
			return
		}
		if forROI {
			intbytes, statuscode, _, err := ce.GetIntelligence(profilemap["base"], profilemap["auth"])
			if err != nil || statuscode != 200 {
				log.Println("Unable to retrieve intelligence - please check your role")
				return
			}
			roibytes, err := output.ElementsForROICalculator(bodybytes, intbytes)
			if err != nil {
				log.Println("Unable to format onto JSON for the ROI Calculator", err.Error())
				return
			}
			fmt.Printf("%s\n", roibytes)
			return
		}
		// output
		ce.OutputElementsTable(bodybytes, orderBy, filterBy)
	},
}

var lbdocsForce bool
var lbdocsVersion string

var elementLBDocsCmd = &cobra.Command{
	Use:    "lbdocs <id|key>",
	Short:  "Output the IBM LoopBack Model document for the Element",
	Long:   "Outputs the IBM LoopBack Model document for the Element",
	Hidden: true,
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

		// Get element LBDocs
		bodybytes, statuscode, curlcmd, err := ce.GetElementLBDocs(profilemap["base"], profilemap["auth"], strconv.Itoa(elementid), lbdocsForce, lbdocsVersion)
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

// elementExportCmd exports an Element given its id
var elementExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export the Element JSON",
	Long:  `Export an Element JSON given the ID of an Element`,
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

		// Get element
		bodybytes, statuscode, curlcmd, err := ce.GetExportElement(profilemap["base"], profilemap["auth"], strconv.Itoa(elementid))
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
		var element interface{}
		err = json.Unmarshal(bodybytes, &element)
		if err != nil {
			fmt.Println("Can't unmarshal")
			os.Exit(1)
		}
		formattedbytes, err := json.MarshalIndent(element, "", "    ")
		if err != nil {
			fmt.Println("Can't format json")
			os.Exit(1)
		}
		fmt.Printf("%s", formattedbytes)
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

var deleteElementCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes an Element by ID",
	Long:  `Given an Element ID, deletes the Element on the Platform`,
	Run: func(cmd *cobra.Command, args []string) {
		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// check for Element ID
		if len(args) < 1 {
			fmt.Println("Please provide an Element ID")
			return
		}
		elementID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Element ID must be a number")
			os.Exit(1)
		}

		bodybytes, statuscode, curlcmd, err := ce.DeleteElement(profilemap["base"], profilemap["auth"], elementID)
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
			fmt.Printf("%s\n", bodybytes)
			// handle this nicely, show error description
		}
		if statuscode == 200 {
			fmt.Printf("Deleted Element ID %v\n", elementID)
		}

	},
}

var cheatsheetsServerPort int

var cheatsheetsCmd = &cobra.Command{
	Use:    "cheatsheets",
	Short:  "Start a cheatsheets server",
	Long:   `Element Cheatsheets server`,
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		bodybytes, statuscode, _, err := ce.GetIntelligence(profilemap["base"], profilemap["auth"])
		if err != nil {
			if statuscode == -1 {
				fmt.Println("Unable to reach CE API. Please check your configuration / profile.")
			}
			fmt.Println(err)
			os.Exit(1)
		}

		if statuscode != 200 {
			fmt.Println("Unable to connect", statuscode)
			os.Exit(1)
		}
		var intelligence ce.Intelligence
		err = json.Unmarshal(bodybytes, &intelligence)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		portstr := fmt.Sprintf(":%v", cheatsheetsServerPort)

		log.Printf("Starting Element Cheatsheet server on port %v ... Please open a web browser\n", cheatsheetsServerPort)
		r := mux.NewRouter()
		r.HandleFunc("/", handleCheatsheetIndex(intelligence))
		r.HandleFunc("/{id}", handleCheatsheet(profilemap))
		http.Handle("/", r)
		err = http.ListenAndServe(portstr, nil)
		if err != nil {
			fmt.Println("Unable to start cheatsheet server on port", cheatsheetsServerPort)
			os.Exit(1)
		}
	},
}

func handleCheatsheet(profilemap map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		url := fmt.Sprintf("%s/elements/%s/cheat-sheet", profilemap["base"], id)
		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			http.Error(w, "Can't form url", http.StatusInternalServerError)
			return
		}
		req.Header.Add("Authorization", profilemap["auth"])
		req.Header.Add("Accept", "text/html")
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Can't get cheat sheet", http.StatusNotFound)
			return
		}
		bodybytes, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		w.Header().Set("content-type", "text/html")
		w.Write(bodybytes)
	}
}

func handleCheatsheetIndex(intelligence ce.Intelligence) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.New("foo").Parse(output.CheatsheetIndexTemplate)
		err = t.Execute(w, intelligence)
		if err != nil {
			log.Println(err)
		}
		//w.Write([]byte("hi"))
	}
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
	elementsCmd.PersistentFlags().BoolVarP(&outputCSV, "csv", "", false, "output as CSV")

	elementsCmd.AddCommand(listElementsCmd)
	listElementsCmd.Flags().BoolVarP(&forROI, "roi", "", false, "Output as JSON for ROI calculator")
	elementsCmd.AddCommand(elementMetadataCmd)
	elementsCmd.AddCommand(elementDocsCmd)
	elementsCmd.AddCommand(elementInstancesCmd)
	elementsCmd.AddCommand(elementExportCmd)
	elementsCmd.AddCommand(importElementCmd)
	//elementsCmd.AddCommand(elementModelValidation)
	elementsCmd.AddCommand(elementLBDocsCmd)
	elementLBDocsCmd.Flags().BoolVarP(&lbdocsForce, "force", "", false, "Force refresh of current LBDocs version")
	elementLBDocsCmd.Flags().StringVar(&lbdocsVersion, "version", "", "LBDocs specific version")

	elementsCmd.AddCommand(transformationsForElementCmd)
	elementsCmd.AddCommand(cheatsheetsCmd)
	cheatsheetsCmd.Flags().IntVarP(&cheatsheetsServerPort, "port", "p", 12001, "optional port for cheatsheetserver")
	elementsCmd.AddCommand(deleteElementCmd)

	// order-by flag: Order element list by
	// --order-by key|name|id|hub
	listElementsCmd.Flags().StringVarP(&orderBy, "order", "", "", "order element (hub, name)")
	// filter-by flag: Show only elements where filter is true
	// --filter-by active|deleted|private|beta|cloneable|extendable
	// TODO this isn't implemented in ce.elements#OutputElementsTable
	//listElementsCmd.Flags().StringVarP(&filterBy, "filter", "", "", "elements where filter is true")

}

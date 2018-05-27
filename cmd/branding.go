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
	"os"

	"github.com/ghchinoy/ce-go/ce"
	"github.com/spf13/cobra"
)

// brandingCmd represents the root command for managing Platform brandking
var brandingCmd = &cobra.Command{
	Use:   "branding",
	Short: "Manage Branding of the Platform",
	Long:  `Manage Branding of the Platform`,
}

// getBrandingCmd is the cli command to return branding for the Platform profile
// usage:
// branding get
var getBrandingCmd = &cobra.Command{
	Use:   "get",
	Short: "Retrieve current Platform branding",
	Long:  "Returns JSON representation of current Platform branding",
	Run: func(cmd *cobra.Command, args []string) {
		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// Get branding definition for Platform
		bodybytes, statuscode, curlcmd, err := ce.GetBranding(profilemap["base"], profilemap["auth"], debug)
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
		if statuscode == 404 {
			fmt.Println("No branding on this account.")
			os.Exit(0)
		}
		if statuscode != 200 {
			log.Printf("HTTP Error: %v\n", statuscode)
			// handle this nicely, show error description
			fmt.Printf("%s\n", bodybytes)
		}
		var pretty bytes.Buffer
		err = json.Indent(&pretty, bodybytes, "", "  ")
		if err != nil {
			fmt.Printf("%s\n", bodybytes)
			return
		}
		fmt.Printf("%s\n", pretty.Bytes())
	},
}

var brandingJSONFile string

// setBrandingCmd is the cli command to set branding for the Platform profile
// there are two options, "set <attribute> <value>" where a specific attribute
// can be set to a specific attribute
// or the flag "--file" which will read in a JSON file to set the branding
// usage:
// branding set [<attribute> <value>] [--file <branding.json>]
var setBrandingCmd = &cobra.Command{
	Use:   "set [<attribute> <value>] [--file <branding.json>]",
	Short: "Sets Platform branding either by attribute or by file",
	Long: `Sets Platform branding one of two ways, either by providing
an attribute and corresponding value, like so:
set branding <attribute> <value>
or by specifying a JSON file containing branding configuration like so:
set branding --file <branding.json>
`,
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// check for file
		if brandingJSONFile != "" {
			// read in file
			filebytes, err := ioutil.ReadFile(brandingJSONFile)
			if err != nil {
				fmt.Println("unable to read file", brandingJSONFile, err.Error())
				os.Exit(1)
			}
			var brandingobject interface{}
			err = json.Unmarshal(filebytes, &brandingobject)
			if err != nil {
				fmt.Println("Unable to parse JSON file", err.Error())
				os.Exit(1)
			}
			// invoke ce branding API with file contents as json
			bodybytes, statuscode, curlcmd, err := ce.SetBranding(profilemap["base"], profilemap["auth"], brandingobject, debug)
			if err != nil {
				log.Println(err.Error())
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
				fmt.Printf("%s\n", bodybytes)
			}
			// TODO return response

			return
		}

		// check for attribute & value
		if len(args) < 2 {
			cmd.Usage()
			os.Exit(1)
		}
		// TODO construct JSON object
		var brandingobject interface{}

		// TODO determine if existing branding, add modification to that
		// TODO if not, utilize ce.DefaultBranding object

		// TODO accommodate for image values by reading in and base64 encodding
		// attributes: favicon, logo
		// should be pngs
		if args[0] == "favicon" || args[0] == "logo" {

		}

		// invoke ce branding API with JSON object
		bodybytes, statuscode, curlcmd, err := ce.SetBranding(profilemap["base"], profilemap["auth"], brandingobject, debug)
		if err != nil {
			log.Println(err.Error())
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
			fmt.Printf("%s\n", bodybytes)
		}
		// TODO return response

	},
}

var resetBrandingCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset Platform branding to default",
	Long:  "Resets Platform branding to Cloud Elements default theme",
	Run: func(cmd *cobra.Command, args []string) {
		// check for profile
		profilemap, err := getAuth(profile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// invoke reset branding
		bodybytes, statuscode, curlcmd, err := ce.ResetBranding(profilemap["base"], profilemap["auth"], debug)
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
		if statuscode == 404 {
			fmt.Println("No branding on this account.")
			os.Exit(0)
		}
		if statuscode != 200 {
			log.Printf("HTTP Error: %v\n", statuscode)
			// handle this nicely, show error description
			fmt.Printf("%s\n", bodybytes)
		}
		var pretty bytes.Buffer
		err = json.Indent(&pretty, bodybytes, "", "  ")
		if err != nil {
			fmt.Printf("%s\n", bodybytes)
			return
		}
		fmt.Printf("%s\n", pretty.Bytes())
		fmt.Printf("%v\n", statuscode)
	},
}

func init() {
	RootCmd.AddCommand(brandingCmd)
	brandingCmd.AddCommand(getBrandingCmd)
	brandingCmd.AddCommand(setBrandingCmd)
	brandingCmd.AddCommand(resetBrandingCmd)

	brandingCmd.PersistentFlags().StringVar(&profile, "profile", "default", "profile name")
	brandingCmd.PersistentFlags().BoolVarP(&outputJSON, "json", "j", false, "output as json")
	brandingCmd.PersistentFlags().BoolVarP(&showCurl, "curl", "c", false, "show curl command")
	brandingCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "print debug info")

	setBrandingCmd.PersistentFlags().StringVar(&brandingJSONFile, "file", "", "branding configuration json file")
}

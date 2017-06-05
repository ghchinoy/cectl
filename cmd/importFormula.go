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
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"encoding/json"

	"github.com/ghchinoy/cectl/ce"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var importFormulaPath string

// importFormulaCmd represents the importFormula command
var importFormulaCmd = &cobra.Command{
	Use:   "import <filepath>",
	Short: "Imports a Formula to the platform",
	Long:  `Providing a Formula JSON, this command adds a Formula template to the platform`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("must supply a path to a Formula JSON")
			os.Exit(1)
		}

		if !viper.IsSet(profile + ".base") {
			fmt.Println("Can't find info for profile", profile)
			os.Exit(1)
		}

		// read in file
		filebytes, err := ioutil.ReadFile(args[0])
		if err != nil {
			fmt.Println("unable to read file", args[0], err.Error())
			os.Exit(1)
		}

		// Check if can decode into formula struct
		var f ce.Formula
		err = json.Unmarshal(filebytes, &f)
		if err != nil {
			fmt.Println(args[0], "doesn't seem like a Formula")
			os.Exit(1)
		}

		base := viper.Get(profile + ".base")
		user := viper.Get(profile + ".user")
		org := viper.Get(profile + ".org")

		url := fmt.Sprintf("%s%s",
			base,
			"/formulas",
		)
		auth := fmt.Sprintf("User %s, Organization %s", user, org)

		client := &http.Client{}
		req, err := http.NewRequest("POST", url, bytes.NewReader(filebytes))
		if err != nil {
			fmt.Println("Can't construct request", err.Error())
			os.Exit(1)
		}
		req.Header.Add("Authorization", auth)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Cannot process response", err.Error())
			os.Exit(1)
		}
		bodybytes, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			fmt.Println(resp.Status)
			fmt.Printf("%s\n", bodybytes)
		}

		if outputJSON {
			fmt.Printf("%s\n", bodybytes)
			return
		}

		if resp.StatusCode == 200 {
			fmt.Println("Formula template added to Platform.")
			var f ce.Formula
			err = json.Unmarshal(bodybytes, &f)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			err = ce.FormulaDetailsTableOutput(f)
			if err != nil {
				fmt.Println("Unable to render Formula details")
				os.Exit(1)
			}
		}

	},
}

func init() {
	formulasCmd.AddCommand(importFormulaCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importFormulaCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importFormulaCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	//importFormulaCmd.Flags().StringVarP(&importFormulaPath, "file", "f", "", "file path location of Formula JSON to import")
}

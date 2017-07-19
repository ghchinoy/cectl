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
	"os"

	"github.com/ghchinoy/cectl/ce"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// formulaDetailsCmd represents the formulaDetails command
var formulaDetailsCmd = &cobra.Command{
	Use:   "details <id>",
	Short: "Output details of a Formula template",
	Long:  `Given a Formula ID, print out details`,
	Run: func(cmd *cobra.Command, args []string) {

		//formulaformat := "/formulas/%s"

		if len(args) < 1 {
			fmt.Println("must supply an ID of a Formula")
			os.Exit(1)
		}

		if !viper.IsSet(profile + ".base") {
			fmt.Println("Can't find info for profile", profile)
			os.Exit(1)
		}

		base := viper.Get(profile + ".base")
		user := viper.Get(profile + ".user")
		org := viper.Get(profile + ".org")

		/*
			url := fmt.Sprintf("%s%s",
				base,
				fmt.Sprintf(formulaformat, args[0]),
			)
		*/
		auth := fmt.Sprintf("User %s, Organization %s", user, org)

		bodybytes, statuscode, err := ce.FormulaDetailsAsBytes(args[0], fmt.Sprintf("%s", base), auth)
		if err != nil {
			fmt.Println("unable to retrieve formula", err.Error())
			os.Exit(1)
		}

		if outputJSON {
			fmt.Printf("%s\n", bodybytes)
			return
		}

		if statuscode != 200 {
			fmt.Println(statuscode)
			var ficr ce.FormulaInstanceCreationResponse
			err = json.Unmarshal(bodybytes, &ficr)
			if err != nil {
				fmt.Println("Cannot process response, tried error message")
				os.Exit(1)
			}
			fmt.Println(ficr.Message)
			os.Exit(1)
		}

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

	},
}

func init() {
	formulasCmd.AddCommand(formulaDetailsCmd)

}

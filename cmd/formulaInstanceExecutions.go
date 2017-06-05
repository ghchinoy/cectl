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
	"github.com/spf13/cobra"
)

// formulaInstanceExecutionsCmd represents the formula-instance-executions command
var formulaInstanceExecutionsCmd = &cobra.Command{
	Use:   "executions",
	Short: "Manage Formula Instance Executions",
	Long:  `Manage the Executions of a Formula Instance`,
}

func init() {
	RootCmd.AddCommand(formulaInstanceExecutionsCmd)
	formulaInstanceExecutionsCmd.PersistentFlags().StringVar(&profile, "profile", "default", "profile name")
	formulaInstanceExecutionsCmd.PersistentFlags().BoolVarP(&outputJSON, "json", "j", false, "output as json")
}

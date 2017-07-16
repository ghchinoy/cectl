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

// profilesCmd represents the profile command
var profilesCmd = &cobra.Command{
	Use:   "profiles",
	Short: "Manage profiles",
	Long:  `Add, remove, list profiles to manage Cloud Elements access`,
}

func init() {
	RootCmd.AddCommand(profilesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// profileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// profileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	profilesCmd.PersistentFlags().StringVar(&cfgFile, "config", "", cfgHelp)
	profilesCmd.PersistentFlags().StringVar(&profile, "profile", "default", "profile name")
	// Set bash-completion
	validConfigFilenames := []string{"toml", ""}
	profilesCmd.PersistentFlags().SetAnnotation("config", cobra.BashCompFilenameExt, validConfigFilenames)

}
